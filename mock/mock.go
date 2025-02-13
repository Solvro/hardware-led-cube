package mock

import (
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	ledSize   = 0.125
	cubeWidth = 4
)

var (
	cameraPos   = mgl32.Vec3{1.5, 1.5, 8}
	cameraFront = mgl32.Vec3{0, 0, -1}
	cameraRight = mgl32.Vec3{0, 1, 0}
	cameraUp    = mgl32.Vec3{0, -1, 0}
	deltaTime   float32
)

type Led struct {
	X, Y, Z float32
	R, G, B float32
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

func render(window *glfw.Window, program uint32, cameraUniform, modelUniform, ledColorUniform int32, leds [][][]Led, sphereVertices []float32, ledUpdates chan [][][]Led) {
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

		select {
		case newLeds := <-ledUpdates:
			leds = newLeds
		default:
		}

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

func Setup(ledUpdates chan [][][]Led) {
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
	if ledColorUniform == -1 {
		log.Fatalln("failed to get uniform location for ledColor")
	}
	sphereVertices := generateSphereVertices(1.0, 30, 30)

	render(window, program, cameraUniform, modelUniform, ledColorUniform, leds, sphereVertices, ledUpdates)
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
