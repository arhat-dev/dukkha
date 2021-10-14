package rshelper

import (
	"fmt"
	"os"

	"arhat.dev/rs"
)

var _ rs.RenderingHandler = (*FileHandler)(nil)

// FileHandler treats rawData as file path
type FileHandler struct{}

func (h *FileHandler) RenderYaml(
	_ string, rawData interface{},
) (interface{}, error) {
	path := ""
	switch t := rawData.(type) {
	case string:
		path = t
	case []byte:
		path = string(t)
	default:
		return nil, fmt.Errorf("invalid raw data type, want string or []byte, got %T", t)
	}

	return os.ReadFile(path)
}
