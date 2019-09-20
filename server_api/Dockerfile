FROM golang:1.13

WORKDIR /go/src/api_avito
COPY ./api_avito .

EXPOSE 9000

RUN go get github.com/lib/pq

ENTRYPOINT [ "go", "run", "api_avito" ]