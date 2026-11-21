//go:build windows

// NMS Personal Atlas — minimal native Windows GUI launcher.
// Runs the bundled Python privately, with no terminal and no admin escalation.
package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	WM_CREATE          = 0x0001
	WM_DESTROY         = 0x0002
	WM_CLOSE           = 0x0010
	WM_COMMAND         = 0x0111
	WM_TIMER           = 0x0113
	WM_CTLCOLORSTATIC  = 0x0138
	WM_CTLCOLOREDIT    = 0x0133
	WM_SETFONT         = 0x0030
	WS_CHILD           = 0x40000000
	WS_VISIBLE         = 0x10000000
	WS_CAPTION         = 0x00C00000
	WS_SYSMENU         = 0x00080000
	WS_MINIMIZEBOX     = 0x00020000
	WS_TABSTOP         = 0x00010000
	WS_BORDER          = 0x00800000
	WS_VSCROLL         = 0x00200000
	WS_EX_CLIENTEDGE   = 0x00000200
	ES_MULTILINE       = 0x0004
	ES_READONLY        = 0x0800
	ES_AUTOVSCROLL     = 0x0040
	SS_LEFT            = 0
	SS_CENTER          = 1
	SW_SHOW            = 5
	SW_HIDE            = 0
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	MB_OK              = 0
	MB_ICONWARNING     = 0x30
	MB_ICONQUESTION    = 0x20
	MB_YESNO           = 4
	IDYES              = 6
	PBS_MARQUEE        = 0x0008
	PBM_SETMARQUEE     = 0x040A
	PBM_SETPOS         = 0x0402
	WM_USER            = 0x0400
	TRANSPARENT        = 1
	ICC_PROGRESS_CLASS = 0x20
	IDC_DETAIL         = 401
	IDC_LOG            = 402
	IDC_CLOSE          = 403
	COLOR_BG           = 0x00211713 // Win32 COLORREF = 0x00BBGGRR
	COLOR_FG           = 0x00F4F1ED
	COLOR_MUTED        = 0x00BFB4AB
	COLOR_EDIT         = 0x002C211D
)

type wndClassEx struct {
	Size       uint32
	Style      uint32
	Proc       uintptr
	ExtraCls   int32
	ExtraWnd   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}
type point struct{ X, Y int32 }
type msgStruct struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
	Private uint32
}
type initCommonControlsEx struct {
	Size  uint32
	Flags uint32
}

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	gdi32             = syscall.NewLazyDLL("gdi32.dll")
	comctl32          = syscall.NewLazyDLL("comctl32.dll")
	shell32           = syscall.NewLazyDLL("shell32.dll")
	pRegisterClass    = user32.NewProc("RegisterClassExW")
	pCreateWindow     = user32.NewProc("CreateWindowExW")
	pDefWindowProc    = user32.NewProc("DefWindowProcW")
	pShowWindow       = user32.NewProc("ShowWindow")
	pUpdateWindow     = user32.NewProc("UpdateWindow")
	pGetMessage       = user32.NewProc("GetMessageW")
	pTranslateMessage = user32.NewProc("TranslateMessage")
	pDispatchMessage  = user32.NewProc("DispatchMessageW")
	pPostQuit         = user32.NewProc("PostQuitMessage")
	pSetTimer         = user32.NewProc("SetTimer")
	pKillTimer        = user32.NewProc("KillTimer")
	pSetWindowText    = user32.NewProc("SetWindowTextW")
	pSendMessage      = user32.NewProc("SendMessageW")
	pDestroyWindow    = user32.NewProc("DestroyWindow")
	pMessageBox       = user32.NewProc("MessageBoxW")
	pSetWindowPos     = user32.NewProc("SetWindowPos")
	pGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	pSetTextColor     = gdi32.NewProc("SetTextColor")
	pSetBkColor       = gdi32.NewProc("SetBkColor")
	pSetBkMode        = gdi32.NewProc("SetBkMode")
	pCreateBrush      = gdi32.NewProc("CreateSolidBrush")
	pCreateFont       = gdi32.NewProc("CreateFontW")
	pDeleteObject     = gdi32.NewProc("DeleteObject")
	pGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	pLoadCursor       = user32.NewProc("LoadCursorW")
	pLoadImage        = user32.NewProc("LoadImageW")
	pInitCommon       = comctl32.NewProc("InitCommonControlsEx")
	pShellExecute     = shell32.NewProc("ShellExecuteW")
	hwndMain          uintptr
	hwndTitle         uintptr
	hwndStatus        uintptr
	hwndStep          uintptr
	hwndHint          uintptr
	hwndProgress      uintptr
	hwndDetails       uintptr
	hwndToggle        uintptr
	hwndOpenLog       uintptr
	hwndClose         uintptr
	brushBg           uintptr
	brushEdit         uintptr
	fontBody          uintptr
	fontTitle         uintptr
	fontMono          uintptr
	expanded          bool
	terminalError     string
	started           bool
	runnerDone        bool
	runnerErr         error
	dataFolder        string
	appRoot           string
	logPath           string
	logReadOffset     int64
	logRemainder      string
	detailsText       string
	currentView       phaseState
	oldStatus         string
	oldDetails        string
	activeCmd         *exec.Cmd
	runnerMutex       sync.Mutex
	classString       *uint16
)

