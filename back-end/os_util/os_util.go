package os_util

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// This is dumb but necessary
// https://stackoverflow.com/questions/31059023/how-to-reference-a-relative-file-from-code-and-tests
func RelativeProjectPath(rel string) string {
	currentPath := getCurrentPath()
	baseProjectPath := strings.Split(currentPath, "/back-end/")[0]

	return filepath.Join(baseProjectPath, rel)
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
