/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        areia: '#f7f1e8',
        coqueiro: '#155e63',
        mar: '#0f8b8d',
        sol: '#f4a261',
        coral: '#e76f51',
        noite: '#1f2937',
      },
      boxShadow: {
        mapa: '0 24px 64px rgba(21, 94, 99, 0.18)',
      },
      fontFamily: {
        display: ['Sora', 'sans-serif'],
        body: ['Manrope', 'sans-serif'],
      },
    },
  },
  plugins: [],
};