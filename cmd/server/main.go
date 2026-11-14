package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"portfolio/internal/api"
	"portfolio/internal/db"
	"portfolio/web"
)

func main() {
	// クラウド環境（Render, Fly.io, Cloud Run等）の PORT 環境変数に対応
	defaultPort := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			defaultPort = p
		}
	}

	defaultDB := "./portfolio.db"
	if envDB := os.Getenv("DB_PATH"); envDB != "" {
		defaultDB = envDB
	}

	port := flag.Int("port", defaultPort, "サーバーのリッスンポート (環境変数 PORT でも指定可能)")
	dbPath := flag.String("db", defaultDB, "SQLiteデータベースファイルのパス (環境変数 DB_PATH でも指定可能)")
	flag.Parse()

	log.Println("==================================================")
	log.Println(" FUWAFUWA PORTFOLIO - GAME SERVER APPLICANT ENGINE")
	log.Println("==================================================")

	// 1. データベース接続・初期化
	database, err := db.InitDB(*dbPath)
	if err != nil {
		log.Fatalf("データベース初期化失敗: %v", err)
	}
	defer database.Close()

	// 2. シードマスターデータ投入
	if err := db.SeedData(database); err != nil {
		log.Fatalf("シードデータ投入失敗: %v", err)
	}

	// 3. APIハンドラーの生成
	handler := api.NewHandler(database)

	// 4. ルーティング設定
	mux := http.NewServeMux()

	// REST API エンドポイント
	mux.HandleFunc("/api/profile", handler.GetProfile)
	mux.HandleFunc("/api/works", handler.GetWorks)
	mux.HandleFunc("/api/logs", handler.GetLogs)
	mux.HandleFunc("/api/quiz/questions", handler.GetQuizQuestions)
	mux.HandleFunc("/api/quiz/answer", handler.AnswerQuiz)
	mux.HandleFunc("/api/random", handler.GetRandom)
	mux.HandleFunc("/api/health", handler.GetHealth)

	// 5. 静的アセット配信（ハイブリッド構成: ローカルディレクトリ優先、なければembed.FS）
	var staticFS http.Handler
	var indexHandler http.HandlerFunc

	if _, err := os.Stat(filepath.Join("web", "index.html")); err == nil {
		// ローカルの web/ ディレクトリから配信（開発時）
		staticDir := filepath.Join("web", "static")
		staticFS = http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))
		indexHandler = func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			http.ServeFile(w, r, filepath.Join("web", "index.html"))
		}
	} else {
		// バイナリ内に埋め込まれた embed.FS から配信（本番デプロイ時・完全単一バイナリ動作）
		staticSubFS, err := fs.Sub(web.EmbeddedFiles, "static")
		if err != nil {
			log.Fatalf("静的ファイルシステムの初期化失敗: %v", err)
		}

		staticFS = http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS)))
		indexHandler = func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			data, err := web.EmbeddedFiles.ReadFile("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}
	}

	mux.Handle("/static/", staticFS)
	mux.HandleFunc("/", indexHandler)

	// 6. ミドルウェアチェーンの適用
	var rootHandler http.Handler = mux
	rootHandler = api.LimitBodyMiddleware(1024 * 1024)(rootHandler) // 1MB制限
	rootHandler = api.SecurityMiddleware(rootHandler)
	rootHandler = api.LoggingMiddleware(rootHandler)

	// 7. HTTPサーバー構成
	serverAddr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        rootHandler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 8. 非同期起動 & Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 サーバーが起動しました: http://localhost:%d (PORT=%d)\n", *port, *port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("サーバーエラー: %v", err)
		}
	}()

	// 終了シグナル待機
	<-stopChan
	log.Println("\nシャットダウンシグナルを受信しました。安全に停止しています...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("サーバーの正常停止に失敗しました: %v", err)
	}

	log.Println("サーバーが正常に終了しました。Good bye!")
}
