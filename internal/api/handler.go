package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"portfolio/internal/models"
)

// Handler はAPIリクエストを処理する構造体
type Handler struct {
	db        *sql.DB
	startTime time.Time
}

// NewHandler は新しいHandlerインスタンスを生成します
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:        db,
		startTime: time.Now(),
	}
}

// jsonResponse はJSON形式でレスポンスを返却するヘルパー関数
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"JSONエンコード失敗"}`, http.StatusInternalServerError)
	}
}

// errorResponse はエラーメッセージをJSON形式で返却するヘルパー関数
func errorResponse(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

// GetProfile はプロフィール情報を取得します (GET /api/profile)
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "GETメソッドのみ許可されています")
		return
	}

	var p models.Profile
	query := `SELECT id, name, role, bio, vision, keywords, stats, updated_at FROM profiles LIMIT 1`
	err := h.db.QueryRow(query).Scan(&p.ID, &p.Name, &p.Role, &p.Bio, &p.Vision, &p.Keywords, &p.Stats, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(w, http.StatusNotFound, "プロフィール情報が見つかりません")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "DBエラーが発生しました")
		return
	}

	jsonResponse(w, http.StatusOK, p)
}

// GetWorks は制作物一覧を取得します (GET /api/works?category=web&q=キーワード)
func (h *Handler) GetWorks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "GETメソッドのみ許可されています")
		return
	}

	category := r.URL.Query().Get("category")
	keyword := r.URL.Query().Get("q")

	query := `SELECT id, title, category, tagline, description, technologies, highlights, demo_url, github_url, featured, created_at FROM works WHERE 1=1`
	var args []interface{}

	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	if keyword != "" {
		query += ` AND (title LIKE ? OR description LIKE ? OR technologies LIKE ?)`
		likePattern := "%" + keyword + "%"
		args = append(args, likePattern, likePattern, likePattern)
	}

	query += ` ORDER BY featured DESC, created_at DESC`

	rows, err := h.db.Query(query, args...)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "制作物の取得に失敗しました")
		return
	}
	defer rows.Close()

	works := make([]models.Work, 0)
	for rows.Next() {
		var wItem models.Work
		var demoURL, githubURL sql.NullString
		err := rows.Scan(
			&wItem.ID, &wItem.Title, &wItem.Category, &wItem.Tagline,
			&wItem.Description, &wItem.Technologies, &wItem.Highlights,
			&demoURL, &githubURL, &wItem.Featured, &wItem.CreatedAt,
		)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, "行スキャンエラー")
			return
		}
		if demoURL.Valid {
			wItem.DemoURL = demoURL.String
		}
		if githubURL.Valid {
			wItem.GithubURL = githubURL.String
		}
		works = append(works, wItem)
	}

	jsonResponse(w, http.StatusOK, works)
}

// GetLogs は日々の活動・学習記録を取得します (GET /api/logs?limit=10)
func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "GETメソッドのみ許可されています")
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	query := `SELECT id, date, title, category, content, metrics, created_at FROM activity_logs ORDER BY date DESC, created_at DESC LIMIT ?`
	rows, err := h.db.Query(query, limit)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "活動ログの取得に失敗しました")
		return
	}
	defer rows.Close()

	logs := make([]models.ActivityLog, 0)
	for rows.Next() {
		var item models.ActivityLog
		var metrics sql.NullString
		err := rows.Scan(&item.ID, &item.Date, &item.Title, &item.Category, &item.Content, &metrics, &item.CreatedAt)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, "行スキャンエラー")
			return
		}
		if metrics.Valid {
			item.Metrics = metrics.String
		}
		logs = append(logs, item)
	}

	jsonResponse(w, http.StatusOK, logs)
}

// GetQuizQuestions はクイズ問題一覧を取得します (GET /api/quiz/questions)
// セキュリティ上、正解インデックス (correct_index) は除外してクライアントに返します。
func (h *Handler) GetQuizQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "GETメソッドのみ許可されています")
		return
	}

	query := `SELECT id, question, category, options, points FROM quiz_questions ORDER BY id ASC`
	rows, err := h.db.Query(query)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "クイズ問題の取得に失敗しました")
		return
	}
	defer rows.Close()

	questions := make([]models.QuizQuestionPublic, 0)
	for rows.Next() {
		var q models.QuizQuestionPublic
		if err := rows.Scan(&q.ID, &q.Question, &q.Category, &q.Options, &q.Points); err != nil {
			errorResponse(w, http.StatusInternalServerError, "行スキャンエラー")
			return
		}
		questions = append(questions, q)
	}

	jsonResponse(w, http.StatusOK, questions)
}

// AnswerQuiz はクイズの回答を受け取り、Goサーバー側で採点して判定結果を返します (POST /api/quiz/answer)
func (h *Handler) AnswerQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "POSTメソッドのみ許可されています")
		return
	}

	// Content-Type チェック
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		errorResponse(w, http.StatusBadRequest, "Content-Typeはapplication/jsonである必要があります")
		return
	}

	var req models.QuizAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "不正なJSONリクエストです")
		return
	}

	// 入力値バリデーション
	if req.QuestionID == "" {
		errorResponse(w, http.StatusBadRequest, "question_idは必須です")
		return
	}
	if req.SelectedIndex < 0 || req.SelectedIndex > 3 {
		errorResponse(w, http.StatusBadRequest, "selected_indexは0〜3の範囲である必要があります")
		return
	}

	// サーバーサイドDBで正解と解説を取得
	var correctIndex, points int
	var explanation string
	query := `SELECT correct_index, explanation, points FROM quiz_questions WHERE id = ?`
	err := h.db.QueryRow(query, req.QuestionID).Scan(&correctIndex, &explanation, &points)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(w, http.StatusNotFound, "該当する問題が見つかりません")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "採点処理中にDBエラーが発生しました")
		return
	}

	// 正誤判定
	isCorrect := (req.SelectedIndex == correctIndex)
	pts := 0
	if isCorrect {
		pts = points
	}

	resp := models.QuizAnswerResponse{
		IsCorrect:   isCorrect,
		Explanation: explanation,
		Points:      pts,
	}

	jsonResponse(w, http.StatusOK, resp)
}

// GetRandom はランダムな豆知識・コードスニペットを取得します (GET /api/random)
func (h *Handler) GetRandom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "GETメソッドのみ許可されています")
		return
	}

	// SQLiteのORDER BY RANDOM()でサーバー側から1件ランダム取得
	query := `SELECT id, category, title, content, code_snippet FROM random_cards ORDER BY RANDOM() LIMIT 1`
	var card models.RandomCard
	var snippet sql.NullString
	err := h.db.QueryRow(query).Scan(&card.ID, &card.Category, &card.Title, &card.Content, &snippet)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "ランダムコンテンツの取得に失敗しました")
		return
	}
	if snippet.Valid {
		card.CodeSnippet = snippet.String
	}

	jsonResponse(w, http.StatusOK, card)
}

// GetHealth はサーバーの稼働メトリクスを取得します (GET /api/health)
func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).Truncate(time.Second).String()
	resp := models.HealthResponse{
		Status:     "OK",
		Uptime:     uptime,
		Goroutines: runtime.NumGoroutine(),
		GoVersion:  runtime.Version(),
	}
	jsonResponse(w, http.StatusOK, resp)
}
