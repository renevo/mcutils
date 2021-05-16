package java

import (
	"path/filepath"
	"runtime"
)

// ExecPath returns the executable path for java based on the supplied Java Home
func ExecPath(jh string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(jh, "bin", "java.exe")
	}

	return filepath.Join(jh, "bin", "java")
}
