// 代码高亮辅助（Prism.js）
import Prism from 'prismjs'
import 'prismjs/components/prism-go'

// 将代码字符串转换为带高亮的 HTML
export function highlightGo(code) {
  if (!code) return ''
  const normalized = code.replace(/\\n/g, '\n')
  return Prism.highlight(normalized, Prism.languages.go, 'go')
}
