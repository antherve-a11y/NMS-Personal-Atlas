# Security and Antivirus Transparency

**Project:** NMS Personal Atlas  
**Current release:** V12.7.257  
**Last updated:** November 24, 2026

NMS Personal Atlas is an independent Windows companion application for No Man's Sky. This document describes the distributed application's security characteristics, the scope of available verification, and how to report concerns privately.

## Supported releases

| Version | Security and compatibility attention |
| --- | --- |
| **V12.7.257** | Current release for new bug and vulnerability reports. |
| V12.7.255 and earlier | Historical releases; fixes are not necessarily backported. |
| V12.7.256-DEV | Development build, not the recommended public release. |

Security updates are best-effort. No formal response-time or support guarantee is offered.

## Antivirus report — exact V12.7.257 archive

The user-supplied VirusTotal screenshot dated **November 24, 2026** shows **0/47 security vendors flagging the uploaded ZIP** at the time of the screenshot. The displayed result and participating engine count may change if VirusTotal reanalyzes the file.

- **Filename:** `NMS_Personal_Atlas_V12.7.257.zip`
- **SHA-256:** `b830aab0da0a442431dfe15da1f2fc0a61f5384e5240878550f1a941d48bf180`
- **VirusTotal:** [View the public detection report](https://www.virustotal.com/gui/file/b830aab0da0a442431dfe15da1f2fc0a61f5384e5240878550f1a941d48bf180/detection)
- **Repository checksum list:** [SHA256SUMS.txt](SHA256SUMS.txt)

**Scope:** This result applies **only** to the ZIP whose SHA-256 matches the value above. It is not a separate antivirus scan of every executable inside the archive. Recompressing, altering or rebuilding the ZIP generally changes its checksum and invalidates the identity check.

**Important limitations:** A 0/47 result does **not** certify the program as malware-free or guarantee that every security product will accept it. The project does not have a paid code-signing certificate. Windows SmartScreen and third-party antivirus products may still warn about unknown or unsigned software. Never disable antivirus or SmartScreen protections to install Atlas.

Earlier releases had different antivirus results, including heuristic flags on a historical native Go launcher. Those results concern different files and must not be represented as the scan result for V12.7.257.

## Verify your download on Windows

In PowerShell, open the directory containing the ZIP and run:

```powershell
Get-FileHash -Algorithm SHA256 ".\NMS_Personal_Atlas_V12.7.257.zip"
```

Compare the **entire** SHA-256 against:

```text
b830aab0da0a442431dfe15da1f2fc0a61f5384e5240878550f1a941d48bf180
```

If it differs, do **not** assume the published VirusTotal report applies. Confirm the download source and contact the maintainer if needed.

## What Atlas does

### Local execution

- The standard Windows distribution bundles a private Python runtime and its required libraries. Ordinary startup does not install Python globally or download Python dependencies.
- `LANCER_ATLAS.cmd` starts the local Streamlit application, normally on `127.0.0.1:8532`, and opens the local interface in a web browser.
- Administrator privileges are not normally required.
- The distributed configuration disables Streamlit's usage-statistics collection. This does not independently prove the absence of network traffic from Windows, the browser or other software components.

### Game saves and local data

- Atlas is designed to read supported original No Man's Sky save files without rewriting them.
- Atlas stores its own history, SQLite database, caches, configuration and diagnostic material in `NMSFaunaAnalyzerData` and, for selected workflows, other specifically designated local folders.
- Save-key decoding normally uses bundled NMS JSON mapping data; a missing or incomplete local mapping can fall back to the packaged resources.
- An optional **Repair JSON keys (Internet)** operation can retrieve a mapping resource when explicitly selected by the user.
- Exports, support links and other externally linked features may open file dialogs or web pages after a user action.

### Memory inspection and diagnostics

- Ordinary exploration is save-driven. The regular Alliance Observatory does not require automatic NMS process-memory inspection.
- Some optional technical or legacy equipment diagnostics may contain read-only process-memory inspection capabilities. They can use Windows APIs such as `OpenProcess`, `VirtualQueryEx` and `ReadProcessMemory` when explicitly invoked.
- Such capabilities can trigger security heuristics. The use of read-only APIs does not, by itself, constitute a comprehensive guarantee about every executable behavior.

## Optional local verification tools

Depending on the distributed package, the following tools may be available:

- `ATLAS_Save_Diagnostic.cmd` — reports save detection and JSON mapping health.
- `VERIFIER_INTEGRITE_ATLAS.cmd` — checks distributed files against a local integrity manifest.
- `AUDIT_SAVES_AVANT.cmd` and `AUDIT_SAVES_APRES.cmd` — compare save-file hashes around a user-performed test.
- `AUDIT_RESEAU_WINDOWS.cmd` — monitors selected network connections during a limited observation period.

These tools provide diagnostics, **not security certification**. Logs may expose file paths, profile identifiers or other personal metadata. Redact them before sharing.

## Known limits and claims we do not make

- A static source or package integrity check is **not** a comprehensive dynamic security audit of all bundled binaries, libraries and behavior on every Windows machine.
- Source and technical-review materials published on GitHub are not necessarily complete, reproducible source builds of the exact Windows ZIP.
- No VirusTotal result proves the absence of harmful behavior or ensures acceptance by all antivirus products.
- Save compatibility and file-write behavior have not been independently certified for every Windows and game configuration.
- Some game-derived assets and data require continued provenance and redistribution-rights review. Availability in the package does not grant downstream reuse rights.

## Reporting a vulnerability privately

**Security contact:** personnalatlas@gmail.com

Please report security or privacy concerns by email rather than publishing exploits, private player saves or sensitive logs in a public issue. When possible, include:

1. Exact Atlas release, ZIP filename and SHA-256.
2. Windows version and No Man's Sky platform/storefront.
3. Reproduction steps and the expected versus observed behavior.
4. Sanitized diagnostic details, without credentials, unrelated personal paths or private saves.
5. Antivirus vendor, detection name and report link, if applicable.

Do not disable security protections or run destructive tests against other users. For non-sensitive bugs and suggestions, the project's public GitHub issue tracker can be used.

---

NMS Personal Atlas is independent and is **not affiliated with or endorsed by Hello Games**.


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

---

## Security information — V12.7.250

### Local operation

NMS Personal Atlas is designed to run locally on Windows 10/11 x64.

- Python and its required libraries are bundled with the distribution.
- No automatic dependency installation is required at startup.
- The application serves its interface locally at `127.0.0.1:8532`.
- Administrator privileges are not normally required.
- Original No Man's Sky save files are intended to remain read-only.

### Memory access

Certain optional technical diagnostic tools can read the No Man's Sky process memory when explicitly requested.

The regular Alliance Observatory no longer requires process memory inspection. Alliance names can be manually recorded and associated with their native identifiers.

### Network access

The application is designed for offline operation.

Some optional functions may contact external services when explicitly requested, such as repairing save-file key mappings. External links may also open in the user's browser.

Network behavior and privacy verification remain under review.

### Antivirus verification

**Version:** V12.7.250  
**Scan date:** November 21, 2026  
**VirusTotal result:** 0/51 detections

[View the VirusTotal report](https://www.virustotal.com/gui/file/3065c250089ebd5d98f865fa9784548217c563c6ebee193fa224a2b23b3a5927/detection)

**Archive SHA-256:**

`3065c250089ebd5d98f865fa9784548217c563c6ebee193fa224a2b23b3a5927`

This result applies specifically to the V12.7.250 archive. It does not constitute a security certification or guarantee that future versions will receive the same result.

### Reporting security issues


**NMS Personal Atlas — V12.7.255**  
Last updated: **November 22, 2026**

NMS Personal Atlas is a free, independent Windows companion application for No Man's Sky. This document explains the security posture of the distributed application, what is and is not verified, and how to report a concern responsibly.

## Supported releases

| Version | Security and compatibility attention |
| --- | --- |
| **V12.7.255** | Current release for new bug and vulnerability reports. |
| V12.7.254 and earlier | Historical versions. Update to the current release where practical; fixes may not be backported. |

Security updates are best-effort for this independent community project. No formal response-time or support-service guarantee is offered.

## Antivirus report — exact release file

The **author-repacked V12.7.255 ZIP archive** was submitted to VirusTotal on **November 22, 2026**. It showed **0 / 53 detections** among participating antivirus engines at the time of the scan.

- [VirusTotal detection report](https://www.virustotal.com/gui/file/c87435acb6f2ff7c7a235f700aeec7ebac72c66d5c5f21ab8d9b8591914b883f/detection)
- **Reported filename:** `NMS_Personal_Atlas_V12.7.255.zip`
- **Reported SHA-256:**

```text
c87435acb6f2ff7c7a235f700aeec7ebac72c66d5c5f21ab8d9b8591914b883f
```

This hash refers **only to the ZIP that was scanned**. Repacking, modifying or rebuilding a ZIP results in a different hash. The exact GitHub/Nexus Mods download must match before this scan is cited as evidence for that download.

**A 0/53 result is not a guarantee of safety, a malware-free certification or an endorsement by the scanning vendors.** Heuristic antivirus and SmartScreen warnings may still occur, particularly for independent/unsigned software.

### Verify a download on Windows

Open PowerShell in the folder containing the downloaded ZIP:

```powershell
Get-FileHash -Algorithm SHA256 ".\NMS_Personal_Atlas_V12.7.255.zip"
```

Compare the full output with the SHA-256 above. If it differs, do **not** assume the VirusTotal report applies. Recheck the download source and release notes.

Do **not** disable Microsoft Defender, Avast, SmartScreen or other protections to run Atlas. Report a quarantine or detection to the developer, including the filename, SHA-256 and vendor's detection name.

## What the application does

**Local execution**

- The standard Windows distribution includes an embedded Python 3.13 runtime and required third-party packages. Ordinary startup does not install Python globally or fetch Python dependencies.
- The `LANCER_ATLAS.cmd` launcher starts a local Streamlit application on **`127.0.0.1:8532`** and opens a browser page on the user's PC.
- Streamlit usage-statistics collection is disabled in the distributed configuration. This does **not** prove the absence of all possible network traffic by browsers, Windows or third-party components.

**File and save access**

- Atlas is designed to **read supported original NMS saves without rewriting them**.
- Atlas stores its own database, history, caches and diagnostic outputs in **`NMSFaunaAnalyzerData`** and, for selected technical workflows, other explicitly named local directories.
- The V12.7.255 save-key decoder uses a bundled NMS JSON mapping if the local map is missing or incomplete. It does not need an automatic online mapping download during ordinary startup and does not silently overwrite an existing custom map.
- The optional **Repair JSON keys (Internet)** operation explicitly requests a resource from GitHub when the user chooses it.
- Export actions and selected external links can open local save dialogs or web pages following user action.

**Memory and diagnostic capabilities**

- Ordinary Atlas exploration is save-driven. The Alliance Observatory no longer uses automatic NMS process-memory inspection.
- **Some optional, legacy Equipment technical diagnostics still contain Windows APIs for read-only inspection of `NMS.exe` process memory** (including `OpenProcess`, `VirtualQueryEx` and `ReadProcessMemory`) when invoked by the user. Such behavior can trigger security heuristics.
- The presence of read-only memory APIs does not establish that every invocation succeeds or that every executable behavior has been independently verified. Do not describe Atlas categorically as incapable of accessing process memory.

## Optional verification tools

The standalone distribution includes tools that users can run locally:

- **`ATLAS_Save_Diagnostic.cmd`** — reports detected save files and the health of local/bundled/effective JSON mappings; profile identities are anonymized in its intended output.
- **`VERIFIER_INTEGRITE_ATLAS.cmd`** — checks distributed files against the supplied integrity manifest.
- **`AUDIT_SAVES_AVANT.cmd` / `AUDIT_SAVES_APRES.cmd`** — compare save-file hashes before and after a user-run test. They do not attribute a change to a particular process.
- **`AUDIT_RESEAU_WINDOWS.cmd`** — observes selected Atlas Python TCP connections for a limited period. It is **not** a complete audit of browser traffic, DNS, UDP or very short-lived connections.

These tools produce **local reports**, not security certifications. Logs can still contain user paths, local network details or other sensitive metadata. Review and redact before sharing.

## Known limits and claims we do not make

- A static source inspection is **not** a comprehensive Windows behavior audit of every bundled DLL and third-party component.
- A clean VirusTotal report does **not** prove that a program is harmless or that every antivirus product will accept it.
- The public source and review packages are **not necessarily complete, byte-for-byte reproducible builds** of the distributed Windows ZIP. The exact scope must be checked for each version.
- Save compatibility, calculation accuracy and protection from every possible file-write or network defect have not been universally certified across all Windows and NMS configurations.
- Some game-derived data, the bundled JSON-key map and some visual assets require further **provenance and redistribution-rights review**. Publication does not automatically grant downstream reuse rights.

See the release's bundled `AUDIT_SECURITE_LIRE_EN_PREMIER.md`, `THIRD_PARTY_NOTICES.txt`, provenance reports and licence documents for further context.

## Reporting a vulnerability privately

**Security contact: `personnalatlas@gmail.com`**

Please use email for suspected vulnerabilities or privacy issues rather than posting exploits, player saves or sensitive logs in a public GitHub issue. Include, where possible:

1. Atlas version and the exact release ZIP filename and SHA-256.
2. Windows version and relevant NMS storefront (Steam, GOG or Game Pass).
3. Clear reproduction steps, expected/observed behavior and any security impact.
4. Sanitized log excerpts or screenshots; avoid sharing credentials, personal paths, game saves or unrelated personal information unless a secure transfer has been agreed.
5. Antivirus product, detection name and report link, if relevant.

Do not run destructive tests against other users or share private data without authorization. The developer can investigate and publish a fix or advisory where warranted; there is no guaranteed response time.

For ordinary feature requests and non-sensitive bugs, use the project's public issue/discussion channels instead.

---

NMS Personal Atlas is an independent, community-developed companion application and is **not affiliated with, supported by or endorsed by Hello Games**.

Security concerns and potential vulnerabilities are taken seriously.

Please avoid publicly disclosing sensitive exploit details or personal save data. Contact the project maintainer through an appropriate private channel before publishing technical vulnerability details.
