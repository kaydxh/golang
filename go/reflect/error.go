package reflect

import "fmt"

// ErrUnsupportedType ...
type ErrUnsupportedType struct {
	valueType string
}

// NewErrUnsupportedType returns new ErrUnsupportedType
func NewErrUnsupportedType(valueType string) ErrUnsupportedType {
	return ErrUnsupportedType{valueType}
}

// Error method so we implement the error interface
func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("%v is not one of supported types", e.valueType)
}
