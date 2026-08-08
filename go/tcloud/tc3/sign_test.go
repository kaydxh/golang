package tc3

import "testing"

// TestAuthorization_Golden 用固定输入校验 TC3-HMAC-SHA256 签名与开源社区 CLI
// (palm_openapi_community_cli.build_authorization) 输出完全一致。
//
// 该期望值由以下脚本生成（Python 参考实现）：
//
//	import palm_openapi_community_cli as cli
//	cli.build_authorization(
//	    "AKID_test_secret_id_12345", "test_secret_key_67890", 223,
//	    "CreateAccessToken", "fixed_nonce_for_test_abcdef", 1700000000,
//	    "open.intl.palm.tencent.com",
//	    `{"AppId": 223, "GrantType": "client_credential", "SecretId": "AKID_test_secret_id_12345", "SecretKeyHash": "deadbeef", "UserId": ""}`,
//	)["Authorization"]
func TestAuthorization_Golden(t *testing.T) {
	payload := `{"AppId": 223, "GrantType": "client_credential", "SecretId": "AKID_test_secret_id_12345", "SecretKeyHash": "deadbeef", "UserId": ""}`
	got := Authorization(Options{
		SecretID:      "AKID_test_secret_id_12345",
		SecretKey:     "test_secret_key_67890",
		Service:       "palm",
		Host:          "open.intl.palm.tencent.com",
		Timestamp:     1700000000,
		Payload:       payload,
		SignedHeaders: []string{"content-type", "host", "x-palm-appid", "x-tc-nonce", "x-tc-timestamp"},
		HeaderValues: map[string]string{
			"content-type":   "application/json",
			"host":           "open.intl.palm.tencent.com",
			"x-palm-appid":   "223",
			"x-tc-nonce":     "fixed_nonce_for_test_abcdef",
			"x-tc-timestamp": "1700000000",
		},
	})
	want := "TC3-HMAC-SHA256 Credential=AKID_test_secret_id_12345/2023-11-14/palm/tc3_request, SignedHeaders=content-type;host;x-palm-appid;x-tc-nonce;x-tc-timestamp, Signature=eb939c307be85de81ed91a0eeee18d897d575a7e584ab5253a4080df355b0d71"
	if got != want {
		t.Errorf("Authorization mismatch:\n got: %s\nwant: %s", got, want)
	}
}

// TestAuthorization_Deterministic 相同输入应产生相同签名；任一参数变化签名应不同。
func TestAuthorization_Deterministic(t *testing.T) {
	opts := Options{
		SecretID:      "id",
		SecretKey:     "key",
		Service:       "palm",
		Host:          "h.example.com",
		Timestamp:     1700000000,
		Payload:       `{}`,
		SignedHeaders: []string{"content-type", "host"},
		HeaderValues: map[string]string{
			"content-type": "application/json",
			"host":         "h.example.com",
		},
	}
	a := Authorization(opts)
	b := Authorization(opts)
	if a != b {
		t.Fatalf("same inputs produced different signatures:\n %s\n %s", a, b)
	}

	opts.Payload = `{"x":1}`
	if c := Authorization(opts); c == a {
		t.Fatalf("signature should change when payload changes")
	}
}
