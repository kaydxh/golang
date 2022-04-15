package logrus

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	runtime_ "github.com/kaydxh/golang/go/runtime"
	"github.com/sirupsen/logrus"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

var (
	baseTimestamp time.Time
	pid           int
)

func init() {
	baseTimestamp = time.Now()
	pid = os.Getpid()
}

/*
https://medium.com/technical-tips/google-log-glog-output-format-7eb31b3f0ce5
[IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg
IWEF — Log Levels, I for INFO, W for WARNING, E for ERROR and `F` for FATAL.
yyyymmdd — Year, Month and Date.
hh:mm:ss.uuuuuu — Hours, Minutes, Seconds and Microseconds.
threadid — PID/TID of the process/thread.
file:line — File name and line number.
msg — Actual user-specified log message.
*/

// GlogFormatter formats logs into text
type GlogFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	terminalInitOnce sync.Once

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int

	// Disables the glog style ：[IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg msg...
	// replace with ：[IWEF] [yyyymmdd] [hh:mm:ss.uuuuuu] [threadid] [file:line] msg msg...
	//EnablePrettyLog   bool
	EnableGoroutineId bool
}

func (f *GlogFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
	// Get the max length of the level text
	for _, level := range logrus.AllLevels {
		levelTextLength := utf8.RuneCount([]byte(level.String()))
		if levelTextLength > f.levelTextMaxLength {
			f.levelTextMaxLength = levelTextLength
		}
	}
}

func (f *GlogFormatter) isColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		switch force, ok := os.LookupEnv("CLICOLOR_FORCE"); {
		case ok && force != "0":
			isColored = true
		case ok && force == "0", os.Getenv("CLICOLOR") == "0":
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}

// Format renders a single log entry
func (f *GlogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	fixedKeys := make([]string, 0, 4+len(data))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyMsg))
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		} else {
			if !f.isColored() {
				fixedKeys = append(fixedKeys, keys...)
				f.SortingFunc(fixedKeys)
			} else {
				f.SortingFunc(keys)
			}
		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	levelText := f.level(entry)
	b.Write(f.formatHeader(entry, levelText))

	for _, key := range fixedKeys {
		var value interface{}
		switch {
		case key == f.FieldMap.resolve(FieldKeyMsg):
			value = entry.Message
			f.appendValue(b, value)
			//continue means msg not need call f.appendKeyValue(b, key, value)
			//or will duplicate write message
			continue
		case key == f.FieldMap.resolve(fieldKey(logrus.ErrorKey)):
			value = data[logrus.ErrorKey]
		default:
			value = data[key]
		}
		key = fmt.Sprintf("%s", key)
		f.appendKeyValue(b, key, value)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *GlogFormatter) level(entry *logrus.Entry) string {
	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation && !f.PadLevelText {
		levelText = levelText[0:4]
	}
	if f.PadLevelText {
		// Generates the format string used in the next line, for example "%-6s" or "%-7s".
		// Based on the max level text length.
		formatString := "%-" + strconv.Itoa(f.levelTextMaxLength) + "s"
		// Formats the level text by appending spaces up to the max length, for example:
		// 	- "INFO   "
		//	- "WARNING"
		levelText = fmt.Sprintf(formatString, levelText)
	}

	return levelText
}

func (f *GlogFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *GlogFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *GlogFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

// Log line format: [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg
func (f *GlogFormatter) formatHeader(entry *logrus.Entry, levelText string) []byte {

	var buf bytes.Buffer
	switch {
	case f.DisableTimestamp:
		buf.WriteString(fmt.Sprintf("[%s]", levelText))
	case !f.FullTimestamp:
		buf.WriteString(fmt.Sprintf("[%s] [%04d]", levelText, int(entry.Time.Sub(baseTimestamp)/time.Second)))
	default:
		buf.WriteString(fmt.Sprintf("[%s] [%s]", levelText,
			entry.Time.Format(f.TimestampFormat),
		))
	}

	if f.EnableGoroutineId {
		buf.WriteString(fmt.Sprintf(" [%d]", runtime_.GoroutineID()))
	} else {
		//use pid instead of goroutine id
		buf.WriteString(fmt.Sprintf(" [%d]", pid))
	}

	var (
		function string
		fileline string = "???:-1"
	)

	if entry.HasCaller() {
		if f.CallerPrettyfier != nil {
			function, fileline = f.CallerPrettyfier(entry.Caller)
		} else {
			function = entry.Caller.Function
			fileline = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}

		buf.WriteString(fmt.Sprintf(" [%s](%s)", fileline, function))
	}

	//split head and body
	buf.WriteString(" ")

	return buf.Bytes()
}
