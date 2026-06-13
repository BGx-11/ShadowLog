import Nav from '@/components/Nav';
import Footer from '@/components/Footer';
import Link from 'next/link';

export const metadata = {
  title: 'Terms of Service',
  description: 'Terms of Service for ShadowLog, a discrete activity analytics framework for authorized cybersecurity research.',
  robots: 'noindex, nofollow',
};

export default function TermsPage() {
  return (
    <>
      <Nav alwaysOpaque />

      <main className="legalPage">
        <div className="legalContainer">
          <Link href="/" className="backLink">← Back to Home</Link>
          <h1>Terms of Service</h1>
          <p className="legalDate">Last updated: April 23, 2026</p>

          <h2>1. Acceptance of Terms</h2>
          <p>
            By accessing this website, downloading, installing, copying, or using ShadowLog (&ldquo;the Software&rdquo;) in
            any manner, you unconditionally agree to be bound by these Terms of Service (&ldquo;Terms&rdquo;). If you do not
            agree to every provision of these Terms, you are strictly prohibited from downloading, installing, possessing,
            or using the Software in any capacity. These Terms constitute a legally binding agreement between you
            (&ldquo;User&rdquo;, &ldquo;You&rdquo;) and Devansh Agarwal (&ldquo;Developer&rdquo;, &ldquo;Author&rdquo;, &ldquo;We&rdquo;).
          </p>

          <h2>2. Nature of the Software</h2>
          <p>
            ShadowLog is a <strong>proof-of-concept security research tool</strong> developed for the sole purpose of
            demonstrating system-level monitoring techniques in controlled, authorized environments. It is an
            <strong> educational and research instrument</strong>, not a commercial product. The Software is provided as
            open-source code to advance understanding of cybersecurity attack vectors and defensive strategies.
          </p>

          <h2>3. Authorized Use Only</h2>
          <p>
            The Software is developed and distributed <strong>exclusively for the following lawful purposes</strong>:
          </p>
          <ul>
            <li>Authorized penetration testing and red team exercises conducted by credentialed security professionals with explicit written authorization from system owners</li>
            <li>Academic and educational study of system-level hooks, monitoring techniques, and defensive countermeasures within controlled laboratory or virtual environments</li>
            <li>System administrators monitoring their own infrastructure or infrastructure they have documented, explicit authorization to monitor</li>
            <li>Security researchers analyzing activity capture, exfiltration methodologies, and encryption implementations in sandboxed environments</li>
            <li>Development of detection signatures, antivirus heuristics, and defensive security tools</li>
          </ul>
          <p>
            <strong>You represent and warrant</strong> that your use of the Software falls within one of the above categories
            and that you possess all necessary authorizations, licenses, and legal right to deploy monitoring software on
            any system where you use ShadowLog.
          </p>

          <h2>4. Strictly Prohibited Use</h2>
          <p>You expressly agree <strong>never</strong> to use ShadowLog for any of the following purposes:</p>
          <ul>
            <li>Monitoring, surveilling, intercepting, recording, or capturing any activity on any system, mobile device, network, or account that you do not own or for which you lack explicit, documented, written authorization from the legal owner</li>
            <li>Stalking, harassment, intimidation, domestic surveillance, or any form of unauthorized surveillance of any individual via their PC or Android device</li>
            <li>Corporate espionage, theft of trade secrets, intellectual property theft, unauthorized competitive intelligence gathering, or any unauthorized data collection</li>
            <li>Identity theft, credential harvesting, financial fraud, or any form of unauthorized access to accounts or systems</li>
            <li>Distribution of the Software bundled with malware, ransomware, trojans, rootkits, or other malicious payloads</li>
            <li>Use against minors or vulnerable individuals in any context</li>
            <li>Any activity that violates applicable local, state, national, or international laws, statutes, regulations, or ordinances</li>
            <li>Circumventing security measures, access controls, or authorization mechanisms on systems you are not authorized to test</li>
          </ul>
          <p>
            <strong>Any use of the Software outside the authorized purposes listed in Section 3 is a direct violation of
            these Terms and may constitute a criminal offense.</strong> The Developer will cooperate fully with law enforcement
            authorities in the investigation of any suspected illegal use of the Software.
          </p>

          <h2>5. Complete Disclaimer of Warranty</h2>
          <p>
            THE SOFTWARE IS PROVIDED <strong>&ldquo;AS IS&rdquo;</strong> AND <strong>&ldquo;AS AVAILABLE&rdquo;</strong> WITHOUT
            WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE IMPLIED WARRANTIES OF
            MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, TITLE, NON-INFRINGEMENT, AND ACCURACY. THE
            DEVELOPER MAKES NO WARRANTY THAT THE SOFTWARE WILL MEET YOUR REQUIREMENTS, THAT ITS OPERATION WILL
            BE UNINTERRUPTED, ERROR-FREE, OR SECURE, OR THAT DEFECTS WILL BE CORRECTED.
          </p>
          <p>
            THE DEVELOPER DOES NOT WARRANT, ENDORSE, GUARANTEE, OR ASSUME RESPONSIBILITY FOR ANY THIRD-PARTY
            APPLICATIONS, SERVICES, OR PLATFORMS THAT INTERACT WITH THE SOFTWARE.
          </p>

          <h2>6. Absolute Limitation of Liability</h2>
          <p>
            TO THE MAXIMUM EXTENT PERMITTED BY APPLICABLE LAW, IN NO EVENT SHALL THE DEVELOPER, HIS AFFILIATES,
            AGENTS, LICENSORS, OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, CONSEQUENTIAL,
            EXEMPLARY, OR PUNITIVE DAMAGES OF ANY KIND, INCLUDING WITHOUT LIMITATION:
          </p>
          <ul>
            <li>Loss of data, profits, revenue, business, goodwill, or anticipated savings</li>
            <li>Personal injury, emotional distress, or property damage</li>
            <li>Legal consequences, criminal charges, fines, penalties, or incarceration arising from use or misuse of the Software</li>
            <li>Damages arising from unauthorized access to or alteration of your transmissions or data</li>
            <li>Cost of procurement of substitute goods or services</li>
            <li>Any matter beyond the Developer&apos;s reasonable control</li>
          </ul>
          <p>
            THIS LIMITATION APPLIES REGARDLESS OF THE LEGAL THEORY (CONTRACT, TORT, STRICT LIABILITY, NEGLIGENCE,
            OR OTHERWISE) AND EVEN IF THE DEVELOPER HAS BEEN ADVISED OF THE POSSIBILITY OF SUCH DAMAGES. IN
            JURISDICTIONS THAT DO NOT ALLOW THE EXCLUSION OR LIMITATION OF INCIDENTAL OR CONSEQUENTIAL DAMAGES,
            THE DEVELOPER&apos;S LIABILITY SHALL BE LIMITED TO THE MAXIMUM EXTENT PERMITTED BY LAW.
          </p>

          <h2>7. Assumption of Risk</h2>
          <p>
            <strong>You expressly acknowledge and assume all risks</strong> associated with downloading, installing,
            possessing, and using the Software. You understand that:
          </p>
          <ul>
            <li>The Software is a powerful systems and mobile monitoring tool that can capture sensitive data, including keystrokes, notifications, and location data on Android</li>
            <li>Misuse of the Software may result in severe civil and criminal liability under both computer crime and wiretapping laws</li>
            <li>You are solely and exclusively responsible for ensuring that your use complies with all applicable laws</li>
            <li>The Developer has no control over and assumes no responsibility for how the Software is deployed or used after download</li>
          </ul>

          <h2>8. Indemnification</h2>
          <p>
            You agree to <strong>fully indemnify, defend, and hold harmless</strong> the Developer, his affiliates, agents,
            contributors, and licensors from and against any and all claims, liabilities, damages, losses, costs, expenses,
            and fees (including reasonable attorney&apos;s fees and court costs) arising from or related to:
          </p>
          <ul>
            <li>Your use, misuse, or deployment of the Software</li>
            <li>Your violation of these Terms</li>
            <li>Your violation of any applicable law, regulation, or third-party right</li>
            <li>Any claim by a third party arising from your use of the Software</li>
            <li>Any data captured, stored, transmitted, or disclosed through your use of the Software</li>
          </ul>
          <p>
            This indemnification obligation shall survive the termination of these Terms and your cessation of use of the Software.
          </p>

          <h2>9. Legal Compliance</h2>
          <p>
            You acknowledge that the use of monitoring software is regulated by numerous laws, including but not limited to:
          </p>
          <ul>
            <li>The Computer Fraud and Abuse Act (CFAA) — 18 U.S.C. § 1030</li>
            <li>The Electronic Communications Privacy Act (ECPA) — 18 U.S.C. §§ 2510–2523</li>
            <li>The Stored Communications Act (SCA) — 18 U.S.C. §§ 2701–2712</li>
            <li>The Wiretap Act — 18 U.S.C. §§ 2511–2522</li>
            <li>The General Data Protection Regulation (GDPR) — EU Regulation 2016/679</li>
            <li>The California Consumer Privacy Act (CCPA) and CPRA</li>
            <li>Equivalent federal, state, and international privacy and computer crime statutes in your jurisdiction</li>
          </ul>
          <p>
            <strong>It is your sole and exclusive responsibility</strong> to understand and comply with all applicable laws
            before downloading, installing, or using the Software. Ignorance of the law is not a defense.
          </p>

          <h2>10. Intellectual Property</h2>
          <p>
            ShadowLog is the intellectual property of the Developer. You are granted a limited, non-exclusive,
            non-transferable, revocable license to use the Software strictly in accordance with these Terms. You may not
            sublicense, sell, redistribute, or commercially exploit the Software without explicit written permission from
            the Developer.
          </p>

          <h2>11. Termination</h2>
          <p>
            The Developer may terminate or suspend your right to use the Software at any time, for any reason, without
            notice. Upon termination, you must immediately cease all use of the Software and destroy all copies in
            your possession. Sections 5, 6, 7, 8, and 9 shall survive termination.
          </p>

          <h2>12. Severability</h2>
          <p>
            If any provision of these Terms is held to be invalid, illegal, or unenforceable by a court of competent
            jurisdiction, the remaining provisions shall continue in full force and effect. The invalid provision shall
            be modified to the minimum extent necessary to make it valid and enforceable while preserving its original intent.
          </p>

          <h2>13. Entire Agreement</h2>
          <p>
            These Terms, together with the Privacy Policy, constitute the entire agreement between you and the Developer
            regarding the Software and supersede all prior agreements, understandings, and communications.
          </p>

          <h2>14. Modifications to Terms</h2>
          <p>
            The Developer reserves the right to modify these Terms at any time without prior notice. Changes become
            effective immediately upon posting to this page. Continued use of the Software after modifications
            constitutes acceptance of the updated Terms. You are responsible for reviewing these Terms periodically.
          </p>

          <h2>15. Governing Law &amp; Jurisdiction</h2>
          <p>
            These Terms shall be governed by and construed in accordance with the laws of India, without regard to
            conflict of law principles. Any disputes arising under these Terms shall be subject to the exclusive
            jurisdiction of the courts in Uttar Pradesh, India. You irrevocably consent to such jurisdiction and venue.
          </p>

          <h2>16. Contact</h2>
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
