FROM --platform=$BUILDPLATFORM gcr.io/distroless/static:latest
ARG BUILDPLATFORM

WORKDIR /
COPY $BUILDPLATFORM .
ENTRYPOINT ["/wavefront-adapter", "--logtostderr=true"]
