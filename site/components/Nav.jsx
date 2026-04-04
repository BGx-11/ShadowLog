'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';

export default function Nav({ alwaysOpaque = false }) {
  const [scrolled, setScrolled] = useState(alwaysOpaque);
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    if (alwaysOpaque) return;
    function onScroll() {
      setScrolled(window.scrollY > 20);
    }
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, [alwaysOpaque]);

  return (
    <nav className={`nav${scrolled ? ' scrolled' : ''}`} role="navigation" aria-label="Main navigation">
      <div className="navInner">
        <Link href="/" className="navLogo" aria-label="ShadowLog Home" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <div className="navLogoIcon" aria-hidden="true" style={{ display: 'flex', width: '24px', height: '24px', position: 'relative' }}>
            <Image src="/logo.png" alt="ShadowLog Logo" fill style={{ objectFit: 'contain' }} />
          </div>
          ShadowLog
        </Link>

        <ul className={`navLinks${menuOpen ? ' navLinksOpen' : ''}`} role="list">
          <li><a href="/#features" onClick={() => setMenuOpen(false)}>Features</a></li>
          <li><a href="/#setup" onClick={() => setMenuOpen(false)}>Setup</a></li>
          <li><a href="/#tools" onClick={() => setMenuOpen(false)}>Tools</a></li>
          <li><a href="/#download" onClick={() => setMenuOpen(false)}>Download</a></li>
        </ul>

        <div className="navActions">
          <a href="https://github.com/BGx-11/ShadowLog" target="_blank" rel="noopener noreferrer" className="btn btnGhost btnSm">
            <svg className="btnIcon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true"><path d="M18 13v6a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
            Source Code
          </a>
          <button
            className="navToggle"
            onClick={() => setMenuOpen(!menuOpen)}
            aria-label="Toggle navigation menu"
            aria-expanded={menuOpen}
          >
            {menuOpen ? (
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            ) : (
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            )}
          </button>
        </div>
      </div>
    </nav>
  );
}
