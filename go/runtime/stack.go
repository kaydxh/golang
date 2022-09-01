package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

func FormatStack() ([]byte, error) {
	debugStack := debug.Stack()
	return prettyStack{}.parse(debugStack)
}

type prettyStack struct {
}

func (s prettyStack) parse(debugStack []byte) ([]byte, error) {
	var err error
	useColor := true
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, " panic: ")
	fmt.Fprintf(buf, "\n")

	// process debug stack info
	stack := strings.Split(string(debugStack), "\n")
	lines := []string{}

	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(") {
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	// decorate
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, useColor, i)
		if err != nil {
			return nil, err
		}
	}

	for _, l := range lines {
		fmt.Fprintf(buf, "%s", l)
	}
	return buf.Bytes(), nil
}

func (s prettyStack) decorateLine(line string, useColor bool, num int) (string, error) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:") {
		return s.decorateSourceLine(line, useColor, num)
	} else if strings.HasSuffix(line, ")") {
		return s.decorateFuncCallLine(line, useColor, num)
	} else {
		if strings.HasPrefix(line, "\t") {
			return strings.Replace(line, "\t", "      ", 1), nil
		} else {
			return fmt.Sprintf("    %s\n", line), nil
		}
	}
}

func (s prettyStack) decorateFuncCallLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("not a func call line")
	}

	buf := &bytes.Buffer{}
	pkg := line[0:idx]
	// addr := line[idx:]
	method := ""

	if idx := strings.LastIndex(pkg, string(os.PathSeparator)); idx < 0 {
		if idx := strings.Index(pkg, "."); idx > 0 {
			method = pkg[idx:]
			pkg = pkg[0:idx]
		}
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0 : idx+1]
		if idx := strings.Index(method, "."); idx > 0 {
			pkg += method[0:idx]
			method = method[idx:]
		}
	}

	if num == 0 {
		fmt.Fprintf(buf, " -> ")
	} else {
		fmt.Fprintf(buf, "    ")
	}
	fmt.Fprintf(buf, "%s", pkg)
	fmt.Fprintf(buf, "%s\n", method)
	return buf.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("not a source line")
	}

	buf := &bytes.Buffer{}
	path := line[0 : idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0 : idx+1]
	file := path[idx+1:]

	idx = strings.Index(lineno, " ")
	if idx > 0 {
		lineno = lineno[0:idx]
	}

	if num == 1 {
		fmt.Fprintf(buf, " -> ")
	} else {
		fmt.Fprintf(buf, "    ")
	}
	fmt.Fprintf(buf, "%s", dir)
	fmt.Fprintf(buf, "%s", file)
	fmt.Fprintf(buf, "%s", lineno)

	if num == 1 {
		fmt.Fprintf(buf, "\n")
	}
	fmt.Fprintf(buf, "\n")

	return buf.String(), nil
}
