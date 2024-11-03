package assert

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mattn/go-isatty"
)

var noColor bool

func init() {
	isTerminal := !isatty.IsTerminal(os.Stdout.Fd())
	// prefer https://no-color.org (with any value)
	var envNoColor bool
	if envColorStr := os.Getenv("NO_COLOR"); envColorStr != "" {
		envNoColor, _ = strconv.ParseBool(envColorStr)
	}
	noColor = isTerminal && !envNoColor
}

type color string

// ANSI color codes
const (
	colorRed    color  = "\033[31m"
	colorGreen  color  = "\033[32m"
	colorYellow color  = "\033[33m"
	endColor    string = "\033[0m"
)

func colorize(c color, s string) string {
	if noColor {
		return s
	}
	return string(c) + s + endColor
}

// Equal checks if two objects of any type are equal and reports an error if they are not.
func Equal[T any](t testing.TB, actual, expected T, msgAndArgs ...any) bool {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	result := [][2]string{
		{"Error", colorize(colorRed, "Not equal")},
	}
	extra := messageFromMsgAndArgs(msgAndArgs...)
	if extra != "" {
		result = append(result, [2]string{"Message", colorize(colorYellow, extra)})
	}
	result = append(result, [2]string{"Expected", colorize(colorGreen, fmt.Sprintf("%#v", expected))})
	result = append(result, [2]string{"Actual", colorize(colorRed, fmt.Sprintf("%#v", actual))})
	diffS := diff(expected, actual)
	if diffS != "" {
		result = append(result, [2]string{"Diff", diffS})
	}
	t.Log(sprintList(result))
	t.Fail()
	return false
}

// NotEqual checks if two objects of any type are not equal and reports an error if they are.
func NotEqual[T any](t testing.TB, actual, expected T, msgAndArgs ...any) bool {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		return true
	}
	result := [][2]string{
		{"Error", "Values should not be equal"},
	}
	extra := messageFromMsgAndArgs(msgAndArgs...)
	if extra != "" {
		result = append(result, [2]string{"Message", colorize(colorYellow, extra)})
	}
	result = append(result, [2]string{"Value", colorize(colorRed, fmt.Sprintf("%#v", expected))})
	t.Log(sprintList(result))
	t.Fail()

	return false
}

// MustEqual asserts that two objects are equal and stops test execution if they are not.
func MustEqual[T any](t testing.TB, actual, expected T) {
	t.Helper()
	if !Equal(t, actual, expected) {
		t.FailNow()
	}
}

// MustNotEqual asserts that two objects are not equal and stops test execution if they are not.
func MustNotEqual[T any](t testing.TB, actual, expected T) {
	t.Helper()
	if !NotEqual(t, actual, expected) {
		t.FailNow()
	}
}

func diff[T any](actual, expected T) string {
	et := reflect.TypeOf(expected)
	at := reflect.TypeOf(actual)
	if et == nil || at == nil {
		return ""
	}
	ek := et.Kind()
	if ek != reflect.Struct && ek != reflect.Map && ek != reflect.Slice && ek != reflect.Array && ek != reflect.String {
		return ""
	}
	diff := cmp.Diff(expected, actual)
	return colorizeDiff(diff)
}

func colorizeDiff(diff string) string {
	var output strings.Builder
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "-"):
			output.WriteString(colorize(colorRed, line+"\n"))
		case strings.HasPrefix(line, "+"):
			output.WriteString(colorize(colorGreen, line+"\n"))
		default:
			output.WriteString(line + "\n")
		}
	}
	return strings.TrimSpace(output.String())
}

func sprintList(input [][2]string) string {
	if len(input) == 0 {
		return ""
	}
	var output strings.Builder
	firstPartLength := 8
	formatSecondPart := func(input string) string {
		if input == "" {
			return ""
		}
		lines := strings.Split(input, "\n")
		if len(lines) == 1 {
			return input
		}
		for i := 1; i < len(lines); i++ {
			lines[i] = fmt.Sprintf("%*s  %s", firstPartLength, "", lines[i])
		}
		return strings.Join(lines, "\n")
	}
	for _, item := range input {
		output.WriteString(fmt.Sprintf("\n%s:%s %s", item[0], strings.Repeat(" ", firstPartLength-len(item[0])), formatSecondPart(item[1])))
	}
	return output.String()
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if len(msgAndArgs) == 1 {
		if msg, ok := msgAndArgs[0].(string); ok {
			return msg
		}
		return fmt.Sprintf("%+v", msgAndArgs[0])
	}
	return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
}
