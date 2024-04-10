# golang:alpine3.19
FROM golang:alpine@sha256:2523a6f68a0f515fe251aad40b18545155135ca6a5b2e61da8254df9153e3648 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001 
ENV GID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN addgroup --gid $GID appuser
RUN adduser \    
  --disabled-password \    
  --gecos "" \    
  --home "/nonexistent" \    
  --shell "/sbin/nologin" \    
  --no-create-home \    
  --uid "${UID}" \    
  -G "appuser" \
  "${USER}"

RUN mkdir /var/log/lists && touch /var/log/lists/debug.log
RUN chown -R appuser:appuser /var/log/lists/debug.log

WORKDIR /src
COPY . .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

# Build the binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" ./cmd/lists/

FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /var/log/lists/debug.log /var/log/lists/

# Copy our static executable
COPY --from=builder /src/lists /usr/local/bin/lists

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 4322/tcp

ENTRYPOINT ["/usr/local/bin/lists"]