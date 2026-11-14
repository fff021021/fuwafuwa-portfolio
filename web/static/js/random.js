// Random Generator 制御スクリプト（ガチャ風アイデア＆トリビアディスペンサー）
class RandomGenerator {
  constructor() {
    this.isSpinning = false;
    this.displayTitle = null;
    this.displayBody = null;
    this.displayCode = null;
    this.triggerBtn = null;
  }

  init() {
    this.displayCategory = document.getElementById('generator-category');
    this.displayTitle = document.getElementById('generator-title');
    this.displayBody = document.getElementById('generator-body');
    this.displayCode = document.getElementById('generator-code');
    this.triggerBtn = document.getElementById('generator-btn');

    if (!this.triggerBtn) return;

    this.triggerBtn.addEventListener('click', () => this.generate());

    // 初期表示として1件取得
    this.generate(true);
  }

  async generate(silent = false) {
    if (this.isSpinning) return;
    this.isSpinning = true;

    if (!silent && window.soundController) {
      window.soundController.playClick();
    }

    // スピン演出（シャッフルアニメーション）
    this.triggerBtn.disabled = true;
    this.displayTitle.innerHTML = '<span class="anim-float">🎲 GENERATING...</span>';
    this.displayBody.textContent = 'ランダムコンテンツをサーバーから引いています...';
    if (this.displayCode) this.displayCode.style.display = 'none';

    // ルーレット音演出
    let beepCount = 0;
    const beepInterval = setInterval(() => {
      if (!silent && window.soundController && beepCount < 5) {
        window.soundController.playBeep(400 + beepCount * 80, 0.04, 'square');
      }
      beepCount++;
    }, 80);

    try {
      const data = await window.API.getRandomContent();
      clearInterval(beepInterval);

      setTimeout(() => {
        if (!silent && window.soundController) {
          window.soundController.playSuccess();
        }

        if (this.displayCategory) {
          this.displayCategory.textContent = data.category || 'TECH TRIVIA';
        }
        this.displayTitle.textContent = data.title;
        this.displayBody.textContent = data.content;

        if (data.code_snippet && this.displayCode) {
          this.displayCode.textContent = data.code_snippet;
          this.displayCode.style.display = 'block';
        } else if (this.displayCode) {
          this.displayCode.style.display = 'none';
        }

        this.isSpinning = false;
        this.triggerBtn.disabled = false;
      }, 500);

    } catch (err) {
      clearInterval(beepInterval);
      this.displayTitle.textContent = '通信エラー';
      this.displayBody.textContent = 'コンテンツの取得に失敗しました。サーバー起動を確認してください。';
      this.isSpinning = false;
      this.triggerBtn.disabled = false;
    }
  }
}

window.RandomGen = new RandomGenerator();
