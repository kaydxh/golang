package filetransfer

import "time"

func WithDownloadTimeout(downloadTimeout time.Duration) FileTransferOption {
	return FileTransferOptionFunc(func(r *FileTransfer) {
		r.opts.downloadTimeout = downloadTimeout
	})
}
