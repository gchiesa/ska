FROM alpine:3.19
COPY ska /usr/bin/ska
ENTRYPOINT ["/usr/bin/ska"]