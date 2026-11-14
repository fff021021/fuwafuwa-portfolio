// イースターエッグ・隠し要素・実績アンロック機能（Webネイティブ操作版）
class EasterEggSystem {
  constructor() {
    this.secretsFound = new Set(JSON.parse(localStorage.getItem('fuwa_secrets') || '[]'));
    this.totalSecrets = 3;
    this.logoClickCount = 0;
  }

  init() {
    this.updateSecretCounter();
    this.bindLogoClick();
    this.bindToastSystem();
  }

  updateSecretCounter() {
    const el = document.getElementById('secrets-counter');
    if (el) {
      el.textContent = `SECRETS: ${this.secretsFound.size}/${this.totalSecrets} FOUND`;
      if (this.secretsFound.size === this.totalSecrets) {
        el.style.background = 'var(--color-pink)';
        el.style.color = '#FFFFFF';
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

    this.showToast(`🏆 実績解除！ (${this.secretsFound.size}/${this.totalSecrets})`, message);
  }

  bindLogoClick() {
    const logo = document.getElementById('fuwa-logo-click');
    if (!logo) return;

    logo.addEventListener('click', () => {
      this.logoClickCount++;
      if (window.soundController) window.soundController.playClick();

      // バウンスエフェクト
      logo.style.transform = `scale(1.15) rotate(${this.logoClickCount % 2 === 0 ? 6 : -6}deg)`;
      setTimeout(() => {
        logo.style.transform = 'none';
      }, 180);

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
      <div style="font-size: 1.6rem;">🌟</div>
      <div>
        <div style="font-size: 0.8rem; color: var(--color-pink); font-weight: 800; letter-spacing: 0.04em;">${title}</div>
        <div style="font-size: 0.9rem; color: #171717; font-weight: 700; margin-top: 0.2rem;">${body}</div>
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
