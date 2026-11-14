package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"portfolio/internal/api"
	"portfolio/internal/db"
)

func main() {
	port := flag.Int("port", 8080, "サーバーのリッスンポート")
	dbPath := flag.String("db", "./portfolio.db", "SQLiteデータベースファイルのパス")
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

	// 静的アセット配信 (/static/ -> web/static/)
	staticDir := filepath.Join("web", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// ルート (/ -> web/index.html)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join("web", "index.html"))
	})

	// 5. ミドルウェアチェーンの適用
	// LimitBody (1MB) -> Security/CORS -> Logging
	var rootHandler http.Handler = mux
	rootHandler = api.LimitBodyMiddleware(1024 * 1024)(rootHandler) // 1MB制限
	rootHandler = api.SecurityMiddleware(rootHandler)
	rootHandler = api.LoggingMiddleware(rootHandler)

	// 6. HTTPサーバー構成
	serverAddr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        rootHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 7. 非同期起動 & Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 サーバーが起動しました: http://localhost:%d\n", *port)
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
