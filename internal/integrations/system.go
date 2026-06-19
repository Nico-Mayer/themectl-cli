package integrations

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/nico-mayer/themectl-cli/internal/model"
	winReg "golang.org/x/sys/windows/registry"
)

const (
	HWND_BROADCAST   = 0xFFFF
	WM_SETTINGCHANGE = 0x001A
	WM_THEMECHANGED  = 0x031A
)

type SystemTheme struct{}

func init() {
	Register(SystemTheme{})
}

func (SystemTheme) Name() string {
	return "system-theme"
}

func (i SystemTheme) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	logger.Debug("applying", "appearance", themeInfo.Appearance, "os", runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		return i.setMacOSTheme(themeInfo)
	case "windows":
		return i.setWindowsTheme(themeInfo)
	default:
		return fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
}

func (i SystemTheme) setMacOSTheme(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	mode := strings.ToLower(themeInfo.Appearance)

	var script string
	switch mode {
	case "dark":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to true`
	case "light":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to false`
	default:
		return fmt.Errorf("unsupported appearance %q: expected \"dark\" or \"light\"", mode)
	}

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript failed: %w (output: %s)", err, strings.TrimSpace(string(output)))
	}

	logger.Info("applied", "appearance", mode)
	return nil
}

func (i SystemTheme) setWindowsTheme(themeInfo model.ThemeInfo) error {
	// logger := integrationLogger(i)
	mode := strings.ToLower(themeInfo.Appearance)

	var value uint32
	switch mode {
	case "dark":
		value = 0
	case "light":
		value = 1
	default:
		return fmt.Errorf("unsupported appearance %q: expected \"dark\" or \"light\"", mode)
	}

	key, err := winReg.OpenKey(
		winReg.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`,
		winReg.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.SetDWordValue("AppsUseLightTheme", value); err != nil {
		return err
	}

	if err := key.SetDWordValue("SystemUsesLightTheme", value); err != nil {
		return err
	}

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	param, _ := syscall.UTF16PtrFromString("ImmersiveColorSet")

	sendMessageTimeout.Call(
		HWND_BROADCAST,
		WM_SETTINGCHANGE,
		0,
		uintptr(unsafe.Pointer(param)),
		0,
		100,
		0,
	)

	sendMessageTimeout.Call(
		HWND_BROADCAST,
		WM_THEMECHANGED,
		0,
		0,
		0,
		100,
		0,
	)

	return nil
}
