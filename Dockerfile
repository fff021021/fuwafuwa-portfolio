
# ===================================================
# Stage 1: Build stage
# ===================================================
FROM golang:1.27.4 AS builder

WORKDIR /app

# CA証明書・Git・タイムゾーン
# Debian系イメージなのでapt-getを使用
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Goの依存関係を先にコピーしてキャッシュを活用
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Linux AMD64向けにビルド
# Pure-Go SQLiteを想定してCGOを無効化
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server


# ===================================================
# Stage 2: Runtime stage
# ===================================================
FROM alpine:3.20

WORKDIR /app

# 実行環境に必要な証明書とタイムゾーン
RUN apk add --no-cache ca-certificates tzdata

# ビルド成果物をコピー
COPY --from=builder /app/server /app/server

# タイムゾーン
ENV TZ=Asia/Tokyo

# アプリケーション設定
ENV DB_PATH=/app/data/portfolio.db
ENV PORT=8080

# データ永続化ディレクトリ
RUN mkdir -p /app/data

# 非rootユーザーで実行
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

# サーバー起動
CMD ["/app/server"]
