FROM alpine:3.12.0

RUN apk update && \
apk add --no-cache tar qemu-img qemu-block-curl xen-libs libaio \
capstone skopeo libc6-compat

WORKDIR /usr/bin
ADD kubevirt-image-service-exporter .

ENTRYPOINT ["./kubevirt-image-service-exporter"]