# -- Stage 1: Build the Go Binary---
FROM golang:alpine AS builder

WORKDIR /app

#Copy dependency files first to leverage docker layer caching
COPY go.mod go.sum ./
RUN go mod download 

#copy the rest of your source code 
COPY . .

#Compile a statistic linux binary with debug symbols stripped (-ldflags="-w -s")
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o scraper-api main.go

#-- Stage 2: create tiny productuin image
FROM alpine:latest

# Install SSL certificates so your scraper can make HTTPS requests (e.g. https://github.com)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy ONLY the compiled binary from Stage 1
COPY --from=builder /app/scraper-api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the compiled binary
CMD ["./scraper-api"]