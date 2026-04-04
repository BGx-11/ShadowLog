'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import Link from 'next/link';

const TIMER_DURATION = 12;
const CIRCUMFERENCE = 2 * Math.PI * 45;
const DOWNLOAD_URL = 'https://github.com/BGx-11/ShadowLog/releases/latest/download/ShadowLog_Release.zip';

export default function DownloadSection() {
  const [modalOpen, setModalOpen] = useState(false);
  const [timerValue, setTimerValue] = useState(TIMER_DURATION);
  const [timerDone, setTimerDone] = useState(false);
  const [checks, setChecks] = useState([false, false, false]);
  const intervalRef = useRef(null);

  const allChecked = checks.every(Boolean);
  const ready = timerDone && allChecked;

  const openModal = useCallback(() => {
    setModalOpen(true);
    setTimerValue(TIMER_DURATION);
    setTimerDone(false);
    setChecks([false, false, false]);
    document.body.style.overflow = 'hidden';
  }, []);

  const closeModal = useCallback(() => {
    setModalOpen(false);
    document.body.style.overflow = '';
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }, []);

  useEffect(() => {
    if (!modalOpen) return;

    intervalRef.current = setInterval(() => {
      setTimerValue((prev) => {
        if (prev <= 1) {
          clearInterval(intervalRef.current);
          intervalRef.current = null;
          setTimerDone(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [modalOpen]);

  useEffect(() => {
    function onKeyDown(e) {
      if (e.key === 'Escape' && modalOpen) closeModal();
    }
    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [modalOpen, closeModal]);

  function toggleCheck(index) {
    setChecks((prev) => {
      const copy = [...prev];
      copy[index] = !copy[index];
      return copy;
    });
  }

  function handleDownload() {
    if (!ready) return;
    window.location.href = DOWNLOAD_URL;
    setTimeout(closeModal, 2000);
  }

  const progress = (TIMER_DURATION - timerValue) / TIMER_DURATION;
  const dashOffset = progress * CIRCUMFERENCE;

  let btnText = 'Waiting for acknowledgement...';
  if (ready) btnText = 'Download ShadowLog_Release.zip';
  else if (timerDone && !allChecked) btnText = 'Please check all boxes above';
  else if (!timerDone && allChecked) btnText = `Waiting for timer (${timerValue}s)...`;
  else if (!timerDone) btnText = 'Waiting for acknowledgement...';

  return (
    <>
      {/* DOWNLOAD SECTION */}
      <section className="downloadSection" id="download">
        <div className="container">
          <div className="sectionHeader animateIn">
            <span className="sectionLabel">Release</span>
            <h2 className="sectionTitle">Download ShadowLog</h2>
          </div>

          <div className="downloadCard animateIn">
            <div className="downloadBlob"></div>
            <div className="downloadContent">
              <div className="downloadIcon" aria-hidden="true">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              </div>
              <h2 className="downloadHeading">ShadowLog Release</h2>
              <p className="downloadDesc">
                Contains the compiled monitor, forensic decryptor, and system uninstaller. Windows 10/11 required.
              </p>

              <div className="downloadMeta">
                <div className="downloadMetaItem"><span>Format</span> ZIP Archive</div>
                <div className="downloadMetaItem"><span>Platform</span> Windows x64</div>
                <div className="downloadMetaItem"><span>Includes</span> 3 executables</div>
              </div>

              <div className="downloadBtnWrapper">
                <div className="downloadBtnGlow"></div>
                <button className="btn btnPrimary downloadBtn" onClick={openModal} type="button">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{marginRight: '8px'}}><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                  Download Master Archive
                </button>
              </div>

              <p className="downloadLegal">
                By downloading, you acknowledge the <Link href="/terms">Terms of Service</Link>{' '}
                and <Link href="/privacy">Privacy Policy</Link>. This tool is for authorized use only.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* MODAL */}
      <div
        className={`modalOverlay${modalOpen ? ' modalOverlayActive' : ''}`}
        role="dialog"
        aria-modal="true"
        aria-label="Security Acknowledgement"
        onClick={(e) => { if (e.target === e.currentTarget) closeModal(); }}
      >
        <div className="modal">
          <div className="modalHeader">
            <div className="modalWarningIcon" aria-hidden="true">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
            <h2>Security Acknowledgement Required</h2>
          </div>

          <div>
            <div className="modalWarningText">
              <p>
                <strong>This is a powerful systems monitoring tool.</strong> Before proceeding, you must confirm
                that you understand the legal implications of downloading and using this software. Unauthorized
                surveillance is a serious criminal offense in most jurisdictions.
              </p>
            </div>

            <div className="modalCheckboxes">
              {[
                <>I confirm that I will <strong>only</strong> use this tool on systems I own or have explicit, documented written authorization to monitor.</>,
                <>I understand that unauthorized deployment may violate federal and state laws, including the <strong>Computer Fraud and Abuse Act</strong>, and I accept full personal responsibility for my use of this software.</>,
                <>I have read and agree to the <Link href="/terms" target="_blank">Terms of Service</Link> and <Link href="/privacy" target="_blank">Privacy Policy</Link>.</>,
              ].map((label, i) => (
                <label className="modalCheckbox" key={i}>
                  <input type="checkbox" checked={checks[i]} onChange={() => toggleCheck(i)} />
                  <span>{label}</span>
                </label>
              ))}
            </div>

            <div className="modalTimer">
              <div className="timerRing" role="timer" aria-label="Countdown timer">
                <svg viewBox="0 0 100 100" aria-hidden="true">
                  <circle className="ringBg" cx="50" cy="50" r="45"/>
                  <circle
                    className={`ringProgress${timerValue <= 5 && !timerDone ? ' ringProgressDanger' : ''}`}
                    cx="50" cy="50" r="45"
                    style={{ strokeDasharray: CIRCUMFERENCE, strokeDashoffset: dashOffset }}
                  />
                </svg>
                <span className="timerText">{timerDone ? '✓' : timerValue}</span>
              </div>
              <span className="timerLabel">{timerDone ? 'Timer complete' : 'Please read the terms above'}</span>
            </div>

            <button
              className={`btn btnPrimary modalDownloadBtn${ready ? ' modalDownloadBtnReady' : ''}`}
              disabled={!ready}
              onClick={handleDownload}
              type="button"
            >
              {btnText}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
