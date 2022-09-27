/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
