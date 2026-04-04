import Nav from '@/components/Nav';
import Footer from '@/components/Footer';
import Link from 'next/link';

export const metadata = {
  title: 'Privacy Policy — ShadowLog',
  description: 'Privacy Policy for the ShadowLog distribution website.',
  robots: 'noindex',
};

export default function PrivacyPage() {
  return (
    <>
      <Nav alwaysOpaque />

      <main className="legalPage">
        <div className="legalContainer">
          <Link href="/" className="backLink">← Back to Home</Link>
          <h1>Privacy Policy</h1>
          <p className="legalDate">Last updated: April 4, 2026</p>

          <h2>1. Overview</h2>
          <p>
            This Privacy Policy describes how information is handled in relation to the ShadowLog distribution website
            (&ldquo;the Site&rdquo;) and the ShadowLog software (&ldquo;the Software&rdquo;). The Developer, Devansh Agarwal (&ldquo;BGx&rdquo;), is committed
            to transparency about data practices.
          </p>

          <h2>2. Information Collected by the Website</h2>
          <p>
            <strong>This website does not collect, store, or process any personal data.</strong> There are no cookies,
            no analytics trackers, no sign-up forms, and no user accounts. The Site is a static distribution page hosted
            on Vercel.
          </p>
          <p>
            Vercel, as the hosting provider, may collect standard server logs (IP address, browser user agent, request
            timestamps) as part of its infrastructure. This data is subject to{' '}
            <a href="https://vercel.com/legal/privacy-policy" target="_blank" rel="noopener noreferrer">Vercel&apos;s Privacy Policy</a>{' '}
            and is not accessed or controlled by the Developer.
          </p>

          <h2>3. Information Collected by the Software</h2>
          <p>
            ShadowLog is a systems monitoring tool designed for authorized research. When deployed by an authorized user,
            the Software captures the following types of data on the host system:
          </p>
          <ul>
            <li>Keystroke data correlated with active window titles and process metadata</li>
            <li>Screenshots of the active window triggered by keyword detection or user interaction</li>
            <li>Encrypted local log files stored on the host machine</li>
          </ul>
          <p>
            <strong>This data is generated and controlled entirely by the user who deploys the Software.</strong>{' '}
            The Developer does not receive, access, or have any visibility into data captured by deployed instances
            of ShadowLog. The Developer has no telemetry, phone-home functionality, or remote data collection
            capabilities embedded in the Software.
          </p>

          <h2>4. Data Exfiltration Channels</h2>
          <p>
            The Software provides optional integrations with third-party services (Discord, Telegram) configured
            entirely by the deploying user. Data transmitted through these channels is subject to the respective
            privacy policies of those platforms:
          </p>
          <ul>
            <li><a href="https://discord.com/privacy" target="_blank" rel="noopener noreferrer">Discord Privacy Policy</a></li>
            <li><a href="https://telegram.org/privacy" target="_blank" rel="noopener noreferrer">Telegram Privacy Policy</a></li>
          </ul>
          <p>
            The Developer has no access to or control over data sent through these user-configured channels.
          </p>

          <h2>5. Data Security</h2>
          <p>
            Local log backups created by the Software are encrypted using AES-256-GCM with a user-defined password.
            The security of this data depends on the strength of the password chosen by the deploying user. The
            Developer is not responsible for data breaches resulting from weak passwords or improper handling of
            encrypted files.
          </p>

          <h2>6. Third-Party Links</h2>
          <p>
            This Site may contain links to external websites (GitHub, portfolio, third-party privacy policies).
            The Developer is not responsible for the content or privacy practices of external sites. Users should
            review the privacy policies of any linked websites independently.
          </p>

          <h2>7. Children&apos;s Privacy</h2>
          <p>
            This Software is not intended for use by individuals under the age of 18. The Developer does not
            knowingly distribute software to minors.
          </p>

          <h2>8. User Responsibility</h2>
          <p>
            Users who deploy ShadowLog are solely responsible for compliance with applicable data protection laws
            (including GDPR, CCPA, and local equivalents) in relation to any data captured by the Software.
            The Developer assumes no responsibility for how captured data is stored, transmitted, or used
            by the deploying user.
          </p>

          <h2>9. Changes to This Policy</h2>
          <p>
            The Developer may update this Privacy Policy at any time. Changes will be reflected by an updated
            &ldquo;Last updated&rdquo; date. Continued use of the Site or Software after changes constitutes acceptance of
            the updated policy.
          </p>

          <h2>10. Contact</h2>
          <p>
            For privacy-related inquiries, please reach out via the Developer&apos;s portfolio at{' '}
            <a href="https://iambgx.in" target="_blank" rel="noopener noreferrer">iambgx.in</a>.
          </p>
        </div>
      </main>

      <Footer />
    </>
  );
}
