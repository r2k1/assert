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

func (mock *mockT) Log(args ...any) {
	mock.logs.WriteString(args[0].(string))
}

func (mock *mockT) Logf(format string, args ...any) {
	mock.Log(fmt.Sprintf(format, args...))
}

func (mock *mockT) Logs() string {
	return mock.logs.String()
}

func (mock *mockT) Fail() {
	mock.failed = true
}

func (mock *mockT) Failed() bool {
	return mock.failed
}

func (mock *mockT) FailNow() {
	mock.failed = true
}

func (mock *mockT) AssertFailed(t *testing.T, expected string) {
	t.Helper()
	if expected == "" {
		if mock.failed {
			t.Logf("shouldn't fail")
			t.Fail()
		}
		return
	}
	if !mock.failed {
		t.Logf("should fail")
		t.Fail()
	}
	actual := strings.ReplaceAll(mock.Logs(), "\u00a0", " ")
	if expected != actual {
		Equal(t, expected, actual, "mismatched t.Log output")
	}

	return
}

type testStruct struct {
	Name string
}

var anError = errors.New("test error")
var nilError error

var nilMap map[string]int
var emptyMap = make(map[string]int)
var nilSlice []string
var emptySlice = make([]string, 0)
var anyA, anyB any

func TestMain(m *testing.M) {
	originalNoColor := noColor
	noColor = true
	defer func() { noColor = originalNoColor }()
	anyA = 1
	anyB = 1
	m.Run()
}

