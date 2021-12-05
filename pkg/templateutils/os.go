package templateutils

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
)

func createOSNS(rc dukkha.RenderingContext) *osNS {
	return &osNS{rc: rc}
}

type osNS struct {
	rc dukkha.RenderingContext
}

func (ns *osNS) ReadFile(filename string) (string, error) {
	data, err := ns.rc.FS().ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ns *osNS) WriteFile(filename string, d interface{}) error {
	var data []byte
	switch dt := d.(type) {
	case string:
		data = []byte(dt)
	case []byte:
		data = dt
	default:
		return fmt.Errorf("invalid non string nor bytes data: %T", d)
	}

	return ns.rc.FS().WriteFile(filename, data, 0640)
}

func (ns *osNS) MkdirAll(path string) error {
	return ns.rc.FS().MkdirAll(path, 0755)
}
