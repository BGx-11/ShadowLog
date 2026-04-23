import Link from 'next/link';
import Image from 'next/image';

export default function Footer() {
  return (
    <footer className="footer" role="contentinfo">
      <div className="footerInner">
        <div className="footerTop">
          <div className="footerBrand">
            <Link href="/" className="navLogo" style={{ textDecoration: 'none', color: 'var(--text-primary)' }}>
              <div className="navLogoIcon" aria-hidden="true">
                <Image src="/logo.png" alt="ShadowLog Logo" width={24} height={24} />
              </div>
              ShadowLog
            </Link>
            <p>A native systems monitoring framework for authorized cybersecurity research. Built with precision in Go.</p>
          </div>

          <div className="footerLinksGroup">
            <div className="footerCol">
              <h4>Project</h4>
              <ul role="list">
                <li><a href="/#features">Features</a></li>
                <li><a href="/#setup">Setup Guide</a></li>
                <li><a href="/#download">Download</a></li>
                <li><a href="https://github.com/BGx-11/ShadowLog" target="_blank" rel="noopener noreferrer">Source Code</a></li>
              </ul>
            </div>
            <div className="footerCol">
              <h4>Legal</h4>
              <ul role="list">
                <li><Link href="/terms">Terms of Service</Link></li>
                <li><Link href="/privacy">Privacy Policy</Link></li>
              </ul>
            </div>
            <div className="footerCol">
              <h4>Author</h4>
              <ul role="list">
                <li><a href="https://iambgx.in" target="_blank" rel="noopener noreferrer">Portfolio</a></li>
                <li><a href="https://github.com/BGx-11" target="_blank" rel="noopener noreferrer">GitHub</a></li>
              </ul>
            </div>
          </div>
        </div>

        <div className="footerBottom">
          <p className="footerCopyright">
            &copy; {new Date().getFullYear()}{' '}
            <a href="https://iambgx.in" target="_blank" rel="noopener noreferrer">Devansh Agarwal (BGx)</a>. All rights reserved.
          </p>
          <div className="footerLegal">
            <Link href="/terms">Terms</Link>
            <Link href="/privacy">Privacy</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
