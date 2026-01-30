# Stage 1: Base image
FROM golang:1.25.6 AS base

# Install air
RUN go install github.com/air-verse/air@latest

# Copy all the folders
WORKDIR /app
COPY . .


# Run Server
FROM base AS server
CMD ["make", "server"]

# Run Database
FROM base AS database
CMD ["make", "database"]

# Run Server
FROM base AS kafka
CMD ["make", "kafka"]

FROM base AS all
CMD ["make", "all", "-j", "4"]

