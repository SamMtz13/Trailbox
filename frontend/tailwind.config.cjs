/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{svelte,ts,js}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Space Grotesk"', 'Inter', 'system-ui', 'sans-serif'],
      },
      colors: {
        forest: '#1f3d2b',
        moss: '#4a6b47',
        mist: '#e9f1ed',
        clay: '#c4622d'
      }
    }
  },
  plugins: [require('@tailwindcss/forms')],
};
