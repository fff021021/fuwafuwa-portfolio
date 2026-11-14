// バックエンドREST API通信クライアント
const API = {
  baseUrl: '/api',

  async get(endpoint, params = {}) {
    try {
      const url = new URL(this.baseUrl + endpoint, window.location.origin);
      Object.keys(params).forEach(key => {
        if (params[key] !== undefined && params[key] !== null && params[key] !== '') {
          url.searchParams.append(key, params[key]);
        }
      });

      const res = await fetch(url.toString(), {
        headers: {
          'Accept': 'application/json'
        }
      });

      if (!res.ok) {
        throw new Error(`API Error: ${res.status} ${res.statusText}`);
      }
      return await res.json();
    } catch (err) {
      console.error(`Failed to GET ${endpoint}:`, err);
      throw err;
    }
  },

  async post(endpoint, data = {}) {
    try {
      const res = await fetch(this.baseUrl + endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify(data)
      });

      if (!res.ok) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error || `API Error: ${res.status}`);
      }
      return await res.json();
    } catch (err) {
      console.error(`Failed to POST ${endpoint}:`, err);
      throw err;
    }
  },

  // 各種エンドポイント
  getProfile() {
    return this.get('/profile');
  },

  getWorks(category = '', query = '') {
    return this.get('/works', { category, q: query });
  },

  getLogs(limit = 10, category = '') {
    return this.get('/logs', { limit, category });
  },

  getQuizQuestions() {
    return this.get('/quiz/questions');
  },

  submitQuizAnswer(questionId, selectedIndex) {
    return this.post('/quiz/answer', {
      question_id: questionId,
      selected_index: selectedIndex
    });
  },

  getRandomContent() {
    return this.get('/random');
  },

  getHealth() {
    return this.get('/health');
  }
};

window.API = API;
