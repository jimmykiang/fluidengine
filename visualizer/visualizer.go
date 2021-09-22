package visualizer

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"jimmykiang/fluidengine/Vector3D"
	"time"
)

func Visualize(vector []*Vector3D.Vector3D, n int64) {

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Spheres
	sphereGeometry := geometry.NewSphere(0.02, 32, 16)

	d2 := light.NewDirectional(&math32.Color{0.6, 1, 0.6}, 1.0)
	d2.SetPosition(0, 1, 0)
	scene.Add(d2)

	for i := int64(0); i < n; i++ {

		pbrMat := material.NewPhysical()
		pbrMat.SetMetallicFactor(0.4)
		pbrMat.SetRoughnessFactor(0.2)
		//pbrMat.SetEmissiveFactor(&math32.Color{0.3, 0.3, 0.6})
		pbrMat.SetBaseColorFactor(&math32.Color4{0, 0, 1, 1})
		sphere := graphic.NewMesh(sphereGeometry, pbrMat)
		sphere.SetPosition(float32(vector[i].X), float32(vector[i].Y), float32(vector[i].Z))
		scene.Add(sphere)

	}

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
