# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
RUN apk add --no-cache git

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s" -o app  .

# Final stage
FROM scratch

COPY --from=builder /build/app /app
COPY --from=builder /usr/bin/wget /bin/wget
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /lib/ld-musl-* /lib/
COPY --from=builder /lib/libc.* /lib/

ENV PORT=8080
HEALTHCHECK --interval=5s --timeout=3s \
CMD ["wget", "-qO-", "http://localhost:8080/health"]

EXPOSE 8080
ENTRYPOINT ["/app"]
