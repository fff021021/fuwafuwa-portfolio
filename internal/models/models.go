package models

import "time"

// Profile は自己紹介・クリエイター情報を表す構造体
type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Bio       string    `json:"bio"`
	Vision    string    `json:"vision"`
	Keywords  string    `json:"keywords"` // JSON文字列 (例: ["Go", "Game Server"])
	Stats     string    `json:"stats"`    // JSON文字列 (パラメータマップ)
	UpdatedAt time.Time `json:"updated_at"`
}

// Work は制作物・アーカイブ作品を表す構造体
type Work struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"` // "web", "bot", "tool", "lab"
	Tagline      string    `json:"tagline"`
	Description  string    `json:"description"`
	Technologies string    `json:"technologies"` // JSON文字列 (例: ["Go", "WebSocket"])
	Highlights   string    `json:"highlights"`   // こだわりや試行錯誤
	DemoURL      string    `json:"demo_url,omitempty"`
	GithubURL    string    `json:"github_url,omitempty"`
	Featured     int       `json:"featured"`
	CreatedAt    time.Time `json:"created_at"`
}

// ActivityLog は日々の開発・学習活動記録を表す構造体
type ActivityLog struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	Metrics   string    `json:"metrics,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// QuizQuestion はサーバー内部で管理するクイズ問題構造体
type QuizQuestion struct {
	ID           string `json:"id"`
	Question     string `json:"question"`
	Category     string `json:"category"`
	Options      string `json:"options"` // JSON文字列 (4つの選択肢配列)
	CorrectIndex int    `json:"-"`       // サーバー側のみ保持し、クライアントには渡さない
	Explanation  string `json:"explanation"`
	Points       int    `json:"points"`
}

// QuizQuestionPublic はクライアントに安全に送信するクイズ問題DTO
type QuizQuestionPublic struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Category string `json:"category"`
	Options  string `json:"options"`
	Points   int    `json:"points"`
}

// QuizAnswerRequest はクライアントからのクイズ回答リクエスト
type QuizAnswerRequest struct {
	QuestionID    string `json:"question_id"`
	SelectedIndex int    `json:"selected_index"`
}

// QuizAnswerResponse はサーバーの採点結果レスポンス
type QuizAnswerResponse struct {
	IsCorrect   bool   `json:"is_correct"`
	Explanation string `json:"explanation"`
	Points      int    `json:"points"`
}

// RandomCard はランダム生成器用のコンテンツ構造体
type RandomCard struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	CodeSnippet string `json:"code_snippet,omitempty"`
}

// HealthResponse はサーバーの稼働状態を表す構造体
type HealthResponse struct {
	Status     string `json:"status"`
	Uptime     string `json:"uptime"`
	Goroutines int    `json:"goroutines"`
	GoVersion  string `json:"go_version"`
}
