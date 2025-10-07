FROM golang:1.25

# Install tools with correct repositories
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate docs
RUN swag init -g ./cmd/api/main.go -o ./docs

EXPOSE 5000
CMD ["air"]