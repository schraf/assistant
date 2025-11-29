FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/service ./cmd/service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/job ./cmd/job

FROM gcr.io/distroless/static-debian12:nonroot

# Copy both binaries
COPY --from=builder /app/bin/service /service
COPY --from=builder /app/bin/job /job

EXPOSE 8080

# Default entrypoint for Cloud Run service
ENTRYPOINT [ "/service" ]
