# Security and Antivirus Transparency

NMS Personal Atlas is an independent Windows application, distributed without a digital code-signing certificate.

## Antivirus detections
The Windows launcher in V12.7.241 has been reported by VirusTotal as detected by 6 of 71 vendors. This report must be taken seriously. A detection is not by itself proof of malware, but a 6/71 result is **not proof of safety** either.

Please do not disable antivirus protection to run Atlas. If security software quarantines the file, report the detection with the exact SHA-256 and vendor name so it can be investigated.

## Relevant behaviors
- The native Go GUI launcher starts the bundled Python runtime.
- Initial dependency installation downloads pip and packages from external services over HTTPS.
- Atlas creates and removes temporary installation files and caches.
- The uninstaller removes Atlas-owned files after a confirmation dialog. It may schedule removal of its own executable after the next Windows reboot.
- The Streamlit interface is bound to the local loopback interface (`127.0.0.1`).
- Technical features may read NMS process memory; the Alliance Observatory may trigger a read-only background memory probe automatically when opened.

## Current limitations requiring follow-up
- This is **partial-source transparency**: the Go GUI launcher's original source was not present in the examined distribution archive, so a byte-for-byte reproducible build is not available.
- The first-run Python dependency installation is not fully reproducible: version ranges are used and a separately verified, pinned dependency wheelhouse is not included.
- The static examination does not establish the complete runtime behavior on every Windows system.
- VirusTotal sandbox classifications such as process injection, sandbox evasion, and data destruction must be investigated against detailed behavioral events. This repository does not claim that those classifiers have been disproved.

## Responsible reporting
For a suspected security issue, include the exact filename, SHA-256, detection name, version, Windows version, and reproduction steps. Do not publish private game saves or other users' data in a public issue.

Contact: personnalatlas@gmail.com
