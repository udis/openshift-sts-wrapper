package util

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens a URL in the user's default browser
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return nil // Unsupported platform, silently fail
	}

	return cmd.Start()
}
