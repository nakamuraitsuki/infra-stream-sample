# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simulator ./cmd

# Runtime Stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/simulator /simulator

ENV SIMULATOR_USERS=10
ENV SIMULATOR_BASE_URL=http://frontend

ENTRYPOINT ["/simulator"]