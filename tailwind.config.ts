import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './lib/**/*.{js,ts,jsx,tsx,mdx}'
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', '-apple-system', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial', 'sans-serif'],
      },
      colors: {
        primary: {
          50:  '#E5F1FF',
          100: '#B3D7FF',
          200: '#66B0FF',
          300: '#3D9AFF',
          400: '#1A88FF',
          500: '#0A84FF',
          600: '#0870DB',
          700: '#065CB7',
          800: '#044893',
          900: '#02346F',
          DEFAULT: '#0A84FF',
        },
        accent: {
          50:  '#FFF0EC',
          100: '#FFD6CC',
          200: '#FFB09E',
          300: '#FF8B70',
          400: '#FF6B4A',
          500: '#FF6B4A',
          600: '#E55535',
          700: '#CC4020',
          DEFAULT: '#FF6B4A',
        },
        neutral: {
          0:   '#FFFFFF',
          50:  '#FAFAF8',
          100: '#F5F5F3',
          200: '#E8E5E0',
          300: '#D5D3CC',
          400: '#AEAB9F',
          500: '#9B9B9B',
          600: '#6B685C',
          700: '#55534A',
          800: '#37352F',
          900: '#1A1A1A',
          950: '#0D0D0D',
        },
        success: {
          50:  'rgba(52,199,89,0.12)',
          500: '#34C759',
          600: '#248A3D',
          DEFAULT: '#34C759',
        },
        warning: {
          50:  'rgba(255,179,64,0.15)',
          500: '#FFB340',
          600: '#A05A00',
          DEFAULT: '#FFB340',
        },
        danger: {
          50:  'rgba(255,59,48,0.12)',
          500: '#FF3B30',
          600: '#D70015',
          DEFAULT: '#FF3B30',
        },
        info: {
          50:  'rgba(10,132,255,0.1)',
          500: '#007AFF',
          600: '#0A84FF',
          DEFAULT: '#007AFF',
        },
        /* Semantic aliases */
        background: '#FAFAF8',
        surface: '#F5F5F3',
        border: '#E8E5E0',
        muted: '#9B9B9B',
        body: '#37352F',
        heading: '#1A1A1A',
      },
      spacing: {
        '1':  '0.25rem',
        '2':  '0.5rem',
        '3':  '0.75rem',
        '4':  '1rem',
        '6':  '1.5rem',
        '8':  '2rem',
        '12': '3rem',
        '16': '4rem',
      },
      borderRadius: {
        sm:   '8px',
        md:   '12px',
        lg:   '16px',
        xl:   '20px',
        full: '999px',
      },
      boxShadow: {
        subtle:  '0 1px 3px rgba(0,0,0,0.06)',
        medium:  '0 4px 12px rgba(0,0,0,0.08)',
        raised:  '0 8px 24px rgba(0,0,0,0.12)',
        focus:   '0 0 0 4px rgba(10,132,255,0.12)',
      },
      fontSize: {
        'display': ['2rem', { lineHeight: '2.5rem', fontWeight: '600' }],
        'h1':      ['1.75rem', { lineHeight: '2.25rem', fontWeight: '600' }],
        'h2':      ['1.375rem', { lineHeight: '1.75rem', fontWeight: '600' }],
        'h3':      ['1.125rem', { lineHeight: '1.5rem', fontWeight: '500' }],
        'body':    ['1rem', { lineHeight: '1.5rem', fontWeight: '400' }],
        'small':   ['0.875rem', { lineHeight: '1.25rem', fontWeight: '400' }],
        'caption': ['0.75rem', { lineHeight: '1rem', fontWeight: '400' }],
      },
    }
  },
  plugins: []
};

export default config;
