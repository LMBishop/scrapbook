FROM golang:1.25 AS build
WORKDIR /build
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN make SYS_CONF_DIR='/etc/scrapbook' SYS_DATA_DIR='/var/lib/scrapbook'

FROM scratch
COPY --from=build /build/scrapbook /bin/scrapbook
EXPOSE 80
CMD ["/bin/scrapbook"]
