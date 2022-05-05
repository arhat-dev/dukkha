package dukkha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskExecOptions(t *testing.T) {
	t.Parallel()

	opts := CreateTaskExecOptions(1, 10).(*taskExecOpts)
	assert.Equal(t, opts.id, 1)
	assert.Equal(t, opts.seq, -1)
	assert.Equal(t, opts.total, 10)

	for i := 0; i < opts.total; i++ {
		mOpts := opts.NextMatrixExecOptions().(*taskMatrixExecOpts)
		assert.Equal(t, opts.id, 1)
		assert.Equal(t, opts.seq, i)
		assert.Equal(t, opts.total, 10)

		assert.Equal(t, mOpts.id, mOpts.ID())
		assert.Equal(t, mOpts.ID(), opts.id)

		assert.Equal(t, mOpts.seq, mOpts.Seq())
		assert.Equal(t, mOpts.Seq(), opts.seq)

		assert.Equal(t, mOpts.total, mOpts.Total())
		assert.Equal(t, mOpts.Total(), opts.total)

		assert.Equal(t, i == opts.total-1, mOpts.IsLast())
	}
}
