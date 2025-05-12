package main

import (
	"reflect"
	"testing"

	"github.com/go-gl/gl/v4.6-core/gl"	// Import the OpenGL implementation for Go, used for graphics rendering
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

func TestInitOpenGL(t *testing.T) {
	var got any = InitOpenGL()

	assertType(t, got, "uint32")
}

func TestMakeVao(t *testing.T) {
	// define some generic points for drawing a trianle
	points := []float32{	// should draw a traingle approx. mid-screen
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
	var got any = MakeVao(points)

	assertType(t, got, "uint32")
}

func TestCompileShader(t *testing.T) {
	_, err := CompileShader(VertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		t.Error("Shader compilation returned err: ", err)
	}
}

func TestDrawGame(t *testing.T) {	// not entirely sure how to test this one
	// define some generic points for drawing a trianle
	points := []float32{	// should draw a traingle approx. mid-screen
		0.0, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
	}

	// ensure all arguments passed are of the correct type
	
	win := InitGlfw()
	
	assertType(t, win, "*glfw.Window")

	prog := InitOpenGL()

	assertType(t, prog, "uint32")

	vao := MakeVao(points)

	assertType(t, vao, "uint32")

	DrawGame(vao, win, prog)
}