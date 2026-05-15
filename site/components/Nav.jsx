'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import Image from 'next/image';

export default function Nav({ alwaysOpaque = false }) {
  const [scrolled, setScrolled] = useState(alwaysOpaque);
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    if (alwaysOpaque) return;
    const onScroll = () => setScrolled(window.scrollY > 20);
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, [alwaysOpaque]);

  useEffect(() => {
    document.body.style.overflow = menuOpen ? 'hidden' : '';
    return () => { document.body.style.overflow = ''; };
  }, [menuOpen]);

  const closeMenu = useCallback(() => setMenuOpen(false), []);

  return (
    <nav className={`nav ${scrolled ? 'scrolled' : ''} ${menuOpen ? 'menuActive' : ''}`} role="navigation">
      <div className="navInner">
        <Link href="/" className="navLogo" onClick={closeMenu}>
          <div className="navLogoIcon">
            <Image src="/logo.png" alt="ShadowLog" width={20} height={20} priority />
          </div>
          <span>ShadowLog</span>
        </Link>

        <ul className={`navLinks ${menuOpen ? 'menuActive' : ''}`}>
          <li><a href="/#features" onClick={closeMenu}>Features</a></li>
          <li><a href="/#changelog" onClick={closeMenu}>Changelog</a></li>
          <li><a href="/#setup" onClick={closeMenu}>Documentation</a></li>
          <li><a href="/#faq" onClick={closeMenu}>FAQ</a></li>
          <li><a href="/#download" onClick={closeMenu}>Download</a></li>
        </ul>

        <div className="navActions">
          <a href="https://github.com/BGx-11/ShadowLog" target="_blank" rel="noopener noreferrer" className="btn btn-ghost btn-sm">
            GitHub
          </a>
          <button className={`navToggle ${menuOpen ? 'navToggleActive' : ''}`} onClick={() => setMenuOpen(!menuOpen)}>
            <span className="hamburger">
              <span className="hamburgerLine"></span>
              <span className="hamburgerLine"></span>
              <span className="hamburgerLine"></span>
            </span>
          </button>
        </div>
      </div>
    </nav>
  );
}
