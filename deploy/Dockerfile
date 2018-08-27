from alpine
RUN apk --no-cache add curl
COPY wavefront-adapter /
ENTRYPOINT ["/wavefront-adapter", "--logtostderr=true"]
