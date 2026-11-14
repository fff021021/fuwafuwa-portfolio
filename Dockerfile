# ===================================================
# Stage 1: Build stage
# ===================================================
FROM golang:1.27.4 AS builder

WORKDIR /app

# CA証明書等のビルド用ツール
RUN apk add --no-cache git ca-certificates tzdata

# 依存パッケージのキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# CGO不要（Pure-Go SQLite）でLinux向けに最適化ビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# ===================================================
# Stage 2: Minimal runtime image (~20MB)
# ===================================================
FROM alpine:3.20

WORKDIR /app

# CA証明書とタイムゾーンデータのコピー
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/server /app/server

# データ永続化ディレクトリの作成
RUN mkdir -p /app/data && chmod 777 /app/data
ENV DB_PATH=/app/data/portfolio.db
ENV PORT=8080

EXPOSE 8080

# サーバー起動
CMD ["/app/server"]
