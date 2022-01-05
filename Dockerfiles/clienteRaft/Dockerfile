FROM alpine

COPY --from=golang:1.17.1-alpine /usr/local/go/ /usr/local/go/
 
ENV PATH="/usr/local/go/bin:${PATH}"

RUN set -ex && \
    apk add build-base
RUN apk add vim

COPY cliente /usr/local/bin/cliente
ADD CodigoEsqueleto /

EXPOSE 7000