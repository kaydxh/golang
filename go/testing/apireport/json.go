package apireport

import "encoding/json"

// RenderJSON 渲染 JSON 报告。
func (c *Collector) RenderJSON() string {
	pass, fail, skip := c.counts()
	out := map[string]any{
		"total":   len(c.Results),
		"passed":  pass,
		"failed":  fail,
		"skipped": skip,
		"results": c.Results,
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b)
}