func wide(s string) *uint16   { return syscall.StringToUTF16Ptr(s) }
func strptr(s string) uintptr { return uintptr(unsafe.Pointer(wide(s))) }
func showText(h uintptr, s string) {
	if h != 0 {
		pSetWindowText.Call(h, strptr(s))
	}
}
func setFont(h, font uintptr) { pSendMessage.Call(h, WM_SETFONT, font, 1) }
func create(parent uintptr, class, text string, style uint32, ex uint32, x, y, w, h int32, id uintptr) uintptr {
	instance, _, _ := pGetModuleHandle.Call(0)
	handle, _, _ := pCreateWindow.Call(uintptr(ex), strptr(class), strptr(text), uintptr(style),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h), parent, id, instance, 0)
	return handle
}
func createFont(size int32, weight uintptr, face string) uintptr {
	h, _, _ := pCreateFont.Call(uintptr(size), 0, 0, 0, weight, 0, 0, 0, 1, 0, 0, 0, 0, strptr(face))
	return h
}
func setWindowHeight(n int32) {
	pSetWindowPos.Call(hwndMain, 0, 0, 0, 690, uintptr(n), SWP_NOMOVE|SWP_NOZORDER)
}
func shellOpenFile(path string) {
	pShellExecute.Call(hwndMain, strptr("open"), strptr(path), 0, 0, SW_SHOW)
}
func showMessage(title, text string, flags uintptr) uintptr {
	r, _, _ := pMessageBox.Call(hwndMain, strptr(text), strptr(title), flags)
	return r
}

