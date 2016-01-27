# Netuitive StatsD Agent
==========================
A Docker image and example code for testing tagged metrics with Netuitive's StatsD backend.

## To Use:

### 1. Build and run the Docker container
**1.1** Clone this repo and change into the netuitive-statsd-agent directory.

**1.2** Build the Docker image: 

```sh
$ sudo docker build --rm=true -t netuitive-statsd-agent .
```

**1.3** Run a container from this image:

```sh
$ sudo docker run -d --name netuitive-statsd-agent -v /proc:/host_proc:ro -v /var/run/docker.sock:/var/run/docker.sock:ro -p 8125:8125/udp -d netuitive-statsd-agent
```

This will run the agent inside the container, exposing UDP port 8125. You'll need to determine your
container's IP address in order to send metrics to it.

### 2. Configure StatsD

The **config.js** file will need to be modified to ensure a valid API key for Netuitive is in use:

```js
{
    backends:["./backends/netuitive"],
    netuitive: {
        apiKey: "<API-KEY>",	<------ PUT YOUR API KEY HERE!
        apiHost: "api.app.netuitive.com",
        apiPort: 443,
        mappings: [
            {
      pattern: "(test.rob-egan)\\.(head.requests)\\.counter",
      element: {
        type: "SERVER",
        name: "$1",
        metric: {
          name: "$2"
        }
      }
    },
...
```
> You can obtain the API Key by viewing the StatsD datasource you wish to use in the Netuitive UI

### 3. Run the example code.

Two (overly simplified) code snippets have been provided that will create generic web services that will
generate metrics and send them to the StatsD agent.

**3.1** Python example

This example assumes you have Python 2.6 or greater installed, with the time, BaseHTTPServer, and statsd
modules installed.

```sh
$ ./tests/webserver.py &
```
This will run a web server, listening on port 8000. You can hit the URL **http://localhost:8000** in order
to begin generating metrics.

**3.2** Golang example

This example assumes you have Golang installed, and uses a thrid party Statsd client library [quipo/statsd]

```sh
$ go run ./tests/webserver.go &
```

This will run a webserver, listening on port 8000. Hit the URL http://localhost:8000/help for details on what 
paths are valid. Each valid path will generate metrics when hit.

**3.3** Golang/Datadog Example

This example requires Golang plus the datadog-go statsd libraries. It was added to show a comparison between
Datadog and Netuitive with regards to adding tags within the code being used to instrument metrics via StatsD.

```sh
$ go run ./tests/webserver_dd.go &
```
> This example also assumes that a datadog-statsd agent is running.

This will create a web service on localhost, listening on TCP port 8000. The URL Path "/help" will explain
details on what valid paths can be access, each of which produce various metrics.
