FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/crowdreview ./cmd/api

FROM alpine:3.19
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/crowdreview /app/crowdreview
EXPOSE 8080
CMD ["/app/crowdreview"]
