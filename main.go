package main

import (
	"fmt"
	"log"
	// "math/rand"
	"runtime"
	"strings"
	// "sync"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"	// Import the OpenGL implementation for Go, used for graphics rendering
	"github.com/go-gl/glfw/v3.3/glfw"	// Import the GLFW implementation for Go, simplifies creating a window with OpenGL
)

const (
	Fps       = 30 // frames per second
	Radius	  = 0.2

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
								// indices of the vertices
		0.0, 0.965, 0.0,		// 0
		0.965, 0.125, 0.0,		// 1
		0.5, -0.965, 0.0,		// 2
		-0.5, -0.965, 0.0,		// 3
		-0.965, 0.125, 0.0,		// 4
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

		// Inner Pentagon Points
		0.25, 0.4275, 0.0,		// 15
		-0.25, 0.4275, 0.0,		// 16
		-0.5, -0.125, 0.0,		// 17
		0.0, -0.625, 0.0,		// 18
		0.5, -0.125, 0.0,		// 19
	}

	Indices = []uint32{	// these are the endpoints for each line that's drawn with gl.DrawElements
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

	isFullscreen = true // window starts in windowed mode, but for first toggle, this must be set true
)

func main() {
	runtime.LockOSThread()	// main OS thread has to be locked for OpenGL rendering, though threading is possible for anything not involving OpenGL

	window := InitGlfw()	// initiate the render window
	defer glfw.Terminate()	// tells the program to close the window when it reaches the end of the main function

	program := InitOpenGL()	// create the shader by combining the vertex and fragment shaders
	
	// code found by searching 'fetch desktop screen size golang opengl glfw'
	// This code gets the size of the current primary monitor in your display settings to allow the window to be rendered at the size of the monitor
	mainMonitor := glfw.GetPrimaryMonitor()
	if mainMonitor == nil {
		panic("Failed to get primary monitor size.")
	}
	videoMode := mainMonitor.GetVideoMode()
	if videoMode == nil {
		panic("Failed to get the video mode of the primary monitor.")
	}

	// code found by searching 'allow swapping between fullscreen and windowed mode opengl glfw golang'
	window.SetPos((videoMode.Width - 800) / 2, (videoMode.Height - 800) / 2)

	for !window.ShouldClose() {	// while the window is not closed
		t := time.Now()
	
		// This code took a lot of research to get to, ultimately pieces of it were each sourced from different places
		// I realized the version of glfw and OpenGL I was using weren't the most recent and some of the functions
		// I needed to do this were only in the most recent versions, so I had to update from OpenGL 4.1 to 4.6
		// and Glfw 3.1 to 3.3
		// searched 'golang SetPrimaryMonitor toggle fullscreen modern opengl glfw'
		// searched 'opengl glfw golang key callback to toggle fullscreen'
		// 'https://www.glfw.org/docs/3.3/input_guide.html#input_key'
		// 'https://www.glfw.org/docs/latest/group__window.html#ga81c76c418af80a1cce7055bccb0ae0a7' were the most helpful

		window.SetKeyCallback(KeyCallback)

		log.Println(window.GetCursorPos())

		vao := MakeVao(Dodecahedron)	// this creates and returns the vertex array object (vao) for drawing
		DrawGame(vao, window, program)		// this function takes the shader program and vao and draws the shape (pentagon right now)

		time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))	// this locks the framerate of the game, currently at 30fps
	}
}

// initGlfw initializes glfw and returns a Window object that can be used to render graphics.
func InitGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)	// tells GLFW the window will not be resizable
	glfw.WindowHint(glfw.ContextVersionMajor, 4)	// tells GLFW the major version being used
	glfw.WindowHint(glfw.ContextVersionMinor, 1)	// tells GLFW the minor version being used
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)	// tells GLFW to use the default configuration settings
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)	// tells GLFW that this program will be compatible with newer versions of OpenGL (I think?)

	window, err := glfw.CreateWindow(800, 800, "Hunt the Wumpus", nil, nil)	// creates the window
	if err != nil {
		panic(err)
	}
	
	window.MakeContextCurrent()	// makes the window the current context to display
	glfw.SwapInterval(glfw.True)

	return window
}

// initOpenGL initializes OpenGL and returns an initialized program
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {	// this if structure initializes in the first step and then does the check whether to run the conditional code in the second
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))	// OpenGL actually uses a *uint8 type for strings, so the string needs to be converted to that type
	log.Println("OpenGL version", version)	// log seems to do the same thing as fmt, though I think log might work better with logging to a file, may create a log file

	// compiles the vertex shader, telling the GPU what shape will be drawn through its vertices and their locations relative to the center of the cell
	vertexShader, err := CompileShader(VertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	// compiles the fragment shader, telling the GPU what color the shape drawn in the cell will be
	fragmentShader, err := CompileShader(FragmentShaderSource, gl.FRAGMENT_SHADER)
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
func MakeVao(points []float32) uint32 {
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
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)	// creates the shader of either type vertex or fragment

	csources, free := gl.Strs(source)	// re-types the shader source from a string type to a *uint8 type for use by the GPU
	gl.ShaderSource(shader, 1, csources, nil)	// tells OpenGL what the shader source code is
	free()	// honestly not sure why free needs to be called
	gl.CompileShader(shader)	// compiles the shader source code on the GPU

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)	// makes sure the shader source code compiles successfully
	if status == gl.FALSE {	// status will store gl.FALSE if the compilation was unsuccessful
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)	// this grabs paramaters from the shader source code

		log := strings.Repeat("\x00", int(logLength+1))	// this divides the two shaders (vertex and fragment) apart
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))	// gets the log error message from the GPU about the shader compilation, where the error was

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)	// prints the GPU's compilation log to the CPU console window
	}

	return shader, nil	// if status is gl.TRUE, compilation successful, so return the compiled shader as a program to use in the GPU
}

// draw clears anything that's on the screen before drawing new objects
// Cannot parallelize draws as OpenGL requires operations to happen on a single thread
func DrawGame(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)	// clears the drawn colors between each frame
	gl.UseProgram(program)	// tells the GPU what shader to use in the generated program

	gl.BindVertexArray(vao)	// bind the vertex array object to the enabled vertix attribute
	gl.DrawElements(gl.LINES, 60, gl.UNSIGNED_INT, nil)	// tells the GPU to draw the vertex array object in the order specified by the Indices array that was passed to the ELEMENT_ARRAY_BUFFER
	
	gl.PointSize(10.0)	// tells the GPU how large to draw the points
	gl.DrawArrays(gl.POINTS, 0, 20)	// draws the vertices at the size of what was specified by the line above
	gl.BindVertexArray(0)

	glfw.PollEvents()
	window.SwapBuffers()
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {	// this code will be useful for any other keypresses I want to implement
	if action == glfw.Press && key == glfw.KeyF11 {
		ToggleFullscreen(w)
	}
}

func ToggleFullscreen(w *glfw.Window) {
	// if isFullscreen is true, set the monitor reference to the current monitor and update the viewport so the draw ratio is correct
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()
	if isFullscreen {
		w.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		gl.Viewport(0, 0, int32(mode.Width), int32(mode.Height))

	} else {
		w.SetMonitor(nil, (mode.Width - 800) / 2, (mode.Height - 800) / 2, 800, 800, mode.RefreshRate)
		gl.Viewport(0, 0, int32(800), int32(800))
	}
	isFullscreen = !isFullscreen
}