func wndProc(hwnd uintptr, message uint32, wp, lp uintptr) uintptr {
	switch message {
	case WM_CREATE:
		hwndMain = hwnd
		fontBody = createFont(-18, 400, "Segoe UI")
		fontTitle = createFont(-27, 650, "Segoe UI")
		fontMono = createFont(-14, 400, "Consolas")
		brushBg, _, _ = pCreateBrush.Call(COLOR_BG)
		brushEdit, _, _ = pCreateBrush.Call(COLOR_EDIT)
		hwndTitle = create(hwnd, "STATIC", "NMS  PERSONAL ATLAS", WS_CHILD|WS_VISIBLE|SS_LEFT, 0, 26, 22, 632, 42, 100)
		hwndHint = create(hwnd, "STATIC", "Préparation de votre Atlas personnel", WS_CHILD|WS_VISIBLE, 0, 28, 76, 620, 25, 101)
		hwndStatus = create(hwnd, "STATIC", "Vérification de l'environnement…", WS_CHILD|WS_VISIBLE, 0, 28, 125, 620, 26, 102)
		hwndStep = create(hwnd, "STATIC", "Vérification des composants", WS_CHILD|WS_VISIBLE, 0, 28, 154, 622, 23, 103)
		hwndProgress = create(hwnd, "msctls_progress32", "", WS_CHILD|WS_VISIBLE|PBS_MARQUEE, 0, 28, 193, 620, 18, 104)
		hwndToggle = create(hwnd, "BUTTON", "Afficher les détails", WS_CHILD|WS_VISIBLE|WS_TABSTOP, 0, 28, 233, 179, 33, IDC_DETAIL)
		hwndOpenLog = create(hwnd, "BUTTON", "Ouvrir le journal", WS_CHILD|WS_VISIBLE|WS_TABSTOP, 0, 219, 233, 178, 33, IDC_LOG)
		hwndClose = create(hwnd, "BUTTON", "Fermer", WS_CHILD|WS_VISIBLE|WS_TABSTOP, 0, 512, 233, 135, 33, IDC_CLOSE)
		hwndDetails = create(hwnd, "EDIT", "", WS_CHILD|ES_MULTILINE|ES_READONLY|ES_AUTOVSCROLL|WS_VSCROLL|WS_BORDER, WS_EX_CLIENTEDGE, 28, 286, 620, 230, 105)
		for _, h := range []uintptr{hwndHint, hwndStatus, hwndStep, hwndToggle, hwndOpenLog, hwndClose} {
			setFont(h, fontBody)
		}
		setFont(hwndTitle, fontTitle)
		setFont(hwndDetails, fontMono)
		pSendMessage.Call(hwndProgress, PBM_SETMARQUEE, 1, 35)
		pSetTimer.Call(hwnd, 1, 280, 0)
		go launchAtlas()
		return 0
	case WM_CTLCOLORSTATIC:
		pSetTextColor.Call(wp, COLOR_FG)
		if lp == hwndHint || lp == hwndStep {
			pSetTextColor.Call(wp, COLOR_MUTED)
		}
		pSetBkMode.Call(wp, TRANSPARENT)
		return brushBg
	case WM_CTLCOLOREDIT:
		pSetTextColor.Call(wp, COLOR_FG)
		pSetBkColor.Call(wp, COLOR_EDIT)
		return brushEdit
	case WM_COMMAND:
		switch wp & 0xffff {
		case IDC_DETAIL:
			expanded = !expanded
			if expanded {
				pShowWindow.Call(hwndDetails, SW_SHOW)
				showText(hwndToggle, "Masquer les détails")
				setWindowHeight(568)
			}
			if !expanded {
				pShowWindow.Call(hwndDetails, SW_HIDE)
				showText(hwndToggle, "Afficher les détails")
				setWindowHeight(319)
			}
			return 0
		case IDC_LOG:
			shellOpenFile(logPath)
			return 0
		case IDC_CLOSE:
			pSendMessage.Call(hwnd, WM_CLOSE, 0, 0)
			return 0
		}
	case WM_TIMER:
		updateFromLog()
		if currentView.Error != "" && terminalError == "" {
			terminalError = currentView.Error
		}
		if terminalError != "" {
			showText(hwndStatus, "Installation interrompue")
			showText(hwndStep, "Consultez les détails. "+terminalError)
			pSendMessage.Call(hwndProgress, PBM_SETMARQUEE, 0, 0)
		}
		if started && isAtlasHealthy() {
			pKillTimer.Call(hwnd, 1)
			pDestroyWindow.Call(hwnd)
		}
		runnerMutex.Lock()
		done, runErr := runnerDone, runnerErr
		runnerMutex.Unlock()
		if done && runErr != nil && terminalError == "" {
			terminalError = fmt.Sprintf("Atlas n'a pas démarré : %v", runErr)
			showText(hwndStatus, "Le démarrage a échoué")
			showText(hwndStep, terminalError)
			pSendMessage.Call(hwndProgress, PBM_SETMARQUEE, 0, 0)
		}
		return 0
	case WM_CLOSE:
		runnerMutex.Lock()
		isDone, process := runnerDone, activeCmd
		runnerMutex.Unlock()
		if !isDone && !started {
			answer := showMessage("Fermer Atlas ?", "La préparation est en cours. L'interrompre nécessitera de reprendre l'installation au prochain lancement. Voulez-vous fermer ?", MB_YESNO|MB_ICONQUESTION)
			if answer != IDYES {
				return 0
			}
			if process != nil && process.Process != nil {
				_ = process.Process.Kill()
			}
		}
		pDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		pKillTimer.Call(hwnd, 1)
		for _, h := range []uintptr{fontBody, fontTitle, fontMono, brushBg, brushEdit} {
			if h != 0 {
				pDeleteObject.Call(h)
			}
		}
		pPostQuit.Call(0)
		return 0
	}
	result, _, _ := pDefWindowProc.Call(hwnd, uintptr(message), wp, lp)
	return result
}

