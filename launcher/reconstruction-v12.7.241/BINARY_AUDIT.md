# V12.7.241 native Go launcher — source provenance and binary comparison

## Evidence reviewed

| Artifact | SHA-256 |
|---|---|
| Corrected launcher V12.7.209 (Nov. 13, 2026) | `0b28c10b4b33ec455c5eec07f712a5b8d3c5940017ca79642d46208f6b5b8915` |
| Official launcher V12.7.241 (Nov. 20, 2026) | `6393663ef96910565d80ce8fc78190e1c2805b4cc16fa8e707ae6de663e76595` |
| Reconstructed launcher compiled for comparison (not distributed) | `f2e408b1dd09759f9e12ca85073f3c2a3ec98bd873c5a73bfeb95c1a3e8c3240` |

The official V12.7.241 SHA-256 also matches the executable fingerprint in the provided VirusTotal screenshots. This identifies the file, not its safety.

## Exact binary comparison: corrected V12.7.209 vs. V12.7.241

The following PE sections are byte-for-byte **identical**: `.text`, `.rdata`, `.data`, `.pdata`, `.xdata`, `.idata`, `.reloc` and `.symtab`.

V12.7.241 additionally contains a `.rsrc` section holding Windows icon resources. This confirms the current executable still uses the **same compiled Go logic as the corrected 13-Nov-2026 launcher**.

A machine-readable record of section sizes and SHA-256 values is available in `binary_section_evidence.json`; use `compare_pe_sections.py` to repeat the comparison on files you legally possess.

## Source reconstruction work performed

- Recovered `main_windows.go`, `status.go`, `status_test.go` and `go.mod` from the original V12.7.209 distribution archive.
- The older source rebuilt its original 13-Nov-2026 launcher exactly; its hash does **not** match the corrected 13-Nov binary.
- Reconstructed later behavior: Win32 progress timer at 450 ms, log reading limited to 48 KiB per update, detail text tail kept to 8000 bytes, throttled visible-detail refresh, asynchronous server readiness polling, and the corrected `runtime/launcher/demarrer_atlas.py` path.
- Recompiled with Go 1.23.2 (`GOOS=windows`, `GOARCH=amd64`, `CGO_ENABLED=0`, `-trimpath`, `-H=windowsgui -s -w`).
- **Two Go status tests passed**; `go vet` for Windows passed; Windows test compilation succeeded.
- Reconstructed code has the same `main.*` function address layout and sizes as the corrected binary, as verified from Go's internal function metadata.
- The compiled `main.wndProc` function (4,416 bytes), which controls the Windows GUI, is **byte-identical** to the corrected original.
- The compiled `main.updateFromLog` (1,920 bytes) and `main.launchAtlas` (3,552 bytes) functions each differ at only **three individual bytes**, all associated with code-relative references; remaining differences include the background polling routine and build metadata. See `function_compare.json`.

## Remaining discrepancy

The rebuilt reconstruction has SHA-256 `f2e408b1dd09759f9e12ca85073f3c2a3ec98bd873c5a73bfeb95c1a3e8c3240`, which is **not** `0b28c10b4b33ec455c5eec07f712a5b8d3c5940017ca79642d46208f6b5b8915`. Exact source-level equivalence and byte-for-byte reproducibility are **not confirmed**. Differences remain in compiled instructions and build metadata, even though the reconstructed source is functionally very close.

There has been **no full dynamic test of this reconstructed launcher on Windows**, no analysis of all dependency download behavior, and no reproduction of the antivirus sandbox events. Do not conclude that VirusTotal findings are false positives based on this audit alone.

## Reproducible section check (no execution)

```bash
python compare_pe_sections.py "Corrected-209.exe" "NMS Personal Atlas.exe"
```

The historical corrected launcher must be extracted from the archived V12.7.209 Windows installation patch, and the official launcher must be extracted from the actual V12.7.241 release ZIP. Both filenames and hashes should be verified before comparison.
