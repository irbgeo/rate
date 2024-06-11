FROM golang:1.22-alpine3.20 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go install /app/cmd/rate

FROM alpine:3.20

COPY --from=builder /go/bin/rate /usr/local/bin/app

ENTRYPOINT ["app"]