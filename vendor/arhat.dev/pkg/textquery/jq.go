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
	"bytes"
	"encoding/json"
	"io"
)

// JQ runs query over json data
func JQ(query, data string) (string, error) {
	return JQBytes(query, []byte(data))
}

// JQ runs query over json data bytes
func JQBytes(query string, dataBytes []byte) (string, error) {
	return Query(query, NewJSONIterator(bytes.NewReader(dataBytes)), json.Marshal)
}

func NewJSONIterator(r io.Reader) func() (interface{}, bool) {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	exit := false
	return func() (interface{}, bool) {
		if exit {
			return nil, false
		}

		var data interface{}
		err := dec.Decode(&data)
		if err != nil {
			if err == io.EOF {
				return nil, false
			}

			buffered, _ := io.ReadAll(dec.Buffered())
			remainder, _ := io.ReadAll(r)
			exit = true

			return append(buffered, remainder...), true
		}

		return data, true
	}
}
