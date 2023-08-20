# scrapbook
A web service which accepts zipped websites and deploys them somewhere (i.e.
unzips them and shoves them into another directory).

Might expand to serving said files as well, but for now it is intended to
drop files into directories served by other web servers, like nginx.

Used in conjunction with my own [static site generator](https://github.com/LMBishop/panulat)
in a GitHub workflow.

## Example

```
curl -X POST -u username:password -F content=@file.zip https://publish.example.com/:site/upload
```
