package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"codeberg.org/bit101/go-ansi"
)

const (
	SYM_YELLOW_BOX = "🟨"
	SYM_GREEN_BOX  = "🟩"
	SYM_BLUE_BOX   = "🟦"
	SYM_RED_BOX    = "🟥"
)

func PrintOk(s string) {
	ansi.Printf(ansi.Default, SYM_GREEN_BOX)
	ansi.Printf(ansi.Green, "[%s]", s)
}

func PrintfWarn(s string) {
	ansi.Printf(ansi.Default, SYM_YELLOW_BOX)
	ansi.Printf(ansi.Yellow, "[%s]", s)
}
func PrintErr(s string) {
	ansi.Printf(ansi.Default, SYM_RED_BOX)
	ansi.Printf(ansi.Red, "[%s]", s)
}

func Print(s string) {
	ansi.Printf(ansi.Default, SYM_BLUE_BOX)
	ansi.Printf(ansi.Blue, "[%s]", s)
}

func Sum(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func IsValidSum(path string, expected string) (bool, error) {
	sum, err := Sum(path)
	if err != nil {
		return false, err
	}
	return fmt.Sprintf("%x", sum) == expected, nil
}
