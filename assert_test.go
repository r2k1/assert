package assert

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type mockT struct {
	*testing.T
	logs   bytes.Buffer
	failed bool
}

func NewMockT() *mockT {
	return &mockT{T: &testing.T{}}
}

func (tb *mockT) Log(args ...interface{}) {
	tb.logs.WriteString(args[0].(string))
}

func (tb *mockT) Logf(format string, args ...interface{}) {
	tb.Log(fmt.Sprintf(format, args...))
}

func (tb *mockT) Logs() string {
	return tb.logs.String()
}

func (tb *mockT) Fail() {
	tb.failed = true
	tb.T.Fail()
}

func assertEqual[T any](t *testing.T, actual, expected T, expectedError string) {
	t.Helper()
	mockT := NewMockT()
	Equal(mockT, actual, expected)

	if expectedError == "" && mockT.failed {
		t.Log("shouldn't fail")
		t.Fail()
		return
	}

	if expectedError != "" && !mockT.failed {
		t.Log("should fail")
		t.Fail()
		return
	}

	actualError := strings.ReplaceAll(mockT.Logs(), "\u00a0", " ")
	if expectedError != actualError {
		t.Log("unexpected log message", actualError)
		t.Fail()
		Equal(t, actualError, expectedError, "mismatched t.Log output")
	}
}

type testStruct struct {
	Name string
}

var anError = errors.New("test error")
var nilError error

func TestMain(m *testing.M) {
	m.Run()
}

func setNoColor(t *testing.T) {
	noColor = true
	t.Cleanup(func() {
		noColor = false
	})
}

func TestEqual(t *testing.T) {
	setNoColor(t)
	assertEqual(t, 1, 1, "")
	assertEqual(t, 1.1, 1.1, "")
	assertEqual(t, true, true, "")
	assertEqual(t, false, false, "")
	assertEqual(t, "test", "test", "")
	assertEqual(t, []byte("test"), []byte("test"), "")
	assertEqual(t, []int{1, 2, 3}, []int{1, 2, 3}, "")
	assertEqual(t, map[string]int{"one": 1}, map[string]int{"one": 1}, "")
	assertEqual(t, testStruct{Name: "John"}, testStruct{Name: "John"}, "")
	assertEqual(t, &testStruct{Name: "John"}, &testStruct{Name: "John"}, "")
	var anyA, anyB any
	anyA = 1
	anyB = 1
	assertEqual(t, anyA, anyB, "")
	anyA = nil
	anyB = nil
	assertEqual(t, anyA, anyB, "")
	assertEqual(t, 1, 2, `
Error:    Not equal
Expected: 2
Actual:   1`)

	assertEqual(t, testStruct{Name: "John1"}, testStruct{Name: "John2"}, `
Error:    Not equal
Expected: assert.testStruct{Name:"John2"}
Actual:   assert.testStruct{Name:"John1"}
Diff:     assert.testStruct{
          - 	Name: "John1",
          + 	Name: "John2",
            }`)

	assertEqual(t, anError, anError, ``)
	assertEqual(t, errors.New("test1"), errors.New("test1"), ``)
	assertEqual(t, errors.New("test1"), errors.New("test2"), `
Error:    Not equal
Expected: &errors.errorString{s:"test2"}
Actual:   &errors.errorString{s:"test1"}`)
	assertEqual(t, anError, nilError, `
Error:    Not equal
Expected: <nil>
Actual:   &errors.errorString{s:"test error"}`)
	assertEqual(t, nilError, anError, `
Error:    Not equal
Expected: &errors.errorString{s:"test error"}
Actual:   <nil>`)
	var nilMap map[string]int
	emptyMap := make(map[string]int)
	assertEqual(t, nilMap, emptyMap, `
Error:    Not equal
Expected: map[string]int{}
Actual:   map[string]int(nil)
Diff:     map[string]int(
          - 	nil,
          + 	{},
            )`)
	assertEqual(t, nilMap, map[string]int{"val": 1}, `
Error:    Not equal
Expected: map[string]int{"val":1}
Actual:   map[string]int(nil)
Diff:     map[string]int(
          - 	nil,
          + 	{"val": 1},
            )`)
	var nilSlice []string
	emptySlice := make([]string, 0)
	assertEqual(t, nilSlice, emptySlice, `
Error:    Not equal
Expected: []string{}
Actual:   []string(nil)
Diff:     []string(
          - 	nil,
          + 	{},
            )`)
}

