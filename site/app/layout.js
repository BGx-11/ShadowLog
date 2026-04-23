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

const SITE_URL = 'https://shadowlog.iambgx.in';
const SITE_TITLE = 'ShadowLog — Discrete Activity Analytics Framework';
const SITE_DESCRIPTION =
  'ShadowLog is a high-performance, native systems monitoring framework built in Go for authorized cybersecurity research. Zero dependencies, AES-256 encryption, and enterprise-grade capture in a single binary under 1MB.';

export const viewport = {
  themeColor: [
    { media: '(prefers-color-scheme: dark)', color: '#09090b' },
    { media: '(prefers-color-scheme: light)', color: '#09090b' },
  ],
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
  viewportFit: 'cover',
  colorScheme: 'dark',
};

export const metadata = {
  metadataBase: new URL(SITE_URL),
  title: {
    default: SITE_TITLE,
    template: '%s — ShadowLog',
  },
  description: SITE_DESCRIPTION,
  icons: {
    icon: [
      { url: '/logo.png', type: 'image/png', sizes: '512x512' },
    ],
    apple: [
      { url: '/logo.png', type: 'image/png', sizes: '512x512' },
    ],
  },
  keywords: [
    'ShadowLog',
    'Shadow Log',
    'shadow log tool',
    'shadowlog download',
    'shadowlog github',
    'keylogger Go',
    'activity monitor',
    'cybersecurity research',
    'security tool',
    'Go monitoring framework',
    'system analytics',
    'stealth monitoring',
    'Windows keylogger',
    'AES-256 encryption tool',
    'BGx',
    'iambgx',
    'iambgx.in',
    'shadowlog.iambgx.in',
    'Devansh Agarwal',
  ],
  authors: [{ name: 'Devansh Agarwal (BGx)', url: 'https://iambgx.in' }],
  creator: 'Devansh Agarwal (BGx)',
  publisher: 'BGx',
  category: 'Technology',
  classification: 'Cybersecurity',
  robots: {
    index: true,
    follow: true,
    'max-image-preview': 'large',
    'max-snippet': -1,
    'max-video-preview': -1,
    googleBot: {
      index: true,
      follow: true,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  openGraph: {
    title: SITE_TITLE,
    description: SITE_DESCRIPTION,
    url: SITE_URL,
    type: 'website',
    locale: 'en_US',
    siteName: 'ShadowLog',
    images: [
      {
        url: '/logo.png',
        width: 512,
        height: 512,
        alt: 'ShadowLog — Activity Analytics Framework',
        type: 'image/png',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: SITE_TITLE,
    description: SITE_DESCRIPTION,
    images: ['/logo.png'],
    creator: '@iambgx',
  },
  alternates: {
    canonical: SITE_URL,
  },
  other: {
    'msapplication-TileColor': '#09090b',
    'apple-mobile-web-app-capable': 'yes',
    'apple-mobile-web-app-status-bar-style': 'black-translucent',
    'apple-mobile-web-app-title': 'ShadowLog',
    'format-detection': 'telephone=no',
  },
};

// JSON-LD Structured Data for search engines
const jsonLd = {
  '@context': 'https://schema.org',
  '@type': 'SoftwareApplication',
  name: 'ShadowLog',
  alternateName: 'Shadow Log',
  description: SITE_DESCRIPTION,
  url: SITE_URL,
  applicationCategory: 'SecurityApplication',
  operatingSystem: 'Windows 10, Windows 11',
  offers: {
    '@type': 'Offer',
    price: '0',
    priceCurrency: 'USD',
  },
  author: {
    '@type': 'Person',
    name: 'Devansh Agarwal',
    alternateName: 'BGx',
    url: 'https://iambgx.in',
  },
  programmingLanguage: 'Go',
  downloadUrl: 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Release.zip',
  softwareVersion: 'latest',
  fileSize: '<1MB',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <head>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body>{children}</body>
    </html>
  );
}
