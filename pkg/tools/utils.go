package tools

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"arhat.dev/pkg/hashhelper"
)

func GetScriptCache(cacheDir, script string) (string, error) {
	scriptName := hex.EncodeToString(hashhelper.Sha256Sum([]byte(script)))
	scriptPath := filepath.Join(cacheDir, scriptName)

	_, err := os.Stat(scriptPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to check existence of script cache: %w", err)
		}

		err = os.WriteFile(scriptPath, []byte(script), 0600)
		if err != nil {
			return "", fmt.Errorf("failed to write script cache: %w", err)
		}
	}

	return scriptPath, nil
}
