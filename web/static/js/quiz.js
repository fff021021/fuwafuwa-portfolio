// Mini Quiz 制御スクリプト（サーバー採点連携）
class QuizController {
  constructor() {
    this.questions = [];
    this.currentIndex = 0;
    this.score = 0;
    this.answers = [];
    this.isLoading = false;
  }

  async init() {
    this.container = document.getElementById('quiz-container');
    if (!this.container) return;

    this.renderLoading();
    try {
      this.questions = await window.API.getQuizQuestions();
      this.currentIndex = 0;
      this.score = 0;
      this.answers = [];
      this.renderQuestion();
    } catch (err) {
      this.container.innerHTML = `
        <div style="padding: 1.5rem; text-align: center; color: var(--color-pink); font-family: var(--font-mono);">
          ⚠️ クイズデータの読み込みに失敗しました (サーバーが起動しているか確認してください)
          <br><br>
          <button class="brutal-btn brutal-btn-yellow" onclick="window.QuizApp.init()">再試行</button>
        </div>
      `;
    }
  }

  renderLoading() {
    this.container.innerHTML = `
      <div style="padding: 2rem; text-align: center; font-family: var(--font-mono);">
        <div class="anim-float" style="font-size: 2rem; margin-bottom: 0.5rem;">🎮</div>
        <div>LOADING QUIZ CHALLENGE...</div>
      </div>
    `;
  }

  renderQuestion() {
    if (!this.questions || this.questions.length === 0) {
      this.container.innerHTML = '<div>問題がありません</div>';
      return;
    }

    if (this.currentIndex >= this.questions.length) {
      this.renderResult();
      return;
    }

    const q = this.questions[this.currentIndex];
    const progress = `${this.currentIndex + 1} / ${this.questions.length}`;

    let options = [];
    try {
      options = typeof q.options === 'string' ? JSON.parse(q.options) : q.options;
    } catch (e) {
      options = q.options || [];
    }

    this.container.innerHTML = `
      <div class="quiz-box">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; border-bottom: var(--border-thin); padding-bottom: 0.5rem;">
          <span class="brutal-badge brutal-badge-pink">QUESTION ${progress}</span>
          <span class="brutal-badge brutal-badge-yellow">SCORE: ${this.score} PTS</span>
        </div>

        <h3 style="font-family: var(--font-display); font-size: 1.3rem; margin-bottom: 1.5rem; line-height: 1.4;">
          ${q.question}
        </h3>

        <div class="quiz-options-list" id="quiz-options">
          ${options.map((opt, idx) => `
            <button class="quiz-option-btn" data-index="${idx}">
              <span style="display: inline-block; width: 26px; font-family: var(--font-mono); font-weight: 800;">${['A', 'B', 'C', 'D'][idx]}.</span>
              ${this.escapeHTML(opt)}
            </button>
          `).join('')}
        </div>

        <div id="quiz-feedback" style="display: none; margin-top: 1.5rem; padding: 1rem; border: var(--border-thick); font-family: var(--font-mono); font-size: 0.9rem;">
        </div>

        <div id="quiz-action-area" style="display: none; margin-top: 1rem; text-align: right;">
          <button class="brutal-btn brutal-btn-pink" id="quiz-next-btn">次の問題へ ➔</button>
        </div>
      </div>
    `;

    // 選択肢のクリックイベント
    const optionBtns = this.container.querySelectorAll('.quiz-option-btn');
    optionBtns.forEach(btn => {
      btn.addEventListener('click', () => {
        const selectedIdx = parseInt(btn.dataset.index, 10);
        this.submitAnswer(q.id, selectedIdx, optionBtns, btn);
      });
    });
  }

