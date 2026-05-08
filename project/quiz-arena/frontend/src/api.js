// API 客户端：统一封装后端接口调用
const API_BASE = import.meta.env.VITE_API_BASE || ''

async function request(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${await res.text()}`)
  }
  return res.json()
}

// 构造 querystring
function buildQuery(params) {
  const parts = []
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== null && v !== '' && v !== 'all') {
      parts.push(`${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
    }
  }
  return parts.length ? `?${parts.join('&')}` : ''
}

export const api = {
  // 获取题目列表（支持筛选）
  // filters: { category, difficulty, ids }
  getQuestions: (filters = {}) => {
    const qs = buildQuery({
      category: filters.category,
      difficulty: filters.difficulty,
      ids: Array.isArray(filters.ids) ? filters.ids.join(',') : filters.ids,
    })
    return request(`/api/questions${qs}`)
  },

  // 获取所有分类
  getCategories: () => request('/api/categories'),

  // 提交答案
  submitAnswer: (id, choice) =>
    request('/api/submit', {
      method: 'POST',
      body: JSON.stringify({ id, choice }),
    }),

  // 健康检查
  health: () => request('/api/health'),

  // ========== 编码题 ==========
  // 获取编码题列表
  getCodingQuestions: () => request('/api/coding/questions'),

  // 执行编码题代码
  runCodingCode: (id, code) =>
    request('/api/coding/run', {
      method: 'POST',
      body: JSON.stringify({ id, code }),
    }),
}
