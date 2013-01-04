# Poller

A simple and composable HTTP monitoring application written in Go.

It's very alpha for now, so please bear with it.

## What is it?

Poller's job is to monitor http application by submitting GET requests to URL
you define in a config file. These URLs are called "Checks".

Once a check is done, the result is sent to one or many backend of your choice. 

Current supported backend are stdout, syslog,
[librato](http://metrics.librato.com/) and [statsd](://github.com/etsy/statsd).

## How to build/install it

As the software is alpha now, I'll assume that you know how to build go
softwares. For the others, I'll try to provide binaries soon enough once
stability has been proven for 1 check to 1000s of checks.

If you want to run it on heroku, I've made a sample repository for that
available here:
[https://github.com/marcw/poller-heroku](https://github.com/marcw/poller-heroku)

## How to configure it

Configuration is done is a json file (for now). A sample configuration file
looks like this:

    {
        "timeout": "5s",
        "backends": ["stdout", "syslog"],
        "checks": [
            {
                "key": "com_google",
                "url": "http://google.com",
                "interval": "10s"
            },
            {
                "key": "fr_yahoo",
                "url": "http://yahoo.fr",
                "interval": "10s"
            }
        ]
    }

This config file defines 2 backends and 2 checks. These two checks will be
executes every 10 seconds. Configuration for the backend is achieved with
environment variables. The connection timeout is set to 5 second.

The `key` is the identifier that will be used in the output.

## How to run it

Running poller is really simple and is a metter of setting a few environment
variables (if needed) and passing the binary the path of the configuration
file.

### Backends configuration

Here is a list of supported backend and how to configure them with environment
variables.

#### Librato

- `LIBRATO_USER` (required): Username of your librato account
- `LIBRATO_TOKEN` (required): API token of your librato account
- `LIBRATO_SOURCE` (optional): Source name for your metrics. Defaults to `poller`

#### Stdout

- No configuration is necessary.

#### Statsd

Statsd backend uses of these environment variables:

- `STATSD_HOST` (required): Host of your statsd instance
- `STATSD_PORT` (optional): Port of your statsd instance. Defaults to 8125.
- `STATSD_PROTOCOL` (optional): Either `tcp` or `udp`. Defaults to `udp`.

#### Syslog

You can configure the syslog backend with these:

- `SYSLOG_NETWORK` (optional): "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
  "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6"
  (IPv6-only), "unix" and "unixpacket". Defaults to nothing.
- `SYSLOG_ADDRESS` (optional): Address of your syslog daemon. Defaults to nothing.
- `SYSLOG_PREFIX` (optional): This will be added to your log. Defaults to "poller".

## What is buggy or likely to be changed/added in future release?

- Metrics might not be stored in a rightful way. Feel free to provide feedback.
- I'm currently not sure if I'll keep using json as a configuration format as
  [s-expr](http://en.wikipedia.org/wiki/S-expression) are a much better thing
  (and they supports COMMENTS)
- Per-check configuration for each backend (i.e: the statsd metric name)
- Poller will need to supports live-updating of the checks list
- Maintenance mode support.
- Customization of the user agent.

## License

The Uptime code is free to use and distribute, under the [MIT license](https://github.com/marcw/poller/blob/master/LICENSE).
