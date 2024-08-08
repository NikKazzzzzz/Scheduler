FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY cmd/ .
COPY config/scheduler.yaml ./config/scheduler.yaml

RUN go mod tidy
RUN go build -o scheduler

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/scheduler .
COPY --from=builder /app/config/scheduler.yaml ./config/scheduler.yaml

CMD ["./scheduler"]
