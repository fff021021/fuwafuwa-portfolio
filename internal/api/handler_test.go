package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"portfolio/internal/db"
	"portfolio/internal/models"
)

// setupTestDB はテスト用のインメモリSQLiteとシードデータを初期化します
func setupTestDB(t *testing.T) *sql.DB {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("テスト用DB初期化失敗: %v", err)
	}

	if err := db.SeedData(database); err != nil {
		t.Fatalf("テスト用シードデータ投入失敗: %v", err)
	}

	return database
}

// TestGetProfile はプロフィール取得APIの正常動作を検証します
func TestGetProfile(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	handler := NewHandler(database)
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()

	handler.GetProfile(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ステータスコード不一致: 期待値 200, 実際 %d", resp.StatusCode)
	}

	var p models.Profile
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatalf("JSONデコード失敗: %v", err)
	}

	if p.Name != "FUWAFUWA DEV" {
		t.Errorf("名前の不一致: %s", p.Name)
	}
}

// TestGetWorks は制作物検索とカテゴリフィルターを検証します
func TestGetWorks(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	handler := NewHandler(database)

	// 全件取得テスト
	req := httptest.NewRequest(http.MethodGet, "/api/works", nil)
	w := httptest.NewRecorder()
	handler.GetWorks(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("全件取得失敗: %d", w.Result().StatusCode)
	}

	var works []models.Work
	json.NewDecoder(w.Result().Body).Decode(&works)
	if len(works) == 0 {
		t.Fatal("制作物が1件も取得できませんでした")
	}

	// カテゴリフィルタテスト ("bot")
	reqBot := httptest.NewRequest(http.MethodGet, "/api/works?category=bot", nil)
	wBot := httptest.NewRecorder()
	handler.GetWorks(wBot, reqBot)

	var botWorks []models.Work
	json.NewDecoder(wBot.Result().Body).Decode(&botWorks)
	for _, bw := range botWorks {
		if bw.Category != "bot" {
			t.Errorf("カテゴリフィルタ失敗: 期待 'bot', 実際 '%s'", bw.Category)
		}
	}
}

// TestGetQuizQuestionsSecurity はクイズ問題一覧で正解が伏せられているかを検証します
func TestGetQuizQuestionsSecurity(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	handler := NewHandler(database)
	req := httptest.NewRequest(http.MethodGet, "/api/quiz/questions", nil)
	w := httptest.NewRecorder()

	handler.GetQuizQuestions(w, req)

	var rawList []map[string]interface{}
	json.NewDecoder(w.Result().Body).Decode(&rawList)

	for _, item := range rawList {
		if _, exists := item["correct_index"]; exists {
			t.Errorf("セキュリティ警告: クライアント向けレスポンスに correct_index が露出しています！")
		}
	}
}

// TestAnswerQuiz はサーバー側でのクイズ採点処理を検証します
func TestAnswerQuiz(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	handler := NewHandler(database)

	// 1. 正解パターンの検証 (q-01の正解は 0)
	correctBody, _ := json.Marshal(models.QuizAnswerRequest{
		QuestionID:    "q-01",
		SelectedIndex: 0,
	})
	reqCorrect := httptest.NewRequest(http.MethodPost, "/api/quiz/answer", bytes.NewReader(correctBody))
	reqCorrect.Header.Set("Content-Type", "application/json")
	wCorrect := httptest.NewRecorder()

	handler.AnswerQuiz(wCorrect, reqCorrect)

	var respCorrect models.QuizAnswerResponse
	json.NewDecoder(wCorrect.Result().Body).Decode(&respCorrect)

	if !respCorrect.IsCorrect {
		t.Errorf("正解判定失敗: 期待 true, 実際 false")
	}
	if respCorrect.Points != 10 {
		t.Errorf("獲得ポイント不一致: 期待 10, 実際 %d", respCorrect.Points)
	}

	// 2. 不正解パターンの検証 (q-01でインデックス 1 を選択)
	wrongBody, _ := json.Marshal(models.QuizAnswerRequest{
		QuestionID:    "q-01",
		SelectedIndex: 1,
	})
	reqWrong := httptest.NewRequest(http.MethodPost, "/api/quiz/answer", bytes.NewReader(wrongBody))
	reqWrong.Header.Set("Content-Type", "application/json")
	wWrong := httptest.NewRecorder()

	handler.AnswerQuiz(wWrong, reqWrong)

	var respWrong models.QuizAnswerResponse
	json.NewDecoder(wWrong.Result().Body).Decode(&respWrong)

	if respWrong.IsCorrect {
		t.Errorf("不正解判定失敗: 期待 false, 実際 true")
	}
	if respWrong.Points != 0 {
		t.Errorf("不正解時のポイント不正: %d", respWrong.Points)
	}

	// 3. バリデーションエラー検証 (選択肢が範囲外: 99)
	invalidBody, _ := json.Marshal(models.QuizAnswerRequest{
		QuestionID:    "q-01",
		SelectedIndex: 99,
	})
	reqInvalid := httptest.NewRequest(http.MethodPost, "/api/quiz/answer", bytes.NewReader(invalidBody))
	reqInvalid.Header.Set("Content-Type", "application/json")
	wInvalid := httptest.NewRecorder()

	handler.AnswerQuiz(wInvalid, reqInvalid)

	if wInvalid.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("バリデーション失敗時のステータス不一致: 期待 400, 実際 %d", wInvalid.Result().StatusCode)
	}
}

// TestGetHealth はサーバーメトリクスAPIの動作を検証します
func TestGetHealth(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	handler := NewHandler(database)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	handler.GetHealth(w, req)

	var h models.HealthResponse
	json.NewDecoder(w.Result().Body).Decode(&h)

	if h.Status != "OK" {
		t.Errorf("ヘルスチェック不一致: %s", h.Status)
	}
	if h.Goroutines <= 0 {
		t.Errorf("goroutine数が不正です: %d", h.Goroutines)
	}
}
