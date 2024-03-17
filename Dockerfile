# syntax=docker/dockerfile:1

FROM golang:1.22 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY ./cmd ./cmd
COPY ./pkg ./pkg

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /sealet-secret-cert-hatch cmd/sealet-secret-cert-hatch.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

FROM scratch

COPY --from=builder /sealet-secret-cert-hatch /sealet-secret-cert-hatch
# Run
CMD ["/sealet-secret-cert-hatch"]