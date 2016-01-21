// webserver.go - Web server for playing with various statsd libraries
//              - This server will have some statsd metrics built into it.
//              - Various URL paths will do some random stuff, i.e.:
//              - http://localhost:8000/count displays the requests count
//              - http://localhost:8000/debug shows header/request info
//              - http://localhost:8000/lissajous displays 801x801 Lissajous figures
//              - http://localhost:8000/help describes how to use the site
//
// Author: Rob Egan
// Updated: January 7, 2016
//
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd" // Data Dog import
)

var (
	mu      sync.Mutex
	count   int
	palette = []color.Color{color.White, color.Black}
	dd      *statsd.Client
	err     error
)

const (
	whiteIndex    = 0
	blackIndex    = 1
	ddAgentServer = "192.168.99.100:8125" // IP/hostname and port of the DataDog Agent
	webServerHost = "localhost:8000"      // IP/hostname and port of the web server
)

func main() {
	// Connect to the statsd/datadogd server
	dd, err = statsd.New(ddAgentServer)
	if err != nil {
		log.Fatalf("Error connecting to Datadog Agent: %s", err.Error())
	}

	// The default prefix for my test metrics
	dd.Namespace = "test.rob-egan."

	// Handler functions
	http.HandleFunc("/", handler)
	http.HandleFunc("/lissajous", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	})
	http.HandleFunc("/debug", debug)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/help", help)

	// The actual web server
	log.Print("Starting web server...\n")
	log.Print("Point your browser to http://", webServerHost, "/help for more info...\n")
	log.Fatal(http.ListenAndServe(webServerHost, nil))
}

// Handler to display 'usage' options for the web service
func help(w http.ResponseWriter, r *http.Request) {
	hit() // Incement the 'total hits' metric
	fmt.Fprintf(w, "Usage -  the following URL paths are available:\n")
	fmt.Fprintf(w, "      - \"/help\" to print this usage page...\n")
	fmt.Fprintf(w, "      - \"/debug\" to display request header info...\n")
	fmt.Fprintf(w, "      - \"/count\" to see the site's hit counter...\n")
	fmt.Fprintf(w, "      - \"/lissajous\" to see Lissajous figures...\n")
	dd.Count("help.pageview.count", 1, nil, 1) // Metric to record hits per cycle for this handler
}

// handler to echo back the Path component of the requested URL. This is the default handler.
func handler(w http.ResponseWriter, r *http.Request) {
	hit() // Increment 'total hits' metric
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	dd.Count("default.pageview.count", 1, nil, 1) // Metric to record hits per cycle for this handler
}

// debug handler function prints elements of the HTTP request including:
// method, url, protocol, headers, form data, and more
func debug(w http.ResponseWriter, r *http.Request) {
	hit() // Increment 'total hits' metric
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
	dd.Count("debug.pageview.count", 1, nil, 1) // Metric to record hits per cycle for this handler
}

// counter handler echoes the number of non-count requests so far
func counter(w http.ResponseWriter, r *http.Request) {
	hit() // Increment 'total hits' metric
	mu.Lock()
	fmt.Fprintf(w, "Count: %d\n", count)
	mu.Unlock()
	dd.Count("counter.pageview.count", 1, nil, 1) // Metric to record hits per cycle for this handler
}

// handler that displays lissajous figures in the browser
func lissajous(out io.Writer) {
	hit() // Increment 'total hits' metric
	l_start := time.Now()
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 400   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
	l_secs := time.Since(l_start).Seconds()
	dd.Histogram("lissajous.load.time", l_secs, nil, 1) // Metric to record time to load a figure
	dd.Count("lissajous.pageview.count", 1, nil, 1)     // Metric to record hits per cycle for this handler
}

// Increment the 'hits' counter anytime a request is received
func hit() {
	dd.Count("total.pageview.count", 1, nil, 1)
}
