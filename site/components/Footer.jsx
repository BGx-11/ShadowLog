import Link from 'next/link';
import Image from 'next/image';

export default function Footer() {
  return (
    <footer className="footer" role="contentinfo">
      <div className="container">
        <div className="footerGrid">
          <div className="footerBrand">
            <Link href="/" className="navLogo" style={{ marginBottom: '16px', display: 'flex' }}>
              <div className="navLogoIcon">
                <Image src="/logo.png" alt="ShadowLog Logo" width={20} height={20} />
              </div>
              <span>ShadowLog</span>
            </Link>
            <p>A native, high-performance systems monitoring framework engineered in Go for authorized cybersecurity research, red teaming, and system analytics.</p>
          </div>

          <div className="footerCol">
            <h4>Project Links</h4>
            <ul role="list">
              <li><a href="/#features">Architecture</a></li>
              <li><a href="/#setup">Deployment Guide</a></li>
              <li><a href="https://github.com/BGx-11/ShadowLog" target="_blank" rel="noopener noreferrer">Source Code</a></li>
              <li><a href="https://go.dev/" target="_blank" rel="noopener noreferrer">Go Programming</a></li>
            </ul>
          </div>
          
          <div className="footerCol">
            <h4>Creator</h4>
            <ul role="list">
              <li><a href="https://iambgx.in" target="_blank" rel="noopener noreferrer">Devansh Agarwal</a></li>
              <li><a href="https://github.com/BGx-11" target="_blank" rel="noopener noreferrer">GitHub Profile</a></li>
            </ul>
          </div>
        </div>

        <div className="footerBottom">
          <p className="footerCopy">
            &copy; {new Date().getFullYear()} <a href="https://iambgx.in" target="_blank" rel="noopener noreferrer">BGx</a>. All rights reserved.
          </p>
          <div className="footerLegal">
            <a href="https://github.com/BGx-11/ShadowLog/blob/main/LICENSE" target="_blank" rel="noopener noreferrer">License</a>
            <Link href="/terms">Terms of Service</Link>
            <Link href="/privacy">Privacy Policy</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
