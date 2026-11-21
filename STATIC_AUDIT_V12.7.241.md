# V12.7.241 — limited static inspection

**Examined archive:** `NMS_Personal_Atlas_V12.7.241_STABILITE_COSMOS_CATALOGUE_COMPLET_20261120.zip`

**Archive SHA-256:** `05ddf50f9f0fe0fd341f8db0effa6c93508ba1394082561dff9f8e9c4acc3afb`

**Windows launcher:** `NMS Personal Atlas.exe`

**Launcher SHA-256:** `6393663ef96910565d80ce8fc78190e1c2805b4cc16fa8e707ae6de663e76595`

**Launcher size:** 5,600,768 bytes

**Build metadata observed:** Go 1.23.2; module `atlasnative/launcher`; Windows amd64 GUI executable; `CGO_ENABLED=0`; `-trimpath=true`.

**Important:** The executable SHA-256 matches the one in the supplied VirusTotal screenshots. That match does **not** establish that the complete Nexus Mods ZIP is identical to this V12.7.241 archive.

## Findings in the distributed sources
- `runtime/launcher/demarrer_atlas.py` launches local installation and Streamlit processes using Windows `CREATE_NO_WINDOW`.
- On first run, the launcher may fetch `get-pip.py` from `bootstrap.pypa.io` and install packages using pip with `https://pypi.org/simple`.
- The Python runtime is bundled; dependencies are installed into `runtime/site-packages`.
- `atlas_core/engines/live_memory.py` includes read-only `OpenProcess`, `VirtualQueryEx`, and `ReadProcessMemory` Win32 bindings.
- `ui/pages/equipment_page.py` includes explicitly triggered technical RAM capture tools.
- `ui/pages/alliance_observatory_page.py` automatically starts a bounded NMS RAM alliance probe when that page is opened and relevant station identifiers are available.
- In the examined Python sources, no literal use of `WriteProcessMemory` or `CreateRemoteThread` was identified. This is a **limited static observation**, not proof that the full application cannot perform an injection.
- `OUTILS/uninstaller_source.c` includes confirmation dialogs and schedules self-removal after uninstall through `MoveFileExW(..., MOVEFILE_DELAY_UNTIL_REBOOT)`.
- No original Go `main_windows.go` source is present in this distribution archive. The Go executable has **not** been rebuilt from public source in this audit.

## What this audit does NOT prove
This inspection did not execute the Windows binaries, test the full first-run network installation, independently reproduce the Go build, or trace the specific VirusTotal MITRE behavioral events. It therefore **does not certify the binaries as malware-free**, and does not establish that all reported detections are false positives.

## Suggested next security work
1. Locate or reconstruct and publish the authentic Go launcher source, with a documented build recipe.
2. Add deterministic dependency versions and integrity checks for the first-run dependency installation.
3. Make automatic NMS RAM inspection explicit to users, preferably opt-in, and remove unused historical capture modules from the player edition.
4. Review behavioral sandbox events and submit carefully documented false-positive reports only where evidence supports them.
