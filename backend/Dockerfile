# -- Build --
FROM golang:1.26-alpine AS builder

WORKDIR /app

# CGO=0 because alpine doesn't have the c libs
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build exe
RUN go build -trimpath -ldflags="-s -w" -o oth ./cmd/api/

# -- Runtime --
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/oth .

EXPOSE 4000
CMD ["./oth"]