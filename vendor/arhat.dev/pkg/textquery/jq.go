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

package textquery

import (
	"encoding/json"

	"arhat.dev/pkg/decodehelper"
)

func JQ(query, data string) (string, error) {
	return JQBytes(query, []byte(data))
}

func JQBytes(query string, input []byte) (string, error) {
	return Query(query, input, decodehelper.UnmarshalJSON, json.Marshal)
}
