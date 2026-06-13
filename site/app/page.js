import Nav from '@/components/Nav';
import ScrollAnimations from '@/components/ScrollAnimations';
import SetupTabs from '@/components/SetupTabs';
import DownloadSection from '@/components/DownloadSection';
import Footer from '@/components/Footer';

export default function Home() {
  return (
    <>
      <ScrollAnimations />
      <Nav />

      {/* ===== HERO ===== */}
      <section className="hero" id="hero">
        <div className="heroBackground"></div>
        <div className="heroContent">
          <h1 className="heroTitle reveal delay-1">
            <span className="highlight">ShadowLog</span>
          </h1>
          
          <p className="heroSubtitle reveal delay-2">
            A high-performance, native systems monitoring framework engineered for 
            authorized security research. Built in Go (Windows) and Kotlin (Android) with zero runtime dependencies, 
            featuring quad-channel exfiltration and AES-256-GCM encryption.
          </p>

          <div className="heroActions reveal delay-3">
            <a href="#setup" className="btn btn-primary btn-lg">
              Deployment Guide
            </a>
            <a href="#download" className="btn btn-ghost btn-lg">
              <svg className="btn-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Download Archive
            </a>
          </div>
        </div>
      </section>

      {/* ===== WARNING BANNER ===== */}
      <section className="warningBanner" id="legal-warning" role="alert" style={{ background: '#fff', borderTop: '1px solid var(--glass-stroke)', borderBottom: '1px solid var(--glass-stroke)', padding: '24px 0' }}>
        <div className="container" style={{ display: 'flex', alignItems: 'flex-start', gap: '16px', maxWidth: '800px', margin: '0 auto' }}>
          <div style={{ flexShrink: 0, width: '24px', height: '24px', color: 'var(--danger)', marginTop: '4px' }}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
          </div>
          <div>
            <h3 style={{ fontSize: '1rem', color: 'var(--text-primary)', marginBottom: '4px', fontWeight: 600 }}>Legal &amp; Ethical Notice</h3>
            <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', lineHeight: 1.5 }}>
              This software is strictly for <strong>authorized educational and security research purposes</strong>.
              Deploying monitoring tools on systems without explicit, documented authorization violates federal laws including the CFAA.
            </p>
          </div>
        </div>
      </section>

      {/* ===== FEATURES ===== */}
      <section className="section" id="features">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Architecture</span>
            <h2 className="sectionTitle">Core Capabilities</h2>
            <p className="sectionDesc">Designed for stealth and reliability with a highly optimized system footprint.</p>
          </div>

          <div className="bentoGrid">
            <div className="bentoItem col-span-2 reveal">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              </div>
              <h3>Native Zero-Dependency</h3>
              <p>Compiled down to a single, standalone executable (~10.3 MB). By avoiding external runtimes (.NET, Java, etc.) and bloated frameworks, it guarantees a hyper-optimized execution profile directly on the Windows API.</p>
            </div>

            <div className="bentoItem col-span-2 reveal delay-1">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="17 1 21 5 17 9"/><path d="M3 11V9a4 4 0 014-4h14"/><polyline points="7 23 3 19 7 15"/><path d="M21 13v2a4 4 0 01-4 4H3"/></svg>
              </div>
              <h3>Quad-Channel Exfiltration</h3>
              <p>Multiplexes telemetry across Discord, Telegram Bot API, SMTP email, and DNS-over-HTTPS (DoH) fallback securely.</p>
            </div>

            <div className="bentoItem col-span-1 reveal">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>
              </div>
              <h3>Rich Contextual Telemetry</h3>
              <p>Every atomic event is tagged in real-time with the foreground Process ID, application name, and precise window title.</p>
            </div>

            <div className="bentoItem col-span-1 reveal delay-1">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
              </div>
              <h3>Volatile Screen Capture</h3>
              <p>Instantly captures the active UI state purely in RAM to leave absolutely zero forensic disk footprint.</p>
            </div>

            <div className="bentoItem col-span-2 reveal delay-2">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
              </div>
              <h3>Clipboard Monitoring</h3>
              <p>Automatically intercepts and securely logs text-based clipboard events to capture sensitive copy-paste actions.</p>
            </div>

            <div className="bentoItem col-span-4 reveal">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
              </div>
              <h3>Automated USB Exfiltration</h3>
              <p>On systems with isolated or heavily monitored networks, ShadowLog can automatically dump encrypted telemetry arrays to an authorized, pre-configured USB drive upon insertion, bypassing network firewalls entirely.</p>
            </div>
          </div>
        </div>
      </section>

      {/* ===== DOWNLOAD ===== */}
      <DownloadSection />

      {/* ===== WHAT'S NEW ===== */}
      <section className="section section-alt" id="changelog">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Latest Release</span>
            <h2 className="sectionTitle">What&apos;s New in v4.0</h2>
          </div>

          <div className="whatsNewGrid reveal delay-1">
            <div className="wnCard">
              <div className="wnIcon wnIconCore">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>
              </div>
              <h3>Intelligent Screenshot Batching</h3>
              <p>3-shot burst for sensitive windows (login, banking), single shot for everything else. Automatic context-aware capture.</p>
            </div>
            <div className="wnCard">
              <div className="wnIcon wnIconCore">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M1 6v16l7-4 8 4 7-4V2l-7 4-8-4-7 4z"/><line x1="8" y1="2" x2="8" y2="18"/><line x1="16" y1="6" x2="16" y2="22"/></svg>
              </div>
              <h3>Wi-Fi Password Extraction</h3>
              <p>Automatically retrieves and logs plaintext passwords for all connected wireless networks on the target machine.</p>
            </div>
            <div className="wnCard">
              <div className="wnIcon wnIconTool">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
              </div>
              <h3>Smart Decryptor Auto-Detection</h3>
              <p>The forensic decryptor now automatically detects if you&apos;re on the same machine and decrypts without requiring a password.</p>
            </div>
            <div className="wnCard">
              <div className="wnIcon wnIconPerf">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
              </div>
              <h3>30s Telegram Polling</h3>
              <p>Long-polling for near-instant remote command execution. Kill switch, pause, resume, and status checks respond in seconds.</p>
            </div>
          </div>


        </div>
      </section>

      {/* ===== SETUP GUIDE ===== */}
      <section className="section" id="setup">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Deployment</span>
            <h2 className="sectionTitle">Getting Started</h2>
            <p className="sectionDesc">Step-by-step instructions for authorized deployment on Windows and Android.</p>
          </div>
          <div className="reveal delay-1">
            <SetupTabs />
          </div>
        </div>
      </section>

      {/* ===== TOOLS & SUPPORT ===== */}
      <section className="section section-alt" id="tools">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Resources</span>
            <h2 className="sectionTitle">Tools &amp; Support</h2>
            <p className="sectionDesc">Companion utilities included in every release, plus solutions for common deployment issues.</p>
          </div>

          {/* Companion Tools */}
          <div className="toolsGrid reveal delay-1">
            <div className="toolCard">
              <h3>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                Forensic Decryptor
              </h3>
              <p>Local web-based dashboard for reconstructing and analyzing AES-encrypted telemetry.</p>
              <ul>
                <li>Premium analysis dashboard with filtering</li>
                <li>Auto-detection of machine-local encryption keys</li>
                <li>Visual screenshot event identification</li>
              </ul>
            </div>

            <div className="toolCard">
              <h3>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                System Uninstaller
              </h3>
              <p>Complete removal utility that cleanly decommissions the monitoring framework.</p>
              <ul>
                <li>Terminates all active monitor processes</li>
                <li>Unregisters Windows hook persistence</li>
                <li>Purges encrypted telemetry from disk</li>
              </ul>
            </div>
          </div>

          {/* Troubleshooting */}
          <div className="supportGrid reveal delay-2">
            <h3 className="supportHeading">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
              Common Issues
            </h3>
            <div className="supportCards">
              <div className="supportCard">
                <h4>
                  <span className="supportDot supportDotWarn"></span>
                  Windows Defender / SmartScreen
                </h4>
                <p>Due to its use of low-level system monitoring APIs (like keyboard hooks) and the absence of an EV code-signing certificate, antiviruses will flag ShadowLog as a Trojan. This is a <strong>false positive</strong>. Add the folder as a <strong>Windows Defender exclusion</strong> before running. For SmartScreen, click <strong>More info → Run anyway</strong>.</p>
              </div>
              <div className="supportCard">
                <h4>
                  <span className="supportDot supportDotAndroid"></span>
                  Play Protect Blocking APK
                </h4>
                <p>Temporarily disable Play Protect via <strong>Play Store → Profile → Play Protect → Settings → Toggle OFF</strong>. Re-enable after install. If still blocked, download the archive to a laptop and run <code>install-mobile.bat</code> via USB.</p>
              </div>
              <div className="supportCard">
                <h4>
                  <span className="supportDot supportDotKey"></span>
                  Decryptor Password Issues
                </h4>
                <p>On the same machine, decryption is automatic (v4.0+). On a different machine, you <strong>must</strong> enter the custom password set during setup. Machine GUID-based keys cannot be transferred.</p>
              </div>
              <div className="supportCard">
                <h4>
                  <span className="supportDot"></span>
                  3-Minute Startup Delay
                </h4>
                <p>Intentional by design — avoids heuristic-based detection that flags immediate system hook initialization at boot. Telemetry collection begins automatically after the delay.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ===== FOOTER ===== */}
      <Footer />
    </>
  );
}
