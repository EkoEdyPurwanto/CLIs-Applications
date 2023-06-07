### Build Stage ###
FROM golang:1.19.1-alpine3.16 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary Go modules files
COPY go.mod go.sum ./

# Download and tidy Go dependencies
RUN go mod download && go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

### Final Image Stage ###
FROM alpine:3.16 AS app

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the build stage
COPY --from=build /app/app .
COPY --from=build /app/internal/database/migrations/ ./internal/database/migrations/

# Remove any unnecessary tools or dependencies from the final stage
# (You can customize this based on your application's requirements)

# Set the command to run when the container starts
CMD ["./app", "worker"]
