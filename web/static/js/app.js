// アプリケーション初期化・統合スクリプト
document.addEventListener('DOMContentLoaded', async () => {
  // 1. サウンドボタンの初期化
  const soundBtn = document.getElementById('sound-toggle-btn');
  if (soundBtn) {
    if (window.soundController.enabled) {
      soundBtn.classList.add('active');
      soundBtn.innerHTML = '🔊 SOUND [ON]';
    } else {
      soundBtn.classList.remove('active');
      soundBtn.innerHTML = '🔇 SOUND [OFF]';
    }

    soundBtn.addEventListener('click', () => {
      const isEnabled = window.soundController.toggle();
      if (isEnabled) {
        soundBtn.classList.add('active');
        soundBtn.innerHTML = '🔊 SOUND [ON]';
      } else {
        soundBtn.classList.remove('active');
        soundBtn.innerHTML = '🔇 SOUND [OFF]';
      }
    });
  }

  // 2. 全クリック要素への効果音バインド
  document.addEventListener('click', (e) => {
    const target = e.target.closest('button, a, .sticker, .work-card, .category-tab-btn');
    if (target && window.soundController) {
      // クイズボタンやサウンドトグル自身は個別で処理するためスキップ
      if (!target.classList.contains('quiz-option-btn') && target.id !== 'sound-toggle-btn') {
        window.soundController.playClick();
      }
    }
  });

  // 3. 各マネージャーの初期化
  if (window.StickerManager) window.StickerManager.init();
  if (window.EasterEggs) window.EasterEggs.init();
  if (window.QuizApp) window.QuizApp.init();
  if (window.RandomGen) window.RandomGen.init();

  // 4. データ読み込みとレンダリング
  loadProfile();
  initArchiveSection();
  loadLogs();
});

// プロフィール情報の取得と反映
async function loadProfile() {
  try {
    const profile = await window.API.getProfile();
    if (!profile) return;

    const bioEl = document.getElementById('profile-bio');
    const quoteEl = document.getElementById('profile-quote');
    const keywordsEl = document.getElementById('profile-keywords');

    if (bioEl && profile.bio) bioEl.textContent = profile.bio;
    if (quoteEl && profile.vision) quoteEl.textContent = `“${profile.vision}”`;

    if (keywordsEl && profile.keywords) {
      const keywords = typeof profile.keywords === 'string' ? JSON.parse(profile.keywords) : profile.keywords;
      keywordsEl.innerHTML = keywords.map(kw => `<span class="keyword-chip">${kw}</span>`).join('');
    }

    // ステータスバーの反映
    if (profile.stats) {
      const stats = typeof profile.stats === 'string' ? JSON.parse(profile.stats) : profile.stats;
      Object.keys(stats).forEach(key => {
        const bar = document.getElementById(`stat-bar-${key}`);
        const val = document.getElementById(`stat-val-${key}`);
        if (bar) bar.style.width = `${stats[key]}%`;
        if (val) val.textContent = `${stats[key]}`;
      });
    }
  } catch (err) {
    console.warn('Profile load fallback:', err);
  }
}

