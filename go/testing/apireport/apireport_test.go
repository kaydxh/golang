package apireport

import (
	"strings"
	"testing"
)

func TestExitCode(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Status: "pass"})
	if c.ExitCode() != 0 {
		t.Errorf("all pass → 0, got %d", c.ExitCode())
	}
	c.Add(Result{Status: "fail"})
	if c.ExitCode() != 1 {
		t.Errorf("has fail → 1, got %d", c.ExitCode())
	}
}

func TestJUnitRender_ContainsCounts(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "v1", Name: "CreateUser", Status: "pass", DurationMs: 10})
	c.Add(Result{Suite: "v1", Name: "BadCase", Status: "fail", Reason: "boom", DurationMs: 5})
	xml := c.RenderJUnit()

	if !strings.Contains(xml, `tests="2"`) {
		t.Errorf("junit missing tests count: %s", xml)
	}
	if !strings.Contains(xml, `failures="1"`) {
		t.Errorf("junit missing failures count: %s", xml)
	}
	if !strings.Contains(xml, "BadCase") {
		t.Errorf("junit missing case name: %s", xml)
	}
}

func TestJSONRender_Valid(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "web", Name: "GetTag", Status: "pass"})
	js := c.RenderJSON()
	if !strings.Contains(js, "\"GetTag\"") {
		t.Errorf("json missing case: %s", js)
	}
}
