package db

import (
	"database/sql"
	"fmt"
	"log"
)

// SeedData はデータベースが空の場合に初期マスターデータをトランザクションで投入します
func SeedData(db *sql.DB) error {
	// プロフィールの存在確認
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&count)
	if err != nil {
		return fmt.Errorf("データ件数確認エラー: %w", err)
	}

	if count > 0 {
		return nil // 既にデータが存在する場合はスキップ
	}

	log.Println("シードデータの投入を開始します...")

	// 整合性を保つためトランザクションを使用
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", err)
	}
	defer tx.Rollback()

	// 1. プロフィール投入
	profileQuery := `INSERT INTO profiles (id, name, role, bio, vision, keywords, stats) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(profileQuery,
		"default_profile",
		"FUWAFUWA DEV",
		"Game Server Engineer Applicant",
		"ゲーム開発会社でのサーバーサイドエンジニアを志望しています。Go言語による並行処理（goroutine/channel）や低フットプリントなAPIサーバー設計、分散システム、リアルタイム通信に関心を持っています。大規模なトラフィック下でも淀みなくパケットを捌き切るアーキテクチャの美しさに惹かれ、日々自作ツールやプロトタイプの開発、インフラ技術の学習に取り組んでいます。",
		"高い並行性と低レイテンシを追求し、プレイヤーがストレスを感じない盤石な世界基盤を組む。",
		`["Go言語 (net/http)", "並行並列処理", "リアルタイム通信", "SQLite / RDBMS", "REST API設計", "低レイテンシ", "ゲームループ設計"]`,
		`{"go": 88, "concurrency": 85, "db": 80, "network": 82, "curiosity": 98}`,
	)
	if err != nil {
		return fmt.Errorf("プロフィールシードエラー: %w", err)
	}

	// 2. 制作物・活動アーカイブ (Works) 投入
	workQuery := `INSERT INTO works (id, title, category, tagline, description, technologies, highlights, demo_url, github_url, featured) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	works := []struct {
		id, title, category, tagline, desc, techs, highlights, demo, github string
		featured                                                            int
	}{
		{
			id:         "work-01",
			title:      "FUWAFUWA PORTFOLIO",
			category:   "web",
			tagline:    "Digital Playground × 雑誌エディトリアルの自己紹介Webアプリ",
			desc:       "ゲーム開発会社のサーバーエンジニア向けインターンシップ応募作品として自作した、フルスタックWebポートフォリオです。Goのnet/httpとPure-Go SQLiteで構築されたバックエンドAPIと、Vanilla JSによるリッチなインタラクティブUIが協調動作します。",
			techs:      `["Go", "net/http", "SQLite", "Vanilla JS", "Web Audio API", "CSS3"]`,
			highlights: "【試行錯誤・こだわり】外部の重厚なフレームワークに頼らず、Go標準ライブラリのHTTPサーバーとPure-Go SQLiteを採用。Cgo不要でポータブルに起動でき、メモリ効率と起動速度を極限まで高めました。クイズ採点をサーバーサイドで厳格に判定するREST API設計を取り入れています。",
			demo:       "#00-index",
			github:     "https://github.com/example/fuwafuwa-portfolio",
			featured:   1,
		},
		{
			id:         "work-02",
			title:      "GameEvent Dispatcher Bot",
			category:   "bot",
			tagline:    "ゲームコミュニティ向けリアルタイム告知＆イベント管理Discord Bot",
			desc:       "ゲームギルド内のレイドバトルや対戦トーナメントの進行を自動管理するDiscord Botです。非同期イベントキューを用いて数千件の通知リクエストを詰まらせずに配信します。",
			techs:      `["Go", "DiscordGo", "Goroutine", "Channel", "Worker Pool"]`,
			highlights: "【試行錯誤・こだわり】急激なスパイクアクセス時にDiscord APIのレートリミット（429 Too Many Requests）に抵触しないよう、Goのチャネルを活用したトークンバケットアルゴリズムによる流量制限ワーカープールを自作しました。",
			demo:       "",
			github:     "https://github.com/example/discord-event-bot",
			featured:   1,
		},
		{
			id:         "work-03",
			title:      "Micro-Room WebSocket Server",
			category:   "tool",
			tagline:    "複数プレイヤーの位置同期を検証する軽量WebSocketルームサーバー",
			desc:       "数万ルームの同時生成を想定した、Go製のリアルタイムマッチング・位置同期検証プロトタイプです。クライアント間の座標差分を高頻度（30fps相当）でブロードキャストします。",
			techs:      `["Go", "gorilla/websocket", "Sync.Map", "Atomic", "Benchmark"]`,
			highlights: "【試行錯誤・こだわり】ルームごとのロック競合を最小化するため、sync.RWMutexの粒度を細分化し、sync.Poolによるバッファ再利用でGC（ガベージコレクション）によるレイテンシスパイクを90%削減しました。",
			demo:       "",
			github:     "https://github.com/example/ws-room-server",
			featured:   1,
		},
		{
			id:         "work-04",
			title:      "Fast Log Aggregator CLI",
			category:   "lab",
			tagline:    "GB単位のアクセスログを並列解析するストリーミングCLIツール",
			desc:       "ゲームサーバーのアクセスログ（JSONL / Apache形式）から特定のエラーレスポンスや高レイテンシエンドポイントを瞬時に抽出する高速CLIツールです。",
			techs:      `["Go", "bufio.Scanner", "Pipeline Pattern", "CLI"]`,
			highlights: "【試行錯誤・こだわり】シングルスレッド処理では数分かかっていた巨大ログ解析を、プロデューサー・コンシューマー・パイプラインパターンで分割し、マルチコアを使い切ることで解析時間を85%短縮しました。",
			demo:       "",
			github:     "https://github.com/example/fast-log-cli",
			featured:   0,
		},
	}

	for _, w := range works {
		_, err = tx.Exec(workQuery, w.id, w.title, w.category, w.tagline, w.desc, w.techs, w.highlights, w.demo, w.github, w.featured)
		if err != nil {
			return fmt.Errorf("制作物シードエラー (%s): %w", w.id, err)
		}
	}

	// 3. 活動・学習記録 (Activity Logs) 投入
	logQuery := `INSERT INTO activity_logs (id, date, title, category, content, metrics) VALUES (?, ?, ?, ?, ?, ?)`
	logs := []struct {
		id, date, title, category, content, metrics string
	}{
		{
			id:       "log-01",
			date:     "2026.11",
			title:    "FUWAFUWA PORTFOLIOの設計・実装",
			category: "Project",
			content:  "Go + SQLite + Vanilla JSによるゲーム会社インターン向け自己紹介Webアプリをゼロから構築。雑誌デザインと遊び場の融合を追求。",
			metrics:  "Pure-Go / 0 npm deps",
		},
		{
			id:       "log-02",
			date:     "2026.10",
			title:    "WebSocket同時接続ベンチマークテスト",
			category: "Benchmark",
			content:  "自作ルームサーバーにて10,000クライアント同時接続時のメモリフットプリントとCPU使用率を検証。sync.Poolの効果を確認。",
			metrics:  "10k conns / 45MB RAM",
		},
		{
			id:       "log-03",
			date:     "2026.09",
			title:    "ゲームネットコードのクライアント予測・補間アルゴリズム学習",
			category: "Study",
			content:  "GDCの講演アーカイブやネットコード論文を読み込み、遅延環境下での位置補正ロジック（Dead Reckoning）をプロトタイプ実装。",
			metrics:  "Jitter 50ms test",
		},
		{
			id:       "log-04",
			date:     "2026.08",
			title:    "Go言語のメモリモデルとエスケープ解析の深掘り",
			category: "Study",
			content:  "go build -gcflags='-m' を活用し、スタック領域とヒープ領域の割り当ての違いやポインタ渡しによるオーバーヘッドを測定・ドキュメント化。",
			metrics:  "Zero-Allocation test",
		},
	}

	for _, l := range logs {
		_, err = tx.Exec(logQuery, l.id, l.date, l.title, l.category, l.content, l.metrics)
		if err != nil {
			return fmt.Errorf("活動記録シードエラー (%s): %w", l.id, err)
		}
	}

	// 4. クイズ問題 (Quiz Questions) 投入
	quizQuery := `INSERT INTO quiz_questions (id, question, category, options, correct_index, explanation, points) VALUES (?, ?, ?, ?, ?, ?, ?)`
	quizzes := []struct {
		id, question, category, options, explanation string
		correctIndex, points                         int
	}{
		{
			id:           "q-01",
			question:     "Go言語で数万件の同時接続を扱う際、goroutineがOSスレッドより軽量な主な理由は？",
			category:     "Go",
			options:      `["初期スタックサイズが数KB程度で可変だから", "GCが完全に無効化されているから", "CPUキャッシュを直接バイパスするから", "マルチスレッドではなくシングルスレッドで動くから"]`,
			correctIndex: 0,
			explanation:  "OSスレッドのスタックサイズが通常1MB〜2MBであるのに対し、goroutineの初期スタックはわずか約2KB〜4KBで、必要に応じて動的に拡張・縮小されるため、数万個同時に起動しても省メモリで動作します。",
			points:       10,
		},
		{
			id:           "q-02",
			question:     "高頻度なゲームパケット通信で、TCPではなくUDP（またはKCP/QUIC）が好まれる最大の理由は？",
			category:     "Network",
			options:      `["UDPは暗号化が自動でかかるから", "Head-of-Line Blocking（先頭パケット詰まり）による遅延を防ぎたいから", "UDPの方がパケットヘッダーが巨大だから", "パケット到達順序の保証が必須だから"]`,
			correctIndex: 1,
			explanation:  "TCPはパケット欠落時に再送が完了するまで後続パケットの処理が停止する（Head-of-Line Blocking）ため、リアルタイム対戦ゲームでは直前のパケット欠落を無視して最新の座標を優先できるUDPベースのプロトコルが重宝されます。",
			points:       10,
		},
		{
			id:           "q-03",
			question:     "ゲームサーバー内で一時的なバッファを頻繁に生成・破棄する場合、GC負荷を抑えるのに最適なGoの機能は？",
			category:     "Go",
			options:      `["runtime.GC()を毎フレーム呼ぶ", "sync.Pool を使用してオブジェクトを再利用する", "グローバル変数を1つだけ使い回す", "defer文をすべて削除する"]`,
			correctIndex: 1,
			explanation:  "sync.Poolは割り当て済みのメモリブロックをスレッドセーフにプールして再利用できるため、毎秒数万回発生するパケットシリアライズ用バッファのGCアロケーションを劇的に削減できます。",
			points:       10,
		},
		{
			id:           "q-04",
			question:     "本ポートフォリオ「FUWAFUWA PORTFOLIO」のクイズ採点処理はどこで実行されている？",
			category:     "Architecture",
			options:      `["クライアント側のJavaScriptコード内", "Firebaseのクラウド関数", "Go言語の自作バックエンドサーバー (/api/quiz/answer)", "外部のサードパーティクイズAPI"]`,
			correctIndex: 2,
			explanation:  "ゲーム開発会社のサーバーエンジニア応募作品として、重要な採点ロジックや正解データはクライアントに一切露出させず、GoサーバーサイドのREST API (/api/quiz/answer) で厳格に判定しています！",
			points:       10,
		},
	}

	for _, q := range quizzes {
		_, err = tx.Exec(quizQuery, q.id, q.question, q.category, q.options, q.correctIndex, q.explanation, q.points)
		if err != nil {
			return fmt.Errorf("クイズシードエラー (%s): %w", q.id, err)
		}
	}

	// 5. ランダム生成器コンテンツ (Random Cards) 投入
	randomQuery := `INSERT INTO random_cards (id, category, title, content, code_snippet) VALUES (?, ?, ?, ?, ?)`
	randoms := []struct {
		id, category, title, content, snippet string
	}{
		{
			id:       "rand-01",
			category: "SERVER TIP",
			title:    "Contextによるタイムアウト伝播",
			content:  "ネットワークI/Oや外部DBアクセス時は、context.WithTimeoutで必ずデッドラインを設定し、サーバーリソースの枯渇を防ぎましょう。",
			snippet:  "ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)\ndefer cancel()",
		},
		{
			id:       "rand-02",
			category: "PERFORMANCE",
			title:    "sync.RWMutexの使い分け",
			content:  "読み取りが9割を占めるゲームマスタデータへのアクセスは、通常のMutexではなくRWMutexのRLock()を使うことでスループットが大幅に向上します。",
			snippet:  "mu.RLock()\nval := cache[key]\nmu.RUnlock()",
		},
		{
			id:       "rand-03",
			category: "ALGORITHM",
			title:    "空間インデックス（Quadtree）",
			content:  "広大な2Dマップ上の全プレイヤー間の距離を総当たり（O(N^2)）で計算するとパンクします。四分木で空間を分割すればO(N log N)で近傍探索が可能です。",
			snippet:  "type Quadtree struct {\n    Boundary Rect\n    Objects  []Entity\n    Children [4]*Quadtree\n}",
		},
		{
			id:       "rand-04",
			category: "PHILOSOPHY",
			title:    "ゲームサーバーの美学",
			content:  "プレイヤーが『ラグを感じない』こと。それは、裏側で緻密な予測補正と堅牢なパケット処理が1秒間に数十回、完璧に回り続けている証拠です。",
			snippet:  "// Silent Guardian: The Game Server Engine",
		},
	}

	for _, r := range randoms {
		_, err = tx.Exec(randomQuery, r.id, r.category, r.title, r.content, r.snippet)
		if err != nil {
			return fmt.Errorf("ランダムシードエラー (%s): %w", r.id, err)
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("シードデータコミットエラー: %w", err)
	}

	log.Println("シードデータの投入が成功しました！")
	return nil
}
