package base

import (
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
)

//Generate random strings with given length
func GenerateRandomName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Execute : Run apictl command
//
func Execute(t *testing.T, args ...string) (string, error) {
	cmd := exec.Command(RelativeBinaryPath+BinaryName, args...)

	t.Log("base.Execute() - apkctl command:", cmd.String())
	// run command
	output, err := cmd.Output()

	t.Log("base.Execute() - apkctl command output:", string(output))
	return string(output), err
}

func GetExportedPathFromOutput(output string) string {
	//Check directory path to omit changes due to OS differences
	if strings.Contains(output, ":\\") {
		arrayOutput := []rune(output)
		extractedPath := string(arrayOutput[strings.Index(output, ":\\")-1:])
		return strings.ReplaceAll(strings.ReplaceAll(extractedPath, "\n", ""), " ", "")
	} else {
		return strings.ReplaceAll(strings.ReplaceAll(output[strings.Index(output, string(os.PathSeparator)):], "\n", ""), " ", "")
	}
}

// IsFileAvailable checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func IsFileAvailable(t *testing.T, filepath string) bool {
	t.Log("base.IsFileAvailable() - API file path:", filepath)

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
