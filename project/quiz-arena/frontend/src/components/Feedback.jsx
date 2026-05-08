// 答题反馈组件
export default function Feedback({ correct, explain }) {
  const wrapCls = correct
    ? 'bg-emerald-50/80 dark:bg-emerald-500/10 border-emerald-400/80 dark:border-emerald-500/40'
    : 'bg-rose-50/80 dark:bg-rose-500/10 border-rose-400/80 dark:border-rose-500/40'

  const iconBg = correct
    ? 'bg-emerald-500 shadow-glow-green'
    : 'bg-rose-500 shadow-glow-red'

  const titleColor = correct
    ? 'text-emerald-800 dark:text-emerald-200'
    : 'text-rose-800 dark:text-rose-200'

  return (
    <div
      className={`relative rounded-xl p-4 mb-5 border-l-4 animate-slide-down ${wrapCls}`}
    >
      <div className="flex items-start gap-3">
        <div
          className={`w-8 h-8 shrink-0 rounded-full flex items-center justify-center text-white font-bold ${iconBg}`}
        >
          {correct ? '✓' : '✗'}
        </div>
        <div className="flex-1 min-w-0">
          <div className={`font-bold text-[15px] mb-1.5 ${titleColor}`}>
            {correct ? '回答正确！+10 分' : '回答错误'}
          </div>
          <div className="text-sm text-slate-700 dark:text-slate-300 leading-relaxed">
            <strong className="text-slate-800 dark:text-slate-200">💡 解析：</strong>
            <span className="ml-1">{explain}</span>
          </div>
        </div>
      </div>
    </div>
  )
}
