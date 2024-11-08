/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: [
    './web/template/**/*.templ',
  ],
  darkMode: 'selector',
  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      black: '#000',
      white: '#fff',
      red: 'rgb(var(--color-red) / <alpha-value>)',
      green: 'rgb(var(--color-green) / <alpha-value>)',
      gray: {
        DEFAULT: 'rgb(var(--color-gray) / <alpha-value>)',
        light: 'rgb(var(--color-gray-light) / <alpha-value>)',
        dark: 'rgb(var(--color-gray-dark) / <alpha-value>)',
      },
      background: 'rgb(var(--color-background) / <alpha-value>)',
      foreground: 'rgb(var(--color-foreground) / <alpha-value>)',
      error: 'rgb(var(--color-error) / <alpha-value>)',
      success: 'rgb(var(--color-success) / <alpha-value>)',
    },
    fontFamily: {
      sans: ['Figtree', 'sans-serif'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
};
