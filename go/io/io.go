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
package io

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
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
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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

//ReadFile read data from file
func ReadFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, err
}

func WriteFile(filePath string, content []byte, appended bool) error {
	file, err := os_.OpenFile(filePath, appended)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	_, err = buf.Write(content)
	if err != nil {
		return err
	}

	if err := buf.Flush(); err != nil {
		return err
	}

	return nil
}

// WriteLine join all line to file.
func WriteFileLines(filePath string, lines []string, appended bool) (err error) {

	file, err := os_.OpenFile(filePath, appended)
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
	file, err := os_.OpenFile(filePath, appended)
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
	file, err := os_.OpenAll(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = file.WriteAt(bytes, n)
	if err != nil {
		return err
	}

	return nil
}

func WriteReader(filepath string, r io.Reader) error {
	file, err := os_.OpenAll(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	//	file, err := os_.OpenFile(filepath, true)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func WriteReaderAt(filepath string, r io.Reader, offset, length int64) error {
	file, err := os_.OpenAll(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	//	file, err := os_.OpenFile(filepath, true)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	var total int64

	for total < length {
		nr, err := r.Read(buf)
		if err == nil || err == io.EOF {
			if nr > 0 {
				n, err := file.Seek(offset, io.SeekStart)
				if err != nil {
					return err
				}
				_, err = file.WriteAt(buf[:nr], n)
				if err != nil {
					return err
				}

				offset += int64(nr)
				total += int64(nr)
			}

			if err == io.EOF {
				break
			}

		} else {
			return err
		}

	}

	return nil
}
