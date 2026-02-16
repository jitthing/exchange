import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './lib/**/*.{js,ts,jsx,tsx,mdx}'
  ],
  theme: {
    extend: {
      colors: {
        ink: '#0b1324',
        mist: '#f4f7fb',
        accent: '#0ea5a4',
        warn: '#d97706',
        danger: '#dc2626'
      }
    }
  },
  plugins: []
};

export default config;
