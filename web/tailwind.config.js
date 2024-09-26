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
        light: '#b0cf78',
        dark: '#7e9142',
      },
      gray: {
        DEFAULT: '#333333',
        light: '#f5f5f5',
        lighter: '#f9f9f9',
        dark: '#1a1a1a',
        darker: '#111111',
      },
    },
    fontFamily: {
      sans: ['Poppins', 'sans-serif'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
};
