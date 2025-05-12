# Hunt-the-Wumpus
This repo contains code for the game 'Hunt the Wumpus', written in Go.

# Dependencies
Note    - To install all dependencies when cloning this project, just run 'go mod tidy' in your terminal from the project's directory (assuming you already have Go installed, if not, why did you clone a Go project?)
- OpenGL  - github.com/go-gl/gl
- GLFW    - github.com/go-gl/glfw

# OpenGL
- A lot of the boilerplate is the same as from Conway's, as that same code is required in each OpenGL project to get a render window up-and-running, so it's honestly just copied over (https://github.com/SmoCloud/Conway-GOL).

- Any differences in the OpenGL are added/changed and will be listed here under the OpenGL section of this Readme

- I tried the coordinates found here - (https://www.reddit.com/r/opengl/comments/mqp63d/creating_dodecahedron/). They're not in the format OpenGL needs. I tried working with them for a while, probably too long, they weren't giving me what I needed, which made me realize I didn't fully understand how OpenGL uses vertices.

- With the help of AI (Copilot), I was given some example code of what to do with the buffers by searching 'golang glfw opengl dodecahedron 2D', though I didn't end up using this, either, as the use-case it was presenting was lacking the information I needed, further cementing the feeling that I was not understanding what was happening with the vertices.

- After messing around with the code example, changing this variable and that, messing with the coordinates, changing X and Y values of the different coordinates, I stumbled upon this article by searching 'drawing a polygon with OpenGL' here (https://open.gl/drawing), which gave me an excellent rundown of OpenGL's graphics pipeline.

- Skimming through this Graphics Pipeline article made me realize that I needed to scrap the coordinates and make my own. I needed to play with drawing more from the ground up, starting by drawing a basic pentagon. 

- Thanks to that documentation for drawing shapes in OpenGL, despite being for C++, was incredibly helpful, and I was able to draw a basic pentagon and gain a deeper understanding of how the vertices are used and how I can use LINE_LOOP to draw any shape I wanted to draw, and I was able to draw the starting shape of a pentagon.

- The plan is to draw the entire dodecahedron one line at a time, connecting the vertices as needed, when needed. 

- I can further utilize something I found in this article I previously was unaware of, the ELEMENT_ARRAY_BUFFER, to have it draw lines from a specific vertex to another specific vertex, as, without this, it draws to the next vertex from the previous one that it drew to. 

- With an ELEMENT_ARRAY_BUFFER, you can specify the starting and ending vertices, which gives more control and reduces the size of the vertex array object being used, ultimately meaning I will be able to draw the dodecahedron.

- Using gl.LINES combined with an element array that specifies the two vertices of each drawn line, I am able to draw the connected map as it should be.

- The values used for the coordinates were calculated in 0.075 increments. I really just played with them, changing them in those increments, until I got something that resembled what I wanted (a pentagon inside of a decagon inside of a pentagon).

- Redrawing the vertices in the VBO using drawArrays and specifying them with gl.POINTS and setting gl.PointSize to 10.0 allows the points to be drawn larger than the lines themselves

- Using a keyCallback, I was able to implement the ability for the render window to swap between fullscreen and windowed modes

# Design
- The design of the map came from the image provided under the Development tab on the wikipedia page for Hunt The Wumpus, that shows the different vertices of a dodecahedron as each room (https://en.wikipedia.org/wiki/Hunt_the_Wumpus).

- The plan for designing the game will be to use each vertex as the rooms, and each connected vertex will be a room that can be travelled to. I will learn to draw circles over top of these vertices, and whichever room the player is in will change color to indicate to the player their current position in the map. 

- The Wumpus, the sinkholes, and the bats will not be visible in any way, to stay somewhat true to the original text-based game, and the player will use input from the console to specify which "room" or vertex they wish to travel to next. 

- If I can figure out how to implement the use of the mouse, I would also like to add the ability for the player to just click on the room to go to it, so long as it is a room connected to the room they are in, but the console will still be an option, with each vertex having an 'id' of sorts to denote which 'room' it is in the map (the dodecahedron).

# Unit Testing
- Courtesy of the go.dev documentation, I found that go has a "testing" module that makes writing unit tests very easy and convenient. https://go.dev/doc/tutorial/add-a-test

- I also found by searching google "get variable type in golang" that there's another useful module, "reflect", that can report what a variable's type is, and using this, I'm able to test that the correct type is returned from functions inside of my unit tests (i.e., check that InitGlfw returns a pointer to a glfw.Window object, which is the first unit test I've written for this project).

# Workflows
- The .yaml file inside of the .github/workflows directory pulls a docker image for an ubuntu environment running amd64 architecture. It installs the necessary software dependencies (gcc, lib1gl-mesa-dev, libx11, xorg, xvfb), exports the pkgconfig/gl.pc file to the PKG_CONFIG_PATH, all to be able to build the project on the container. These dependencies are needed to ensure the environment has a compatible version of OpenGL and GLFW, and to ensure that the cgo toolchain can compile the C libraries that OpenGL and GLFW use.

- The workflow also runs the tests on a xvfb server. Since the container is headless, this is needed to test anything that uses OpenGL.

- All thanks for getting this workflow working goes to Copilot's ability to diagnose the errors and report their fixes to me (I take minimal credit for the workflow file).