package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskHeaderLineData_json(t *testing.T) {
	tests := []struct {
		name string
		data TaskHeaderLineData
		want string
	}{
		{
			name: "Normal",
			data: TaskHeaderLineData{
				ToolKind: "a",
				ToolName: "b",
				TaskKind: "c",
				TaskName: "d",
			},
			want: `{ "kind": "a:c", "tool_name": "b", "name": "d" }`,
		},
		{
			name: "No Tool Name",
			data: TaskHeaderLineData{
				ToolKind: "a",
				ToolName: "",
				TaskKind: "c",
				TaskName: "d",
			},
			want: `{ "kind": "a:c", "tool_name": "", "name": "d" }`,
		},
		{
			name: "No Tool and Task Name",
			data: TaskHeaderLineData{
				ToolKind: "a",
				ToolName: "",
				TaskKind: "c",
				TaskName: "",
			},
			want: `{ "kind": "a:c", "tool_name": "" }`,
		},
		{
			name: "Matrix",
			data: TaskHeaderLineData{
				ToolKind: "",
				ToolName: "",
				TaskKind: "",
				TaskName: "",
				Matrix: map[string]string{
					"a": "b",
				},
			},
			want: `{ "kind": ":", "tool_name": "", "matrix": { "a": "b" } }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.data.json())
		})
	}
}
