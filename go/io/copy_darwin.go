package io

import "os"

func CopyPath(srcPath, dstPath string, f os.FileInfo, copyMode Mode) error {
	return nil
}

func doCopyWithFileClone(srcFile, dstFile *os.File) error {
	return nil
}

func doCopyWithFileRange(srcFile, dstFile *os.File, fileinfo os.FileInfo) error {
	return nil
}
