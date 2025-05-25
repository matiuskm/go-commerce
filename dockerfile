# Gunakan image Go 1.24 terbaru (alpine supaya ringan)
FROM golang:1.24-alpine

# Set working directory di dalam container
WORKDIR /app

# Copy file go.mod dan go.sum terlebih dahulu (biar cache dependency)
COPY go.mod go.sum ./

# Download semua dependencies
RUN go mod download

# Copy seluruh source code ke dalam container
COPY . .

# Build binary Go
RUN go build -o main .

# Buka port 8080 (atau ganti sesuai port backend lo)
EXPOSE 8080

# Jalankan binary hasil build
CMD ["./main"]
