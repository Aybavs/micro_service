FROM golang:1.23

WORKDIR /app

# Go mod dosyalarını kopyala, bağımlılıkları indir
COPY go.mod go.sum ./
RUN go mod download

# Proje içeriğini kopyala
COPY . .

# Binary üret
RUN go build -o main .

# Uygulamayı başlat
CMD ["./main"]
