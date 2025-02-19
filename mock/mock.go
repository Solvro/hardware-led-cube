package mock

import (
	"hardware-led-cube/frames"
	_ "image/png"
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	ledSize = 0.125
)

var (
	cameraPos       = mgl32.Vec3{1.5, 1.5, 8}
	cameraFront     = mgl32.Vec3{0, 0, -1}
	cameraRight     = mgl32.Vec3{0, 1, 0}
	cameraUp        = mgl32.Vec3{0, -1, 0}
	previousTime    float64
	deltaTime       float32
	window          *glfw.Window
	program         uint32
	cameraUniform   int32
	modelUniform    int32
	ledColorUniform int32
	sphereVertices  []float32
	angle           float64
)

type led struct {
	X, Y, Z float32
	R, G, B float32
}

type Cube struct {
	leds [][][]led
}

func InitCube(width, height, depth int) *Cube {
	// we only want to lock the os thread if we are mocking
	runtime.LockOSThread()
	window = initializeGLFW()

	window.SetCursorPosCallback(mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	program = initializeOpenGL()
	_, cameraUniform, modelUniform = configureShaders(program)
	configureVertexData(program)
	configureGlobalSettings()

	ledColorUniform = gl.GetUniformLocation(program, gl.Str("ledColor\x00"))
	if ledColorUniform == -1 {
		log.Fatalln("failed to get uniform location for ledColor")
	}
	sphereVertices = generateSphereVertices(1.0, 30, 30)
	angle = 0.0
	previousTime = 0.0
	deltaTime = 0.0

	leds := createLeds(width, height, depth)
	return &Cube{leds}
}

func (c *Cube) SetLeds(f frames.Frame) {
	newLeds := f.ToXYZ()
	for x := range len(newLeds) {
		for y := range len(newLeds[x]) {
			for z := range len(newLeds[x][y]) {
				led := &c.leds[x][y][z]
				led.R = float32(newLeds[x][y][z] >> 16 & 0xFF)
				led.G = float32(newLeds[x][y][z] >> 8 & 0xFF)
				led.B = float32(newLeds[x][y][z] & 0xFF)
			}
		}
	}
}

func (c *Cube) Render() error {
	if !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		if previousTime != 0.0 {
			deltaTime = float32(time - previousTime)
		}
		previousTime = time

		angle += float64(deltaTime)
		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		processInput(window)

		camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraRight)
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		gl.UseProgram(program)

		for x := range len(c.leds) {
			for y := range len(c.leds[x]) {
				for z := range len(c.leds[x][y]) {
					led := c.leds[x][y][z]

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
	return nil
}

func (c *Cube) Fini() {
	glfw.Terminate()
}

func createLeds(width, height, depth int) [][][]led {
	leds := make([][][]led, width)
	for x := 0; x < width; x++ {
		leds[x] = make([][]led, height)
		for y := 0; y < height; y++ {
			leds[x][y] = make([]led, depth)
			for z := 0; z < depth; z++ {
				leds[x][y][z] = led{
					X: float32(x),
					Y: float32(y),
					Z: float32(z),
					R: 0.0,
					G: 0.0,
					B: 0.0,
				}
			}
		}
	}
	return leds
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
