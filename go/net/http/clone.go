package http

import (
	"mime/multipart"
	"net/http"
	"net/url"

	_ "unsafe" // for go:linkname
)

//go:linkname CloneURLValues net/http.cloneURLValues
func CloneURLValues(v url.Values) url.Values

//go:linkname CloneURL net/http.cloneURL
func CloneURL(u *url.URL) *url.URL

//go:linkname CloneMultipartForm net/http.cloneMultipartForm
func CloneMultipartForm(f *multipart.Form) *multipart.Form

//go:linkname CloneMultipartFileHeader net/http.cloneMultipartFileHeader
func CloneMultipartFileHeader(fh *multipart.FileHeader) *multipart.FileHeader

// CloneOrMakeHeader invokes Header.Clone but if the
// result is nil, it'll instead make and return a non-nil Header.
//go:linkname CloneOrMakeHeader net/http.cloneOrMakeHeader
func CloneOrMakeHeader(hdr http.Header) http.Header
