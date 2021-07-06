package tools

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"arhat.dev/pkg/hashhelper"
	"gopkg.in/yaml.v3"
)

type TaskReference struct {
	ToolKind string
	ToolName string
	TaskKind string
	TaskName string

	MatrixFilter map[string][]string
}

func (r *TaskReference) HasToolName() bool {
	return len(r.ToolName) != 0
}

// ParseTaskReference parse task ref
//
// <tool-kind>{:<tool-name>}:<task-kind>(<task-name>, ...)
//
// e.g. buildah:bud(dukkha) # use default matrix
// 		buildah:bud(dukkha, {kernel: [linux]}) # use custom matrix
//		buildah:in-docker:bud(dukkha, {kernel: [linux]}) # with tool-name
func ParseTaskReference(taskRef string) (*TaskReference, error) {
	callStart := strings.IndexByte(taskRef, '(')
	if callStart < 0 {
		return nil, fmt.Errorf("missing task call `()`")
	}

	call, err := ParseBrackets(taskRef[callStart+1:])
	if err != nil {
		return nil, fmt.Errorf("invalid task call: %w", err)
	}

	ref := &TaskReference{}
	callArgs := strings.SplitN(call, ",", 2)
	ref.TaskName = strings.TrimSpace(callArgs[0])

	switch len(callArgs) {
	case 1:
		// default matrix spec
	case 2:
		ref.MatrixFilter = make(map[string][]string)
		err = yaml.Unmarshal([]byte(callArgs[1]), &ref.MatrixFilter)
		if err != nil {
			return nil, fmt.Errorf("invalid matrix arg %q: %w", callArgs[1], err)
		}
	default:
		return nil, fmt.Errorf(
			"invalid number of task call args, expecting 1 or 2 args, got %q (%d)",
			call, len(callArgs),
		)
	}

	parts := strings.Split(taskRef[:callStart], ":")
	ref.ToolKind = parts[0]

	switch len(parts) {
	case 2:
		ref.TaskKind = parts[1]
	case 3:
		ref.ToolName = parts[1]
		ref.TaskKind = parts[2]
	default:
		return nil, fmt.Errorf("invalid prefix %q", taskRef)
	}

	return ref, nil
}

// ParseBrackets `()`
func ParseBrackets(s string) (string, error) {
	leftBrackets := 0
	for i := range s {
		switch s[i] {
		case '(':
			leftBrackets++
		case ')':
			if leftBrackets == 0 {
				return s[:i], nil
			}
			leftBrackets--
		}
	}

	// invalid data
	return "", fmt.Errorf("unexpected non-terminated brackets")
}

func GetScriptCache(cacheDir, script string) (string, error) {
	scriptName := hex.EncodeToString(hashhelper.Sha256Sum([]byte(script)))
	scriptPath := filepath.Join(cacheDir, scriptName)

	_, err := os.Stat(scriptPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to check existence of script cache: %w", err)
		}

		err = os.WriteFile(scriptPath, []byte(script), 0600)
		if err != nil {
			return "", fmt.Errorf("failed to write script cache: %w", err)
		}
	}

	return scriptPath, nil
}

func CartesianProduct(m map[string][]string) []map[string]string {
	names := make([]string, 0)
	mat := make([][]string, 0)
	for k, v := range m {
		names = append(names, k)
		mat = append(mat, v)
	}

	sort.Slice(names, func(i, j int) bool {
		ok := names[i] < names[j]
		if ok {
			mat[i], mat[j] = mat[j], mat[i]
		}
		return ok
	})

	listCart := cartNext(mat)

	result := make([]map[string]string, 0)
	for _, list := range listCart {
		vMap := make(map[string]string)
		for i, v := range list {
			vMap[names[i]] = v
		}
		result = append(result, vMap)
	}

	return result
}

func cartNext(mat [][]string) [][]string {
	if len(mat) == 0 {
		return nil
	}

	tupleCount := 1
	for _, list := range mat {
		if len(list) == 0 {
			// ignore empty list
			continue
		}

		tupleCount *= len(list)
	}

	result := make([][]string, tupleCount)

	buf := make([]string, tupleCount*len(mat))
	indexPerList := make([]int, len(mat))

	start := 0
	for i := range result {
		end := start + len(mat)

		tuple := buf[start:end]
		result[i] = tuple

		start = end

		for j, idx := range indexPerList {
			// mat[j] is the list

			tuple[j] = mat[j][idx]
		}

		for j := len(indexPerList) - 1; j >= 0; j-- {
			indexPerList[j]++
			if indexPerList[j] < len(mat[j]) {
				break
			}

			// reset for next tuple
			indexPerList[j] = 0
		}
	}

	return result
}
