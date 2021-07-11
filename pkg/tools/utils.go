package tools

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

func getFieldNamesToResolve(typ reflect.Type) []string {
	var ret []string
	for i := 1; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !(f.Name[0] >= 'A' && f.Name[0] <= 'Z') {
			// unexported, ignore
			continue
		}

		if f.Anonymous && f.Name == "BaseTask" {
			// it's me
			continue
		}

		val, ok := f.Tag.Lookup("yaml")
		if !ok {
			// no yaml field, ignore
			continue
		}

		if strings.Contains(val, "-") {
			// ignored explicitly
			continue
		}

		ret = append(ret, f.Name)
	}

	return ret
}
