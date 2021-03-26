package io

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	os_ "github.com/kaydxh/golang/go/os"
)

// ReadFileLines read line from file divide with \n
func ReadFileLines(filepath string) ([]string, error) {
	file, err := os_.OpenAll(filepath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) { //文件已经结束
				break
			}

			return nil, err
		}
		lines = append(lines, string(line))
	}

	return lines, nil
}

// ReadLines read line from byteArray divide with \n
func ReadLines(byteArray []byte) []string {
	var lines []string
	index := 0
	readIndex := 0
	for ; readIndex < len(byteArray); index++ {
		line, n := ReadLineAt(readIndex, byteArray)
		readIndex = n

		lines = append(lines, string(line))
	}

	return lines
}

func ReadLineAt(readIndex int, byteArray []byte) ([]byte, int) {
	currentReadIndex := readIndex

	// consume left spaces
	for currentReadIndex < len(byteArray) {
		if byteArray[currentReadIndex] == ' ' {
			currentReadIndex++
		} else {
			break
		}
	}

	// leftTrimIndex stores the left index of the line after the line is left-trimmed
	leftTrimIndex := currentReadIndex

	// rightTrimIndex stores the right index of the line after the line is right-trimmed
	// it is set to -1 since the correct value has not yet been determined.
	rightTrimIndex := -1

	for ; currentReadIndex < len(byteArray); currentReadIndex++ {
		if byteArray[currentReadIndex] == ' ' {
			// set rightTrimIndex
			if rightTrimIndex == -1 {
				rightTrimIndex = currentReadIndex
			}
		} else if (byteArray[currentReadIndex] == '\n') || (currentReadIndex == (len(byteArray) - 1)) {
			// end of line or byte buffer is reached
			if currentReadIndex <= leftTrimIndex {
				return nil, currentReadIndex + 1
			}
			// set the rightTrimIndex
			if rightTrimIndex == -1 {
				rightTrimIndex = currentReadIndex
				if currentReadIndex == (len(byteArray)-1) && (byteArray[currentReadIndex] != '\n') {
					// ensure that the last character is part of the returned string,
					// unless the last character is '\n'
					rightTrimIndex = currentReadIndex + 1
				}
			}
			// Avoid unnecessary allocation.
			return byteArray[leftTrimIndex:rightTrimIndex], currentReadIndex + 1
		} else {
			// unset rightTrimIndex
			rightTrimIndex = -1
		}
	}

	return nil, currentReadIndex
}

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

	buf := bytes.NewBuffer(nil)
	err = WriteLine(buf, words...)
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

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

func WriteBytesAt(filepath string, bytes []byte, offset int64) error {
	file, err := os_.OpenAll(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
