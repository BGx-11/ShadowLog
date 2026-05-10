import { Inter, JetBrains_Mono, Outfit } from 'next/font/google';
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

const outfit = Outfit({
  subsets: ['latin'],
  variable: '--font-outfit',
  display: 'swap',
});

const SITE_URL = 'https://shadowlog.iambgx.in';
const SITE_TITLE = 'ShadowLog | Advanced Systems Monitoring & Analytics Framework';
const SITE_DESCRIPTION = 'ShadowLog is a state-of-the-art, high-performance native systems monitoring framework engineered in Go. Featuring AES-256-GCM encryption, zero dependencies, and quad-channel exfiltration for authorized cybersecurity research and threat simulation.';

export const viewport = {
  themeColor: [
    { media: '(prefers-color-scheme: dark)', color: '#030303' },
    { media: '(prefers-color-scheme: light)', color: '#030303' },
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
    template: '%s | ShadowLog',
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
    'Shadow Log framework',
    'systems monitoring tool',
    'Windows keylogger',
    'Go security research tool',
    'activity analytics framework',
    'stealth monitoring',
    'cybersecurity tools',
    'AES-256 encryption monitoring',
    'telemetry exfiltration',
    'advanced threat simulation',
    'BGx',
    'Devansh Agarwal',
    'iambgx'
  ],
  authors: [{ name: 'Devansh Agarwal (BGx)', url: 'https://iambgx.in' }],
  creator: 'Devansh Agarwal (BGx)',
  publisher: 'BGx Cybersecurity',
  category: 'Cybersecurity Software',
  classification: 'Security Analytics',
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
        width: 1200,
        height: 630,
        alt: 'ShadowLog - Advanced Systems Monitoring',
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
    'msapplication-TileColor': '#030303',
    'apple-mobile-web-app-capable': 'yes',
    'apple-mobile-web-app-status-bar-style': 'black-translucent',
    'apple-mobile-web-app-title': 'ShadowLog',
    'format-detection': 'telephone=no',
  },
};

const jsonLd = {
  '@context': 'https://schema.org',
  '@type': 'SoftwareApplication',
  name: 'ShadowLog',
  alternateName: 'Shadow Log Activity Framework',
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
    sameAs: [
      'https://github.com/BGx-11',
      'https://www.linkedin.com/in/devansha25/'
    ]
  },
  programmingLanguage: 'Go',
  downloadUrl: 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Release.zip',
  softwareVersion: 'v2.3',
  fileSize: '1MB',
  featureList: [
    'Zero dependencies',
    'AES-256-GCM encryption',
    'Quad-channel exfiltration',
    'Volatile Screen Capture',
    'Advanced stealth OPSEC'
  ]
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable} ${outfit.variable}`}>
      <head>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body>
        <div className="noise-bg"></div>
        {children}
      </body>
    </html>
  );
}
