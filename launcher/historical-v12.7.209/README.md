# NMS Personal Atlas — Historical Go Launcher

## Overview

This directory contains the original Go source code for the native Windows launcher of **NMS Personal Atlas**, first introduced in version 12.7.209 (November 2026).

The launcher was developed to provide a simple, user-friendly way to start Atlas without requiring users to manually install Python or execute command-line scripts.

## How It Works

The launcher provides:

- A native Windows graphical interface.
- A startup progress indicator and installation status messages.
- Automatic initialization of Atlas's private Python environment.
- Execution of the Python startup script without displaying console windows.
- Automatic opening of the local Streamlit interface.
- Error reporting and installation diagnostics.

## Included Source Files

| File | Description |
|---|---|
| `main_windows.go` | Native Windows GUI and application startup logic |
| `status.go` | Startup status and progress message handling |
| `status_test.go` | Go unit tests |
| `go.mod` | Go module configuration |

## Technical Information

- **Language:** Go
- **Target platform:** Windows x64
- **Original toolchain:** Go 1.23.2
- **Application interface:** Native Windows GUI
- **Python environment:** Bundled privately with NMS Personal Atlas
- **Web interface:** Streamlit, running locally

## Source Code Transparency

These files represent the historical implementation of the launcher.

The original V12.7.209 executable was successfully reproduced byte-for-byte from the recovered historical source code.

However, subsequent updates introduced additional launcher improvements and fixes.

**Important:** These historical sources should not be considered an exact source-code representation of the launcher distributed with NMS Personal Atlas V12.7.241.

The goal of publishing this code is to provide transparency about the architecture and development of the application.

## Security Notice

NMS Personal Atlas is independent, fan-made software and is distributed without a commercial code-signing certificate.

Some antivirus products may flag unfamiliar or unsigned executables. Such detections must be investigated and should not automatically be dismissed as false positives.

Publishing this source code does not constitute a security certification of the current application.

For more information, see the repository's main `SECURITY.md` and static audit documentation.

---

**NMS Personal Atlas** — A free companion application created by a passionate No Man's Sky player, with extensive AI-assisted development.
