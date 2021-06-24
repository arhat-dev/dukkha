/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package constant

import "context"

type ContextKey string

// nolint:revive
const (
	ContextKeyConfig ContextKey = "config"

	contextKeyWorkerCount  ContextKey = "worker_count"
	contextKeyMatrixFilter ContextKey = "matrix_filter"
)

func WithWorkerCount(ctx context.Context, n int) context.Context {
	if n < 1 {
		n = 1
	}

	return context.WithValue(ctx, contextKeyWorkerCount, n)
}

func GetWorkerCount(ctx context.Context) int {
	v, ok := ctx.Value(contextKeyWorkerCount).(int)
	if !ok {
		// default in serial mode
		return 1
	}

	return v
}

func WithMatrixFilter(ctx context.Context, filter map[string][]string) context.Context {
	return context.WithValue(ctx, contextKeyMatrixFilter, filter)
}

func GetMatrixFilter(ctx context.Context) map[string][]string {
	v, ok := ctx.Value(contextKeyMatrixFilter).(map[string][]string)
	if !ok {
		// default in serial mode
		return nil
	}

	return v
}
