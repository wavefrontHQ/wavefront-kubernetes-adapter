from scratch
COPY wavefront-adapter /
COPY ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/wavefront-adapter", "--logtostderr=true"]
