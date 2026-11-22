# 🚀 NMS Personal Atlas

### Explore • Scan • Identify

**Your personal universe, organized.**

NMS Personal Atlas is a free, independent Windows companion application for **No Man's Sky**.

It transforms your game saves into a comprehensive exploration database, combining discovery tracking, procedural research, starship analysis, equipment optimization, industrial management and galactic trading intelligence.

**Read-only by design — your original No Man's Sky saves are never modified.**

---

## 🌌 Main Features

- **Galactic Atlas:** Track star systems, planets, discoveries, biomes, resources and exploration progress.
- **Interactive Galaxy Map:** Visualize visited systems, bases, stations and exploration routes.
- **Planetary Analysis:** Explore detailed planetary records and interactive 3D holograms.
- **Fauna & Species:** Catalogue discoveries, analyze statistics and preserve historical records.
- **Companions & Genetics:** Compare companions and plan egg-sequencing experiments.
- **Fleet Intelligence:** Analyze starships, frigates, squadrons and Corvette configurations.
- **Equipment Optimizer:** Explore improved technology layouts, adjacency bonuses and supercharged slots.
- **Infrastructure:** Monitor bases, extractors, silos and resource production networks.
- **Galactic Commerce:** Calculate supported market-price scenarios and compare trading opportunities.
- **Explorer Progression:** Track achievements, records and personal exploration milestones.

Some calculations are estimates based on reconstructed game mechanics. Atlas distinguishes observed data from calculated, estimated and illustrative results.

---

## 💻 Installation

**Platform:** Windows 64-bit

1. Download a complete release archive when available.
2. Extract the ZIP into a writable folder.
3. Launch `NMS Personal Atlas.exe`.
4. Allow the first-run installer to prepare its private Python dependencies.
5. Open the local Atlas interface in your browser.

Atlas normally uses `http://127.0.0.1:8532`.

An embedded Python runtime is included. The initial setup may download dependencies from official Python package repositories.

**No mod manager is required.**

### Security Notice

The Windows executable is currently unsigned. Microsoft Defender SmartScreen or third-party antivirus software may display warnings.

Always verify downloaded files and keep your antivirus protections enabled.

---

## 🔬 Technology

NMS Personal Atlas uses:

- Python 3.13
- Streamlit
- SQLite
- NumPy and pandas
- Procedural reconstruction algorithms
- No Man's Sky save-data analysis
- Native game-reference research

The project was developed under human direction with extensive AI assistance in coding, debugging, research and documentation.

---

## 📖 Documentation

A comprehensive illustrated English user guide has been prepared, covering the main application modules and their relationship to No Man's Sky.

Documentation and download links will be published as the project repository is populated.

---

## 🛡️ Project Status

NMS Personal Atlas is under active development.

The repository is being prepared for source-code publication, technical review and Windows release distribution.

Some advanced native reconstruction features remain experimental.

---

## 🤝 Feedback

Bug reports, suggestions and constructive feedback are welcome.

**Contact:** personnalatlas@gmail.com

---

## Disclaimer

NMS Personal Atlas is an independent fan-made application.

No Man's Sky is developed by Hello Games. This project is not affiliated with or endorsed by Hello Games.

Third-party components retain their respective licences. Redistribution rights for included game-derived assets must be verified before public release.

---

## Latest release — V12.7.250

**NMS Personal Atlas — Standalone Windows Edition**

NMS Personal Atlas is a free, independent companion application for No Man's Sky.

### Security verification

**VirusTotal result: 0/51 detections**

Scanned on November 21, 2026.

