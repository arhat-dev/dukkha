package tools

import (
	"fmt"
	"io"
	"strconv"

	"arhat.dev/dukkha/pkg/dukkha"
)

type Actions []*Action

func ResolveActions(
	rc dukkha.TaskExecContext,
	parent dukkha.Resolvable,
	actions *Actions,
	actionsTagName string,
) ([]dukkha.TaskExecSpec, error) {
	jobIndex := make(map[string]int)
	mCtx := rc.DeriveNew()
	err := parent.DoAfterFieldsResolved(mCtx, 1, true, func() error {
		jobs := *actions

		for i := range jobs {
			name := jobs[i].Name
			if len(name) == 0 {
				name = strconv.FormatInt(int64(i), 10)
			}

			jobIndex[name] = i
		}

		return nil
	}, actionsTagName)

	if err != nil {
		return nil, err
	}

	if len(jobIndex) == 0 {
		return nil, nil
	}

	return next(mCtx,
		parent, actions, actionsTagName,
		jobIndex, 0,
	)
}

func next(
	mCtx dukkha.TaskExecContext,
	x dukkha.Resolvable,
	actions *Actions,
	actionsTagName string,

	// data
	jobIndex map[string]int,
	index int,
) ([]dukkha.TaskExecSpec, error) {
	var (
		thisAction dukkha.RunTaskOrRunCmd
		thisJob    *Action

		skip bool
	)

	var err error
	// depth = 1 to get job list only, DO NOT render inner actions for now
	err = x.DoAfterFieldsResolved(mCtx, 1, true, func() error {
		jobs := *actions

		if index >= len(jobs) {
			return nil
		}

		thisJob = jobs[index]

		// resolve single action
		return thisJob.DoAfterFieldResolved(mCtx, func(run bool) error {
			if !run {
				skip = true
				return nil
			}

			thisAction, err = thisJob.GenSpecs(mCtx.DeriveNew(), index)
			return err
		})
	}, actionsTagName)

	if err != nil || thisJob == nil {
		return nil, err
	}

	if skip {
		// not running this action, continue to next
		// DO NOT depend on thisAction value, as it can be nil when using idle
		return next(mCtx, x, actions, actionsTagName, jobIndex, index+1)
	}

	return []dukkha.TaskExecSpec{
		{
			// respect continue_on_error, when not set or set to false
			// then we SHOULD not go further

			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				return thisAction, nil
			},
		},
		{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				var ni int

				// we will dead lock self when *next is self and calling
				// next() directly inside
				// so DO NOT do it
				err = thisJob.DoAfterFieldResolved(mCtx, func(run bool) error {
					if nj := thisJob.Next; nj != nil {
						var ok bool
						ni, ok = jobIndex[*nj]
						if !ok {
							return fmt.Errorf("unknown next job reference %q", *nj)
						}
						return nil
					}

					ni = index + 1
					return nil
				}, "next")

				if err != nil {
					return nil, err
				}

				return next(
					mCtx,
					x, actions, actionsTagName,
					jobIndex, ni,
				)
			},
		},
	}, nil
}