func assertNotEqual[T any](t *testing.T, actual, expected T, expectedError string) {
	t.Helper()
	mockT := NewMockT()
	NotEqual(mockT, actual, expected)

	if expectedError == "" && mockT.failed {
		t.Log("shouldn't fail")
		t.Fail()
		return
	}

	if expectedError != "" && !mockT.failed {
		t.Log("should fail")
		t.Fail()
		return
	}

	actualError := strings.ReplaceAll(mockT.Logs(), "\u00a0", " ")
	if expectedError != actualError {
		t.Log("unexpected log message", actualError)
		t.Fail()
		Equal(t, actualError, expectedError, "mismatched t.Log output")
	}
}

func TestNotEqual(t *testing.T) {
	setNoColor(t)
	assertNotEqual(t, 1, 2, "")
	assertNotEqual(t, 1.1, 2.2, "")
	assertNotEqual(t, true, false, "")
	assertNotEqual(t, false, true, "")
	assertNotEqual(t, "test", "test1", "")
	assertNotEqual(t, []byte("test"), []byte("test1"), "")
	assertNotEqual(t, []int{1, 2, 3}, []int{4, 5, 6}, "")
	assertNotEqual(t, map[string]int{"one": 1}, map[string]int{"two": 2}, "")
	assertNotEqual(t, testStruct{Name: "John"}, testStruct{Name: "Doe"}, "")
	assertNotEqual(t, &testStruct{Name: "John"}, &testStruct{Name: "Doe"}, "")
	var anyA, anyB any
	anyA = 1
	anyB = 2
	assertNotEqual(t, anyA, anyB, "")
	anyA = nil
	anyB = 1
	assertNotEqual(t, anyA, anyB, "")
	assertNotEqual(t, 1, 1, `
Error:    Values should not be equal
Value:    1`)

	assertNotEqual(t, testStruct{Name: "John"}, testStruct{Name: "John"}, `
Error:    Values should not be equal
Value:    assert.testStruct{Name:"John"}`)

	assertNotEqual(t, anError, nilError, ``)
	assertNotEqual(t, errors.New("test1"), errors.New("test2"), ``)
	assertNotEqual(t, errors.New("test1"), errors.New("test1"), `
Error:    Values should not be equal
Value:    &errors.errorString{s:"test1"}`)
	assertNotEqual(t, anError, anError, `
Error:    Values should not be equal
Value:    &errors.errorString{s:"test error"}`)
	var nilMap map[string]int
	emptyMap := make(map[string]int)
	assertNotEqual(t, nilMap, map[string]int{"val": 1}, "")
	assertNotEqual(t, emptyMap, map[string]int{"val": 1}, "")
	assertNotEqual(t, emptyMap, nilMap, "")
	assertNotEqual(t, nilMap, nilMap, `
Error:    Values should not be equal
Value:    map[string]int(nil)`)

	var nilSlice []string
	emptySlice := make([]string, 0)
	assertNotEqual(t, nilSlice, nilSlice, `
Error:    Values should not be equal
Value:    []string(nil)`)
	assertNotEqual(t, emptySlice, emptySlice, `
Error:    Values should not be equal
Value:    []string{}`)
	assertNotEqual(t, nilSlice, emptySlice, ``)
	assertNotEqual(t, nilError, nil, `
Error:    Values should not be equal
Value:    <nil>`)
	assertNotEqual(t, anError, nil, ``)
}
