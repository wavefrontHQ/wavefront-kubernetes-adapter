from scratch

COPY _output/amd64/wavefront-adapter-linux wavefront-adapter
COPY ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/wavefront-adapter", "--logtostderr=true"]
