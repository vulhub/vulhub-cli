/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      keyframes: {
        "accordion-down": {
          from: { height: "0" },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: "0" },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
      typography: {
        DEFAULT: {
          css: {
            '--tw-prose-body': 'hsl(215 20% 75%)',
            '--tw-prose-headings': 'hsl(210 40% 96%)',
            '--tw-prose-lead': 'hsl(215 20% 65%)',
            '--tw-prose-links': 'hsl(217 91% 60%)',
            '--tw-prose-bold': 'hsl(210 40% 96%)',
            '--tw-prose-counters': 'hsl(215 20% 55%)',
            '--tw-prose-bullets': 'hsl(217 91% 60%)',
            '--tw-prose-hr': 'hsl(217 33% 25%)',
            '--tw-prose-quotes': 'hsl(215 20% 75%)',
            '--tw-prose-quote-borders': 'hsl(217 91% 60%)',
            '--tw-prose-captions': 'hsl(215 20% 55%)',
            '--tw-prose-code': 'hsl(199 89% 60%)',
            '--tw-prose-pre-code': 'hsl(215 20% 85%)',
            '--tw-prose-pre-bg': 'hsl(222 47% 9%)',
            '--tw-prose-th-borders': 'hsl(217 33% 30%)',
            '--tw-prose-td-borders': 'hsl(217 33% 20%)',
            maxWidth: 'none',
            color: 'var(--tw-prose-body)',
            lineHeight: '1.8',
            // Headings
            h1: {
              color: 'var(--tw-prose-headings)',
              fontWeight: '700',
              fontSize: '1.875rem',
              marginTop: '2rem',
              marginBottom: '1rem',
              lineHeight: '1.3',
              borderBottom: '1px solid hsl(217 33% 25%)',
              paddingBottom: '0.5rem',
            },
            h2: {
              color: 'var(--tw-prose-headings)',
              fontWeight: '600',
              fontSize: '1.5rem',
              marginTop: '2rem',
              marginBottom: '0.75rem',
              lineHeight: '1.4',
              borderBottom: '1px solid hsl(217 33% 20%)',
              paddingBottom: '0.375rem',
            },
            h3: {
              color: 'var(--tw-prose-headings)',
              fontWeight: '600',
              fontSize: '1.25rem',
              marginTop: '1.5rem',
              marginBottom: '0.5rem',
              lineHeight: '1.5',
            },
            h4: {
              color: 'var(--tw-prose-headings)',
              fontWeight: '600',
              fontSize: '1.125rem',
              marginTop: '1.25rem',
              marginBottom: '0.5rem',
            },
            // Paragraphs
            p: {
              marginTop: '1rem',
              marginBottom: '1rem',
            },
            // Links
            a: {
              color: 'var(--tw-prose-links)',
              textDecoration: 'none',
              fontWeight: '500',
              borderBottom: '1px solid transparent',
              transition: 'border-color 0.2s',
              '&:hover': {
                borderBottomColor: 'var(--tw-prose-links)',
              },
            },
            // Strong/Bold
            strong: {
              color: 'var(--tw-prose-bold)',
              fontWeight: '600',
            },
            // Code (inline)
            code: {
              color: 'var(--tw-prose-code)',
              backgroundColor: 'hsl(222 47% 13%)',
              padding: '0.2em 0.4em',
              borderRadius: '0.25rem',
              fontSize: '0.875em',
              fontFamily: "'Fira Code', monospace",
              fontWeight: '400',
            },
            'code::before': {
              content: '""',
            },
            'code::after': {
              content: '""',
            },
            // Pre/Code blocks
            pre: {
              backgroundColor: 'var(--tw-prose-pre-bg)',
              color: 'var(--tw-prose-pre-code)',
              borderRadius: '0.5rem',
              padding: '1rem',
              overflowX: 'auto',
              border: '1px solid hsl(217 33% 25%)',
              marginTop: '1.5rem',
              marginBottom: '1.5rem',
            },
            'pre code': {
              backgroundColor: 'transparent',
              padding: '0',
              fontSize: '0.875rem',
              color: 'inherit',
            },
            // Lists
            ul: {
              marginTop: '1rem',
              marginBottom: '1rem',
              paddingLeft: '1.5rem',
            },
            ol: {
              marginTop: '1rem',
              marginBottom: '1rem',
              paddingLeft: '1.5rem',
            },
            li: {
              marginTop: '0.375rem',
              marginBottom: '0.375rem',
            },
            'ul > li::marker': {
              color: 'var(--tw-prose-bullets)',
            },
            'ol > li::marker': {
              color: 'var(--tw-prose-counters)',
            },
            // Blockquotes
            blockquote: {
              borderLeftWidth: '4px',
              borderLeftColor: 'var(--tw-prose-quote-borders)',
              paddingLeft: '1rem',
              fontStyle: 'italic',
              color: 'var(--tw-prose-quotes)',
              marginTop: '1.5rem',
              marginBottom: '1.5rem',
              backgroundColor: 'hsl(222 47% 9%)',
              padding: '1rem',
              borderRadius: '0 0.5rem 0.5rem 0',
            },
            'blockquote p:first-of-type::before': {
              content: '""',
            },
            'blockquote p:last-of-type::after': {
              content: '""',
            },
            // Horizontal rule
            hr: {
              borderColor: 'var(--tw-prose-hr)',
              marginTop: '2rem',
              marginBottom: '2rem',
            },
            // Tables
            table: {
              width: '100%',
              marginTop: '1.5rem',
              marginBottom: '1.5rem',
              borderCollapse: 'collapse',
              fontSize: '0.875rem',
            },
            thead: {
              borderBottomWidth: '2px',
              borderBottomColor: 'var(--tw-prose-th-borders)',
            },
            'thead th': {
              color: 'var(--tw-prose-headings)',
              fontWeight: '600',
              padding: '0.75rem',
              textAlign: 'left',
              backgroundColor: 'hsl(222 47% 11%)',
            },
            'tbody tr': {
              borderBottomWidth: '1px',
              borderBottomColor: 'var(--tw-prose-td-borders)',
            },
            'tbody td': {
              padding: '0.75rem',
            },
            // Images
            img: {
              borderRadius: '0.5rem',
              marginTop: '1.5rem',
              marginBottom: '1.5rem',
            },
          },
        },
      },
    },
  },
  plugins: [
    require("tailwindcss-animate"),
    require("@tailwindcss/typography"),
  ],
}
