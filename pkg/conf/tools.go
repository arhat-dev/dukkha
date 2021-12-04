package conf

import (
	"fmt"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

type Tools struct {
	rs.BaseField `yaml:"-"`

	Tools map[string][]dukkha.Tool `yaml:",inline"`
}

func (m *Tools) Merge(a *Tools) error {
	err := m.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("inherit tools config: %w", err)
	}

	if len(a.Tools) != 0 {
		if m.Tools == nil {
			m.Tools = make(map[string][]dukkha.Tool)
		}

		for k := range a.Tools {
			m.Tools[k] = append(m.Tools[k], a.Tools[k]...)
		}
	}

	return nil
}
