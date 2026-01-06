package engine

import (
	"encoding/json"
	"os"
)

func MergePackageJSON(targetPath string, patchPath string) error {
	targetBytes, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}
	patchBytes, err := os.ReadFile(patchPath)
	if err != nil {
		return err
	}

	var target map[string]any
	var patch map[string]any
	if err := json.Unmarshal(targetBytes, &target); err != nil {
		return err
	}
	if err := json.Unmarshal(patchBytes, &patch); err != nil {
		return err
	}

	mergeMap(target, patch)

	out, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(targetPath, out, 0644)
}

func mergeMap(dst, src map[string]any) {
	for k, v := range src {
		// Special-case known keys that are maps
		if m2, ok := v.(map[string]any); ok {
			if m1, ok := dst[k].(map[string]any); ok {
				mergeMap(m1, m2)
				dst[k] = m1
				continue
			}
		}
		dst[k] = v
	}
}
