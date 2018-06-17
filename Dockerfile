ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates

COPY ./WHMCS-currency-update /usr/loca/bin/WHMCS-currency-update

CMD ["WHMCS-currency-update"]