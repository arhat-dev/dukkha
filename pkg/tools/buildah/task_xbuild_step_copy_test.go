package buildah

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestStepCopy_genSpec(t *testing.T) {
	const (
		baseCmd = "buildah"
	)

	tests := []struct {
		name     string
		spec     *stepCopy
		expected []dukkha.TaskExecSpec
	}{
		{
			name: "From local To workdir",
			spec: &stepCopy{
				From: copyFromSpec{
					Local: &copyFromLocalSpec{
						Path: "/foo",
					},
				},
				To: copyToSpec{
					Path: "",
				},
			},
			expected: []dukkha.TaskExecSpec{
				{
					Command: []string{baseCmd, "copy", replace_XBUILD_CURRENT_CONTAINER_ID, "/foo"},
				},
			},
		},
		{
			name: "From local To bar",
			spec: &stepCopy{
				From: copyFromSpec{
					Local: &copyFromLocalSpec{
						Path: "/foo",
					},
				},
				To: copyToSpec{
					Path: "/bar",
				},
			},
			expected: []dukkha.TaskExecSpec{
				{
					Command: []string{baseCmd, "copy", replace_XBUILD_CURRENT_CONTAINER_ID, "/foo", "/bar"},
				},
			},
		},
		{
			name: "From http To bar",
			spec: &stepCopy{
				From: copyFromSpec{
					HTTP: &copyFromHTTPSpec{
						URL: "https://example.com",
					},
				},
				To: copyToSpec{
					Path: "/bar",
				},
			},
			expected: []dukkha.TaskExecSpec{
				{
					Command: []string{baseCmd, "copy", replace_XBUILD_CURRENT_CONTAINER_ID, "https://example.com", "/bar"},
				},
			},
		},
		// TODO: compare FixStdoutValueForReplace will always fail
		// 		{
		// 			name: "From image To bar",
		// 			spec: &stepCopy{
		// 				From: copyFromSpec{
		// 					Image: &copyFromImageSpec{
		// 						Ref:  "some-image",
		// 						Path: "/foo",
		// 					},
		// 				},
		// 				To: copyToSpec{
		// 					Path: "/bar",
		// 				},
		// 			},
		// 			expected: []dukkha.TaskExecSpec{
		// 				{
		// 					StdoutAsReplace:          "<XBUILD_COPY_FROM_IMAGE_ID>",
		// 					FixStdoutValueForReplace: bytes.TrimSpace,
		//
		// 					IgnoreError: false,
		//
		// 					ShowStdout: true,
		// 					Command:    []string{baseCmd, "pull", "some-image"},
		// 				},
		// 				{
		// 					Command: []string{baseCmd, "copy", "--from", "<XBUILD_COPY_FROM_IMAGE_ID>", replace_XBUILD_CURRENT_CONTAINER_ID, "/foo", "/bar"},
		// 				},
		// 			},
		// 		},
		{
			name: "From step To bar",
			spec: &stepCopy{
				From: copyFromSpec{
					Step: &copyFromStepSpec{
						ID:   "some-step",
						Path: "/foo",
					},
				},
				To: copyToSpec{
					Path: "/bar",
				},
			},
			expected: []dukkha.TaskExecSpec{
				{
					Command: []string{baseCmd, "copy", "--from", replace_XBUILD_STEP_CONTAINER_ID("some-step"), replace_XBUILD_CURRENT_CONTAINER_ID, "/foo", "/bar"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := dukkha_test.NewTestContext(context.TODO())
			ret, err := test.spec.genSpec(
				ctx,
				dukkha_test.CreateTaskMatrixExecOptions([]string{baseCmd}),
				false,
			)
			assert.NoError(t, err)
			assert.EqualValues(t, test.expected, ret)
		})
	}
}
