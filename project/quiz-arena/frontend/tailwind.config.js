/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          '-apple-system',
          'BlinkMacSystemFont',
          '"Segoe UI"',
          '"PingFang SC"',
          '"Microsoft YaHei"',
          'sans-serif',
        ],
        mono: [
          '"JetBrains Mono"',
          '"Fira Code"',
          'SFMono-Regular',
          'Menlo',
          'Consolas',
          'monospace',
        ],
      },
      colors: {
        // 品牌色（靛紫）
        brand: {
          50: '#eef2ff',
          100: '#e0e7ff',
          200: '#c7d2fe',
          300: '#a5b4fc',
          400: '#818cf8',
          500: '#6366f1',
          600: '#4f46e5',
          700: '#4338ca',
          800: '#3730a3',
          900: '#312e81',
        },
        // 强调色（青蓝）
        accent: {
          400: '#22d3ee',
          500: '#06b6d4',
          600: '#0891b2',
        },
        // 成功/错误/警告
        neon: {
          green: '#4ade80',
          red: '#fb7185',
          amber: '#fbbf24',
          pink: '#f472b6',
        },
        // 画布背景（暗色）
        canvas: {
          950: '#0b1020',
          900: '#0f172a',
          850: '#111834',
          800: '#1e293b',
        },
      },
      boxShadow: {
        glow: '0 10px 40px -10px rgba(99, 102, 241, 0.55)',
        'glow-lg': '0 20px 60px -15px rgba(99, 102, 241, 0.7)',
        'glow-pink': '0 10px 40px -10px rgba(244, 114, 182, 0.55)',
        'glow-green': '0 10px 40px -10px rgba(74, 222, 128, 0.55)',
        'glow-red': '0 10px 40px -10px rgba(251, 113, 133, 0.55)',
        card: '0 10px 30px -10px rgba(15, 23, 42, 0.5)',
      },
      backgroundImage: {
        'brand-gradient':
          'linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #ec4899 100%)',
        'accent-gradient':
          'linear-gradient(135deg, #06b6d4 0%, #6366f1 100%)',
        'progress-gradient':
          'linear-gradient(90deg, #22d3ee 0%, #6366f1 50%, #ec4899 100%)',
      },
      animation: {
        'fade-in': 'fadeIn 0.45s cubic-bezier(0.16, 1, 0.3, 1)',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-up': 'slideUp 0.35s cubic-bezier(0.16, 1, 0.3, 1)',
        'bounce-slow': 'bounceSlow 1.2s ease-out',
        float: 'float 6s ease-in-out infinite',
        shimmer: 'shimmer 2.4s linear infinite',
        'pulse-glow': 'pulseGlow 2.2s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0', transform: 'translateY(16px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(12px) scale(0.98)' },
          '100%': { opacity: '1', transform: 'translateY(0) scale(1)' },
        },
        bounceSlow: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-14px)' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0) translateX(0)' },
          '50%': { transform: 'translateY(-12px) translateX(6px)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        pulseGlow: {
          '0%, 100%': {
            boxShadow: '0 0 0 0 rgba(99, 102, 241, 0.35)',
          },
          '50%': {
            boxShadow: '0 0 0 14px rgba(99, 102, 241, 0)',
          },
        },
      },
    },
  },
  plugins: [],
}
