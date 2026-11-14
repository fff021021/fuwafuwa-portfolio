// イースターエッグ・隠し要素・実績アンロック機能
class EasterEggSystem {
  constructor() {
    this.secretsFound = new Set(JSON.parse(localStorage.getItem('fuwa_secrets') || '[]'));
    this.totalSecrets = 3;
    this.konamiCode = [
      'ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown',
      'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight',
      'b', 'a'
    ];
    this.konamiIndex = 0;
    this.logoClickCount = 0;
  }

  init() {
    this.updateSecretCounter();
    this.bindKonami();
    this.bindLogoClick();
    this.bindToastSystem();
  }

  updateSecretCounter() {
    const el = document.getElementById('secrets-counter');
    if (el) {
      el.textContent = `SECRETS: ${this.secretsFound.size}/${this.totalSecrets} FOUND`;
      if (this.secretsFound.size === this.totalSecrets) {
        el.style.background = 'var(--color-pink)';
        el.style.color = '#fff';
      }
    }
  }

  unlockSecret(secretId, message) {
    if (this.secretsFound.has(secretId)) return;

    this.secretsFound.add(secretId);
    localStorage.setItem('fuwa_secrets', JSON.stringify(Array.from(this.secretsFound)));
    this.updateSecretCounter();

    if (window.soundController) {
      window.soundController.playSecret();
    }

    this.showToast(`🏆 SECRET UNLOCKED! (${this.secretsFound.size}/${this.totalSecrets})`, message);
  }

  bindKonami() {
    document.addEventListener('keydown', (e) => {
      const expectedKey = this.konamiCode[this.konamiIndex];
      if (e.key.toLowerCase() === expectedKey.toLowerCase()) {
        this.konamiIndex++;
        if (this.konamiIndex === this.konamiCode.length) {
          this.konamiIndex = 0;
          this.triggerDebugOverlay();
          this.unlockSecret('konami_code', 'コナミコマンド発見！サーバーデバッグHUDが有効化されました。');
        }
      } else {
        this.konamiIndex = 0;
      }
    });
  }

  triggerDebugOverlay() {
    const overlay = document.getElementById('debug-overlay');
    if (overlay) {
      overlay.classList.toggle('active');
      if (overlay.classList.contains('active')) {
        this.updateLiveServerStats();
      }
    }
  }

  async updateLiveServerStats() {
    try {
      const health = await window.API.getHealth();
      const statsBody = document.getElementById('debug-stats-body');
      if (statsBody) {
        statsBody.innerHTML = `
          <div>STATUS: <span style="color:#00FF66;">${health.status}</span></div>
          <div>UPTIME: ${health.uptime}</div>
          <div>GOROUTINES: ${health.goroutines}</div>
          <div>GO VERSION: ${health.go_version}</div>
          <div>DB CONNECTION: OK (SQLite)</div>
        `;
      }
    } catch (e) {
      const statsBody = document.getElementById('debug-stats-body');
      if (statsBody) {
        statsBody.innerHTML = `<div>STATUS: OFFLINE</div>`;
      }
    }
  }

  bindLogoClick() {
    const logo = document.getElementById('fuwa-logo-click');
    if (!logo) return;

    logo.addEventListener('click', () => {
      this.logoClickCount++;
      if (window.soundController) window.soundController.playClick();

      // ブルブル震えるエフェクト
      logo.style.transform = `scale(1.15) rotate(${this.logoClickCount % 2 === 0 ? 5 : -5}deg)`;
      setTimeout(() => {
        logo.style.transform = 'none';
      }, 150);

      if (this.logoClickCount === 5) {
        this.unlockSecret('logo_clicker', 'ロゴを5回タップ！FUWAFUWAの真髄に触れました。');
      }
    });
  }

  bindToastSystem() {
    let container = document.getElementById('toast-container');
    if (!container) {
      container = document.createElement('div');
      container.id = 'toast-container';
      container.className = 'toast-container';
      document.body.appendChild(container);
    }
  }

  showToast(title, body) {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = 'brutal-toast';
    toast.innerHTML = `
      <div style="font-size: 1.5rem;">🌟</div>
      <div>
        <div style="font-size: 0.8rem; color: var(--color-pink);">${title}</div>
        <div style="font-size: 0.85rem; color: var(--color-black);">${body}</div>
      </div>
    `;

    container.appendChild(toast);
    requestAnimationFrame(() => toast.classList.add('show'));

    setTimeout(() => {
      toast.classList.remove('show');
      setTimeout(() => toast.remove(), 400);
    }, 4500);
  }
}

window.EasterEggs = new EasterEggSystem();
