package code

import (
	"fmt"
	"strings"
)

func (x *CgoError) Error() string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("cgo (%d)", x.GetErrorCode()))
	if x.GetErrorMessage() != "" {
		msg.WriteString(fmt.Sprintf(", %s", x.GetErrorMessage()))
	}
	if x.GetSdkErrorCode() != 0 {
		msg.WriteString(fmt.Sprintf(": sdk(%d)", x.GetSdkErrorCode()))
		if x.GetSdkErrorMessage() != "" {
			msg.WriteString(fmt.Sprintf(", %s", x.GetSdkErrorMessage()))
		}
	}
	return msg.String()
}
