package tools

import (
	"encoding/hex"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"arhat.dev/pkg/sha256helper"
	"arhat.dev/rs"
)

func GetScriptCache(cacheDir, script string) (string, error) {
	scriptName := hex.EncodeToString(sha256helper.Sum([]byte(script)))
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
		if !token.IsExported(f.Name) {
			// unexported, ignore
			continue
		}

		if f.Anonymous && f.Name == "BaseTask" {
			// it's me
			continue
		}

		dukkhaTags, hasDukkhaTags := f.Tag.Lookup(rs.TagName)
		yamlTags, hasYamlTags := f.Tag.Lookup("yaml")

		switch {
		case hasYamlTags:
			if strings.Contains(yamlTags, "-") {
				// ignored explicitly
				continue
			}
		case hasDukkhaTags:
			if !strings.Contains(dukkhaTags, "other") {
				continue
			}
		default:
			// no yaml tag, not a catch other field, ignore
			continue
		}

		ret = append(ret, f.Name)
	}

	return ret
}
