FROM golang:1.22 as builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/dong

FROM alpine:latest as runner

COPY --from=builder /go/dong /usr/local/bin/dong

EXPOSE 80

ENTRYPOINT ["dong"]