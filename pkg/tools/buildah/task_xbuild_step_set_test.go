package buildah

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
)

func Test_kvArgs(t *testing.T) {
	t.Parallel()

	type args struct {
		flag    string
		entries []*dukkha.EnvEntry
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty",
			args: args{
				flag:    "--empty",
				entries: []*dukkha.EnvEntry{},
			},
			want: nil,
		},
		{
			name: "Remove",
			args: args{
				flag: "--remove",
				entries: []*dukkha.EnvEntry{
					{
						Name:  "key-",
						Value: "",
					},
				},
			},
			want: []string{"--remove", "key-"},
		},
		{
			name: "Remove (value ignored)",
			args: args{
				flag: "--remove",
				entries: []*dukkha.EnvEntry{
					{
						Name:  "key-",
						Value: "value",
					},
				},
			},
			want: []string{"--remove", "key-"},
		},
		{
			name: "Key Only",
			args: args{
				flag: "--key-only",
				entries: []*dukkha.EnvEntry{
					{
						Name:  "key",
						Value: "",
					},
				},
			},
			want: []string{"--key-only", "key="},
		},
		{
			name: "Value Only",
			args: args{
				flag: "--value-only",
				entries: []*dukkha.EnvEntry{
					{
						Name:  "",
						Value: "value",
					},
				},
			},
			want: []string{"--value-only", "=value"},
		},
		{
			name: "Key Value Pair",
			args: args{
				flag: "--key-value",
				entries: []*dukkha.EnvEntry{
					{
						Name:  "key",
						Value: "value",
					},
				},
			},
			want: []string{"--key-value", "key=value"},
		},
		{
			name: "Multiple Key Value Pairs",
			args: args{
				flag: "--key-value",
				entries: []*dukkha.EnvEntry{
					{
						Name:  "key",
						Value: "value",
					},
					{
						Name:  "key",
						Value: "value",
					},
					{
						Name:  "key",
						Value: "value",
					},
				},
			},
			want: []string{"--key-value", "key=value", "--key-value", "key=value", "--key-value", "key=value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.want, kvArgs(tt.args.flag, tt.args.entries))
		})
	}
}
