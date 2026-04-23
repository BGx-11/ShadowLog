import Nav from '@/components/Nav';
import ScrollAnimations from '@/components/ScrollAnimations';
import SetupTabs from '@/components/SetupTabs';
import DownloadSection from '@/components/DownloadSection';
import Footer from '@/components/Footer';

function HeroBackground() {
  return (
    <div className="heroBackground">
      <div className="heroOrb orbPrimary"></div>
      <div className="heroOrb orbSecondary"></div>
    </div>
  );
}

export default function Home() {
  return (
    <>
      <ScrollAnimations />
      <Nav />

      {/* ===== HERO ===== */}
      <section className="hero" id="hero">
        <HeroBackground />
        <div className="heroContent">
          <p className="heroLabel animateIn stagger-1">Activity Analytics Framework</p>
          <h1 className="animateIn stagger-2"><span className="highlight">Shadow Log</span></h1>
          <p className="heroSubtitle animateIn stagger-3">
            A high-performance, native systems monitoring framework engineered for
            authorized security research and demonstrations. Built entirely in Go
            with zero runtime dependencies.
          </p>

          <div className="heroTerminal animateIn stagger-4">
            <div className="terminalHeader">
              <span className="dot dotRed"></span>
              <span className="dot dotYellow"></span>
              <span className="dot dotGreen"></span>
              <span className="terminalTitle">bash</span>
            </div>
            <pre className="terminalBody">
              <code>$ shadowlog --listen :8080 --stealth</code>
            </pre>
          </div>

          <div className="heroCtas animateIn stagger-5">
            <a href="#setup" className="btn btnPrimary btnLg">
              Get Started
            </a>
            <a href="#download" className="btn btnGhost btnLg">
              <svg className="btnIcon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Download
            </a>
          </div>

          <div className="heroStats animateIn stagger-6">
            <div className="statItem">
              <span className="statValue">&lt; 1MB</span>
              <span className="statLabel">Binary Size</span>
            </div>
            <div className="statItem">
              <span className="statValue">0</span>
              <span className="statLabel">Dependencies</span>
            </div>
            <div className="statItem">
              <span className="statValue">AES-256</span>
              <span className="statLabel">Encryption</span>
            </div>
            <div className="statItem">
              <span className="statValue">100%</span>
              <span className="statLabel">Go Native</span>
            </div>
          </div>
        </div>
      </section>

      {/* ===== WARNING BANNER ===== */}
      <section className="warningBanner" id="legal-warning" role="alert">
        <div className="warningInner">
          <div className="warningIcon" aria-hidden="true">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
          </div>
          <div className="warningContent">
            <h3>Legal &amp; Ethical Notice</h3>
            <p>
              This software is intended <strong>exclusively for authorized educational and security research purposes</strong>.
              Deploying monitoring tools on systems you do not own or lack explicit, documented authorization to monitor
              violates federal and state laws including the <strong>Computer Fraud and Abuse Act (CFAA)</strong>.
              The developer provides this tool &ldquo;as-is&rdquo; for security professionals and students and assumes{' '}
              <strong>no liability</strong> for misuse, data loss, or legal consequences.
            </p>
          </div>
        </div>
      </section>

      {/* ===== FEATURES ===== */}
      <section className="section" id="features">
        <div className="container">
          <div className="sectionHeader animateIn">
            <span className="sectionLabel">Technical Overview</span>
            <h2 className="sectionTitle">Engineered for Precision</h2>
            <p className="sectionDesc">A masterclass in operational security, built to deliver enterprise-grade analytics with a microscopic system footprint.</p>
          </div>

          <div className="bentoGrid">
            {/* CARD 1: Wide */}
            <div className="bentoCard bentoWide animateIn" style={{ animationDelay: '0.1s' }}>
              <div className="bentoBlob blobIndigo"></div>
              <div className="bentoContent">
                <div className="bentoHeader">
                  <div className="bentoIconWrapper featIconIndigo" aria-hidden="true">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                  </div>
                  <h3>Native Zero-Dependency</h3>
                </div>
                <p>Compiled down to a single, standalone executable. By avoiding external runtimes and bloated frameworks, it guarantees a hyper-optimized, invisible execution profile.</p>
              </div>
            </div>

            {/* CARD 2: Normal */}
            <div className="bentoCard animateIn" style={{ animationDelay: '0.2s' }}>
              <div className="bentoBlob blobBlue"></div>
              <div className="bentoContent">
                <div className="bentoHeader">
                  <div className="bentoIconWrapper featIconBlue" aria-hidden="true">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>
                  </div>
                </div>
                <h3>Stealth Input Mapping</h3>
                <p>Achieves 100% accuracy across all global keyboard layouts by hooking the low-level Windows ToUnicode API.</p>
              </div>
            </div>

            {/* CARD 3: Normal */}
            <div className="bentoCard animateIn" style={{ animationDelay: '0.3s' }}>
              <div className="bentoBlob blobPurple"></div>
              <div className="bentoContent">
                <div className="bentoHeader">
                  <div className="bentoIconWrapper featIconPurple" aria-hidden="true">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>
                  </div>
                </div>
                <h3>Rich Contextual Telemetry</h3>
                <p>Every atomic event is tagged in real-time with the foreground Process ID, app name, and exact window title.</p>
              </div>
            </div>

            {/* CARD 4: Wide */}
            <div className="bentoCard bentoWide animateIn" style={{ animationDelay: '0.4s' }}>
              <div className="bentoBlob blobPink"></div>
              <div className="bentoContent">
                <div className="bentoHeader">
                  <div className="bentoIconWrapper featIconPink" aria-hidden="true">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="17 1 21 5 17 9"/><path d="M3 11V9a4 4 0 014-4h14"/><polyline points="7 23 3 19 7 15"/><path d="M21 13v2a4 4 0 01-4 4H3"/></svg>
                  </div>
                  <h3>Dual-Stream Exfiltration</h3>
                </div>
                <p>Bypasses the need for vulnerable local disk storage by safely multiplexing telemetry out to Telegram and Discord via rapid, encrypted webhook streams.</p>
              </div>
            </div>

            {/* CARD 5: Wide */}
            <div className="bentoCard bentoWide animateIn" style={{ animationDelay: '0.5s' }}>
              <div className="bentoBlob blobTeal"></div>
              <div className="bentoContent">
                <div className="bentoHeader">
                  <div className="bentoIconWrapper featIconTeal" aria-hidden="true">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
                  </div>
                  <h3>Volatile Screen Capture</h3>
                </div>
                <p>Waits dormant until critical keyword triggers or physical clicks occur, instantly capturing the active UI state purely in RAM to leave zero forensic disk footprint.</p>
              </div>
            </div>

            {/* CARD 6: Normal */}
            <div className="bentoCard animateIn" style={{ animationDelay: '0.6s' }}>
              <div className="bentoBlob blobAmber"></div>
              <div className="bentoContent">
                <div className="bentoHeader">
                  <div className="bentoIconWrapper featIconAmber" aria-hidden="true">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0110 0v4"/></svg>
                  </div>
                </div>
                <h3>Air-Gapped Encryption</h3>
                <p>For scenarios requiring local backups, all data is locked down with AES-256-GCM authenticated encryption, cryptographically bound solely to your custom master key.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ===== SETUP GUIDE ===== */}
      <section className="section sectionAlt" id="setup">
        <div className="container">
          <div className="sectionHeader animateIn">
            <span className="sectionLabel">Getting Started</span>
            <h2 className="sectionTitle">Deployment Guide</h2>
            <p className="sectionDesc">Choose your preferred method to get up and running.</p>
          </div>
          <SetupTabs />
        </div>
      </section>

      {/* ===== TOOLS ===== */}
      <section className="section" id="tools">
        <div className="container">
          <div className="sectionHeader animateIn">
            <span className="sectionLabel">Included</span>
            <h2 className="sectionTitle">Companion Tools</h2>
            <p className="sectionDesc">The release includes utilities for forensic analysis and clean removal.</p>
          </div>

          <div className="toolsGrid">
            <div className="toolCard animateIn">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                Forensic Decryptor
              </h3>
              <p>A local web-based dashboard for reconstructing and analyzing encrypted log backups.</p>
              <ul>
                <li>Spawns a local server with a premium analysis dashboard</li>
                <li>Secure password-based session authentication</li>
                <li>Real-time filtering by window title or keystroke content</li>
                <li>Visual identification of screenshot capture events</li>
                <li>Full JSON export for external forensic tools</li>
              </ul>
            </div>

            <div className="toolCard animateIn">
              <h3>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                System Uninstaller
              </h3>
              <p>Complete removal utility that cleanly decommissions the monitoring framework.</p>
              <ul>
                <li>Terminates all active monitor processes</li>
                <li>Unregisters Windows hook persistence routines</li>
                <li>Removes scheduled tasks and registry entries</li>
                <li>Purges all encrypted telemetry records from disk</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* ===== DOWNLOAD ===== */}
      <DownloadSection />

      {/* ===== FOOTER ===== */}
      <Footer />
    </>
  );
}
