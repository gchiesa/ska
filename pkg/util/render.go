package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aquasecurity/table"
	"github.com/jszwec/csvutil"
)

func RenderJSON(v any) ([]byte, error) {
	output, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func RenderCSV(v any) ([]byte, error) {
	outputCsv, err := csvutil.Marshal(v)
	if err != nil {
		return nil, err
	}
	return outputCsv, nil
}

type TableOptFunc func(t *table.Table)

func RenderTable(v any, opt ...TableOptFunc) ([]byte, error) {
	csvData, err := RenderCSV(v)
	if err != nil {
		return nil, err
	}
	csvBuf := bytes.NewBuffer(csvData)
	var outBuf bytes.Buffer
	t := table.New(&outBuf)
	t.SetAutoMerge(true)
	for _, opt := range opt {
		opt(t)
	}
	if err := t.LoadCSV(csvBuf, true); err != nil {
		return nil, err
	}
	t.Render()
	return outBuf.Bytes(), nil
}

func WithStyleMarkdown() TableOptFunc {
	return func(t *table.Table) {
		t.SetAutoMerge(false)
		t.SetDividers(table.MarkdownDividers)
	}
}

func RenderMarkdown(v any) ([]byte, error) {
	return RenderTable(v, WithStyleMarkdown())
}

func RenderWithOutputFormat(v any, outputFormat string) ([]byte, error) {
	var out []byte
	var err error
	switch outputFormat {
	case "json":
		out, err = RenderJSON(v)
	case "csv":
		out, err = RenderCSV(v)
	case "markdown":
		out, err = RenderMarkdown(v)
	case "table":
		out, err = RenderTable(v)
	default:
		return nil, fmt.Errorf("unknown output format: %s", outputFormat)
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}
