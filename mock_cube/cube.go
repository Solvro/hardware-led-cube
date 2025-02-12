package mock

import (
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
	ledSize      = 0.125
	cubeWidth    = 4
)

var (
	cameraPos   = mgl32.Vec3{1.5, 1.5, 8}
	cameraFront = mgl32.Vec3{0, 0, -1}
	cameraRight = mgl32.Vec3{0, 1, 0}
	cameraUp    = mgl32.Vec3{0, -1, 0}
	deltaTime   float32
)

var (
	yaw         float32 = -90.0
	pitch       float32 = 0.0
	lastX       float32 = windowWidth / 2.0
	lastY       float32 = windowHeight / 2.0
	firstMouse  bool    = true
	sensitivity float32 = 0.1
)

type Led struct {
	X, Y, Z float32
	R, G, B float32
}

func init() {
	runtime.LockOSThread()
}

func initializeOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	program, err := newShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)
	return program
}

func configureShaders(program uint32) (int32, int32, int32) {
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	return projectionUniform, cameraUniform, modelUniform
}

func configureVertexData(program uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	sphereVertices := generateSphereVertices(1.0, 30, 30)
	gl.BufferData(gl.ARRAY_BUFFER, len(sphereVertices)*4, gl.Ptr(sphereVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	return vao
}

func configureGlobalSettings() {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func createLeds(width int) [][][]Led {
	leds := make([][][]Led, width)
	for x := 0; x < width; x++ {
		leds[x] = make([][]Led, width)
		for y := 0; y < width; y++ {
			leds[x][y] = make([]Led, width)
			for z := 0; z < width; z++ {
				leds[x][y][z] = Led{
					X: float32(x),
					Y: float32(y),
					Z: float32(z),
					R: rand.Float32(),
					G: rand.Float32(),
					B: rand.Float32(),
				}
			}
		}
	}
	return leds
}

func render(window *glfw.Window, program uint32, cameraUniform, modelUniform, ledColorUniform int32, leds [][][]Led, sphereVertices []float32) {
	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		deltaTime = float32(time - previousTime)
		previousTime = time

		angle += float64(deltaTime)
		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		processInput(window)

		camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraRight)
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		gl.UseProgram(program)

		for x := 0; x < cubeWidth; x++ {
			for y := 0; y < cubeWidth; y++ {
				for z := 0; z < cubeWidth; z++ {
					led := leds[x][y][z]

					color := mgl32.Vec3{led.R, led.G, led.B}
					gl.Uniform3fv(ledColorUniform, 1, &color[0])

					model = mgl32.Translate3D(led.X, led.Y, led.Z).Mul4(mgl32.Scale3D(ledSize, ledSize, ledSize))
					gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

					gl.DrawArrays(gl.TRIANGLES, 0, int32(len(sphereVertices)/5))
				}
			}
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func main() {
	window := initializeGLFW()
	defer glfw.Terminate()

	window.SetCursorPosCallback(mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	program := initializeOpenGL()
	_, cameraUniform, modelUniform := configureShaders(program)
	configureVertexData(program)
	configureGlobalSettings()

	leds := createLeds(cubeWidth)
	ledColorUniform := gl.GetUniformLocation(program, gl.Str("ledColor\x00"))
	sphereVertices := generateSphereVertices(1.0, 30, 30)

	render(window, program, cameraUniform, modelUniform, ledColorUniform, leds, sphereVertices)
}

func generateSphereVertices(radius float32, latitudeBands int, longitudeBands int) []float32 {
	var vertices []float32

	for latNumber := 0; latNumber < latitudeBands; latNumber++ {
		theta := float32(latNumber) * float32(math.Pi) / float32(latitudeBands)
		nextTheta := float32(latNumber+1) * float32(math.Pi) / float32(latitudeBands)
		sinTheta := float32(math.Sin(float64(theta)))
		cosTheta := float32(math.Cos(float64(theta)))
		sinNextTheta := float32(math.Sin(float64(nextTheta)))
		cosNextTheta := float32(math.Cos(float64(nextTheta)))

		for longNumber := 0; longNumber < longitudeBands; longNumber++ {
			phi := float32(longNumber) * 2 * float32(math.Pi) / float32(longitudeBands)
			nextPhi := float32(longNumber+1) * 2 * float32(math.Pi) / float32(longitudeBands)
			sinPhi := float32(math.Sin(float64(phi)))
			cosPhi := float32(math.Cos(float64(phi)))
			sinNextPhi := float32(math.Sin(float64(nextPhi)))
			cosNextPhi := float32(math.Cos(float64(nextPhi)))

			// First triangle
			vertices = append(vertices, radius*cosPhi*sinTheta, radius*cosTheta, radius*sinPhi*sinTheta, float32(longNumber)/float32(longitudeBands), float32(latNumber)/float32(latitudeBands))
			vertices = append(vertices, radius*cosNextPhi*sinTheta, radius*cosTheta, radius*sinNextPhi*sinTheta, float32(longNumber+1)/float32(longitudeBands), float32(latNumber)/float32(latitudeBands))
			vertices = append(vertices, radius*cosPhi*sinNextTheta, radius*cosNextTheta, radius*sinPhi*sinNextTheta, float32(longNumber)/float32(longitudeBands), float32(latNumber+1)/float32(latitudeBands))

			// Second triangle
			vertices = append(vertices, radius*cosNextPhi*sinTheta, radius*cosTheta, radius*sinNextPhi*sinTheta, float32(longNumber+1)/float32(longitudeBands), float32(latNumber)/float32(latitudeBands))
			vertices = append(vertices, radius*cosNextPhi*sinNextTheta, radius*cosNextTheta, radius*sinNextPhi*sinNextTheta, float32(longNumber+1)/float32(longitudeBands), float32(latNumber+1)/float32(latitudeBands))
			vertices = append(vertices, radius*cosPhi*sinNextTheta, radius*cosNextTheta, radius*sinPhi*sinNextTheta, float32(longNumber)/float32(longitudeBands), float32(latNumber+1)/float32(latitudeBands))
		}
	}

	return vertices
}
