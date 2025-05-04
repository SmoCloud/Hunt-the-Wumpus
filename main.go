package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	// "sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"	// Import the OpenGL implementation for Go, used for graphics rendering
	"github.com/go-gl/glfw/v3.1/glfw"	// Import the GLFW implementation for Go, simplifies creating a window with OpenGL
)

const (
	Fps       = 60 // frames per second

	// vertex shader source code, for telling OpenGL where the vertices of each shape will be relative to the center of each cell
	VertexShaderSource = `
        #version 410
        in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

	// fragment shader source code, tells OpenGL what color the shape that's drawn with the vertex shader will be
	FragmentShaderSource = `
        #version 410
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(0.9, 0.4, 0.7, 1);
        }
    ` + "\x00"
)

var (
	// pre-defined square shape, using two triangles
	Dodecahedron = []float32{
		0, 0.618, 1.618,
		0, -0.618, 1.618,
		0, -0.618, -1.618,
		0, 0.618, -1.618,

		1.618, 0, 0.618,
		-1.618, 0, 0.618,
		-1.618, 0, -0.618,
		1.618, 0, -0.618,

		0.618, 1.618, 0,
		-0.618, 1.618, 0,
		-0.618, -1.618, 0,
		0.618, -1.618, 0,

		1, 1, 1,
		-1, 1, 1,
		-1, -1, 1,
		1, -1, 1,

		1, -1, -1,
		1, 1, -1,
		-1, 1, -1,
		-1, -1, -1,
	}

	Width  = 500	// width of the render window, in pixels
	Height = 500	// height of the render window, in pixels

	// 9*9 grid for a total of 81 cells
	Rows   = 9		// cell count along the width of the window
	Cols   = 9		// cell count along the height of the window
)

func main() {
	runtime.LockOSThread()	// main OS thread has to be locked for OpenGL rendering, though threading is possible for anything not involving OpenGL

	window := initGlfw()	// initiate the render window
	defer glfw.Terminate()	// tells the program to close the window when it reaches the end of the main function

	program := initOpenGL()	// create the shader by combining the vertex and fragment shaders

	for !window.ShouldClose() {	// while the window is not closed
		t := time.Now()
	
		draw(window, program)

		time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))
	}
}

// initGlfw initializes glfw and returns a Window object that can be used to render graphics.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)	// tells GLFW the window will not be resizable
	glfw.WindowHint(glfw.ContextVersionMajor, 4)	// tells GLFW the major version being used
	glfw.WindowHint(glfw.ContextVersionMinor, 1)	// tells GLFW the minor version being used
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)	// tells GLFW to use the default configuration settings
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)	// tells GLFW that this program will be compatible with newer versions of OpenGL (I think?)

	window, err := glfw.CreateWindow(Width, Height, "Hunt the Wumpus", nil, nil)	// creates the window
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()	// makes the window the current context to display
	glfw.SwapInterval(glfw.True)

	return window
}

// initOpenGL initializes OpenGL and returns an initialized program
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {	// this if structure initializes in the first step and then does the check whether to run the conditional code in the second
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))	// OpenGL actually uses a *uint8 type for strings, so the string needs to be converted to that type
	log.Println("OpenGL version", version)	// log seems to do the same thing as fmt, though I think log might work better with logging to a file, may create a log file

	// compiles the vertex shader, telling the GPU what shape will be drawn through its vertices and their locations relative to the center of the cell
	vertexShader, err := compileShader(VertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	// compiles the fragment shader, telling the GPU what color the shape drawn in the cell will be
	fragmentShader, err := compileShader(FragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// creates a full shader, called a program (since a shader is just a program that's ran on the GPU) by attaching the VertexShader and FragmentShader to it
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)	// links the shader program to the buffer for the GPU's use
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 20*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointerWithOffset(0, 4, gl.FLOAT, false, 3*4, 0)

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

// compileShader will send the shader source code to the GPU for compilation on the GPU (shaders handle vertex points of drawn objects as well as their color)
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)	// creates the shader of either type vertex or fragment

	csources, free := gl.Strs(source)	// re-types the shader source from a string type to a *uint8 type for use by the GPU
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// draw clears anything that's on the screen before drawing new objects
// Cannot parallelize draws as OpenGL requires operations to happen on a single thread
func draw(window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	glfw.PollEvents()
	window.SwapBuffers()
}