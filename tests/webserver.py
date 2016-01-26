#!/usr/bin/env python
#
# webserver.py - Creates a basic web server that responds to HEAD and GET requests.
#              - This is for testing the sending of metrics to Netuitive via statsd
#
# Author: Rob Egan
# Updated: January 14, 2016
#

import time, BaseHTTPServer, statsd

HOST_NAME = "localhost"		# Hostname or IP of Web Server
PORT_NUMBER = 8000		# Web Server Port
STATSD_HOST = "172.17.0.16"	# Hostname or IP of Statsd Server
STATSD_PORT = 8125		# Statsd Server Port

class Handler(BaseHTTPServer.BaseHTTPRequestHandler):
    def do_HEAD(s):
        """Invoked if a HEAD request is received."""
        s.send_response(200)
        s.send_header("Content-type", "text/html")
        s.end_headers()
        # Increment a metric for counting HEAD requests in statsd
        c.incr("head.requests.counter")
    def do_GET(s):
        """Invoked if a GET request is received."""
        s.send_response(200)
        s.send_header("Content-type", "text/html")
        s.end_headers()
        s.wfile.write("<p>URL Path: %s</p>" % s.path)
        # Increment a metric for counting GET requests in statsd
        c.incr("get.requests.counter")

if __name__ == '__main__':
    #Initialize the connection to statsd host
    c = statsd.StatsClient(STATSD_HOST, STATSD_PORT, prefix='test.rob-egan')

    # Initialize and start the web server, stopping via keyboard interrupt
    server_class = BaseHTTPServer.HTTPServer
    httpd = server_class((HOST_NAME, PORT_NUMBER), Handler)
    print time.asctime(), "Web Server Startup - %s:%s" % (HOST_NAME, PORT_NUMBER)
    print "Point browser to http://%s:%s" % (HOST_NAME, PORT_NUMBER)
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()
    print time.asctime(), "Web Server Stopped - %s:%s" % (HOST_NAME, PORT_NUMBER) 
