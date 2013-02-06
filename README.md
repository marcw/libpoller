# Poller

[![Build Status](https://travis-ci.org/marcw/poller.png?branch=master)](https://travis-ci.org/marcw/poller)

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
stability has been proven for 1 to 1000s of checks.

If you want to run it on heroku, I've made a sample repository for that matter,
available there:
[https://github.com/marcw/poller-heroku](https://github.com/marcw/poller-heroku)

## How to configure it

Configuration is done via:

- Environment variables for backend.
- JSON for list of checks
- Flags passed when executing the program.

A typical json file for checks looks like this

    [
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

The JSON config file is optional as checks can be added thanks to the HTTP endpoint `/checks`.

Running `./poller --help` will prints a list of available options.

## How to monitor it

A `/health` http endpoint is available. If poller is answering a 200, then all
is good!

## How to add checks while poller is running

Poller supports live configuration changes thanks to the `/checks` http endpoint.
Send a `PUT` request with a valid config JSON in the body of the request and poller
will append the checks to its list.

### Backends configuration

Here is a list of supported backend and how to configure them with environment
variables.

#### Librato

- `LIBRATO_USER` (required): Username of your librato account
- `LIBRATO_TOKEN` (required): API token of your librato account
- `LIBRATO_SOURCE` (optional): Source name for your metrics. Defaults to `poller`
- `LIBRATO_PREFIX` (optional): Prefix of your metrics. Defaults to `poller.checks.`

Given your check's key is `foobar` and `LIBRATO_PREFIX` is `acme.`:

On success:
  - `1` will be sent to metrics `acme.foobar.up`
  - The duration in milliseconds will be sent to `foobar.duration`

On error or timeout:
  - `0` will be sent to metrics `acme.foobar.up`
  - The duration in milliseconds will be sent to `foobar.duration`

#### Stdout

No configuration is necessary.
Output will look like this:

    2012/01/24 11:35:16 com_google UP 345.271ms
    2012/01/24 11:35:17 fr_yahoo DOWN 1.518175s


#### Statsd

Statsd backend uses of these environment variables:

- `STATSD_HOST` (required): Host of your statsd instance
- `STATSD_PORT` (optional): Port of your statsd instance. Defaults to 8125.
- `STATSD_PROTOCOL` (optional): Either `tcp` or `udp`. Defaults to `udp`.
- `STATSD_PREFIX` (optional): Prefix of your metrics. Defaults to `poller.checks.`

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

## Why Go?

A Go application has the advantage of being concurrent, fast, cross-compilable
and easily deployable. I think it's a lot of good reasons to create a
monitoring system out of this language.

### Install go

note: OSX users, do NOT install Go from homebrew. At time of writing, the
package is broken and won't let you cross-compile.

Please follow instructions from:
[http://golang.org/doc/install/source](http://golang.org/doc/install/source)

Please set a `$GOROOT` env var in your `~/.bashrc.` file. If you unpacked the
go source in `/home/you/go`:

    export GOROOT=/home/you/go

Please also set a `$GOPATH` env var in your `~/.bashrc` file. `$GOPATH` is
where go will look for packages and where you clone private projects.  An
example of `$GOPATH` would be `/home/you/work/go`.

    export GOPATH=/home/you/work/go

### Enable cross compilation with Go

Follow this really good blogpost:
[http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go](http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go)

## License

The Poller code is free to use and distribute, under the [MIT license](https://github.com/marcw/poller/blob/master/LICENSE).
