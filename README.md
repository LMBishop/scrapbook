# scrapbook

Scrapbook is a website manager, built to deploy and serve statically generated
web pages. It is designed to be as easy to use as possible, with a web interface
for management and API for programmatic site deployment. It also does not need a
database and uses the file system for organisation, to allow easy interoperability
if you wish to serve sites using another web server.

It was originally built for me to use in conjunction with my own
[static site generator](https://github.com/LMBishop/panulat) in a GitHub workflow.

## Installation

This program is designed to work on any Linux machine. There is a provided Makefile
with `install`, `install-config`, and `install-service` targets. To install the
program itself:

```bash
make
make install
```

If this is the first installation, then you may also want to install the default
configuration with the `install-config` target.

```bash
make install-config
```

If you are on a systemd distribution, then there is also a provided service file which
the `install-service` target will install. This target will also create a `scrapbook`
user on the system.

```bash
make install-service
```

## Configuration

Scrapbook will look for its configuration in `/etc/scrapbook/` by default. You can
run the `install-config` target to install the provided [default configuration](https://github.com/LMBishop/scrapbook/blob/master/dist/config.toml).

You must set a hostname and secret for the web management interface and API. (I
collectively call these the 'Command' interfaces, as it is the way you issue
commands to scrapbook.)

```toml
[Command]
Host = ''
Secret = ''
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
