package mock

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func processInput(window *glfw.Window) {
	var cameraSpeed float32 = 2.5 * deltaTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraRight).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraRight).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraUp.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		cameraPos = cameraPos.Add(cameraUp.Mul(cameraSpeed))
	}

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func mouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	if firstMouse {
		lastX = float32(xpos)
		lastY = float32(ypos)
		firstMouse = false
	}

	xoffset := float32(xpos) - lastX
	yoffset := lastY - float32(ypos)
	lastX = float32(xpos)
	lastY = float32(ypos)

	xoffset *= sensitivity
	yoffset *= sensitivity

	yaw += xoffset
	pitch += yoffset

	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
	}
	cameraFront = front.Normalize()
}
