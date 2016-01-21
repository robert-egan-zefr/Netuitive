FROM ubuntu:14.04

MAINTAINER Robert Egan <robert.egan@zefr.com>

USER root

# Update/upgrade the OS and install dependencies
RUN apt-get -y update\
 && apt-get -y upgrade\
 && apt-get -y --force-yes install\
 tcpdump\
 lsof\
 supervisor\
 git\
 nodejs\
 && apt-get clean\
 && rm -rf /var/lib/apt/lists/* /var/tmp/*

# Install and configure StatsD
RUN git clone -b v0.7.2 --depth 1 https://github.com/etsy/statsd.git /opt/statsd
COPY conf/statsd/config.js /opt/statsd/config.js
COPY conf/statsd/statsd.conf /etc/supervisor/conf.d/statsd.conf

# Install and configure Netuitive StatsD backend
WORKDIR /tmp
RUN git clone --depth 1 https://github.com/Netuitive/statsd-netuitive-backend.git\ 
 && mv statsd-netuitive-backend/netuitive/ /opt/statsd/backends/netuitive\
 && mv statsd-netuitive-backend/netuitive.js /opt/statsd/backends/netuitive.js\
 && rm -rf /tmp/*

# Make sure supervisor kicks off from the statsd directory
WORKDIR /opt/statsd

# Start Supervisor, which will now manage statsd
CMD ["supervisord", "-n"]
