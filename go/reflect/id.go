package reflect

var (
	FieldNameRequestId = "RequestId"
	FieldNameSessionId = "SessionId"
)

type IdRetriever interface {
	GetId() string
}

func RetrieveId(req interface{}, key string) string {
	if req == nil {
		return ""
	}
	// if defined GetRequestId()string method, use it get requestId
	if req, ok := req.(IdRetriever); ok {
		if req != nil {
			return req.GetId()
		}
	}
	return RetrieveStructField(req, key)
}

func TrySetId(req interface{}, key, id string) {
	TrySetStructFiled(req, key, id)
}
