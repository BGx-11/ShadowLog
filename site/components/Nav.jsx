'use client';

import { useState, useEffect, useCallback } from 'react';
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

  // Lock body scroll when mobile menu is open
  useEffect(() => {
    if (menuOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => { document.body.style.overflow = ''; };
  }, [menuOpen]);

  // Close menu on resize past mobile breakpoint
  useEffect(() => {
    function onResize() {
      if (window.innerWidth > 768 && menuOpen) {
        setMenuOpen(false);
      }
    }
    window.addEventListener('resize', onResize, { passive: true });
    return () => window.removeEventListener('resize', onResize);
  }, [menuOpen]);

  // Close on Escape key
  useEffect(() => {
    function onKeyDown(e) {
      if (e.key === 'Escape' && menuOpen) setMenuOpen(false);
    }
    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [menuOpen]);

  const closeMenu = useCallback(() => setMenuOpen(false), []);

  return (
    <>
      <nav
        className={`nav${scrolled ? ' scrolled' : ''}${menuOpen ? ' menuActive' : ''}`}
        role="navigation"
        aria-label="Main navigation"
      >
        <div className="navInner">
          <Link href="/" className="navLogo" aria-label="ShadowLog Home">
            <div className="navLogoIcon" aria-hidden="true">
              <Image src="/logo.png" alt="" width={28} height={28} priority />
            </div>
            <span>ShadowLog</span>
          </Link>

          <ul className={`navLinks${menuOpen ? ' navLinksOpen' : ''}`} role="list">
            <li><a href="/#features" onClick={closeMenu}>Features</a></li>
            <li><a href="/#setup" onClick={closeMenu}>Setup</a></li>
            <li><a href="/#tools" onClick={closeMenu}>Tools</a></li>
            <li><a href="/#changelog" onClick={closeMenu}>Changelog</a></li>
            <li><a href="/#download" onClick={closeMenu}>Download</a></li>
            {/* Mobile-only: show source code link inline */}
            <li className="navLinkMobileOnly">
              <a href="https://github.com/BGx-11/ShadowLog" target="_blank" rel="noopener noreferrer" onClick={closeMenu}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 00-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0020 4.77 5.07 5.07 0 0019.91 1S18.73.65 16 2.48a13.38 13.38 0 00-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 005 4.77a5.44 5.44 0 00-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 009 18.13V22"/></svg>
                Source Code
              </a>
            </li>
          </ul>

          <div className="navActions">
            <a
              href="https://github.com/BGx-11/ShadowLog"
              target="_blank"
              rel="noopener noreferrer"
              className="btn btnGhost btnSm navGithubBtn"
              aria-label="View source code on GitHub"
            >
              <svg className="btnIcon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 00-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0020 4.77 5.07 5.07 0 0019.91 1S18.73.65 16 2.48a13.38 13.38 0 00-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 005 4.77a5.44 5.44 0 00-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 009 18.13V22"/></svg>
              GitHub
            </a>
            <button
              className={`navToggle${menuOpen ? ' navToggleActive' : ''}`}
              onClick={() => setMenuOpen(!menuOpen)}
              aria-label={menuOpen ? 'Close navigation menu' : 'Open navigation menu'}
              aria-expanded={menuOpen}
              type="button"
            >
              <span className="hamburger">
                <span className="hamburgerLine"></span>
                <span className="hamburgerLine"></span>
                <span className="hamburgerLine"></span>
              </span>
            </button>
          </div>
        </div>
      </nav>

      {/* Mobile menu backdrop overlay */}
      {menuOpen && (
        <div className="navBackdrop" onClick={closeMenu} aria-hidden="true" />
      )}
    </>
  );
}
