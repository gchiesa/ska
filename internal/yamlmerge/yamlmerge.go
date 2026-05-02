package yamlmerge

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// EngineID is the engine identifier used in ska-start bracket modifiers.
const EngineID = "yaml-merge"

// MergeYAML deep-merges src YAML content into dst YAML content.
// Keys present in src override those in dst; keys only in dst are preserved.
// Comments in dst are preserved where possible via yaml.v3 Node-based round-trip.
//
// Returns:
//   - merged YAML bytes (with the same leading indentation as dst)
//   - list of dot-separated key paths that were patched (for debug logging)
//   - any error encountered
func MergeYAML(dst, src []byte) ([]byte, []string, error) {
	// Detect and strip common indentation so yaml.v3 can parse without ambiguity.
	dstStripped, dstIndent := stripCommonIndent(dst)
	srcStripped, _ := stripCommonIndent(src)

	// If either side is blank, nothing to merge.
	if len(bytes.TrimSpace(dstStripped)) == 0 || len(bytes.TrimSpace(srcStripped)) == 0 {
		return dst, nil, nil
	}

	var dstDoc yaml.Node
	if err := yaml.Unmarshal(dstStripped, &dstDoc); err != nil {
		return nil, nil, fmt.Errorf("failed to parse destination YAML: %w", err)
	}

	var srcDoc yaml.Node
	if err := yaml.Unmarshal(srcStripped, &srcDoc); err != nil {
		return nil, nil, fmt.Errorf("failed to parse source YAML: %w", err)
	}

	if dstDoc.Kind != yaml.DocumentNode || len(dstDoc.Content) == 0 {
		return dst, nil, nil
	}
	if srcDoc.Kind != yaml.DocumentNode || len(srcDoc.Content) == 0 {
		return dst, nil, nil
	}

	var changedPaths []string
	mergeNodes(dstDoc.Content[0], srcDoc.Content[0], "", &changedPaths)

	// Serialize the mutated dst document back to YAML.
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(dstDoc.Content[0]); err != nil {
		return nil, nil, fmt.Errorf("failed to serialize merged YAML: %w", err)
	}
	_ = enc.Close()

	merged := buf.Bytes()

	// Re-apply the original indentation that was stripped earlier.
	if dstIndent > 0 {
		merged = addIndent(merged, dstIndent)
	}

	return merged, changedPaths, nil
}

// mergeNodes recursively merges src node into dst node, recording changed paths.
func mergeNodes(dst, src *yaml.Node, keyPath string, changedPaths *[]string) {
	if dst.Kind != src.Kind {
		// Incompatible kinds – replace scalar/value and record the path.
		if keyPath != "" {
			*changedPaths = append(*changedPaths, keyPath)
		}
		dst.Kind = src.Kind
		dst.Tag = src.Tag
		dst.Value = src.Value
		dst.Style = src.Style
		dst.Content = src.Content
		return
	}

	switch dst.Kind {
	case yaml.MappingNode:
		mergeMappingNodes(dst, src, keyPath, changedPaths)
	case yaml.SequenceNode:
		// Sequences are replaced wholesale.
		if !nodesEqual(dst, src) {
			if keyPath != "" {
				*changedPaths = append(*changedPaths, keyPath)
			}
			dst.Content = src.Content
			dst.Style = src.Style
		}
	case yaml.ScalarNode:
		if dst.Value != src.Value {
			if keyPath != "" {
				*changedPaths = append(*changedPaths, keyPath)
			}
			dst.Value = src.Value
			dst.Tag = src.Tag
			dst.Style = src.Style
		}
	}
}

// mergeMappingNodes merges a src MappingNode into dst MappingNode.
func mergeMappingNodes(dst, src *yaml.Node, keyPath string, changedPaths *[]string) {
	// Build an index: key string -> index of the value node in dst.Content.
	dstKeyIdx := make(map[string]int, len(dst.Content)/2)
	for i := 0; i < len(dst.Content)-1; i += 2 {
		dstKeyIdx[dst.Content[i].Value] = i + 1
	}

	for i := 0; i < len(src.Content)-1; i += 2 {
		srcKey := src.Content[i]
		srcVal := src.Content[i+1]

		childPath := srcKey.Value
		if keyPath != "" {
			childPath = keyPath + "." + srcKey.Value
		}

		if valIdx, exists := dstKeyIdx[srcKey.Value]; exists {
			dstVal := dst.Content[valIdx]
			if dstVal.Kind == srcVal.Kind {
				mergeNodes(dstVal, srcVal, childPath, changedPaths)
			} else {
				// Kinds differ – replace and record.
				if !nodesEqual(dstVal, srcVal) {
					*changedPaths = append(*changedPaths, childPath)
				}
				dst.Content[valIdx] = srcVal
			}
		} else {
			// New key not present in dst – add it.
			*changedPaths = append(*changedPaths, childPath)
			dst.Content = append(dst.Content, srcKey, srcVal)
		}
	}
}

// nodesEqual performs a deep structural equality check on two yaml.Nodes.
func nodesEqual(a, b *yaml.Node) bool {
	if a.Kind != b.Kind || a.Value != b.Value {
		return false
	}
	if len(a.Content) != len(b.Content) {
		return false
	}
	for i := range a.Content {
		if !nodesEqual(a.Content[i], b.Content[i]) {
			return false
		}
	}
	return true
}

// stripCommonIndent detects the minimum leading whitespace across all non-empty lines,
// strips that many characters from the beginning of every line, and returns the result
// together with the stripped indent width.
func stripCommonIndent(content []byte) ([]byte, int) {
	lines := strings.Split(string(content), "\n")
	minIndent := -1

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent < 0 || n < minIndent {
			minIndent = n
		}
	}

	if minIndent <= 0 {
		return content, 0
	}

	stripped := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			stripped[i] = ""
		} else if len(line) >= minIndent {
			stripped[i] = line[minIndent:]
		} else {
			stripped[i] = line
		}
	}

	return []byte(strings.Join(stripped, "\n")), minIndent
}

// addIndent prepends indent spaces to each non-empty line of content.
func addIndent(content []byte, indent int) []byte {
	prefix := strings.Repeat(" ", indent)
	lines := strings.Split(string(content), "\n")
	indented := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			indented[i] = ""
		} else {
			indented[i] = prefix + line
		}
	}
	return []byte(strings.Join(indented, "\n"))
}
