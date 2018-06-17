ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates

COPY ./WHMCS-currency-update /usr/local/bin/WHMCS-currency-update

CMD ["WHMCS-currency-update"]