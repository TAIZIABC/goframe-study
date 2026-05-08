import { useState, useEffect, useRef } from 'react'
import { api } from '../api'
import { highlightGo } from '../highlight'

const DIFFICULTY_LABELS = {
  1: {
    text: '简单',
    cls: 'bg-emerald-100/80 text-emerald-700 border-emerald-200 dark:bg-emerald-500/15 dark:text-emerald-300 dark:border-emerald-500/30',
  },
  2: {
    text: '中等',
    cls: 'bg-amber-100/80 text-amber-700 border-amber-200 dark:bg-amber-500/15 dark:text-amber-300 dark:border-amber-500/30',
  },
  3: {
    text: '困难',
    cls: 'bg-rose-100/80 text-rose-700 border-rose-200 dark:bg-rose-500/15 dark:text-rose-300 dark:border-rose-500/30',
  },
}

export default function CodingScreen({ onExit }) {
  const [questions, setQuestions] = useState([])
  const [current, setCurrent] = useState(0)
  const [code, setCode] = useState('')
  const [running, setRunning] = useState(false)
  const [result, setResult] = useState(null)
  const [showHint, setShowHint] = useState(false)
  const [showSolution, setShowSolution] = useState(false)
  const [passedIds, setPassedIds] = useState(new Set())
  const [loading, setLoading] = useState(true)
  const textareaRef = useRef(null)

  // 加载编码题
  useEffect(() => {
    api.getCodingQuestions()
      .then((data) => {
        setQuestions(data)
        if (data.length > 0) setCode(data[0].template)
        setLoading(false)
      })
      .catch(() => setLoading(false))
  }, [])

  // 切题时重置
  useEffect(() => {
    if (questions[current]) {
      setCode(questions[current].template)
      setResult(null)
      setShowHint(false)
      setShowSolution(false)
    }
  }, [current, questions])

  // Tab 键支持
  const handleKeyDown = (e) => {
    if (e.key === 'Tab') {
      e.preventDefault()
      const { selectionStart, selectionEnd } = e.target
      const newCode = code.substring(0, selectionStart) + '\t' + code.substring(selectionEnd)
      setCode(newCode)
      requestAnimationFrame(() => {
        e.target.selectionStart = e.target.selectionEnd = selectionStart + 1
      })
    }
  }

  // 执行代码
  const handleRun = async () => {
    if (running) return
    setRunning(true)
    setResult(null)
    try {
      const data = await api.runCodingCode(questions[current].id, code)
      setResult(data)
      if (data.pass) {
        setPassedIds((prev) => new Set([...prev, questions[current].id]))
      }
    } catch (e) {
      setResult({ pass: false, error: '请求失败：' + e.message, output: '' })
    } finally {
      setRunning(false)
    }
  }

  // 重置代码
  const handleReset = () => {
    if (confirm('重置将清除你的代码修改，确定？')) {
      setCode(questions[current].template)
      setResult(null)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <span className="inline-block w-5 h-5 rounded-full border-2 border-brand-500 border-t-transparent animate-spin" />
        <span className="ml-3 text-slate-600 dark:text-slate-300">加载编码题...</span>
      </div>
    )
  }

  if (questions.length === 0) {
    return (
      <div className="card text-center py-10">
        <div className="text-4xl mb-3">📝</div>
        <div className="text-lg font-bold text-slate-700 dark:text-slate-200 mb-2">暂无编码题</div>
        <button onClick={onExit} className="btn btn-ghost mt-4">返回首页</button>
      </div>
    )
  }

  const q = questions[current]
  const diff = DIFFICULTY_LABELS[q.difficulty] || DIFFICULTY_LABELS[2]
  const passed = passedIds.has(q.id)

  return (
    <section className="animate-fade-in">
      {/* 顶部导航 */}
      <div className="flex items-center gap-4 mb-5 px-5 py-3 rounded-2xl bg-white/75 dark:bg-white/[0.04] backdrop-blur-xl border border-slate-200/70 dark:border-white/10 shadow-card">
        <button
          onClick={() => { if (confirm('退出编码练习？')) onExit() }}
          className="text-slate-600 dark:text-slate-300 hover:text-brand-600 dark:hover:text-brand-300 text-sm font-medium shrink-0"
        >
          ← 返回
        </button>

        {/* 题目选择器 */}
        <div className="flex items-center gap-2 flex-1 overflow-x-auto">
          {questions.map((item, i) => (
            <button
              key={item.id}
              onClick={() => setCurrent(i)}
              className={`shrink-0 w-9 h-9 rounded-lg text-sm font-bold transition-all ${
                i === current
                  ? 'text-white shadow-glow'
                  : passedIds.has(item.id)
                    ? 'bg-emerald-100 dark:bg-emerald-500/20 text-emerald-700 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-500/30'
                    : 'bg-white/70 dark:bg-white/5 text-slate-600 dark:text-slate-300 border border-slate-200/70 dark:border-white/10 hover:border-brand-400'
              }`}
              style={i === current ? { backgroundImage: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #ec4899 100%)' } : undefined}
            >
              {passedIds.has(item.id) && i !== current ? '✓' : i + 1}
            </button>
          ))}
        </div>

        <div className="shrink-0 text-sm font-semibold text-slate-700 dark:text-slate-200 tabular-nums">
          <span className="text-emerald-600 dark:text-emerald-300">{passedIds.size}</span>
          <span className="text-slate-400"> / {questions.length}</span>
        </div>
      </div>

      {/* 题目描述 */}
      <div className="card mb-4">
        <div className="flex items-center gap-2 mb-3">
          <span className="chip bg-brand-50 dark:bg-brand-500/15 text-brand-600 dark:text-brand-300 border-brand-200/70 dark:border-brand-500/30">
            {q.category}
          </span>
          <span className={`chip ${diff.cls}`}>{diff.text}</span>
          {passed && (
            <span className="chip bg-emerald-100 dark:bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-500/30">
              ✓ 已通过
            </span>
          )}
          <span className="ml-auto text-slate-400 text-sm font-semibold">#{current + 1}</span>
        </div>

        <h3 className="text-lg font-bold text-slate-800 dark:text-slate-100 mb-3">
          {q.title}
        </h3>
        <div className="text-sm text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-wrap">
          {q.desc}
        </div>

        {/* 提示 */}
        <div className="flex gap-2 mt-4">
          <button
            onClick={() => setShowHint(!showHint)}
            className="text-xs text-brand-500 hover:text-brand-700 dark:text-brand-300 dark:hover:text-brand-200"
          >
            {showHint ? '隐藏提示' : '💡 查看提示'}
          </button>
        </div>
        {showHint && (
          <div className="mt-2 p-3 rounded-lg bg-amber-50/60 dark:bg-amber-500/10 border border-amber-200/60 dark:border-amber-500/20 text-sm text-amber-800 dark:text-amber-200 animate-slide-down">
            💡 {q.hint}
          </div>
        )}
      </div>

      {/* 代码编辑器 */}
      <div className="card mb-4 !p-0 overflow-hidden">
        <div className="flex items-center justify-between px-4 py-2 bg-slate-100/80 dark:bg-white/5 border-b border-slate-200/70 dark:border-white/10">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-rose-400" />
            <div className="w-3 h-3 rounded-full bg-amber-400" />
            <div className="w-3 h-3 rounded-full bg-emerald-400" />
            <span className="ml-2 text-xs text-slate-500 dark:text-slate-400 font-mono">main.go</span>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={handleReset}
              className="text-xs text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 px-2 py-1 rounded hover:bg-slate-200/50 dark:hover:bg-white/10"
            >
              ↺ 重置
            </button>
          </div>
        </div>
        <textarea
          ref={textareaRef}
          value={code}
          onChange={(e) => setCode(e.target.value)}
          onKeyDown={handleKeyDown}
          className="w-full min-h-[350px] p-4 font-mono text-[13px] leading-relaxed resize-y bg-[#0f172a] dark:bg-[#050a1a] text-slate-200 border-0 outline-none"
          spellCheck={false}
          autoCapitalize="off"
          autoCorrect="off"
          placeholder="// 在此编写你的 Go 代码..."
        />
      </div>

      {/* 操作按钮 */}
      <div className="flex items-center gap-3 mb-4">
        <button
          onClick={handleRun}
          disabled={running}
          className="btn btn-primary px-8"
        >
          {running ? (
            <>
              <span className="inline-block w-4 h-4 rounded-full border-2 border-white border-t-transparent animate-spin" />
              运行中...
            </>
          ) : (
            '▶ 运行代码'
          )}
        </button>

        {result && !result.pass && (
          <button
            onClick={() => setShowSolution(!showSolution)}
            className="btn btn-ghost text-sm"
          >
            {showSolution ? '隐藏答案' : '📖 查看参考答案'}
          </button>
        )}
      </div>

      {/* 运行结果 */}
      {result && (
        <div className={`card mb-4 animate-slide-down border-l-4 ${
          result.pass
            ? 'border-l-emerald-500 bg-emerald-50/40 dark:bg-emerald-500/5'
            : 'border-l-rose-500 bg-rose-50/40 dark:bg-rose-500/5'
        }`}>
          <div className="flex items-center gap-2 mb-3">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center text-white font-bold ${
              result.pass ? 'bg-emerald-500 shadow-glow-green' : 'bg-rose-500 shadow-glow-red'
            }`}>
              {result.pass ? '✓' : '✗'}
            </div>
            <span className={`font-bold ${
              result.pass ? 'text-emerald-700 dark:text-emerald-300' : 'text-rose-700 dark:text-rose-300'
            }`}>
              {result.pass ? '所有测试用例通过！' : '未通过'}
            </span>
          </div>

          {result.output && (
            <div className="mb-3">
              <div className="text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1">输出：</div>
              <pre className="p-3 rounded-lg bg-slate-900 text-slate-200 text-xs font-mono leading-relaxed overflow-x-auto max-h-[200px] overflow-y-auto">
                {result.output}
              </pre>
            </div>
          )}

          {result.error && !result.pass && (
            <div>
              <div className="text-xs font-semibold text-rose-600 dark:text-rose-400 mb-1">错误：</div>
              <pre className="p-3 rounded-lg bg-rose-950/30 text-rose-300 text-xs font-mono leading-relaxed overflow-x-auto max-h-[200px] overflow-y-auto">
                {result.error}
              </pre>
            </div>
          )}
        </div>
      )}

      {/* 参考答案 */}
      {showSolution && result && result.solution && (
        <div className="card mb-4 animate-slide-down">
          <div className="text-sm font-bold text-slate-700 dark:text-slate-200 mb-2">
            📖 参考答案
          </div>
          <pre
            className="code-block text-sm"
            dangerouslySetInnerHTML={{ __html: highlightGo(result.solution) }}
          />
        </div>
      )}

      {/* 快捷键说明 */}
      <div className="text-xs text-slate-500 dark:text-slate-500 text-center">
        Tab 缩进 · Ctrl+Enter 运行（即将支持）
      </div>
    </section>
  )
}
