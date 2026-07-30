package apiassert

import "testing"

func TestCode(t *testing.T) {
	if !Code(200, 200).OK {
		t.Error("200==200 should pass")
	}
	if Code(500, 200).OK {
		t.Error("500!=200 should fail")
	}
}

func TestBodyContains(t *testing.T) {
	body := `{"Response":{"UserTagId":"tag_1"}}`
	if !BodyContains(body, "UserTagId").OK {
		t.Error("should contain UserTagId")
	}
	if BodyContains(body, "NotThere").OK {
		t.Error("should not contain NotThere")
	}
}

func TestBodyNotContains(t *testing.T) {
	body := `{"list":["a"]}`
	if !BodyNotContains(body, "other_id").OK {
		t.Error("other_id absent should pass")
	}
	if BodyNotContains(body, "a").OK {
		t.Error("present token should fail not-contains")
	}
}

func TestJSONField_Equal(t *testing.T) {
	body := `{"Response":{"TotalCount":3}}`
	if !JSONField(body, "Response.TotalCount", "3").OK {
		t.Error("TotalCount==3 should pass")
	}
	if JSONField(body, "Response.TotalCount", "5").OK {
		t.Error("TotalCount!=5 should fail")
	}
}

func TestJSONField_NotEmpty(t *testing.T) {
	body := `{"Response":{"RequestId":"rid"}}`
	if !JSONFieldNotEmpty(body, "Response.RequestId").OK {
		t.Error("RequestId non-empty should pass")
	}
	if JSONFieldNotEmpty(body, "Response.Missing").OK {
		t.Error("missing field should fail not-empty")
	}
}
