/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        weenie: {
          red: '#ff6b6b',
          gold: '#feca57',
          dark: '#1a1a2e',
          darker: '#0f0f1a',
          light: '#f8f9fa'
        }
      },
      backgroundImage: {
        'weenie-gradient': 'linear-gradient(135deg, #ff6b6b 0%, #feca57 100%)',
        'weenie-gradient-dark': 'linear-gradient(135deg, #1a1a2e 0%, #0f0f1a 100%)',
      },
      fontFamily: {
        minecraft: ['Minecraft', 'monospace'],
        sans: ['Inter', 'system-ui', 'sans-serif']
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'float': 'float 6s ease-in-out infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        }
      }
    },
  },
  plugins: [],
}
