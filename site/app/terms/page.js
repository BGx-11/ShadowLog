import Nav from '@/components/Nav';
import Footer from '@/components/Footer';
import Link from 'next/link';

export const metadata = {
  title: 'Terms of Service — ShadowLog',
  description: 'Terms of Service for ShadowLog, a discrete activity analytics framework for authorized cybersecurity research.',
  robots: 'noindex',
};

export default function TermsPage() {
  return (
    <>
      <Nav alwaysOpaque />

      <main className="legalPage">
        <div className="legalContainer">
          <Link href="/" className="backLink">← Back to Home</Link>
          <h1>Terms of Service</h1>
          <p className="legalDate">Last updated: April 4, 2026</p>

          <h2>1. Acceptance of Terms</h2>
          <p>
            By accessing, downloading, or using ShadowLog (&ldquo;the Software&rdquo;), you agree to be bound by these Terms of Service (&ldquo;Terms&rdquo;).
            If you do not agree, you must not download, install, or use the Software. These Terms constitute a legally binding
            agreement between you (&ldquo;User&rdquo;) and Devansh Agarwal (&ldquo;Developer&rdquo;, &ldquo;Author&rdquo;).
          </p>

          <h2>2. Intended Use</h2>
          <p>
            ShadowLog is developed and distributed <strong>exclusively for authorized cybersecurity research, educational purposes,
            and legitimate security assessments</strong>. The Software is intended for:
          </p>
          <ul>
            <li>Security professionals conducting authorized penetration testing</li>
            <li>Cybersecurity students studying system-level hooks and monitoring techniques</li>
            <li>System administrators monitoring their own infrastructure with proper authorization</li>
            <li>Researchers analyzing activity capture and exfiltration methodologies in controlled environments</li>
          </ul>

          <h2>3. Prohibited Use</h2>
          <p>You agree <strong>not</strong> to use ShadowLog for any of the following purposes:</p>
          <ul>
            <li>Monitoring, surveilling, or capturing activity on any system, device, or network you do not own or have explicit, documented written authorization to monitor</li>
            <li>Stalking, harassment, or any form of unauthorized surveillance of individuals</li>
            <li>Corporate espionage, theft of trade secrets, or unauthorized data collection</li>
            <li>Any activity that violates applicable local, state, national, or international laws</li>
            <li>Distribution of the Software bundled with malware, ransomware, or other malicious payloads</li>
          </ul>

          <h2>4. No Warranty</h2>
          <p>
            The Software is provided <strong>&ldquo;AS IS&rdquo;</strong> without warranty of any kind, express or implied, including but not
            limited to warranties of merchantability, fitness for a particular purpose, and non-infringement. The Developer does
            not warrant that the Software will meet your requirements or that its operation will be uninterrupted or error-free.
          </p>

          <h2>5. Limitation of Liability</h2>
          <p>
            In no event shall the Developer be liable for any direct, indirect, incidental, special, consequential, or punitive
            damages arising out of or in connection with the use or inability to use the Software. This includes, without
            limitation, damages for loss of data, loss of profits, or any legal consequences resulting from the use or misuse of the Software.
          </p>
          <p>
            <strong>You assume full responsibility</strong> for your use of the Software. The Developer expressly disclaims any
            liability for actions taken by users that violate these Terms or any applicable law.
          </p>

          <h2>6. Indemnification</h2>
          <p>
            You agree to indemnify, defend, and hold harmless the Developer from and against any and all claims, liabilities,
            damages, losses, costs, and expenses (including reasonable attorney&apos;s fees) arising from your use or misuse of the
            Software, your violation of these Terms, or your violation of any applicable law.
          </p>

          <h2>7. Legal Compliance</h2>
          <p>
            You acknowledge that the use of monitoring software may be subject to laws such as the Computer Fraud and Abuse Act
            (CFAA), the Electronic Communications Privacy Act (ECPA), the General Data Protection Regulation (GDPR), and equivalent
            laws in your jurisdiction. <strong>It is your sole responsibility</strong> to understand and comply with all applicable laws
            before using the Software.
          </p>

          <h2>8. Intellectual Property</h2>
          <p>
            ShadowLog is the intellectual property of the Developer. You are granted a limited, non-exclusive, non-transferable
            license to use the Software in accordance with these Terms. You may not sublicense, sell, or redistribute the Software
            for commercial purposes without explicit written permission.
          </p>

          <h2>9. Modifications to Terms</h2>
          <p>
            The Developer reserves the right to modify these Terms at any time. Continued use of the Software after changes
            constitutes acceptance of the updated Terms. Users are responsible for reviewing these Terms periodically.
          </p>

          <h2>10. Governing Law</h2>
          <p>
            These Terms shall be governed by and construed in accordance with applicable laws. Any disputes arising under
            these Terms shall be resolved in the appropriate courts of competent jurisdiction.
          </p>

          <h2>11. Contact</h2>
          <p>
            For questions regarding these Terms, please reach out via the Developer&apos;s portfolio at{' '}
            <a href="https://iambgx.in" target="_blank" rel="noopener noreferrer">iambgx.in</a>.
          </p>
        </div>
      </main>

      <Footer />
    </>
  );
}
