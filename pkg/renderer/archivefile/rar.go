package archivefile

import (
	"fmt"
	"io"

	"github.com/nwaples/rardecode"
)

// nolint:unparam
func unrar(src io.Reader, target, password string) (io.Reader, error) {
	// TODO: implement
	_ = target
	r, err := rardecode.NewReader(src, password)
	if err != nil {
		return nil, err
	}
	_ = r
	return nil, fmt.Errorf("no implementation")
}
