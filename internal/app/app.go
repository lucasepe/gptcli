package app

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	Name               = "gptcli"
	shortcutsFilename  = "shortcuts"
	envGptCliAPIKey    = "GPTCLI_API_KEY"
	envGptCliShortcuts = "GPTCLI_SHORTCUTS"
)

func OpenAIApiKey() string {
	return os.Getenv(envGptCliAPIKey)
}

func ConfigPath() (string, error) {
	if pt := os.Getenv(envGptCliShortcuts); len(pt) > 0 {
		return filepath.Abs(pt)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fn := filepath.Join(cwd, fmt.Sprintf(".%s", shortcutsFilename))
	fp, err := os.Open(fn)
	if err == nil {
		fp.Close()
		return fn, nil
	}

	return "", fmt.Errorf("env variable '%s' is undefined", envGptCliShortcuts)
}
