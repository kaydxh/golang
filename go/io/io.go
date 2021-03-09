package io

import (
	"bytes"
	"os"

	os_ "github.com/kaydxh/golang/go/os"
)

func WriteLines(filePath string, words []string, appended bool) (err error) {

	var file *os.File

	if !appended {
		file, err = os_.OpenAll(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
	} else {
		file, err = os_.OpenAll(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
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
