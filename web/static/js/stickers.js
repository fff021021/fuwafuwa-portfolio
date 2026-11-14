// ドラッグ可能ステッカー機能（Digital Playground - デルタ加算による正確な追従）
class StickerManager {
  constructor() {
    this.stickers = [];
    this.isDragging = false;
    this.currentSticker = null;
    this.startX = 0;
    this.startY = 0;
    this.initialLeft = 0;
    this.initialTop = 0;
    this.dragCount = 0;
  }

  init() {
    this.stickers = document.querySelectorAll('.sticker');
    if (!this.stickers.length) return;

    this.stickers.forEach((sticker) => {
      // 初期パーセント位置をピクセルに変換して固定（初回ドラッグ時のズレ防止）
      const rect = sticker.getBoundingClientRect();
      const parentRect = sticker.offsetParent ? sticker.offsetParent.getBoundingClientRect() : { left: 0, top: 0 };
      const currentLeft = rect.left - parentRect.left;
      const currentTop = rect.top - parentRect.top;

      sticker.style.left = `${currentLeft}px`;
      sticker.style.top = `${currentTop}px`;
      sticker.style.right = 'auto';
      sticker.style.bottom = 'auto';

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

    // ドラッグ開始時のポインタ位置
    this.startX = clientX;
    this.startY = clientY;

    // ドラッグ開始時の要素の現在の位置（parseIntで数値取得）
    this.initialLeft = parseFloat(sticker.style.left) || sticker.offsetLeft || 0;
    this.initialTop = parseFloat(sticker.style.top) || sticker.offsetTop || 0;

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

    // カーソルの移動量（デルタ）
    const deltaX = clientX - this.startX;
    const deltaY = clientY - this.startY;

    // 移動量分だけ正確に加算
    this.currentSticker.style.left = `${this.initialLeft + deltaX}px`;
    this.currentSticker.style.top = `${this.initialTop + deltaY}px`;

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
