/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class', // 👈 activamos modo oscuro por clase .dark
  content: ['./index.html', './src/**/*.{svelte,js,ts}'],
  theme: {
    extend: {
      colors: {
        // Paleta montaña
        rock: {
          900: '#020617', // casi negro azulado
          800: '#0b1120',
          700: '#111827'
        },
        forest: {
          500: '#059669',
          600: '#047857',
          700: '#065f46',
          800: '#064e3b'
        },
        mist: {
          100: '#e5f3ef',
          200: '#d1fae5',
          900: '#020617'
        },
        accent: {
          yellow: '#facc15',
          orange: '#fb923c'
        }
      },
      boxShadow: {
        'elevated': '0 18px 40px rgba(0,0,0,0.45)'
      },
      backgroundImage: {
        'mountain-gradient':
          'radial-gradient(circle at top, rgba(56,189,248,0.25), transparent 55%), radial-gradient(circle at bottom, rgba(34,197,94,0.25), transparent 60%), linear-gradient(to bottom, #020617, #0b1120)'
      }
    }
  },
  plugins: []
};
