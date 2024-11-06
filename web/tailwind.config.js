/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: [
    './web/template/**/*.templ',
  ],
  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      black: '#000',
      white: '#fff',
      green: {
        DEFAULT: '#7e9142',
        light: '#a0cf78',
        dark: '#7e9142',
      },
      gray: {
        DEFAULT: '#333333',
        dark: '#1a1a1a',
        light: '#f0f0f0',
      },
    },
    fontFamily: {
      sans: ['Figtree', 'sans-serif'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
};
