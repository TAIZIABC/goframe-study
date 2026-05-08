import { useState, useEffect } from 'react'
import { api } from '../api'
import { storage } from '../storage'
import { highlightGo } from '../highlight'
import Options from './Options'
import Feedback from './Feedback'

const DIFFICULTY_LABELS = {
  1: {
    text: '简单',
    cls:
      'bg-emerald-100/80 text-emerald-700 border-emerald-200 dark:bg-emerald-500/15 dark:text-emerald-300 dark:border-emerald-500/30',
  },
  2: {
    text: '中等',
    cls:
      'bg-amber-100/80 text-amber-700 border-amber-200 dark:bg-amber-500/15 dark:text-amber-300 dark:border-amber-500/30',
  },
  3: {
    text: '困难',
    cls:
      'bg-rose-100/80 text-rose-700 border-rose-200 dark:bg-rose-500/15 dark:text-rose-300 dark:border-rose-500/30',
  },
}

export default function QuizScreen({
  question,
  current,
  total,
  score,
  onAnswered,
  onNext,
  onExit,
}) {
  const [choice, setChoice] = useState(null)
  const [result, setResult] = useState(null)
  const [submitting, setSubmitting] = useState(false)
  const [favorited, setFavorited] = useState(false)

  useEffect(() => {
    setChoice(null)
    setResult(null)
    setFavorited(storage.isFavorite(question.id))
  }, [question.id])

  const handleSelect = async (idx) => {
    if (choice !== null || submitting) return
    setChoice(idx)
    setSubmitting(true)
    try {
      const data = await api.submitAnswer(question.id, idx)
      setResult(data)
      onAnswered({
        questionId: question.id,
        title: question.title,
        category: question.category,
        difficulty: question.difficulty,
        options: question.options,
        userChoice: idx,
        correctAnswer: data.answer,
        correct: data.correct,
        explain: data.explain,
      })
    } catch (e) {
      alert('提交失败：' + e.message)
      setChoice(null)
    } finally {
      setSubmitting(false)
    }
  }

  const toggleFavorite = () => {
    const fav = storage.toggleFavorite(question.id)
    setFavorited(fav)
  }

  useEffect(() => {
    const handler = (e) => {
      if (['1', '2', '3', '4'].includes(e.key) && choice === null) {
        const idx = parseInt(e.key) - 1
        if (idx < question.options.length) handleSelect(idx)
      } else if (e.key === 'Enter' && result) {
        onNext()
      } else if (e.key === 'f' || e.key === 'F') {
        toggleFavorite()
      } else if (e.key === 'Escape') {
        if (confirm('退出本次试炼？')) onExit()
      }
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [choice, result, question.options.length])

  const progress = ((current + 1) / total) * 100
  const isLast = current === total - 1
  const diff = DIFFICULTY_LABELS[question.difficulty] || DIFFICULTY_LABELS[2]

  return (
    <section className="animate-fade-in">
      {/* 顶部进度条 */}
      <div
        className="flex items-center gap-4 mb-5 px-5 py-3 rounded-2xl
                   bg-white/75 dark:bg-white/[0.04] backdrop-blur-xl
                   border border-slate-200/70 dark:border-white/10 shadow-card"
      >
        <button
          onClick={() => {
            if (confirm('退出本次试炼？')) onExit()
          }}
          className="text-slate-600 dark:text-slate-300 hover:text-brand-600
                     dark:hover:text-brand-300 text-sm font-medium shrink-0"
          title="Esc 退出"
        >
          ← 退出
        </button>

        <div className="flex items-center gap-3 flex-1 min-w-0">
          <div className="flex-1 h-2 rounded-full overflow-hidden bg-slate-200/70 dark:bg-white/10 relative">
            <div
              className="h-full rounded-full transition-all duration-500 bg-progress-gradient relative"
              style={{ width: `${progress}%` }}
            >
              <div className="absolute inset-0 shimmer-bar" />
            </div>
          </div>
          <span className="text-sm font-semibold text-slate-700 dark:text-slate-200 whitespace-nowrap tabular-nums">
            {current + 1} <span className="text-slate-400">/ {total}</span>
          </span>
        </div>

        <div className="flex items-center gap-2 shrink-0">
          <span className="text-xs text-slate-500 dark:text-slate-400">得分</span>
          <span
            className="text-2xl font-extrabold leading-none tabular-nums
                       bg-gradient-to-br from-amber-400 to-pink-500 bg-clip-text text-transparent
                       drop-shadow-[0_2px_8px_rgba(251,191,36,0.35)]"
          >
            {score}
          </span>
        </div>
      </div>

      {/* 题目卡片 */}
      <div className="card-glow min-h-[420px] flex flex-col">
        <div className="flex justify-between items-center mb-5">
          <div className="flex items-center gap-2">
            <span className="chip bg-brand-50 dark:bg-brand-500/15 text-brand-600 dark:text-brand-300 border-brand-200/70 dark:border-brand-500/30">
              {question.category}
            </span>
            <span className={`chip ${diff.cls}`}>{diff.text}</span>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={toggleFavorite}
              className={`text-2xl leading-none transition-all hover:scale-125 ${
                favorited ? 'drop-shadow-[0_0_10px_rgba(251,191,36,0.7)]' : ''
              }`}
              title={favorited ? '取消收藏 (F)' : '收藏本题 (F)'}
              aria-label="收藏"
            >
              {favorited ? '⭐' : '☆'}
            </button>
            <span className="text-slate-400 text-sm font-semibold tabular-nums">
              #{current + 1}
            </span>
          </div>
        </div>

        <h3 className="text-lg sm:text-xl font-semibold text-slate-800 dark:text-slate-100 mb-5 leading-relaxed text-balance">
          {question.title}
        </h3>

        {question.code && (
          <pre
            className="code-block mb-5"
            dangerouslySetInnerHTML={{ __html: highlightGo(question.code) }}
          />
        )}

        <Options
          options={question.options}
          choice={choice}
          result={result}
          onSelect={handleSelect}
        />

        {result && <Feedback correct={result.correct} explain={result.explain} />}

        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mt-auto pt-5 border-t border-slate-200/70 dark:border-white/10">
          <div className="text-xs text-slate-500 dark:text-slate-500 flex flex-wrap gap-x-2 gap-y-1">
            <Kbd>1-4</Kbd>选择
            <Kbd>Enter</Kbd>下一题
            <Kbd>F</Kbd>收藏
            <Kbd>Esc</Kbd>退出
          </div>
          {result && (
            <button onClick={onNext} className="btn btn-primary">
              {isLast ? '查看结果 🏁' : '下一题 →'}
            </button>
          )}
        </div>
      </div>
    </section>
  )
}

function Kbd({ children }) {
  return (
    <kbd className="inline-flex items-center px-1.5 py-0.5 rounded border border-slate-300/70 dark:border-white/15 bg-slate-100/70 dark:bg-white/5 text-[11px] font-mono text-slate-600 dark:text-slate-300 mx-0.5">
      {children}
    </kbd>
  )
}
