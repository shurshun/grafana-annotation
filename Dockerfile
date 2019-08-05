FROM golang:1.12.7 AS BUILD

WORKDIR /app

ENV GO111MODULE=auto
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum main.go ./

RUN go mod vendor

RUN go build -o annotation-poster .

FROM alpine

RUN \
    apk add --no-cache --update \
        ca-certificates \
        tzdata

COPY --from=BUILD /app/annotation-poster /bin/

ENTRYPOINT ["annotation-poster"]
