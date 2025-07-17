# scrapbook

Scrapbook is a website manager, built to deploy and serve statically generated
web pages. It is designed to be as easy to use as possible, with a web interface
for management and API for programmatic site deployment. It also does not need a
database and uses the file system for organisation, to allow easy interoperability
if you wish to serve sites using another web server.

It was originally built for me to use in conjunction with my own
[static site generator](https://github.com/LMBishop/panulat) in a GitHub workflow.

## Example

```
curl -X POST -H "Authorization: Bearer (token)" -F upload=@file.zip https://publish.example.com/api/site/:site/upload
```
