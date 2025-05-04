# Hunt-the-Wumpus
This repo contains code for the game 'Hunt the Wumpus', written in Go.

# OpenGL
A lot of the boilerplate is the same as from Conway's, as that same code is required in each OpenGL project to get a render window up-and-running, so it's honestly just copied over (https://github.com/SmoCloud/Conway-GOL).

Any differences in the OpenGL are added/changed and will be listed here under the OpenGL section of this Readme

Coordinates for the dodecahedron vertices were found on reddit here - (https://www.reddit.com/r/opengl/comments/mqp63d/creating_dodecahedron/). Those vertices are used to draw the dodecahedron

With the help of AI (Copilot), I was given some example code of what to do with the buffers by searching 'golang glfw opengl dodecahedron 2D'

# Design
The design of the map came from the image provided under the Development tab on the wikipedia page for Hunt The Wumpus, that shows the different vertices of a dodecahedron as each room (https://en.wikipedia.org/wiki/Hunt_the_Wumpus)
