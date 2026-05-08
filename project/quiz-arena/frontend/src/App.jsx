import { useState, useEffect, useCallback } from 'react'
import { api } from './api'
import { storage } from './storage'
import Header from './components/Header'
import StartScreen from './components/StartScreen'
import QuizScreen from './components/QuizScreen'
import ResultScreen from './components/ResultScreen'
import CodingScreen from './components/CodingScreen'

const SCREENS = {
  START: 'start',
  QUIZ: 'quiz',
  RESULT: 'result',
  CODING: 'coding',
}

export default function App() {
  const [screen, setScreen] = useState(SCREENS.START)
  const [questions, setQuestions] = useState([])
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // 主题
  const [theme, setTheme] = useState(storage.getTheme())

  // 答题状态
  const [current, setCurrent] = useState(0)
  const [score, setScore] = useState(0)
  const [records, setRecords] = useState([])
  const [startTime, setStartTime] = useState(0)

  // 应用主题到 html 元素
  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark')
    storage.setTheme(theme)
  }, [theme])

  const toggleTheme = () => setTheme((t) => (t === 'dark' ? 'light' : 'dark'))

  // 加载分类
  useEffect(() => {
    api
      .getCategories()
      .then((data) => {
        setCategories(data)
        setLoading(false)
      })
      .catch((err) => {
        setError(err.message)
        setLoading(false)
      })
  }, [])

  // 开始试炼：根据筛选条件加载题目
  const handleStart = useCallback(async (filters = {}) => {
    try {
      setLoading(true)
      const data = await api.getQuestions(filters)
      if (data.length === 0) {
        alert('当前筛选条件下没有题目')
        setLoading(false)
        return
      }
      setQuestions(data)
      setCurrent(0)
      setScore(0)
      setRecords([])
      setStartTime(Date.now())
      setScreen(SCREENS.QUIZ)
    } catch (e) {
      alert('加载题目失败：' + e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  // 进入编码模式
  const handleStartCoding = useCallback(() => {
    setScreen(SCREENS.CODING)
  }, [])

  // 答题（由 QuizScreen 上报）
  const handleAnswered = useCallback((record) => {
    setRecords((prev) => [...prev, record])
    if (record.correct) {
      setScore((s) => s + 10)
      storage.removeWrong(record.questionId)
    } else {
      storage.addWrong(record.questionId)
    }
  }, [])

  // 下一题 / 查看结果
  const handleNext = useCallback(() => {
    if (current < questions.length - 1) {
      setCurrent((c) => c + 1)
    } else {
      const correctCount = records.filter((r) => r.correct).length
      storage.updateHistory(records.length, correctCount)
      setScreen(SCREENS.RESULT)
    }
  }, [current, questions.length, records])

  // 回到开始页
  const handleBackHome = useCallback(() => {
    setScreen(SCREENS.START)
  }, [])

  if (loading && screen === SCREENS.START && categories.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="flex items-center gap-3 text-slate-700 dark:text-slate-200 text-lg">
          <span className="inline-block w-5 h-5 rounded-full border-2 border-brand-500 border-t-transparent animate-spin" />
          正在加载题库...
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6">
        <div className="card max-w-md text-center">
          <div className="text-5xl mb-3">❌</div>
          <div className="font-bold text-lg mb-2 text-slate-800 dark:text-slate-100">
            加载失败
          </div>
          <div className="text-sm text-slate-500 dark:text-slate-400">{error}</div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-6 sm:py-10">
      <Header theme={theme} onToggleTheme={toggleTheme} />

      {screen === SCREENS.START && (
        <StartScreen
          categories={categories}
          onStart={handleStart}
          onStartCoding={handleStartCoding}
        />
      )}

      {screen === SCREENS.QUIZ && (
        <QuizScreen
          question={questions[current]}
          current={current}
          total={questions.length}
          score={score}
          onAnswered={handleAnswered}
          onNext={handleNext}
          onExit={handleBackHome}
        />
      )}

      {screen === SCREENS.RESULT && (
        <ResultScreen
          records={records}
          total={questions.length}
          elapsedSec={Math.round((Date.now() - startTime) / 1000)}
          onRestart={() => handleStart()}
          onHome={handleBackHome}
        />
      )}

      {screen === SCREENS.CODING && (
        <CodingScreen onExit={handleBackHome} />
      )}
    </div>
  )
}
