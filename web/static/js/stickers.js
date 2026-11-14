// ドラッグ可能ステッカー機能（Digital Playground）
class StickerManager {
  constructor() {
    this.stickers = [];
    this.isDragging = false;
    this.currentSticker = null;
    this.offset = { x: 0, y: 0 };
    this.dragCount = 0;
  }

  init() {
    this.stickers = document.querySelectorAll('.sticker');
    if (!this.stickers.length) return;

    this.stickers.forEach((sticker) => {
      // マウスイベント
      sticker.addEventListener('mousedown', (e) => this.onDragStart(e, sticker));
      // タッチイベント
      sticker.addEventListener('touchstart', (e) => this.onDragStart(e, sticker), { passive: false });

      // ダブルクリックで回転リセット/ランダムチルト
      sticker.addEventListener('dblclick', () => {
        const randomDeg = Math.floor(Math.random() * 30) - 15;
        sticker.style.transform = `rotate(${randomDeg}deg) scale(1.1)`;
        if (window.soundController) window.soundController.playClick();
      });
    });

    document.addEventListener('mousemove', (e) => this.onDrag(e));
    document.addEventListener('touchmove', (e) => this.onDrag(e), { passive: false });

    document.addEventListener('mouseup', () => this.onDragEnd());
    document.addEventListener('touchend', () => this.onDragEnd());
  }

  onDragStart(e, sticker) {
    this.isDragging = true;
    this.currentSticker = sticker;
    sticker.style.zIndex = '100';

    const clientX = e.type.startsWith('touch') ? e.touches[0].clientX : e.clientX;
    const clientY = e.type.startsWith('touch') ? e.touches[0].clientY : e.clientY;

    const rect = sticker.getBoundingClientRect();
    this.offset.x = clientX - rect.left;
    this.offset.y = clientY - rect.top;

    if (window.soundController) {
      window.soundController.playClick();
    }

    if (e.type.startsWith('touch')) {
      e.preventDefault();
    }
  }

  onDrag(e) {
    if (!this.isDragging || !this.currentSticker) return;

    const clientX = e.type.startsWith('touch') ? e.touches[0].clientX : e.clientX;
    const clientY = e.type.startsWith('touch') ? e.touches[0].clientY : e.clientY;

    const containerRect = document.querySelector('.hero-cover').getBoundingClientRect();
    let left = clientX - containerRect.left - this.offset.x;
    let top = clientY - containerRect.top - this.offset.y;

    this.currentSticker.style.left = `${left}px`;
    this.currentSticker.style.top = `${top}px`;
    this.currentSticker.style.right = 'auto';
    this.currentSticker.style.bottom = 'auto';

    if (e.type.startsWith('touch')) {
      e.preventDefault();
    }
  }

  onDragEnd() {
    if (this.isDragging && this.currentSticker) {
      this.currentSticker.style.zIndex = '30';
      this.dragCount++;

      // ステッカーを3回以上動かしたら実績解除チェック
      if (this.dragCount === 3 && window.EasterEggs) {
        window.EasterEggs.unlockSecret('sticker_master', 'ステッカーマスター！画面を自由に彩りました。');
      }
    }
    this.isDragging = false;
    this.currentSticker = null;
  }
}

window.StickerManager = new StickerManager();
