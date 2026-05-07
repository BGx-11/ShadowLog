'use client';

import { useState } from 'react';

export default function SetupTabs() {
  const [activeTab, setActiveTab] = useState('quick');

  return (
    <>
      <div className="setupTabs animateIn" role="tablist">
        <button
          className={`setupTab${activeTab === 'quick' ? ' setupTabActive' : ''}`}
          role="tab"
          aria-selected={activeTab === 'quick'}
          onClick={() => setActiveTab('quick')}
        >
          Quick Start
        </button>
        <button
          className={`setupTab${activeTab === 'source' ? ' setupTabActive' : ''}`}
          role="tab"
          aria-selected={activeTab === 'source'}
          onClick={() => setActiveTab('source')}
        >
          Build from Source
        </button>
      </div>

      {/* Quick Start Panel */}
      <div className={`setupPanel${activeTab === 'quick' ? ' setupPanelActive' : ''}`} role="tabpanel">
        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">01</div>
          <div className="stepContent">
            <h4>Download &amp; Extract</h4>
            <p>Download the latest <code>ShadowLog_Release.zip</code> from the download section below. Extract the archive on the target machine.</p>
          </div>
        </div>

        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">02</div>
          <div className="stepContent">
            <h4>Run as Administrator</h4>
            <p>Run the core executable with administrator privileges to start the setup wizard.</p>
            <p className="stepNote">If run without Administrator privileges, the tool will trigger a native Windows warning and terminate to prevent silent hooking failures.</p>
            <p className="stepNote">A terminal window may briefly flash while system hooks are being successfully registered. This is normal behavior.</p>
          </div>
        </div>

        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">03</div>
          <div className="stepContent">
            <h4>Configure Credentials</h4>
            <p>The setup wizard will prompt you to configure the following:</p>
            <ul className="substepList">
              <li><strong>Encryption Password</strong> (Required) — Choose a strong password for securing local encrypted backups.</li>
              <li><strong>Discord Webhook</strong> (Recommended) — Create a private channel, go to Edit Channel → Integrations → Webhooks → New Webhook, and paste the URL.</li>
              <li><strong>Telegram Bot Token</strong> (Optional) — Message @BotFather on Telegram, send /newbot, and follow the prompts to get your API token.</li>
              <li><strong>Telegram Chat ID</strong> (Optional) — Create a group, invite your bot, send /test, then visit the getUpdates API endpoint to find the chat ID.</li>
            </ul>
          </div>
        </div>

        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">04</div>
          <div className="stepContent">
            <h4>Test &amp; Deploy</h4>
            <p>Click <strong>Test Configuration</strong> to verify your setup. Then click <strong>Initialize Monitor</strong> to lock configuration and start the background service.</p>
          </div>
        </div>
      </div>

      {/* Build from Source Panel */}
      <div className={`setupPanel${activeTab === 'source' ? ' setupPanelActive' : ''}`} role="tabpanel">
        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">01</div>
          <div className="stepContent">
            <h4>Prerequisites</h4>
            <ul className="substepList">
              <li><strong>Go 1.21+</strong> — Required for compilation (1.26+ recommended for garble obfuscation)</li>
              <li><strong>Windows 10/11</strong> — Target environment for native system hooks</li>
              <li><strong>garble</strong> (Optional) — <code>go install mvdan.cc/garble@latest</code> for string literal obfuscation</li>
            </ul>
          </div>
        </div>

        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">02</div>
          <div className="stepContent">
            <h4>Clone the Repository</h4>
            <div className="codeBlock">
              <code>git clone https://github.com/BGx-11/ShadowLog.git<br/>cd ShadowLog</code>
            </div>
          </div>
        </div>

        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">03</div>
          <div className="stepContent">
            <h4>Compile Binaries</h4>
            <p>Generate optimized, stripped binaries:</p>
            <div className="codeBlock">
              <code>
                <span className="codeComment"># Core Monitor</span><br/>
                go build -trimpath -ldflags &quot;-H windowsgui -s -w -buildid=&quot; -o WinUpdateSvc.exe main.go<br/><br/>
                <span className="codeComment"># Forensic Decryptor</span><br/>
                go build -trimpath -ldflags &quot;-H windowsgui -s -w -buildid=&quot; -o Decryptor.exe decryptor/main.go<br/><br/>
                <span className="codeComment"># System Uninstaller</span><br/>
                go build -trimpath -ldflags &quot;-H windowsgui -s -w -buildid=&quot; -o Uninstaller.exe uninstaller/main.go
              </code>
            </div>
          </div>
        </div>

        <div className="setupStep animateIn">
          <div className="stepNumber" aria-hidden="true">04</div>
          <div className="stepContent">
            <h4>Run Setup</h4>
            <p>Follow the Quick Start guide from Step 2 onwards to configure and deploy.</p>
          </div>
        </div>
      </div>
    </>
  );
}
