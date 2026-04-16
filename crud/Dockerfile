FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a \
    -ldflags='-w -s -extldflags "-static"' \
    -installsuffix cgo \
    -o main main.go

FROM golang:1.25-alpine

WORKDIR /app

ENV TZ=Europe/Moscow

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]