func setColor(t *testing.T) {
	noColor = false
	t.Cleanup(func() {
		noColor = true
	})
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name        string
		actual      any
		expected    any
		msgAndArgs  []any
		expectedMsg string
	}{
		// Successful tests
		{"integers_equal", 1, 1, nil, ""},
		{"floats_equal", 1.1, 1.1, nil, ""},
		{"booleans_equal", true, true, nil, ""},
		{"booleans_false_equal", false, false, nil, ""},
		{"strings_equal", "test", "test", nil, ""},
		{"bytes_equal", []byte("test"), []byte("test"), nil, ""},
		{"slices_equal", []int{1, 2, 3}, []int{1, 2, 3}, nil, ""},
		{"maps_equal", map[string]int{"one": 1}, map[string]int{"one": 1}, nil, ""},
		{"structs_equal", testStruct{Name: "John"}, testStruct{Name: "John"}, nil, ""},
		{"pointers_equal", &testStruct{Name: "John"}, &testStruct{Name: "John"}, nil, ""},
		{"errors_same_instance", anError, anError, nil, ""},
		{"errors_same_message", errors.New("test1"), errors.New("test1"), nil, ""},
		{"any_equal", anyA, anyB, nil, ""},

		// Failed tests
		{"integers_not_equal", 1, 2, nil, `
Error:    Not equal
Expected: 2
Actual:   1`},
		{"structs_not_equal", testStruct{Name: "John1"}, testStruct{Name: "John2"}, nil, `
Error:    Not equal
Expected: assert.testStruct{Name:"John2"}
Actual:   assert.testStruct{Name:"John1"}
Diff:     assert.testStruct{
          - 	Name: "John1",
          + 	Name: "John2",
            }`},
		{"errors_different", errors.New("test1"), errors.New("test2"), nil, `
Error:    Not equal
Expected: &errors.errorString{s:"test2"}
Actual:   &errors.errorString{s:"test1"}`},
		{"error_vs_nil", anError, nilError, nil, `
Error:    Not equal
Expected: <nil>
Actual:   &errors.errorString{s:"test error"}`},
		{"nil_vs_error", nilError, anError, nil, `
Error:    Not equal
Expected: &errors.errorString{s:"test error"}
Actual:   <nil>`},
		{"custom_message", 1, 2, []any{"custom message"}, `
Error:    Not equal
Message:  custom message
Expected: 2
Actual:   1`},
		{"custom_message_with_args", 1, 2, []any{"custom message %s", "with args"}, `
Error:    Not equal
Message:  custom message with args
Expected: 2
Actual:   1`},
		{"nil_map_vs_empty_map", nilMap, emptyMap, nil, `
Error:    Not equal
Expected: map[string]int{}
Actual:   map[string]int(nil)
Diff:     map[string]int(
          - 	nil,
          + 	{},
            )`},
		{"nil_map_vs_map_with_value", nilMap, map[string]int{"val": 1}, nil, `
Error:    Not equal
Expected: map[string]int{"val":1}
Actual:   map[string]int(nil)
Diff:     map[string]int(
          - 	nil,
          + 	{"val": 1},
            )`},
		{"nil_slice_vs_empty_slice", nilSlice, emptySlice, nil, `
Error:    Not equal
Expected: []string{}
Actual:   []string(nil)
Diff:     []string(
          - 	nil,
          + 	{},
            )`},
		{"non_string_msg", 1, 2, []any{1}, `
Error:    Not equal
Message:  1
Expected: 2
Actual:   1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := NewMockT()
			if tt.msgAndArgs == nil {
				Equal(mockT, tt.actual, tt.expected)
			} else {
				Equal(mockT, tt.actual, tt.expected, tt.msgAndArgs...)
			}
			mockT.AssertFailed(t, tt.expectedMsg)
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name        string
		actual      any
		expected    any
		expectedMsg string
	}{
		// Successful tests
		{"integers_different", 1, 2, ""},
		{"floats_different", 1.1, 2.2, ""},
		{"booleans_different", true, false, ""},
		{"strings_different", "test", "test1", ""},
		{"bytes_different", []byte("test"), []byte("test1"), ""},
		{"slices_different", []int{1, 2, 3}, []int{4, 5, 6}, ""},
		{"maps_different", map[string]int{"one": 1}, map[string]int{"two": 2}, ""},
		{"structs_different", testStruct{Name: "John"}, testStruct{Name: "Doe"}, ""},
		{"pointers_different", &testStruct{Name: "John"}, &testStruct{Name: "Doe"}, ""},
		{"error_vs_nil", anError, nilError, ""},
		{"errors_different_message", errors.New("test1"), errors.New("test2"), ""},
		{"nil_vs_value", nil, 1, ""},
		{"nil_slice_vs_empty_slice", nilSlice, emptySlice, ""},

		// Failed tests
		{"integers_same", 1, 1, `
Error:    Values should not be equal
Value:    1`},
		{"structs_same", testStruct{Name: "John"}, testStruct{Name: "John"}, `
Error:    Values should not be equal
Value:    assert.testStruct{Name:"John"}`},
		{"errors_same_message", errors.New("test1"), errors.New("test1"), `
Error:    Values should not be equal
Value:    &errors.errorString{s:"test1"}`},
		{"errors_same_instance", anError, anError, `
Error:    Values should not be equal
Value:    &errors.errorString{s:"test error"}`},
		{"nil_vs_nil", nilError, nil, `
Error:    Values should not be equal
Value:    <nil>`},
		{"pointers_same", &testStruct{Name: "John"}, &testStruct{Name: "John"}, `
Error:    Values should not be equal
Value:    &assert.testStruct{Name:"John"}`},
		{"slices_same", []int{1, 2, 3}, []int{1, 2, 3}, `
Error:    Values should not be equal
Value:    []int{1, 2, 3}`},
		{"maps_same", map[string]int{"one": 1}, map[string]int{"one": 1}, `
Error:    Values should not be equal
Value:    map[string]int{"one":1}`},
		{"any_same", anyA, anyA, `
Error:    Values should not be equal
Value:    1`},
		{"nil_map_vs_nil_map", nilMap, nilMap, `
Error:    Values should not be equal
Value:    map[string]int(nil)`},
		{"empty_map_vs_empty_map", emptyMap, emptyMap, `
Error:    Values should not be equal
Value:    map[string]int{}`},
		{"nil_slice_vs_nil_slice", nilSlice, nilSlice, `
Error:    Values should not be equal
Value:    []string(nil)`},
		{"empty_slice_vs_empty_slice", emptySlice, emptySlice, `
Error:    Values should not be equal
Value:    []string{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			mockT := NewMockT()
			NotEqual(mockT, tt.actual, tt.expected)

			if tt.expectedMsg == "" {
				MustEqual(t, false, mockT.failed)
			} else {
				MustEqual(t, true, mockT.failed)
				actualError := strings.ReplaceAll(mockT.Logs(), "\u00a0", " ")
				Equal(t, actualError, tt.expectedMsg, "mismatched t.Log output")
			}
		})
	}
}

func TestMustEqual(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		mockT := NewMockT()
		MustEqual(mockT, 1, 1)
		Equal(t, false, mockT.Failed())
	})

	t.Run("failed", func(t *testing.T) {
		mockT := NewMockT()
		MustEqual(mockT, 1, 2)
		Equal(t, true, mockT.Failed())
	})
}

func TestMustNotEqual(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		mockT := NewMockT()
		MustNotEqual(mockT, 1, 2)
		Equal(t, false, mockT.Failed())
	})

	t.Run("failed", func(t *testing.T) {
		mockT := NewMockT()
		MustNotEqual(mockT, 1, 1)
		Equal(t, true, mockT.Failed())
	})
}

func TestColorOutput(t *testing.T) {
	setColor(t)
	mockT := NewMockT()
	Equal(mockT, 1, 2)

	logs := mockT.Logs()
	if !strings.Contains(logs, "\033[31m") {
		t.Error("Expected colored output in error message when colors enabled")
	}

	noColor = true
	mockT = NewMockT()
	Equal(mockT, 1, 2)

	logs = mockT.Logs()
	if strings.Contains(logs, "\033[31m") {
		t.Error("Expected no color codes in output when colors disabled")
	}
}
