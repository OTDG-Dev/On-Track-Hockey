ARG GOLANG_VERSION=1.26-alpine
ARG ALPINE_VERSION=3.23

# =====================
# Stage 1: Build
# =====================
FROM golang:${GOLANG_VERSION} AS builder

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

# =====================
# Stage 2: Runtime
# =====================
FROM alpine:${ALPINE_VERSION}

WORKDIR /app
COPY --from=builder /app/oth .

EXPOSE 4000
CMD ["./oth"]