[View the public VirusTotal report](https://www.virustotal.com/gui/file/3065c250089ebd5d98f865fa9784548217c563c6ebee193fa224a2b23b3a5927/detection)

**SHA-256:**

`3065c250089ebd5d98f865fa9784548217c563c6ebee193fa224a2b23b3a5927`

### Standalone installation

- Bundled Python runtime and dependencies.
- No automatic dependency downloads at startup.
- Local browser-based interface.
- Windows 10/11 x64.
- No separate Python installation required.
- No administrator privileges normally required.

Some optional technical diagnostics can read No Man's Sky process memory when explicitly requested. Optional online features and external links are documented separately.

**Latest documented Windows release: V12.7.255 — November 22, 2026**

NMS Personal Atlas is a **free, independent companion application for [No Man's Sky](https://www.nomanssky.com/)**. It turns supported local game-save records into a personal exploration atlas: known systems and planets, discoveries, ships and equipment, bases, resources, trading information and progression — all accessible through a local browser interface.

It is **not** a mod that runs inside the game, and it is **not affiliated with or endorsed by Hello Games**.

[**Download the latest release**](https://github.com/antherve-a11y/NMS-Personal-Atlas/releases) · [Security & antivirus information](SECURITY.md)

## Features

- **Exploration and discoveries:** catalogues of known galaxies, systems, planets, fauna, flora and minerals; completion indicators and historical records.
- **Galaxy and planetary maps:** interactive system maps and illustrative planetary 3D visualizations.
- **Fleet and starships:** ship and corvette inventories, reconstructed statistics, equipment layouts and comparison tools.
- **Equipment:** multi-tools, weapons, exosuit, exocraft, technology and optimization views.
- **Bases and infrastructure:** resource extraction, storage and production network analysis.
- **Galactic trade:** station and market information, estimated prices, comparisons and route-planning tools.
- **Companions:** companion profiles, traits and genetics tools.
- **Explorer dashboard:** discovery progress, personal history and mission-planning views.
- **Languages:** French and English interface, including the personalized startup screen.

**How to interpret results:** Atlas distinguishes data read from a save from reconstructed, calculated, estimated and illustrative information. A score, predicted terminal price or stylized planet model should not automatically be treated as an exact value displayed by the game. Information depends on what is present in the selected save and Atlas's retained history.

## What's new in V12.7.255

This update incorporates fixes from V12.7.251 through V12.7.255, including changes prompted by community diagnostics.

- **Portable save-key mapping:** fixed a directory-name dependency that could prevent the bundled NMS JSON key map from being loaded after the Atlas folder was renamed.
- **Offline fallback:** when the local key map is absent or incomplete, Atlas can use its bundled, validated mapping without silently overwriting the user's existing file.
- **Better diagnostics:** the save diagnostic reports the status of the local, bundled and effective mapping sources.
- **Startup visibility:** restored the personalized **“Hello, Traveller” / “Bonjour, Voyageur”** screen with real initialization stages, including save processing and catalogue preparation.
- **Save processing:** improved reporting of prolonged synchronization and detection of stale results when the save slot changes.
- **Command Center performance:** reduced repeated calculations and reused cached results where applicable.
- **Planet and creature pages:** improved handling of historical discovery counts, incomplete species profiles and mixed-format table data.
- **Interface compatibility:** replaced a deprecated Streamlit iframe call for the planetary viewer.

These fixes improve diagnosis and compatibility; they do **not** guarantee that every save format, PC configuration or long-running analysis is free of issues.

## Installation — Windows 10/11 (64-bit)

1. Download the **complete V12.7.255 ZIP** from [GitHub Releases](https://github.com/antherve-a11y/NMS-Personal-Atlas/releases).
2. Extract **all** files into a writable folder outside the No Man's Sky game directory (for example, `D:\Games\NMS Personal Atlas`).
3. Open the extracted **`NMS Personal Atlas`** folder and double-click **`LANCER_ATLAS.cmd`**.
4. Keep the console window open while using Atlas. The browser should open automatically at **http://127.0.0.1:8532**.
5. Select or verify the NMS profile and save slot. After saving new progress in the game, use **Reload Save** in Atlas.

The distribution bundles its **own Windows Python runtime and required dependencies**. You should not need to install Python, run `pip`, or download libraries during ordinary startup. A browser is required; a system-wide Python installation is not.

**Updating:** Back up the `NMSFaunaAnalyzerData` folder from your current Atlas installation before replacing the application or moving to a new directory. This folder contains your local history and configuration. **Do not delete or overwrite it** if you want to keep your progress.

**Uninstalling:** Close Atlas and its console, then remove the application folder. Keep or delete `NMSFaunaAnalyzerData` deliberately depending on whether you wish to retain your data. Atlas does not install a Windows service or a system-wide Python environment.

## Troubleshooting

| Symptom | First thing to check |
| --- | --- |
| No save or wrong profile | Check the selected NMS profile, game slot and `.hg` save file; then use **Reload Save**. |
| Save synchronization takes unusually long | Note the last stage shown on the welcome screen. Run **`ATLAS_Save_Diagnostic.cmd`** from the Atlas folder. |
| JSON key mapping warning | In V12.7.255, the bundled mapping should be available even when the folder has been renamed. Check the diagnostic's **effective mapping source** and report the result. |
| Blank or partially rendered page | Save the console traceback and the Atlas version; report the affected page, selected slot and reproduction steps. |
| Windows security / antivirus alert | Verify the download origin and file SHA-256. **Do not disable your security software.** See [SECURITY.md](SECURITY.md). |

The diagnostic writes `ATLAS_SAVE_DIAGNOSTIC.txt` locally. Although it anonymizes save-profile identities, **review any logs or screenshots before sharing them publicly**; paths and other metadata may still be sensitive.

## Security, privacy and the V12.7.255 scan

The **author-submitted, repackaged release archive** was scanned on VirusTotal on **November 22, 2026**: **0 detections from 53 participating antivirus engines**.

[View the VirusTotal report](https://www.virustotal.com/gui/file/c87435acb6f2ff7c7a235f700aeec7ebac72c66d5c5f21ab8d9b8591914b883f/detection)

**SHA-256 of the exact scanned ZIP:**

```text
c87435acb6f2ff7c7a235f700aeec7ebac72c66d5c5f21ab8d9b8591914b883f
```

Check that the **file you download** has the same SHA-256 before relying on this result. Recompressing a ZIP changes its hash. An antivirus scan with no detections is **not a security certification** and does not guarantee that individual binaries are harmless.

Atlas is designed to read original NMS save files without modifying them; it creates its own local data, caches and logs. Its ordinary startup uses a local Streamlit server and does **not** automatically download Python packages or save-key mappings. An **optional, user-initiated** “Repair JSON keys (Internet)” action and external links may access the Internet.

**Important transparency detail:** optional legacy technical equipment diagnostics include code capable of **read-only inspection of the running NMS process memory** when explicitly requested. Routine exploration and the Alliance Observatory are not intended to run an automatic process-memory probe. See [SECURITY.md](SECURITY.md) for the precise scope and audit limitations.

## Source code, licensing and credits

This repository is intended to support code review and transparency. **A published source snapshot does not automatically imply that the full standalone Windows package is reproducibly buildable from the repository.** Verify the version and scope of the available source package before drawing that conclusion.

Atlas is **freeware**, not automatically licensed as an open-source project. The creator's original contributions and third-party components have different rights and terms. Consult the bundled `LICENCE_PERSONAL_ATLAS_FR_EN.txt`, `CREDITS_ET_LICENCES_FR_EN.txt` and `THIRD_PARTY_NOTICES.txt`.

The provenance and redistribution rights of some bundled **save-mapping data and visual assets are still under review**. Do not assume all files, images or data are freely reusable simply because they are included with the application. The developer must resolve those questions before treating public redistribution rights as confirmed.

## Feedback and development

Bug reports and constructive suggestions are welcome. For an ordinary bug, include your Atlas version, Windows version, game storefront (Steam/GOG/Game Pass), expected vs. actual behavior, and a **redacted** diagnostic if relevant. Do not upload complete game saves or personal logs to a public issue without reviewing their contents.

For a suspected **security vulnerability**, follow [SECURITY.md](SECURITY.md) and report it privately instead of opening a public issue.

NMS Personal Atlas has been created with substantial AI assistance under the direction of its author, with research, reverse-engineering, design choices and testing by the creator and community.

**Independent fan project — not affiliated with or endorsed by Hello Games.**


The application is designed to read original No Man's Sky save files without modifying them.

**Transparency:** A clean VirusTotal scan is not a guarantee of complete software security. Additional runtime and privacy verification is ongoing.
