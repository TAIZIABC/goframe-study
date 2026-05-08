// 选项列表组件
export default function Options({ options, choice, result, onSelect }) {
  const getOptionClass = (idx) => {
    const base =
      'group relative flex items-start gap-3 px-4 py-3.5 border rounded-xl text-[15px] ' +
      'transition-all duration-200 select-none'

    if (!result) {
      if (choice === idx) {
        return `${base} border-brand-400 bg-brand-50 dark:bg-brand-500/15 dark:border-brand-400 text-slate-800 dark:text-slate-100 cursor-wait shadow-glow`
      }
      return (
        `${base} border-slate-200/70 dark:border-white/10 bg-white/70 dark:bg-white/[0.03] ` +
        `text-slate-700 dark:text-slate-100 cursor-pointer ` +
        `hover:border-brand-400 hover:bg-brand-50/60 dark:hover:bg-brand-500/10 ` +
        `hover:translate-x-1 hover:shadow-glow`
      )
    }

    // 已出结果
    if (idx === result.answer) {
      return `${base} border-emerald-400 bg-emerald-50 dark:bg-emerald-500/15 dark:border-emerald-400/60 text-emerald-900 dark:text-emerald-200 cursor-not-allowed shadow-glow-green`
    }
    if (idx === choice && !result.correct) {
      return `${base} border-rose-400 bg-rose-50 dark:bg-rose-500/15 dark:border-rose-400/60 text-rose-900 dark:text-rose-200 cursor-not-allowed shadow-glow-red`
    }
    return `${base} border-slate-200/60 dark:border-white/5 bg-white/50 dark:bg-white/[0.02] text-slate-500 dark:text-slate-400 cursor-not-allowed opacity-70`
  }

  const getLabelClass = (idx) => {
    const base =
      'w-8 h-8 rounded-lg flex items-center justify-center font-bold text-sm shrink-0 transition-all'
    if (result) {
      if (idx === result.answer) return `${base} bg-emerald-500 text-white`
      if (idx === choice && !result.correct)
        return `${base} bg-rose-500 text-white`
      return `${base} bg-slate-200/80 dark:bg-white/10 text-slate-500 dark:text-slate-400`
    }
    if (choice === idx) {
      return `${base} bg-gradient-to-br from-brand-500 to-fuchsia-500 text-white`
    }
    return `${base} bg-slate-100 dark:bg-white/10 text-slate-600 dark:text-slate-300 group-hover:bg-brand-500 group-hover:text-white`
  }

  const getIcon = (idx) => {
    if (!result) return null
    if (idx === result.answer) return '✓'
    if (idx === choice && !result.correct) return '✗'
    return null
  }

  return (
    <div className="flex flex-col gap-2.5 mb-5">
      {options.map((text, idx) => {
        const icon = getIcon(idx)
        return (
          <div
            key={idx}
            className={getOptionClass(idx)}
            onClick={() => !result && onSelect(idx)}
          >
            <div className={getLabelClass(idx)}>
              {String.fromCharCode(65 + idx)}
            </div>
            <div className="flex-1 whitespace-pre-wrap leading-relaxed pt-0.5">
              {text.replace(/\\n/g, '\n')}
            </div>
            {icon && (
              <div
                className={`shrink-0 w-7 h-7 rounded-full flex items-center justify-center text-sm font-bold ${
                  icon === '✓'
                    ? 'bg-emerald-500 text-white'
                    : 'bg-rose-500 text-white'
                }`}
              >
                {icon}
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}
