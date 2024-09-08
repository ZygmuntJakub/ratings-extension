FROM golang:1.23.0

# Build the app
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/server ./cmd/server.go

# Run the app
CMD ["./build/server"]
