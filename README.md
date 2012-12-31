# Poller

A simple and composable HTTP monitoring application written in Go.

It's very alpha for now, so please bear with it.

For now, the only supported backends are "statsd" and "stdout". I'll add more
for sure.

## How to build/install it

As the software is alpha now, I'll assume that you know how to build go
softwares. For the others, I'll try to provide binaries soon enough once
stability has been proven for 1 check to 1000s of checks.

## How to run it

If you use a stastd backend, make sure you provide at least the `STATSD_HOST`
environment variable. You can also provide `STATSD_PORT` and `STATSD_PROTOCOL`
but if you don't, they will respectively default to `8125` and `udp`.

poller accepts only one argument which is the filename of a json file. This
file describes the configuration. You can have a sample of this configuration
file in the project source tree.

If you use the statsd backend, the "key" config file for a check will be used
as the metric key.

## What is likely to change in future release?

- I'm currently not sure if I'll use json as a configuration format as
  [s-expr](http://en.wikipedia.org/wiki/S-expression) are a much better thing
  (and they supports COMMENTS)
- Per-check configuration for each backend (i.e: the statsd metric name)
- Poller will need to supports live-updating of the checks list
