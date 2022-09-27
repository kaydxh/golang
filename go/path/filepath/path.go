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
package filepath

import (
	"path"
	"path/filepath"
)

func GetParentRelPath(filePath string) string {
	fileName := path.Base(filePath)
	parentDir := GetParentRelDir(filePath)
	return filepath.Join(parentDir, fileName)
}

func GetParentRelDir(filePath string) string {
	dir := path.Dir(filePath)
	return path.Base(dir)
}

func LastChar(str string) uint8 {
	if str == "" {
		// char '0' means dec 48, 0 means null
		return 0
	}

	return str[len(str)-1]
}

// join paths, keep relativePath suffix
func JoinPaths(rootPath, relativePath string) (string, error) {
	absolutePath, err := filepath.Abs(rootPath)
	if err != nil {
		return "", err
	}

	if relativePath == "" {
		return absolutePath, nil
	}

	finalPath := path.Join(absolutePath, relativePath)
	if LastChar(relativePath) == '/' && LastChar(finalPath) != '/' {
		return finalPath + "/", nil
	}

	return finalPath, nil
}

// CanonicalizePath turns path into an absolute path without symlinks.
func CanonicalizePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	activepath, err := filepath.EvalSymlinks(path)
	if err == nil {
		return activepath, nil
	}

	/*
		// Get a better error if we have an invalid path
		if pathErr, ok := err.(*os.PathError); ok {
			err = errors.Wrap(pathErr.Err, pathErr.Path)
		}
	*/

	return path, nil
}
