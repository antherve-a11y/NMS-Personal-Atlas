# Build and test the reconstructed Go launcher

The historical and reconstructed launchers were built with **Go 1.23.2** for **Windows amd64**. There are no third-party Go module dependencies in this folder.

Open PowerShell inside the folder containing `go.mod`, then run:

```powershell
# Requires Go 1.23.2 on PATH
go version
go test ./...
go vet ./...
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '0'
go build -trimpath -ldflags='-H=windowsgui -s -w' -o 'AtlasLauncher_Reconstruction.exe' .
Get-FileHash .\AtlasLauncher_Reconstruction.exe -Algorithm SHA256
```

For a cross-build from Linux:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags='-H=windowsgui -s -w' -o AtlasLauncher_Reconstruction.exe .
```

## Important distinction

The output is a *reconstruction*. It is not a byte-identical rebuild of the V12.7.241 release EXE, and it does not include its separately inserted Windows icon resources. Do not present the output hash as the hash of the official software.

No release executable is included in this source package. Windows execution tests and third-party security reviews remain necessary.
