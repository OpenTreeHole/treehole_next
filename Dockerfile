FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -tags="release" -ldflags "-s -w" -o treehole

FROM alpine

WORKDIR /app

COPY --from=builder /app/treehole /app/
COPY data data

ENV MODE=production

EXPOSE 8000

ENTRYPOINT ["./treehole"]