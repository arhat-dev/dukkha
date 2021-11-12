package tools

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"arhat.dev/dukkha/pkg/dukkha"
)

type Actions []*Action

func ResolveActions(
	rc dukkha.TaskExecContext,
	x dukkha.Resolvable,
	actionsFieldName string,
	actionsTagName string,
	options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	jobIndex := make(map[string]int)
	mCtx := rc.DeriveNew()
	err := x.DoAfterFieldsResolved(mCtx, 1, true, func() error {
		xv := reflect.ValueOf(x)
		for xv.Kind() == reflect.Ptr {
			xv = xv.Elem()
		}

		jobs := xv.FieldByName(actionsFieldName).Interface().(Actions)

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
		x, actionsFieldName, actionsTagName,
		jobIndex, options, 0,
	)
}

func next(
	mCtx dukkha.TaskExecContext,
	x dukkha.Resolvable,
	actionsFieldName string,
	actionsTagName string,

	// data
	jobIndex map[string]int,
	options dukkha.TaskMatrixExecOptions,
	index int,
) ([]dukkha.TaskExecSpec, error) {
	var (
		thisAction dukkha.RunTaskOrRunCmd
		thisJob    *Action
	)

	var err error
	// depth = 1 to get job list only, DO NOT render inner actions for now
	err = x.DoAfterFieldsResolved(mCtx, 1, true, func() error {
		xv := reflect.ValueOf(x)
		for xv.Kind() == reflect.Ptr {
			xv = xv.Elem()
		}

		jobs := xv.FieldByName(actionsFieldName).Interface().(Actions)

		if index >= len(jobs) {
			return nil
		}

		thisJob = jobs[index]

		// resolve single action
		return thisJob.DoAfterFieldResolved(mCtx, func() error {
			thisAction, err = thisJob.GenSpecs(mCtx.DeriveNew(), index)
			return err
		})
	}, actionsTagName)

	if err != nil || thisJob == nil {
		return nil, err
	}

	return []dukkha.TaskExecSpec{
		{
			// we should respect continune on error settings
			// when continue_on_error is not set or set to false
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
				err = thisJob.DoAfterFieldResolved(mCtx, func() error {
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
					x, actionsFieldName, actionsTagName,
					jobIndex, options, ni,
				)
			},
		},
	}, nil
}
