package main

import (
	"reflect"
	"testing"

	// "github.com/go-gl/gl/v4.6-core/gl"	// Import the OpenGL implementation for Go, used for graphics rendering
	// "github.com/go-gl/glfw/v3.3/glfw"	// Import the GLFW implementation for Go, simplifies creating a window with OpenGL

)

func assertType(t *testing.T, val any, expectedType string) {
	actualType := reflect.TypeOf(val).String()
	if actualType != expectedType {
		t.Errorf("Expected type %s but got %s\n", expectedType, actualType)
	}
}

// BEGIN: Unit tests begin here
func TestInitGlfw(t *testing.T) {
	var got any = InitGlfw()
	
	assertType(t, got, "*glfw.Window")
}