  async submitAnswer(questionId, selectedIdx, allBtns, clickedBtn) {
    if (this.isLoading) return;
    this.isLoading = true;

    // ボタンの無効化
    allBtns.forEach(b => b.disabled = true);
    clickedBtn.classList.add('selected');

    try {
      // Goサーバー側へ回答を送信して採点
      const result = await window.API.submitQuizAnswer(questionId, selectedIdx);

      const feedback = document.getElementById('quiz-feedback');
      const actionArea = document.getElementById('quiz-action-area');
      const nextBtn = document.getElementById('quiz-next-btn');

      if (result.is_correct) {
        clickedBtn.classList.remove('selected');
        clickedBtn.classList.add('correct');
        this.score += result.points || 10;
        if (window.soundController) window.soundController.playSuccess();

        feedback.style.display = 'block';
        feedback.style.background = '#ECFDF5';
        feedback.style.borderColor = '#10B981';
        feedback.innerHTML = `
          <strong style="color: #059669; font-size: 1.1rem;">🎉 正解！ (+${result.points || 10}pts)</strong>
          <p style="margin-top: 0.5rem; color: #1F2937;">${result.explanation}</p>
        `;
      } else {
        clickedBtn.classList.remove('selected');
        clickedBtn.classList.add('wrong');
        if (window.soundController) window.soundController.playFail();

        feedback.style.display = 'block';
        feedback.style.background = '#FEF2F2';
        feedback.style.borderColor = 'var(--color-pink)';
        feedback.innerHTML = `
          <strong style="color: var(--color-pink); font-size: 1.1rem;">❌ 不正解...</strong>
          <p style="margin-top: 0.5rem; color: #1F2937;">${result.explanation}</p>
        `;
      }

      actionArea.style.display = 'block';
      nextBtn.onclick = () => {
        this.currentIndex++;
        this.isLoading = false;
        this.renderQuestion();
      };

    } catch (err) {
      alert('採点中にエラーが発生しました: ' + err.message);
      this.isLoading = false;
      allBtns.forEach(b => b.disabled = false);
    }
  }

  renderResult() {
    const totalQuestions = this.questions.length;
    const maxScore = totalQuestions * 10;
    const percentage = Math.round((this.score / maxScore) * 100);

    let rank = 'C';
    let rankColor = 'var(--color-pink)';
    let comment = 'まだまだ伸びしろがあります！Go言語とサーバーサイドの世界を探求しましょう！';

    if (percentage >= 90) {
      rank = 'S';
      rankColor = 'var(--color-yellow)';
      comment = 'サーバーエンジニア適性MAX！並行処理とネットワークのスペシャリストです！';
      if (window.soundController) window.soundController.playSecret();
      if (window.EasterEggs) {
        window.EasterEggs.unlockSecret('quiz_master', 'クイズマスター！見事に高得点を獲得しました。');
      }
    } else if (percentage >= 70) {
      rank = 'A';
      rankColor = 'var(--color-cyan)';
      comment = '素晴らしい知識量！実戦的なサーバー開発に即戦力として挑めます！';
      if (window.soundController) window.soundController.playSuccess();
    } else if (percentage >= 50) {
      rank = 'B';
      rankColor = '#10B981';
      comment = '基礎はバッチリ！更なるパフォーマンスチューニングを目指そう！';
    }

    this.container.innerHTML = `
      <div class="quiz-box" style="text-align: center; padding: 2.5rem 1.5rem;">
        <span class="brutal-badge brutal-badge-dark" style="margin-bottom: 1rem;">CHALLENGE COMPLETE</span>
        
        <h3 style="font-family: var(--font-display); font-size: 2.4rem; margin-bottom: 0.5rem;">
          FINAL SCORE: <span style="color: var(--color-pink);">${this.score}</span> / ${maxScore}
        </h3>

        <div style="display: inline-block; margin: 1.5rem 0; padding: 1.5rem 2.5rem; background: var(--color-black); color: ${rankColor}; border: var(--border-bold); box-shadow: var(--shadow-yellow); font-family: var(--font-display); font-size: 5rem; line-height: 1;">
          RANK: ${rank}
        </div>

        <p style="font-family: var(--font-mono); font-size: 1rem; color: var(--color-gray-900); max-width: 480px; margin: 0 auto 2rem; line-height: 1.6;">
          ${comment}
        </p>

        <div style="display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap;">
          <button class="brutal-btn brutal-btn-yellow" onclick="window.QuizApp.init()">
            🔄 もう一度挑戦する
          </button>
          <a href="#03-archive" class="brutal-btn brutal-btn-cyan">
            📁 制作物アーカイブを見る ➔
          </a>
        </div>
      </div>
    `;
  }

  escapeHTML(str) {
    if (!str) return '';
    return str.replace(/[&<>'"]/g, 
      tag => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' }[tag] || tag)
    );
  }
}

window.QuizApp = new QuizController();