// 制作物アーカイブ機能（フィルター・検索・詳細モーダル）
function initArchiveSection() {
  const container = document.getElementById('works-grid');
  const searchInput = document.getElementById('archive-search');
  const tabs = document.querySelectorAll('.category-tab-btn');
  const modalBackdrop = document.getElementById('work-modal-backdrop');
  const modalClose = document.getElementById('modal-close-btn');

  let currentCategory = '';
  let searchQuery = '';
  let worksCache = [];

  const renderWorks = (works) => {
    if (!container) return;
    if (!works || works.length === 0) {
      container.innerHTML = `
        <div style="grid-column: 1/-1; padding: 3rem; text-align: center; font-family: var(--font-mono); background: #fff; border: var(--border-bold);">
          該当する制作物が見つかりませんでした。別のキーワードでお試しください。
        </div>
      `;
      return;
    }

    container.innerHTML = works.map(w => {
      let techs = [];
      try {
        techs = typeof w.technologies === 'string' ? JSON.parse(w.technologies) : w.technologies;
      } catch (e) { techs = []; }

      return `
        <div class="work-card" data-id="${w.id}">
          <div class="work-card-header">
            <span class="brutal-badge brutal-badge-cyan">${w.category.toUpperCase()}</span>
            <span style="font-family: var(--font-mono); font-size: 0.75rem; color: #888;">#${w.id}</span>
          </div>
          <h3 class="work-card-title">${w.title}</h3>
          <p class="work-card-tagline">${w.tagline || w.description}</p>
          <div class="work-tech-tags">
            ${techs.map(t => `<span class="brutal-badge brutal-badge-yellow">${t}</span>`).join('')}
          </div>
          <div class="work-card-footer">
            <span style="color: var(--color-pink); font-weight: 700;">EXPLORE ARCHIVE ➔</span>
            <span>READ MORE</span>
          </div>
        </div>
      `;
    }).join('');

    // カードクリックでモーダルを開く
    container.querySelectorAll('.work-card').forEach(card => {
      card.addEventListener('click', () => {
        const id = card.dataset.id;
        const work = worksCache.find(item => item.id === id);
        if (work) openWorkModal(work);
      });
    });
  };

  const fetchAndFilter = async () => {
    try {
      const data = await window.API.getWorks(currentCategory, searchQuery);
      worksCache = data || [];
      renderWorks(worksCache);
    } catch (err) {
      if (container) {
        container.innerHTML = '<div style="padding: 2rem;">データ取得待機中 (サーバー起動待ち)</div>';
      }
    }
  };

  // タブ切り替え
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      tabs.forEach(t => t.classList.remove('active'));
      tab.classList.add('active');
      currentCategory = tab.dataset.category || '';
      fetchAndFilter();
    });
  });

  // 検索入力（デバウンス）
  let debounceTimer;
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        searchQuery = e.target.value.trim();
        fetchAndFilter();
      }, 250);
    });
  }

  // モーダル表示
  const openWorkModal = (work) => {
    if (!modalBackdrop) return;
    document.getElementById('modal-title').textContent = work.title;
    document.getElementById('modal-category').textContent = work.category.toUpperCase();
    document.getElementById('modal-tagline').textContent = work.tagline || '';
    document.getElementById('modal-desc').textContent = work.description;
    document.getElementById('modal-highlights').textContent = work.highlights || '現在ドキュメント整理中';

    let techs = [];
    try {
      techs = typeof work.technologies === 'string' ? JSON.parse(work.technologies) : work.technologies;
    } catch (e) { techs = []; }
    document.getElementById('modal-techs').innerHTML = techs.map(t => `<span class="brutal-badge brutal-badge-pink">${t}</span>`).join('');

    const linkArea = document.getElementById('modal-links');
    linkArea.innerHTML = '';
    if (work.github_url) {
      linkArea.innerHTML += `<a href="${work.github_url}" target="_blank" rel="noopener" class="brutal-btn brutal-btn-yellow" style="font-size: 0.95rem;">GITHUB REPO ↗</a>`;
    }
    if (work.demo_url) {
      linkArea.innerHTML += `<a href="${work.demo_url}" target="_blank" rel="noopener" class="brutal-btn brutal-btn-cyan" style="font-size: 0.95rem;">LIVE DEMO ↗</a>`;
    }

    modalBackdrop.classList.add('active');
  };

  if (modalClose) {
    modalClose.addEventListener('click', () => modalBackdrop.classList.remove('active'));
  }
  if (modalBackdrop) {
    modalBackdrop.addEventListener('click', (e) => {
      if (e.target === modalBackdrop) modalBackdrop.classList.remove('active');
    });
  }

  fetchAndFilter();
}

// 活動記録（Logs）の読み込み
async function loadLogs() {
  const container = document.getElementById('activity-logs-list');
  if (!container) return;

  try {
    const logs = await window.API.getLogs(5);
    if (!logs || logs.length === 0) return;

    container.innerHTML = logs.map(log => `
      <div style="border-bottom: 1px dashed var(--color-gray-300); padding: 1.1rem 0; display: flex; justify-content: space-between; align-items: center; gap: 1.25rem; font-family: var(--font-mono); font-size: 0.9rem; line-height: 1.6;">
        <div>
          <span style="color: var(--color-pink); font-weight: 800;">[${log.date}]</span>
          <span style="font-weight: 700; margin-left: 0.6rem; color: #171717;">${log.title}</span>
          <span style="display: block; font-size: 0.85rem; color: #4B5563; margin-top: 0.35rem; line-height: 1.6;">${log.content}</span>
        </div>
        ${log.metrics ? `<span class="brutal-badge brutal-badge-cyan" style="white-space: nowrap;">${log.metrics}</span>` : ''}
      </div>
    `).join('');
  } catch (e) {
    console.warn('Logs load fallback:', e);
  }
}
