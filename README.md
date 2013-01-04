# Poller

A simple and composable HTTP monitoring application written in Go.

It's very alpha for now, so please bear with it.

## What is it?

Poller's job is to monitor http applications by submitting GET requests to URL
you define in a config file. These URLs are called "Checks".

Once a check is done, the result is sent to one or many backends of your choice. 

Current supported backends are stdout, syslog,
[librato](http://metrics.librato.com/) and [statsd](://github.com/etsy/statsd).

## How to build/install it

As the software is alpha now, I'll assume that you know how to build Go
softwares. For the others, I'll try to provide binaries soon enough once
stability has been proven for 1 check to 1000s of checks.

If you want to run it on heroku, I've made a sample repository for that matter,
available there:
[https://github.com/marcw/poller-heroku](https://github.com/marcw/poller-heroku)

## How to configure it

Configuration is done is a json file. A sample configuration file
looks like this:

    {
        "timeout":  "5s",
        "userAgent": "Poller (https://github.com/marcw/poller)",
        "backends": ["stdout", "librato"],
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
            },
            {
                "key": "connect_sensiolabs_com_api",
                "url": "https://connect.sensiolabs.com/api/",
                "interval": "60s",
                "headers": {
                    "Accept": "application/vnd.com.sensiolabs.connect+xml"
                }
            }
        ]
    }

This config file defines 2 backends and 3 checks.
- The connection timeout is 5s
- The User agent is set to "Poller (https://github.com/marcw/poller)"
- Two backends will be used: stdout and librato
- `key` is the identifier used by backends.
- The first two checks (`com_google` and `fr_yahoo`) will be checked every 10s
- The third check will be checked every 60s
- When checking the third check,the header `Accept:
  application/vnd.com.sensiolabs.connect+xml` will be sent.

## How to run it

Running poller is really simple and is a matter of setting a few environment
variables (if needed) and passing the binary the path of the configuration
file.

## How to monitor it

Before launching poller, export `POLLER_URL` environment variable with something
like this

    export POLLER_URL="localhost:8000"

and poller will launch a webserver. You can then poll
"http://localhost:8000/health" and check if poller is still running

### Backends configuration

Here is a list of supported backend and how to configure them with environment
variables.

#### Librato

- `LIBRATO_USER` (required): Username of your librato account
- `LIBRATO_TOKEN` (required): API token of your librato account
- `LIBRATO_SOURCE` (optional): Source name for your metrics. Defaults to `poller`

Given your check's key is `foobar`:

On success:
  - `1` will be sent to metrics `foobar.up`
  - The duration in milliseconds will be sent to `foobar.duration`

On error:
  - `0` will be sent to metrics `foobar.up`
  - The duration in milliseconds will be sent to `foobar.duration`

On timeout:
  - `0` will be sent to metrics `foobar.up`

#### Stdout

No configuration is necessary.
Output will look like this:

    2012/01/24 11:35:16 com_google 200 345.271ms
    2012/01/24 11:35:17 fr_yahoo 200 1.518175s
    2012/01/24 11:35:16 com_google TIMEOUT


#### Statsd

Statsd backend uses of these environment variables:

- `STATSD_HOST` (required): Host of your statsd instance
- `STATSD_PORT` (optional): Port of your statsd instance. Defaults to 8125.
- `STATSD_PROTOCOL` (optional): Either `tcp` or `udp`. Defaults to `udp`.

The metrics are sent the same way as the Librato backend.

#### Syslog

You can configure the syslog backend with these:

- `SYSLOG_NETWORK` (optional): "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
  "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6"
  (IPv6-only), "unix" and "unixpacket". Defaults to nothing.
- `SYSLOG_ADDRESS` (optional): Address of your syslog daemon. Defaults to nothing.
- `SYSLOG_PREFIX` (optional): This will be added to your log. Defaults to "poller".

Yes, you can send checks results to [loggly](http://www.loggly.com/) by using
the syslog backend.

Output formatting is the same as the stdout backend.

## What is buggy or likely to be changed/added in future release?

- Metrics names for statsd/librato backends might not be named correctly. Feel
  free to provide feedback.
- I'm currently not sure if I'll keep using json as a configuration format as
  [s-expr](http://en.wikipedia.org/wiki/S-expression) are a much better thing
  (and they supports COMMENTS)
- Per-check configuration for each backend ? (i.e: the statsd metric name)
- Poller will need to supports live-updating of the checks list
- Maintenance mode support.

## License

The Uptime code is free to use and distribute, under the [MIT license](https://github.com/marcw/poller/blob/master/LICENSE).
