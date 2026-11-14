package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// InitDB はSQLiteデータベースの初期化とテーブル作成を行います
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	// コネクションプールの設定（SQLiteの単一ファイルアクセスと並行性を考慮）
	db.SetMaxOpenConns(1) // SQLiteは書き込み競合を避けるため1接続推奨
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベース死活監視エラー: %w", err)
	}

	// テーブルスキーマのマイグレーション
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("スキーママイグレーションエラー: %w", err)
	}

	log.Println("SQLiteデータベースの初期化が正常に完了しました。")
	return db, nil
}

// migrate は必要なテーブルとインデックスを作成します
func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			role TEXT NOT NULL,
			bio TEXT NOT NULL,
			vision TEXT NOT NULL,
			keywords TEXT NOT NULL,
			stats TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS works (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			category TEXT NOT NULL,
			tagline TEXT NOT NULL,
			description TEXT NOT NULL,
			technologies TEXT NOT NULL,
			highlights TEXT NOT NULL,
			demo_url TEXT,
			github_url TEXT,
			featured INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_works_category ON works(category);`,
		`CREATE TABLE IF NOT EXISTS activity_logs (
			id TEXT PRIMARY KEY,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			category TEXT NOT NULL,
			content TEXT NOT NULL,
			metrics TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS quiz_questions (
			id TEXT PRIMARY KEY,
			question TEXT NOT NULL,
			category TEXT NOT NULL,
			options TEXT NOT NULL,
			correct_index INTEGER NOT NULL,
			explanation TEXT NOT NULL,
			points INTEGER DEFAULT 10
		);`,
		`CREATE TABLE IF NOT EXISTS random_cards (
			id TEXT PRIMARY KEY,
			category TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			code_snippet TEXT
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
