package main

import (
	"fmt"
	"log"
	// "math/rand"
	"runtime"
	"strings"
	// "sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"	// Import the OpenGL implementation for Go, used for graphics rendering
	"github.com/go-gl/glfw/v3.1/glfw"	// Import the GLFW implementation for Go, simplifies creating a window with OpenGL
)

const (
	Fps       = 30 // frames per second

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
	// Pre-defined dodecahedron shape, using lines
	// Calculations are based on 0.075 increments
	Dodecahedron = []float32{
		// Outer Pentagon points
		0.0, 1.0, 0.0,			// 0
		1.0, 0.125, 0.0,		// 1
		0.5, -1.0, 0.0,			// 2
		-0.5, -1.0, 0.0,		// 3
		-1.0, 0.125, 0.0,		// 4
		// 0.0, 1.0, 0.0,

		// Inner Decagon Points
		0.0, 0.6775, 0.0,		// 5
		0.375, 0.53625, 0.0,	// 6
		0.6925, 0.1925, 0.0,	// 7
		0.6925, -0.225, 0.0,	// 8
		0.375, -0.625, 0.0,		// 9
		0.0, -0.75, 0.0,		// 10
		-0.375, -0.625, 0.0,	// 11
		-0.6925, -0.225, 0.0,	// 12
		-0.6925, 0.1925, 0.0,	// 13
		-0.375, 0.53625, 0.0,	// 14
		// 0.0, 0.6775, 0.0,

		// Inner Pentagon Points
		// 0.375, 0.53625, 0.0,
		0.25, 0.4275, 0.0,		// 15
		-0.25, 0.4275, 0.0,		// 16
		-0.5, -0.125, 0.0,		// 17
		0.0, -0.625, 0.0,		// 18
		0.5, -0.125, 0.0,		// 19
		// 0.25, 0.4275, 0.0,
	}

	Indices = []uint32{
	 	0, 1,
		1, 2, 
		2, 3,
		3, 4,
		4, 0,

		0, 5,
		1, 7,
		2, 9,
		3, 11,
		4, 13,

		5, 6,
		6, 7, 
		7, 8,
		8, 9,
		9, 10,
		10, 11,
		11, 12, 
		12, 13,
		13, 14,
		14, 5,

		6, 15,
		8, 19,
		10, 18,
		12, 17,
		14, 16,

		15, 19,
		19, 18,
		18, 17,
		17, 16,
		16, 15,
		
	}

	Width  = 500	// width of the render window, in pixels
	Height = 500	// height of the render window, in pixels

	// 9*9 grid for a total of 81 cells
	// Rows   = 9		// cell count along the width of the window
	// Cols   = 9		// cell count along the height of the window
)

func main() {
	runtime.LockOSThread()	// main OS thread has to be locked for OpenGL rendering, though threading is possible for anything not involving OpenGL

	window := initGlfw()	// initiate the render window
	defer glfw.Terminate()	// tells the program to close the window when it reaches the end of the main function

	program := initOpenGL()	// create the shader by combining the vertex and fragment shaders

	for !window.ShouldClose() {	// while the window is not closed
		t := time.Now()
	
		vao := makeVao(Dodecahedron)	// this creates and returns the vertex array object (vao) for drawing
		draw(vao, window, program)		// this function takes the shader program and vao and draws the shape (pentagon right now)

		time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))	// this locks the framerate of the game, currently at 30fps
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
	prog := gl.CreateProgram()	// creates a new program object that will be used by the GPU
	gl.AttachShader(prog, vertexShader)	// attaches the vertex shader to the program object
	gl.AttachShader(prog, fragmentShader)	// attaches the fragment shader to the program object
	gl.LinkProgram(prog)	// links the shader program to the buffer for the GPU's use
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	// creates a vertex buffer object (vbo), which will hold the vertex array object (vao) needed to map the vertices to the render window
	var vbo uint32
	gl.GenBuffers(1, &vbo)	// tells the GPU to treat buffer 1 as a vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)	// binds the vbo to the vbo buffer in the GPU
	// 4*len(points) approximates the size (in bytes) of the vertex array, and STATIC_DRAW means it will be uploaded to the GPU once but drawn multiple times
	// may change to a STREAM_DRAW since, once the dodecahedron is drawn, it doesn't need to be re-drawn at any point in time since it won't change, though ClearColor may still clear it out, yet to see
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// creates the vertex array object (vao) that will be stored inside of the vbo
	var vao uint32
	gl.GenVertexArrays(1, &vao)	// tells the GPU to create a vertex array out of the data stored in the vertex buffer at buffer 1
	gl.BindVertexArray(vao)	// binds (or uploads) the vao to the vertex buffer
	gl.EnableVertexAttribArray(0)	// swaps the vertex array in the vertex buffer to buffer 0, which is the primary buffer the GPU reads from
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)	// binds the vao and the vbo together (I think? Anyone correct me if I'm wrong about anything here)
	
	/* 
	First argument can be used to send data to the shader if you want to allow vertices to be added through input
	Second argument specifies how many values there are for each vertex, in our case, 3 (X, Y, Z, though Z won't be used, it's still needed)
		In the case that vertices are added through input, this needs to match what's in the vertex shader (2 for vec2, 3 for vec3, 4 for vec4)
	Third argument specifies the data type of each component piece of the input
	Fourth argument specifies if input should be normalized in the case they aren't floating point inputs
	Fifth argument specifies the stride (or amount of data in between each vertex) in bytes (0 for us since there is no data between each vertex to skip over)
	Final argument specifies the offset, or where in the vao the reading of the vertices begins, as a byte pointer (0 since the vertices start at the beginning of the vao) */
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// this will be needed later for the indices array that tells the GPU what draw order I want to use
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(Indices), gl.Ptr(Indices), gl.STATIC_DRAW)

	// this is used if rasterization is enabled using gl.PointSize(), which allows vertices to be rendered with larger diameters
	// may end up using it to make the vertices larger than the lines drawn between them, to visually indicate each vertex is a room, and is more important than the paths connecting them
	// gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)	//	clears the color buffers (r, g, b, a) and replaces them with the specified input. Needed to 'refresh' the screen between each draw call

	return vao
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
func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.LINES, 60, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	glfw.PollEvents()
	window.SwapBuffers()
}