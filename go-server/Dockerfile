# Use official Go image
FROM golang:1.24 AS builder

WORKDIR /app

# Download go mod dependencies first for cache efficiency
COPY go.mod go.sum ./
RUN go mod download

# Copy go code
COPY . .

# Optional: Static build (no C libs needed)
RUN CGO_ENABLED=0 go build -o main .


# --- Final minimal image ---
FROM debian:bullseye-slim

WORKDIR /app

# Add certificates & curl (optional)
RUN apt-get update && apt-get install -y ca-certificates curl unzip postgresql-client && rm -rf /var/lib/apt/lists/*

# Install Flyway manually
RUN curl -L https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/11.7.0/flyway-commandline-11.7.0-linux-x64.tar.gz | tar xvz && mv flyway-11.7.0 /flyway

ENV PATH="/flyway:$PATH"

# Copy compiled Go binary from build stage
COPY --from=builder /app/main .


# Copy only migration files
COPY db/migrations ./db/migrations

# Copy Flyway config for prod (url/user/pass via ENV vars)
COPY flyway.prod.conf ./flyway.conf

# Expose port
EXPOSE 8080

# Run migrations + start app
CMD flyway migrate && ./main
