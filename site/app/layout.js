import { Inter, JetBrains_Mono } from 'next/font/google';
import './globals.css';

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  variable: '--font-jetbrains',
  display: 'swap',
});

export const metadata = {
  title: 'ShadowLog (Shadow Log) — Discrete Activity Analytics Framework',
  description:
    'ShadowLog (or Shadow Log) is a high-performance, native systems monitoring framework built in Go for authorized cybersecurity research and security demonstrations. Zero dependencies, enterprise-grade capture. The ultimate shadow log tool.',
  keywords: [
    'ShadowLog',
    'Shadow Log',
    'shadow log tool',
    'shadowlog download',
    'activity monitor',
    'cybersecurity research',
    'security tool',
    'Go monitoring framework',
    'system analytics',
    'BGx',
    'iambgx',
    'iambgx.in',
    'shadowlog.iambgx.in',
  ],
  authors: [{ name: 'Devansh Agarwal (BGx)', url: 'https://iambgx.in' }],
  openGraph: {
    title: 'ShadowLog (Shadow Log) — Discrete Activity Analytics Framework',
    description:
      'High-performance native systems monitoring framework for authorized cybersecurity research. Built in Go with zero runtime dependencies. Get the best shadow log experience.',
    type: 'website',
    locale: 'en_US',
    siteName: 'ShadowLog',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'ShadowLog (Shadow Log) — Discrete Activity Analytics',
    description:
      'Native systems monitoring framework for authorized cybersecurity research. Built in Go.',
  },
  alternates: {
    canonical: 'https://shadowlog.iambgx.in', // Adjust base URL as needed or assume self
  },
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <body>{children}</body>
    </html>
  );
}
