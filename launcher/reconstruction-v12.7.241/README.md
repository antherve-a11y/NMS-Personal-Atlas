# NMS Personal Atlas — Go Launcher Source Reconstruction (V12.7.241)

**Status: SOURCE RECONSTRUCTION — NOT YET A BYTE-FOR-BYTE REPRODUCIBLE BUILD.**

This folder contains a reconstructed, compilable Go implementation of the Windows native GUI launcher shipped with NMS Personal Atlas V12.7.241. It is intended for transparency and code review; **the official Windows EXE has not been replaced**.

## Where the source came from

1. The original Go launcher source was recovered from the historical V12.7.209 (November 13, 2026) distribution. That historical source rebuilds its original executable byte-for-byte.
2. The later V12.7.209 maintenance build changed the launcher to fix real Windows installation and GUI responsiveness problems. Its exact revised Go source file was not found in the archived player packages.
3. Static PE comparison confirms that the executable code and data in V12.7.241 are **byte-identical to the corrected V12.7.209 binary**. The newer EXE added an icon resource, not new launcher code.
4. This reconstructed source applies the observed launcher changes to the recovered historical Go source. It compiles on Go 1.23.2 and its function layout matches the corrected Go binary, but its output EXE SHA-256 differs from the original. Therefore it is NOT an exact-source claim.

See `BINARY_AUDIT.md` for evidence and limitations.

## What this source shows

- Native Windows 64-bit GUI and status display (no terminal by default)
- Launch of the bundled private Python interpreter
- Execution of the actual `runtime/launcher/demarrer_atlas.py` script
- Start-up logs and visibility-gated update of detail text
- Bounded log reads, progress state handling and a nonblocking local server health check
- Local Streamlit server at `127.0.0.1:8532`
- No administrator elevation requested by the launcher

The **Python startup script** and the **application's own runtime behavior** require their own separate audit. The Go source alone cannot establish absence of malware.

## Files

- `main_windows.go` — reconstructed Windows launcher logic
- `status.go` and `status_test.go` — recovered historical status parsing and tests
- `go.mod` — module definition
- `BUILD.md` — compilation and testing instructions
- `BINARY_AUDIT.md` — evidence and honest limitations
- `compare_pe_sections.py` and `binary_section_evidence.json` — independent PE section comparison tool and reference comparison output
- `function_compare.json` — function-level comparisons of the reconstructed executable against the corrected historical binary
- `SHA256_SOURCE_FILES.txt` — integrity hashes for published source files

## Safety disclaimer

Do **not** replace the Atlas release EXE with a locally rebuilt version without testing. Do **not** disable antivirus protections. The current VirusTotal detections must still be investigated. This publication is a transparency measure, **not a clean bill of health**.

NMS Personal Atlas is an independent fan project and is not affiliated with Hello Games.

