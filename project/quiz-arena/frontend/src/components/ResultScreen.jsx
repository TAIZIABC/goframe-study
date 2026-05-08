import { useState } from 'react'

// 结果展示组件
export default function ResultScreen({ records, total, elapsedSec, onRestart, onHome }) {
  const [showReview, setShowReview] = useState(false)

  const correct = records.filter((r) => r.correct).length
  const wrong = total - correct
  const rate = total > 0 ? Math.round((correct / total) * 100) : 0

  let emoji, title, tone
  if (rate === 100) {
    emoji = '🏆'
    title = '完美通关！Go 大师就是你！'
    tone = 'from-amber-400 via-pink-500 to-fuchsia-500'
  } else if (rate >= 80) {
    emoji = '🎉'
    title = '表现优秀！基础十分扎实'
    tone = 'from-emerald-400 via-teal-400 to-cyan-500'
  } else if (rate >= 60) {
    emoji = '💪'
    title = '还不错，继续加油！'
    tone = 'from-brand-400 via-fuchsia-400 to-pink-500'
  } else {
    emoji = '📚'
    title = '再接再厉，多练习就能掌握！'
    tone = 'from-slate-400 via-brand-400 to-fuchsia-400'
  }

  const formatTime = (sec) => {
    if (sec < 60) return `${sec}s`
    return `${Math.floor(sec / 60)}m ${sec % 60}s`
  }

  return (
    <section className="animate-fade-in space-y-5">
      <div className="card-glow text-center">
        <div className="text-7xl mb-3 animate-bounce-slow drop-shadow-[0_6px_24px_rgba(0,0,0,0.2)]">
          {emoji}
        </div>
        <h2 className="text-2xl sm:text-3xl font-extrabold mb-6">
          <span className={`bg-gradient-to-r ${tone} bg-clip-text text-transparent`}>
            {title}
          </span>
        </h2>

        {/* 大分数环 */}
        <div className="relative inline-flex flex-col items-center justify-center mb-7">
          <div className="text-6xl sm:text-7xl font-extrabold tracking-tight tabular-nums">
            <span className={`bg-gradient-to-br ${tone} bg-clip-text text-transparent`}>
              {correct}
            </span>
            <span className="text-slate-300 dark:text-slate-600 font-normal mx-3">/</span>
            <span className="text-slate-400 dark:text-slate-500">{total}</span>
          </div>
          <div className="text-slate-600 dark:text-slate-400 text-sm mt-2">
            正确率{' '}
            <span className="text-brand-500 dark:text-brand-300 font-bold text-base">
              {rate}%
            </span>
          </div>
        </div>

        <div className="flex justify-center gap-3 mb-7 flex-wrap">
          <SummaryItem label="答对" icon="✅" value={correct} tone="green" />
          <SummaryItem label="答错" icon="❌" value={wrong} tone="red" />
          <SummaryItem label="耗时" icon="⏱" value={formatTime(elapsedSec)} tone="brand" />
        </div>

        {wrong > 0 && (
          <div className="mb-6 text-sm text-slate-500 dark:text-slate-400 px-3 py-2 rounded-lg bg-amber-50/60 dark:bg-amber-500/10 border border-amber-200/60 dark:border-amber-500/20 inline-block">
            💡 答错的题已加入错题集，可在首页「错题重练」专项提升
          </div>
        )}

        <div className="flex justify-center gap-3 flex-wrap">
          <button onClick={onRestart} className="btn btn-primary">
            再来一次
          </button>
          <button onClick={() => setShowReview((v) => !v)} className="btn btn-ghost">
            {showReview ? '收起回顾' : '查看回顾'}
          </button>
          <button onClick={onHome} className="btn btn-ghost">
            返回首页
          </button>
        </div>
      </div>

      {showReview && (
        <div className="card animate-slide-up">
          <h3 className="text-lg font-bold mb-4 text-slate-800 dark:text-slate-100 flex items-center gap-2">
            📋 答题回顾
            <span className="text-xs font-normal text-slate-500 dark:text-slate-400">
              · 共 {records.length} 题
            </span>
          </h3>
          <div className="space-y-3">
            {records.map((r, i) => (
              <ReviewItem key={i} index={i} record={r} />
            ))}
          </div>
        </div>
      )}
    </section>
  )
}

function SummaryItem({ label, value, icon, tone = 'brand' }) {
  const palettes = {
    brand:
      'bg-brand-50/60 dark:bg-brand-500/10 border-brand-200/70 dark:border-brand-500/20 text-brand-600 dark:text-brand-300',
    green:
      'bg-emerald-50/60 dark:bg-emerald-500/10 border-emerald-200/70 dark:border-emerald-500/20 text-emerald-600 dark:text-emerald-300',
    red:
      'bg-rose-50/60 dark:bg-rose-500/10 border-rose-200/70 dark:border-rose-500/20 text-rose-600 dark:text-rose-300',
  }
  return (
    <div
      className={`px-6 py-3.5 rounded-xl border min-w-[120px] transition-all hover:-translate-y-0.5 ${palettes[tone]}`}
    >
      <div className="text-xs opacity-80 mb-1">
        {icon} {label}
      </div>
      <div className="text-2xl font-extrabold tabular-nums">{value}</div>
    </div>
  )
}

function ReviewItem({ index, record }) {
  const userLabel = String.fromCharCode(65 + record.userChoice)
  const correctLabel = String.fromCharCode(65 + record.correctAnswer)
  const correct = record.correct

  return (
    <div
      className={`p-4 rounded-xl border transition-all hover:-translate-y-0.5 ${
        correct
          ? 'border-emerald-200/70 dark:border-emerald-500/20 bg-emerald-50/40 dark:bg-emerald-500/5'
          : 'border-rose-200/70 dark:border-rose-500/20 bg-rose-50/40 dark:bg-rose-500/5'
      }`}
    >
      <div className="flex justify-between items-start gap-2 mb-2">
        <span className="font-semibold text-slate-800 dark:text-slate-100 text-sm leading-relaxed">
          <span className="text-slate-400 mr-1.5">#{index + 1}</span>
          <span className="chip mr-1.5 !px-2 !py-0.5 text-[11px] bg-brand-50 dark:bg-brand-500/15 text-brand-600 dark:text-brand-300 border-brand-200/70 dark:border-brand-500/30">
            {record.category}
          </span>
          {record.title}
        </span>
        <span
          className={`chip whitespace-nowrap ${
            correct
              ? 'bg-emerald-100 text-emerald-700 border-emerald-200 dark:bg-emerald-500/15 dark:text-emerald-300 dark:border-emerald-500/30'
              : 'bg-rose-100 text-rose-700 border-rose-200 dark:bg-rose-500/15 dark:text-rose-300 dark:border-rose-500/30'
          }`}
        >
          {correct ? '✓ 正确' : '✗ 错误'}
        </span>
      </div>
      <div className="text-sm text-slate-600 dark:text-slate-400 mb-2">
        你的答案：<strong className={correct ? 'text-emerald-600 dark:text-emerald-300' : 'text-rose-600 dark:text-rose-300'}>{userLabel}</strong>
        <span className="mx-2 text-slate-300 dark:text-slate-600">·</span>
        正确答案：
        <strong className="text-emerald-600 dark:text-emerald-300">
          {correctLabel}. {record.options[record.correctAnswer]}
        </strong>
      </div>
      <div className="text-slate-600 dark:text-slate-300 text-sm p-3 rounded-lg bg-white/60 dark:bg-white/5 border-l-2 border-brand-500 leading-relaxed">
        {record.explain}
      </div>
    </div>
  )
}
