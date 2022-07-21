FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -o treehole

FROM alpine

WORKDIR /app

COPY --from=builder /app/treehole /app/
COPY --from=builder /app/data /app/

ENV MODE=production

EXPOSE 8000

ENTRYPOINT ["./treehole"]