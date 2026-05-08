import { useState, useEffect } from 'react'
import { storage } from '../storage'

const DIFFICULTIES = [
  { value: 0, label: '全部难度' },
  { value: 1, label: '简单' },
  { value: 2, label: '中等' },
  { value: 3, label: '困难' },
]

const MODES = {
  NORMAL: 'normal',
  WRONG: 'wrong',
  FAVORITE: 'favorite',
}

export default function StartScreen({ categories, onStart, onStartCoding }) {
  const [category, setCategory] = useState('all')
  const [difficulty, setDifficulty] = useState(0)
  const [mode, setMode] = useState(MODES.NORMAL)
  const [wrongCount, setWrongCount] = useState(0)
  const [favoriteCount, setFavoriteCount] = useState(0)
  const [history, setHistory] = useState({ totalAnswered: 0, totalCorrect: 0, sessions: 0 })

  useEffect(() => {
    setWrongCount(storage.wrongCount())
    setFavoriteCount(storage.favoriteCount())
    setHistory(storage.getHistory())
  }, [])

  const total = categories.reduce((s, c) => s + c.count, 0)
  const accuracy = history.totalAnswered > 0
    ? Math.round((history.totalCorrect / history.totalAnswered) * 100)
    : 0

  const handleStart = () => {
    const filters = {}
    if (mode === MODES.WRONG) {
      filters.ids = storage.getWrongIds()
      if (filters.ids.length === 0) {
        alert('错题集是空的！答错一道题后会自动加入')
        return
      }
    } else if (mode === MODES.FAVORITE) {
      filters.ids = storage.getFavoriteIds()
      if (filters.ids.length === 0) {
        alert('收藏夹是空的！在答题时点击 ⭐ 可收藏题目')
        return
      }
    } else {
      if (category !== 'all') filters.category = category
      if (difficulty > 0) filters.difficulty = difficulty
    }
    onStart(filters)
  }

  return (
    <section className="animate-fade-in space-y-5">
      {/* 个人数据卡片 */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-bold flex items-center gap-2 text-slate-800 dark:text-slate-100">
            <span>📊</span> 我的数据
          </h3>
          <span className="text-xs text-slate-500 dark:text-slate-400">
            持续积累 · 自动记录
          </span>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <StatCard label="完成场次" value={history.sessions} icon="🏁" />
          <StatCard label="累计答题" value={history.totalAnswered} icon="📝" />
          <StatCard
            label="正确率"
            value={`${accuracy}%`}
            icon="🎯"
            tone={accuracy >= 80 ? 'green' : accuracy >= 60 ? 'amber' : 'default'}
          />
          <StatCard
            label="错题"
            value={wrongCount}
            icon="❗"
            tone={wrongCount > 0 ? 'red' : 'default'}
          />
        </div>
      </div>

      {/* 模式选择 */}
      <div className="card">
        <h3 className="text-lg font-bold mb-4 text-slate-800 dark:text-slate-100">
          🎯 选择模式
        </h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6">
          <ModeButton
            active={mode === MODES.NORMAL}
            onClick={() => setMode(MODES.NORMAL)}
            icon="🎲"
            label="标准模式"
            sub={`共 ${total} 题`}
          />
          <ModeButton
            active={mode === MODES.WRONG}
            onClick={() => setMode(MODES.WRONG)}
            icon="📝"
            label="错题重练"
            sub={`${wrongCount} 题`}
            disabled={wrongCount === 0}
          />
          <ModeButton
            active={mode === MODES.FAVORITE}
            onClick={() => setMode(MODES.FAVORITE)}
            icon="⭐"
            label="我的收藏"
            sub={`${favoriteCount} 题`}
            disabled={favoriteCount === 0}
          />
        </div>

        {mode === MODES.NORMAL && (
          <div className="space-y-5 animate-fade-in">
            {/* 分类筛选 */}
            <div>
              <div className="text-sm font-semibold mb-3 text-slate-700 dark:text-slate-300 flex items-center gap-2">
                <span className="w-1 h-4 rounded bg-brand-500" /> 分类
              </div>
              <div className="flex flex-wrap gap-2">
                <FilterChip
                  active={category === 'all'}
                  onClick={() => setCategory('all')}
                  label={`全部 · ${total}`}
                />
                {categories.map((c) => (
                  <FilterChip
                    key={c.name}
                    active={category === c.name}
                    onClick={() => setCategory(c.name)}
                    label={`${c.name} · ${c.count}`}
                  />
                ))}
              </div>
            </div>

            {/* 难度筛选 */}
            <div>
              <div className="text-sm font-semibold mb-3 text-slate-700 dark:text-slate-300 flex items-center gap-2">
                <span className="w-1 h-4 rounded bg-neon-pink" /> 难度
              </div>
              <div className="flex flex-wrap gap-2">
                {DIFFICULTIES.map((d) => (
                  <FilterChip
                    key={d.value}
                    active={difficulty === d.value}
                    onClick={() => setDifficulty(d.value)}
                    label={d.label}
                    difficultyValue={d.value}
                  />
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* 开始按钮 */}
      <div className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-2">
        <button
          onClick={handleStart}
          className="btn btn-primary text-base px-12 py-4 rounded-2xl animate-pulse-glow"
        >
          开始试炼 →
        </button>
        <button
          onClick={onStartCoding}
          className="btn btn-ghost text-base px-10 py-4 rounded-2xl"
        >
          💻 编码练习
        </button>
      </div>
      <div className="text-xs text-slate-500 dark:text-slate-500 mt-3 text-center">
        选择题 · 多分类筛选 · 键盘操作 &nbsp;|&nbsp; 编码题 · 在线编写 Go 代码 · 实时运行验证
      </div>
    </section>
  )
}

/* ---------- 子组件 ---------- */

function StatCard({ label, value, icon, tone = 'default' }) {
  const palettes = {
    default:
      'bg-slate-50/80 dark:bg-white/5 border-slate-200/70 dark:border-white/10 text-brand-600 dark:text-brand-300',
    green:
      'bg-emerald-50 dark:bg-emerald-500/10 border-emerald-200 dark:border-emerald-500/20 text-emerald-600 dark:text-emerald-300',
    amber:
      'bg-amber-50 dark:bg-amber-500/10 border-amber-200 dark:border-amber-500/20 text-amber-600 dark:text-amber-300',
    red:
      'bg-rose-50 dark:bg-rose-500/10 border-rose-200 dark:border-rose-500/20 text-rose-600 dark:text-rose-300',
  }
  return (
    <div
      className={`relative text-center p-4 rounded-xl border transition-all hover:-translate-y-0.5 ${palettes[tone]}`}
    >
      <div className="absolute top-2 right-2 text-xs opacity-70">{icon}</div>
      <div className="text-[1.75rem] font-extrabold leading-none tracking-tight">
        {value}
      </div>
      <div className="text-xs text-slate-500 dark:text-slate-400 mt-1.5">
        {label}
      </div>
    </div>
  )
}

function ModeButton({ active, onClick, icon, label, sub, disabled }) {
  const base =
    'group relative rounded-xl p-4 text-center border transition-all duration-200 cursor-pointer overflow-hidden'
  let cls
  if (disabled) {
    cls = `${base} border-slate-200 dark:border-white/5 bg-slate-50/60 dark:bg-white/[0.02] opacity-55 cursor-not-allowed`
  } else if (active) {
    cls = `${base} border-transparent text-white shadow-glow`
  } else {
    cls = `${base} border-slate-200/80 dark:border-white/10 bg-white/70 dark:bg-white/[0.03] hover:-translate-y-0.5 hover:border-brand-400 dark:hover:border-brand-400/50 hover:shadow-glow`
  }
  return (
    <button
      className={cls}
      onClick={disabled ? undefined : onClick}
      disabled={disabled}
      style={
        active
          ? {
              backgroundImage:
                'linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #ec4899 100%)',
            }
          : undefined
      }
    >
      <div className="text-3xl mb-1.5">{icon}</div>
      <div
        className={`font-bold text-[15px] ${
          active ? 'text-white' : 'text-slate-800 dark:text-slate-100'
        }`}
      >
        {label}
      </div>
      <div
        className={`text-xs mt-1 ${
          active ? 'text-white/85' : 'text-slate-500 dark:text-slate-400'
        }`}
      >
        {sub}
      </div>
    </button>
  )
}

function FilterChip({ active, onClick, label }) {
  const base =
    'chip cursor-pointer select-none hover:-translate-y-0.5 active:translate-y-0'
  const cls = active
    ? `${base} text-white border-transparent shadow-glow`
    : `${base} bg-white/70 dark:bg-white/5 border-slate-200/70 dark:border-white/10 text-slate-700 dark:text-slate-200 hover:border-brand-400 hover:text-brand-600 dark:hover:text-brand-300`
  return (
    <button
      onClick={onClick}
      className={cls}
      style={
        active
          ? {
              backgroundImage:
                'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
            }
          : undefined
      }
    >
      {label}
    </button>
  )
}
