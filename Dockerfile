# Gunakan base image Go
FROM golang:1.22.0-alpine

# Set working directory
WORKDIR /app

# Install git (kalau perlu ambil package dari private repo)
RUN apk update && apk add --no-cache git

# Copy go.mod & go.sum dulu, untuk cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build aplikasi
RUN go build -o service-account .

# Expose port (gunakan APP_PORT dari .env, default 8080)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./service-account"]
