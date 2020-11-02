package io

import "os"

//copyFileContentes copies the contents of the file named src to the file named dst
//The destination file will be created if it does not alreay exist. If the destination
//file exists, all it's contents will be replaced by the contents of the source file
func copyFileContents(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := dstFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return
	}

	err = dstFile.Sync()
	return
}
