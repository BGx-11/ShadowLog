'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import Link from 'next/link';

const TIMER_DURATION = 10;
const CIRCUMFERENCE = 2 * Math.PI * 36;
const DOWNLOAD_URL = 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Release.zip';

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
              <div className="metaBadge">Version 2.3.1</div>
              <div className="metaBadge">Windows x64</div>
              <div className="metaBadge">9.1 MB Size</div>
            </div>

            <button className="btn btn-primary btn-lg" onClick={openModal}>
              <svg className="btn-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Acquire Archive
            </button>
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
