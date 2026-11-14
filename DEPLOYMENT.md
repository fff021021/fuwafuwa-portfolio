# FUWAFUWA PORTFOLIO デプロイ手順書 🚀

本ポートフォリオは、**Go言語（Pure-Go SQLite）とVanilla JS**で構築されており、外部ディレクトリ不要の単一バイナリまたはDockerコンテナとして、世界中のクラウドサービスへ完全無料で簡単にデプロイできます。

---

## おすすめデプロイ方法 3選

| サービス | 費用 | 特徴 | 手間 |
|---|---|---|---|
| **1. Render (推奨)** | **完全無料** | GitHubプッシュだけで自動HTTPS公開。最も簡単！ | ★☆☆ (一番かんたん) |
| **2. Fly.io** | **無料枠あり** | 東京リージョン（nrt）対応で超低遅延。 | ★★☆ (CLIを使用) |
| **3. Railway / Koyeb** | **無料枠あり** | GitHub連携するだけで即時立ち上げ。 | ★☆☆ (簡単) |

---

## 方法 1: Render.com に無料デプロイ（最も簡単・おすすめ）

Renderは、GitHubのリポジトリと連携するだけで、自動でDockerコンテナをビルドし、無料の公開URL（`https://fuwafuwa-portfolio.onrender.com`）を発行してくれます。

### ステップ 1: GitHubにリポジトリを作成してプッシュする
```bash
# GitHubで新規リポジトリ（例: fuwafuwa-portfolio）を作成後、プッシュします
git remote add origin https://github.com/<あなたのユーザー名>/fuwafuwa-portfolio.git
git branch -M main
git push -u origin main
```

### ステップ 2: Render.com でサービスを作成
1. [Render.com](https://render.com/) にアクセスし、無料アカウントを作成（GitHubアカウントでログイン可能）。
2. ダッシュボードで **[New +]** ➔ **[Web Service]** をクリック。
3. **[Build and deploy from a Git repository]** を選択し、先ほどプッシュした `fuwafuwa-portfolio` を選択。
4. 設定画面が開きます（リポジトリ内の `Dockerfile` が自動認識されます）：
   - **Name**: `fuwafuwa-portfolio`
   - **Region**: `Singapore` または `Oregon`
   - **Instance Type**: **Free**（無料）
5. 一番下の **[Create Web Service]** をクリック！

数分待つとビルドが完了し、**`https://fuwafuwa-portfolio.onrender.com`** という全世界からアクセス可能なURLが発行されます！🎉

---

## 方法 2: Fly.io にデプロイ（東京リージョンで爆速）

Fly.ioは、日本の東京リージョン（nrt）が利用でき、ゲームサーバーらしい極限の低レイテンシを体験できます。

### ステップ 1: Fly.io CLI のインストール
```powershell
# PowerShellでインストール
iwr https://fly.io/install.ps1 -useb | iex
```

### ステップ 2: ログイン & デプロイ
```powershell
# ログイン
fly auth login

# アプリの初期化（リポジトリの fly.toml が自動適用されます）
fly launch --no-deploy

# ボリューム（SQLiteの永続化ディスク）の作成（東京リージョン）
fly volumes create portfolio_data --size 1 --region nrt

# デプロイ実行！
fly deploy
```
完了後、`fly open` を実行するとブラウザで公開サイトが立ち上がります。

---

## 方法 3: 自作VPS / Docker Compose でのデプロイ

ConoHa、さくらのVPS、AWS EC2、Google Cloud Compute EngineなどのLinuxサーバーをお持ちの場合：

```bash
git clone <あなたのリポジトリURL>
cd portfolio
docker compose up -d
```
これだけでポート8080でバックグラウンド起動します。NginxやCaddyをリバースプロキシとして前に置くことで、独自ドメインでのHTTPS運用も可能です。

---

## 単一バイナリによる直接実行（Linux / VPS）

本ポートフォリオは `embed.FS` により、HTML/CSS/JS/画像をすべてGoバイナリ内に内包しています。  
そのため、Windows上でLinux向けにクロスコンパイルするだけで、**たった1つの実行ファイルをLinuxサーバーにアップロードして起動するだけ**で動きます！

```powershell
# WindowsからLinux(amd64)向けにコンパイル
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -ldflags="-w -s" -o server_linux ./cmd/server
```
生成された `server_linux` をサーバーに `scp` 等で送り、
```bash
chmod +x server_linux
./server_linux -port 8080
```
と実行するだけで、追加ファイル一切なしで完全動作します。
