// 本地存储管理：错题集、收藏夹、答题历史
const KEYS = {
  WRONG: 'quiz-arena:wrong-ids', // 错题 ID 集合
  FAVORITE: 'quiz-arena:favorite-ids', // 收藏 ID 集合
  HISTORY: 'quiz-arena:history', // 答题历史
  THEME: 'quiz-arena:theme', // 主题：light/dark
}

function readSet(key) {
  try {
    const raw = localStorage.getItem(key)
    return new Set(raw ? JSON.parse(raw) : [])
  } catch {
    return new Set()
  }
}

function writeSet(key, set) {
  localStorage.setItem(key, JSON.stringify([...set]))
}

export const storage = {
  // 错题集
  getWrongIds: () => [...readSet(KEYS.WRONG)],
  addWrong: (id) => {
    const s = readSet(KEYS.WRONG)
    s.add(id)
    writeSet(KEYS.WRONG, s)
  },
  removeWrong: (id) => {
    const s = readSet(KEYS.WRONG)
    s.delete(id)
    writeSet(KEYS.WRONG, s)
  },
  clearWrong: () => writeSet(KEYS.WRONG, new Set()),
  wrongCount: () => readSet(KEYS.WRONG).size,

  // 收藏夹
  getFavoriteIds: () => [...readSet(KEYS.FAVORITE)],
  isFavorite: (id) => readSet(KEYS.FAVORITE).has(id),
  toggleFavorite: (id) => {
    const s = readSet(KEYS.FAVORITE)
    if (s.has(id)) s.delete(id)
    else s.add(id)
    writeSet(KEYS.FAVORITE, s)
    return s.has(id)
  },
  favoriteCount: () => readSet(KEYS.FAVORITE).size,

  // 答题历史（记录总答题数、总正确数）
  getHistory: () => {
    try {
      return (
        JSON.parse(localStorage.getItem(KEYS.HISTORY)) || {
          totalAnswered: 0,
          totalCorrect: 0,
          sessions: 0, // 完成的试炼场次数
        }
      )
    } catch {
      return { totalAnswered: 0, totalCorrect: 0, sessions: 0 }
    }
  },
  updateHistory: (answered, correct) => {
    const h = storage.getHistory()
    h.totalAnswered += answered
    h.totalCorrect += correct
    h.sessions += 1
    localStorage.setItem(KEYS.HISTORY, JSON.stringify(h))
  },
  resetHistory: () => localStorage.removeItem(KEYS.HISTORY),

  // 主题
  getTheme: () => localStorage.getItem(KEYS.THEME) || 'light',
  setTheme: (t) => localStorage.setItem(KEYS.THEME, t),
}
