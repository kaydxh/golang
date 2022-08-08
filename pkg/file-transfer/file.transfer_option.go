package filetransfer

// A FileTransferOption sets options.
type FileTransferOption interface {
	apply(*FileTransfer)
}

// EmptyFileTransferUrlOption does not alter the FileTransferuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyFileTransferOption struct{}

func (EmptyFileTransferOption) apply(*FileTransfer) {}

// FileTransferOptionFunc wraps a function that modifies FileTransfer into an
// implementation of the FileTransferOption interface.
type FileTransferOptionFunc func(*FileTransfer)

func (f FileTransferOptionFunc) apply(do *FileTransfer) {
	f(do)
}

// sample code for option, default for nothing to change
func _FileTransferOptionWithDefault() FileTransferOption {
	return FileTransferOptionFunc(func(*FileTransfer) {
		// nothing to change
	})
}
func (o *FileTransfer) ApplyOptions(options ...FileTransferOption) *FileTransfer {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
