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
            authorized security research. Built entirely in Go with zero runtime dependencies, 
            featuring quad-channel exfiltration and AES-256-GCM encryption.
          </p>

          <div className="heroActions reveal delay-3">
            <a href="#setup" className="btn btn-primary btn-lg">
              Deployment Guide
            </a>
            <a href="#download" className="btn btn-ghost btn-lg">
              <svg className="btn-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Download Archive (10.3 MB)
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
            {/* 1 */}
            <div className="bentoItem col-span-2 reveal">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              </div>
              <h3>Native Zero-Dependency</h3>
              <p>Compiled down to a single, standalone executable (~10.3 MB). By avoiding external runtimes (.NET, Java, etc.) and bloated frameworks, it guarantees a hyper-optimized execution profile directly on the Windows API.</p>
            </div>

            {/* 2 */}
            <div className="bentoItem col-span-2 reveal delay-1">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="17 1 21 5 17 9"/><path d="M3 11V9a4 4 0 014-4h14"/><polyline points="7 23 3 19 7 15"/><path d="M21 13v2a4 4 0 01-4 4H3"/></svg>
              </div>
              <h3>Quad-Channel Exfiltration</h3>
              <p>Multiplexes telemetry across Discord, Telegram Bot API, SMTP email, and DNS-over-HTTPS (DoH) fallback securely.</p>
            </div>

            {/* 3 */}
            <div className="bentoItem col-span-1 reveal">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>
              </div>
              <h3>Rich Contextual Telemetry</h3>
              <p>Every atomic event is tagged in real-time with the foreground Process ID, application name, and precise window title.</p>
            </div>

            {/* 4 */}
            <div className="bentoItem col-span-1 reveal delay-1">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
              </div>
              <h3>Volatile Screen Capture</h3>
              <p>Instantly captures the active UI state purely in RAM to leave absolutely zero forensic disk footprint.</p>
            </div>

            {/* 5 */}
            <div className="bentoItem col-span-2 reveal delay-2">
              <div className="bentoIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
              </div>
              <h3>Clipboard Monitoring</h3>
              <p>Automatically intercepts and securely logs text-based clipboard events to capture sensitive copy-paste actions.</p>
            </div>

            {/* 6 */}
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

      {/* ===== CHANGELOG ===== */}
      <section className="section section-alt" id="changelog">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Updates</span>
            <h2 className="sectionTitle">Release Notes</h2>
            <p className="sectionDesc">Track the evolution and newly added capabilities of the ShadowLog framework.</p>
          </div>

          <div className="changelogTimeline">
            <div className="changeItem reveal">
              <div className="changeDot"></div>
              <div className="changeContent">
                <div className="changeHeader">
                  <h3>Version 4.0 (Current)</h3>
                  <span className="versionBadge">Latest</span>
                </div>
                <ul>
                  <li><strong>Core:</strong> Intelligent Screenshot Batching — 3-shot burst for sensitive windows (login/bank), single shot for others.</li>
                  <li><strong>Core:</strong> Wi-Fi Password Extraction — Automatically retrieves and logs plaintext passwords for connected networks.</li>
                  <li><strong>Core:</strong> Improved Telegram responsiveness — 30s long-polling for near-instant command execution.</li>
                  <li><strong>Decryptor:</strong> Smart auto-detection mode — automatically decrypts logs on the same machine without requiring a password.</li>
                  <li><strong>Decryptor:</strong> Improved password validation with clear error diagnostics for cross-machine usage.</li>
                  <li><strong>Build:</strong> Decryptor binary no longer compiled with hidden window flag, significantly reducing AV false positive rates.</li>
                  <li><strong>Site:</strong> Fixed download button pointing to incorrect release asset filename.</li>
                  <li><strong>Site:</strong> Added comprehensive troubleshooting FAQ for AV exclusions and decryptor setup.</li>
                  <li><strong>Core:</strong> Improved error handling and session management across all companion tools.</li>
                </ul>
              </div>
            </div>

            <div className="changeItem reveal delay-1">
              <div className="changeDot"></div>
              <div className="changeContent">
                <div className="changeHeader">
                  <h3>Version 2.3.1</h3>
                </div>
                <ul>
                  <li><strong>Core:</strong> Implemented automatic native UAC prompt at startup if administrative privileges are missing.</li>
                  <li><strong>Core:</strong> Changed default process initialization to sleep for 3 minutes to avoid heuristic detection at boot.</li>
                  <li><strong>Telemetry:</strong> Implemented active monitoring for Wi-Fi SSID network connections.</li>
                  <li><strong>Network:</strong> Batched screenshot delivery over Telegram to avoid rate limits and reduce CPU overhead.</li>
                  <li><strong>Uninstaller:</strong> Fixed a critical bug where the uninstaller failed to terminate the hidden process cleanly.</li>
                </ul>
              </div>
            </div>

            <div className="changeItem reveal delay-2">
              <div className="changeDot"></div>
              <div className="changeContent">
                <div className="changeHeader">
                  <h3>Version 2.2.0</h3>
                </div>
                <ul>
                  <li><strong>Exfiltration:</strong> Added fully automated offline USB dumping for air-gapped environments.</li>
                  <li><strong>Monitoring:</strong> Added comprehensive clipboard interception.</li>
                  <li><strong>Security:</strong> All screenshot buffers are now scrubbed from memory immediately post-encryption.</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ===== SETUP GUIDE ===== */}
      <section className="section" id="setup">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Deployment</span>
            <h2 className="sectionTitle">Documentation</h2>
            <p className="sectionDesc">Step-by-step instructions for authorized deployment.</p>
          </div>
          <div className="reveal delay-1">
            <SetupTabs />
          </div>
        </div>
      </section>

      {/* ===== TOOLS ===== */}
      <section className="section section-alt" id="tools">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Included</span>
            <h2 className="sectionTitle">Companion Tools</h2>
            <p className="sectionDesc">The release package includes essential utilities for forensic analysis and clean decommission.</p>
          </div>

          <div className="toolsGrid">
            <div className="toolCard reveal">
              <h3>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                Forensic Decryptor
              </h3>
              <p>A local web-based dashboard for reconstructing and analyzing AES-encrypted log backups.</p>
              <ul>
                <li>Spawns a local server with a premium analysis dashboard</li>
                <li>Real-time filtering by window title or keystrokes</li>
                <li>Visual identification of screenshot capture events</li>
              </ul>
            </div>

            <div className="toolCard reveal delay-1">
              <h3>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                System Uninstaller
              </h3>
              <p>Complete removal utility that cleanly decommissions the monitoring framework.</p>
              <ul>
                <li>Terminates active monitor processes instantly</li>
                <li>Unregisters Windows hook persistence routines</li>
                <li>Purges all encrypted telemetry records from disk</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* ===== TROUBLESHOOTING FAQ ===== */}
      <section className="section section-alt" id="faq">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Support</span>
            <h2 className="sectionTitle">Troubleshooting</h2>
            <p className="sectionDesc">Common issues and solutions for deployment and forensic analysis.</p>
          </div>

          <div className="faqGrid">
            <div className="faqCard reveal">
              <div className="faqIcon faqIconWarning">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              </div>
              <h3>Windows Defender & SmartScreen Warnings</h3>
              <p>Because ShadowLog utilizes low-level Windows APIs, lacks a paid EV Authenticode Certificate, and uses process stealth techniques, <strong>Windows Defender and SmartScreen will universally flag it as a Trojan or "Unrecognized App"</strong> (commonly <code>Trojan:Win32/Bearfoos</code> or <code>Wacatac.B!ml</code>). This is mathematically unavoidable for unauthorized security tools without an expensive corporate digital signature.</p>
              <div className="faqSolution">
                <h4>Solution: Add a Windows Defender Exclusion</h4>
                <ol>
                  <li>Open <strong>Windows Security</strong> → Virus &amp; threat protection</li>
                  <li>Click <strong>Manage settings</strong> under Virus &amp; threat protection settings</li>
                  <li>Scroll to <strong>Exclusions</strong> → Add or remove exclusions</li>
                  <li>Click <strong>Add an exclusion</strong> → Folder → Select the extracted ShadowLog directory</li>
                </ol>
                <p className="faqNote">If SmartScreen blocks execution ("Windows protected your PC"), click <strong>More info</strong> and then <strong>Run anyway</strong>. For AV blocks, add an exclusion folder before extracting.</p>
              </div>
            </div>

            <div className="faqCard reveal delay-1">
              <div className="faqIcon faqIconKey">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
              </div>
              <h3>Decryptor Password Not Working</h3>
              <p>The Forensic Decryptor uses the encryption password set during initial setup to decrypt log files. If no password was configured, it derives a key from the machine&apos;s unique hardware GUID.</p>
              <div className="faqSolution">
                <h4>Key Points</h4>
                <ul>
                  <li><strong>Same machine, custom password:</strong> Enter the exact password you set during setup. As of v4.0, the decryptor will auto-detect if you&apos;re on the same machine.</li>
                  <li><strong>Same machine, no password:</strong> The decryptor will automatically unlock using the machine GUID — no password input required.</li>
                  <li><strong>Different machine:</strong> You <strong>must</strong> enter the custom encryption password from setup. Machine GUID-based keys are unique per device and cannot be transferred.</li>
                </ul>
              </div>
            </div>

            <div className="faqCard reveal delay-2">
              <div className="faqIcon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              </div>
              <h3>3-Minute Startup Delay</h3>
              <p>After installation, ShadowLog intentionally waits 3 minutes before initializing. This delay is by design to avoid heuristic-based detection that flags processes which immediately start system hooks at boot time.</p>
              <div className="faqSolution">
                <h4>Expected Behavior</h4>
                <p>The monitor process will appear idle for the first 3 minutes after a reboot. Telemetry collection begins automatically after this initialization period. No action is required.</p>
              </div>
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
