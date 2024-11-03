# Assert Package

The `assert` package provides helper functions for writing tests in Go. It offers detailed error messages and diffs for failed assertions.

## Installation

To install the package, use:

```sh
go get github.com/r2k1/assert
```

## Example

Here's a sample usage of the `assert` package:

```go
// Check for equality
assert.Equal(t, actual, expected, "Person objects should be equal")
// Error:    Not equal
// Message:  Person objects should be equal
// Expected: assert.Person{Name:"Alice", Age:30}
// Actual:   assert.Person{Name:"Alice", Age:25}
// Diff:     {
//         -   Age: 25,
//         +   Age: 30,
//         }

// Stop test execution if not equal
assert.MustEqual(t, doSomething(), nil)
// Error:    Not equal
// Expected: <nil>
// Actual:   &errors.errorString{s:"failed to do something!"}}
```