func isAtlasHealthy() bool {
	if !started {
		return false
	}
	return serverIsReady()
}

func serverIsReady() bool {
	c := http.Client{Timeout: 600 * time.Millisecond}
	resp, err := c.Get("http://127.0.0.1:8532/_stcore/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
func updateFromLog() {
	f, err := os.Open(logPath)
	if err == nil {
		defer f.Close()
		if stat, e := f.Stat(); e == nil && logReadOffset > stat.Size() {
			logReadOffset = 0
			logRemainder = ""
		}
		if _, e := f.Seek(logReadOffset, io.SeekStart); e == nil {
			buf, e := io.ReadAll(io.LimitReader(f, 256*1024))
			if e == nil && len(buf) > 0 {
				logReadOffset += int64(len(buf))
				detailsText += string(buf)
				if len(detailsText) > 16000 {
					detailsText = detailsText[len(detailsText)-16000:]
				}
				if detailsText != oldDetails {
					showText(hwndDetails, detailsText)
					oldDetails = detailsText
				}
				chunk := logRemainder + string(buf)
				lines := strings.Split(chunk, "\n")
				logRemainder = lines[len(lines)-1]
				for _, line := range lines[:len(lines)-1] {
					currentView = parseStatusLine(line, currentView)
					if strings.HasPrefix(line, "@ATLAS_UI\tSTART\t") {
						started = true
					}
				}
				if currentView.Phase > 0 && currentView.Error == "" {
					st := fmt.Sprintf("Étape %d/4 — %s", currentView.Phase, currentView.Text)
					if st != oldStatus {
						showText(hwndStatus, st)
						oldStatus = st
					}
					labels := []string{"", "Données personnelles", "Vérification du runtime", "Téléchargement et installation", "Ouverture d'Atlas"}
					if currentView.Phase <= 4 {
						showText(hwndStep, labels[currentView.Phase])
					}
				}
			}
		}
	}
}
func launchAtlas() {
	executable, err := os.Executable()
	if err != nil {
		runnerMutex.Lock()
		runnerErr = err
		runnerDone = true
		runnerMutex.Unlock()
		return
	}
	appRoot = filepath.Dir(executable)
	dataFolder = filepath.Join(appRoot, "NMSFaunaAnalyzerData")
	logPath = filepath.Join(dataFolder, "runtime", "logs", "demarrage_interface.log")
	if err = os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		runnerMutex.Lock()
		runnerErr = err
		runnerDone = true
		runnerMutex.Unlock()
		return
	}
	python := filepath.Join(appRoot, "runtime", "python", "python.exe")
	script := filepath.Join(appRoot, "tools", "distribution", "demarrer_atlas.py")
	if _, err = os.Stat(python); err != nil {
		runnerMutex.Lock()
		runnerErr = fmt.Errorf("Python embarqué introuvable ; réextraire le ZIP Atlas")
		runnerDone = true
		runnerMutex.Unlock()
		return
	}
	if _, err = os.Stat(script); err != nil {
		runnerMutex.Lock()
		runnerErr = fmt.Errorf("Programme Atlas incomplet ; réextraire le ZIP")
		runnerDone = true
		runnerMutex.Unlock()
		return
	}
	logfile, e := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		runnerMutex.Lock()
		runnerErr = e
		runnerDone = true
		runnerMutex.Unlock()
		return
	}
	_, _ = fmt.Fprintf(logfile, "\n=== Atlas GUI Launcher — %s ===\n", time.Now().Format(time.RFC3339))
	cmd := exec.Command(python, script)
	cmd.Dir = appRoot
	cmd.Stdout = logfile
	cmd.Stderr = logfile
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	env := os.Environ()
	env = append(env,
		"NMS_FAUNA_DATA_DIR="+dataFolder,
		"NMS_RUNTIME_DIR="+filepath.Join(dataFolder, "runtime"),
		"STREAMLIT_BROWSER_GATHER_USAGE_STATS=false", "STREAMLIT_SERVER_HEADLESS=true",
		"PYTHONNOUSERSITE=1", "PYTHONIOENCODING=utf-8", "PIP_NO_INPUT=1", "PIP_DISABLE_PIP_VERSION_CHECK=1",
		"PYTHONPATH=", "PYTHONHOME=", "ATLAS_RUN_PORT=8532")
	cmd.Env = env
	if e = cmd.Start(); e != nil {
		runnerMutex.Lock()
		runnerErr = e
		runnerDone = true
		runnerMutex.Unlock()
		logfile.Close()
		return
	}
	runnerMutex.Lock()
	activeCmd = cmd
	runnerMutex.Unlock()
	e = cmd.Wait()
	logfile.Close()
	runnerMutex.Lock()
	if e != nil {
		runnerErr = e
	} else {
		runnerErr = errors.New("Le serveur Atlas s'est arrêté")
	}
	runnerDone = true
	runnerMutex.Unlock()
}
func main() {
	if exePath, err := os.Executable(); err == nil {
		logPath = filepath.Join(filepath.Dir(exePath), "NMSFaunaAnalyzerData", "runtime", "logs", "demarrage_interface.log")
		if fileInfo, err := os.Stat(logPath); err == nil {
			logReadOffset = fileInfo.Size() // Show this session only; keep the historical log.
		}
	}
	if serverIsReady() {
		pShellExecute.Call(0, strptr("open"), strptr("http://127.0.0.1:8532"), 0, 0, SW_SHOW)
		return
	}
	clsName := "NMSAtlasInstallerWindow"
	classString = wide(clsName)
	instance, _, _ := pGetModuleHandle.Call(0)
	cursor, _, _ := pLoadCursor.Call(0, 32512)
	brush, _, _ := pCreateBrush.Call(COLOR_BG)
	exePath, _ := os.Executable()
	icoPath := filepath.Join(filepath.Dir(exePath), "assets", "atlas_launcher.ico")
	appIcon, _, _ := pLoadImage.Call(0, strptr(icoPath), 1, 0, 0, 0x0050) // IMAGE_ICON | LR_LOADFROMFILE | LR_DEFAULTSIZE
	cls := wndClassEx{Size: uint32(unsafe.Sizeof(wndClassEx{})), Proc: syscall.NewCallback(wndProc), Instance: instance, Cursor: cursor, Background: brush, ClassName: classString, Icon: appIcon, IconSm: appIcon}
	if r, _, _ := pRegisterClass.Call(uintptr(unsafe.Pointer(&cls))); r == 0 {
		return
	}
	icc := initCommonControlsEx{Size: uint32(unsafe.Sizeof(initCommonControlsEx{})), Flags: ICC_PROGRESS_CLASS}
	pInitCommon.Call(uintptr(unsafe.Pointer(&icc)))
	sw, _, _ := pGetSystemMetrics.Call(0)
	sh, _, _ := pGetSystemMetrics.Call(1)
	x := int32(sw)/2 - 345
	y := int32(sh)/2 - 170
	hwnd, _, _ := pCreateWindow.Call(0, uintptr(unsafe.Pointer(classString)), strptr("NMS Personal Atlas — Préparation"),
		WS_CAPTION|WS_SYSMENU|WS_MINIMIZEBOX|WS_VISIBLE, uintptr(x), uintptr(y), 690, 319, 0, 0, instance, 0)
	if hwnd == 0 {
		return
	}
	pShowWindow.Call(hwnd, SW_SHOW)
	pUpdateWindow.Call(hwnd)
	var message msgStruct
	for {
		r, _, _ := pGetMessage.Call(uintptr(unsafe.Pointer(&message)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&message)))
		pDispatchMessage.Call(uintptr(unsafe.Pointer(&message)))
	}
}
