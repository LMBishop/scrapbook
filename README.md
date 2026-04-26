# scrapbook

Scrapbook is a website manager, built to deploy and serve statically generated
web pages. It is designed to be as easy to use as possible, with a web interface
for management and API for programmatic site deployment. It also does not need a
database and uses the file system for organisation, to allow easy interoperability
if you wish to serve sites using another web server.

It was originally built for me to use in conjunction with my own
[static site generator](https://git.leonardobishop.net/panulat) in a CI job.

## Installation

This program is designed to work on any Linux machine. Install to `/usr/local/bin`
with:

```bash
make
make install
```

There is a sample configuration file and service file in `contrib`. By default,
scrapbook will look for its configuration at `/etc/scrapbook/config`.

## Configuration

You must set a hostname and secret for the web management interface and API. (I
collectively call these the 'Control' interfaces, as it is the way you issue
commands to scrapbook.)

```toml
listen "0.0.0.0:80"

control {
  host   ""
  secret ""
}
```

If either values are left blank, then the web management interface and API will
be inaccessible.

## Practical notes and recommended setup

**TLS.** Scrapbook currently has no support for TLS. I would recommend running it
behind a reverse proxy (I use nginx) and terminating TLS connections there before
passing them to scrapbook.

**Certificates / DNS.** On the topic of certificates, I would recommend getting a
wildcard certificate for the (sub-)domain you want to serve scrapbook sites with.
Couple this with a wildcard CNAME pointing to your webserver, and you can very
easily set up a new sites on different subdomains all within the scrapbook web
management interface.
