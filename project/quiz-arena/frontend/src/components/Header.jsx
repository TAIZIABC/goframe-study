export default function Header({ theme, onToggleTheme }) {
  const isDark = theme === 'dark'
  return (
    <header className="relative text-center mb-10 animate-fade-in">
      {/* Logo + 标题 */}
      <div className="inline-flex items-center gap-3 mb-3">
        <span
          className="text-4xl drop-shadow-[0_0_12px_rgba(99,102,241,0.65)]"
          aria-hidden
        >
          ⚔️
        </span>
        <h1 className="text-[2.4rem] sm:text-5xl font-extrabold tracking-tight title-gradient">
          Quiz Arena
        </h1>
      </div>
      <p className="text-sm sm:text-[15px] text-slate-600 dark:text-slate-400">
        Go 并发 · 语法 · GoFrame 三位一体知识挑战场
      </p>

      {/* 主题切换按钮 */}
      <button
        onClick={onToggleTheme}
        className="absolute top-1 right-0 w-11 h-11 rounded-full
                   border border-slate-200/80 bg-white/70 backdrop-blur
                   text-xl flex items-center justify-center
                   hover:scale-110 hover:shadow-glow transition-all
                   dark:bg-white/5 dark:border-white/10"
        title={isDark ? '切换到浅色模式' : '切换到深色模式'}
        aria-label="切换主题"
      >
        {isDark ? '☀️' : '🌙'}
      </button>
    </header>
  )
}
