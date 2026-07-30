package apireport

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// RenderJUnit 渲染 JUnit XML 报告。
func (c *Collector) RenderJUnit() string {
	_, fail, skip := c.counts()
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fmt.Fprintf(&sb, `<testsuite name="palmtest" tests="%d" failures="%d" skipped="%d">`+"\n",
		len(c.Results), fail, skip)
	for _, r := range c.Results {
		name := xmlEscape(r.Suite + "." + r.Name)
		switch r.Status {
		case "fail":
			fmt.Fprintf(&sb, `  <testcase name="%s" time="%.3f"><failure>%s</failure></testcase>`+"\n",
				name, float64(r.DurationMs)/1000, xmlEscape(r.Reason))
		case "skip":
			fmt.Fprintf(&sb, `  <testcase name="%s"><skipped/></testcase>`+"\n", name)
		default:
			fmt.Fprintf(&sb, `  <testcase name="%s" time="%.3f"/>`+"\n", name, float64(r.DurationMs)/1000)
		}
	}
	sb.WriteString(`</testsuite>` + "\n")
	return sb.String()
}

func xmlEscape(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
