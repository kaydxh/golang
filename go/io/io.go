package io

import (
	"bufio"
	"bytes"

	os_ "github.com/kaydxh/golang/go/os"
)

// WriteLine join all line to file.
func WriteFileLines(filePath string, lines []string, appended bool) (err error) {

	file, err := os_.OpenAllAt(filePath, appended)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	for i := range lines {
		if _, err = buf.WriteString(lines[i] + "\n"); err != nil {
			return err
		}
	}

	if err := buf.Flush(); err != nil {
		return err
	}

	return nil
}

// WriteLine join all words with spaces, terminate with newline and
// write to file.
func WriteFileLine(filePath string, words []string, appended bool) (err error) {
	file, err := os_.OpenAllAt(filePath, appended)
	if err != nil {
		return err
	}
	defer file.Close()

	//buf := bufio.NewWriter(file)
	buf := bytes.NewBuffer(nil)
	err = WriteLine(buf, words...)
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	/*
		for i := range lines {
			if _, err = buf.WriteString(lines[i] + "\n"); err != nil {
				return err
			}
		}
	*/

	/*
		if err := buf.Flush(); err != nil {
			return err
		}
	*/

	return nil
}

// WriteLine join all words with spaces, terminate with newline and write to buff.
func WriteLine(buf *bytes.Buffer, words ...string) error {
	// We avoid strings.Join for performance reasons.
	for i := range words {
		_, err := buf.WriteString(words[i])
		if err != nil {
			return err
		}
		if i < len(words)-1 {
			err = buf.WriteByte(' ')
		} else {
			err = buf.WriteByte('\n')
		}
		if err != nil {
			return err
		}

	}

	return nil
}

// WriteBytesLine write bytes to buffer, terminate with newline
func WriteBytesLine(buf *bytes.Buffer, bytes []byte) error {
	_, err := buf.Write(bytes)
	if err != nil {
		return err
	}

	err = buf.WriteByte('\n')
	if err != nil {
		return err
	}

	return nil
}
