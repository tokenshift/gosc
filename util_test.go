package gosc

/**
 * Helper functions for test code.
 */

import (
	"path/filepath"
	"reflect"
	"runtime"
	. "testing"
)

func expectf(t *T, expected, actual interface{}) {
	if _, fname, line, ok := runtime.Caller(2); ok {
		t.Errorf("%s:%d: expected %#v (%T), got %#v (%T)", filepath.Base(fname), line, expected, expected, actual, actual)
	} else {
		t.Errorf("expected %#v (%T), got %#v (%T)", expected, expected, actual, actual)
	}
}

func expectNil(t *T, actual interface{}) bool {
	if actual != nil {
		expectf(t, nil, actual)
		return false
	}

	return true
}

// Test for deep equality (using reflection); returns true if the arguments
// match, otherwise reports an error and returns false.
func expectSame(t *T, expected, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		expectf(t, expected, actual)
		return false
	}

	return true
}
