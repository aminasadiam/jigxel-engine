package engine

import (
	"log"
	"runtime"

	"github.com/aminasadiam/jigxer-engine/internal/ecs"
	"github.com/aminasadiam/jigxer-engine/internal/graphics"
	"github.com/aminasadiam/jigxer-engine/internal/input"
	"github.com/aminasadiam/jigxer-engine/internal/physics"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Engine represents the main game engine
type Engine struct {
	window   *glfw.Window
	title    string
	width    int
	height   int
	running  bool
	lastTime float64

	// Systems
	ecs      *ecs.World
	renderer *graphics.Renderer
	input    *input.Manager
	physics  *physics.World
}

// NewEngine creates a new game engine instance
func NewEngine(title string, width, height int) *Engine {
	return &Engine{
		title:   title,
		width:   width,
		height:  height,
		running: false,
	}
}

// Init initializes the game engine and all its systems
func (e *Engine) Init() error {
	// Lock the main thread for OpenGL
	runtime.LockOSThread()

	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		return err
	}

	// Configure GLFW
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create window
	window, err := glfw.CreateWindow(e.width, e.height, e.title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return err
	}
	e.window = window

	// Make the window's context current
	e.window.MakeContextCurrent()

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		return err
	}

	// Set viewport
	gl.Viewport(0, 0, int32(e.width), int32(e.height))

	// Initialize systems
	e.ecs = ecs.NewWorld()
	e.renderer = graphics.NewRenderer()
	e.input = input.NewManager(e.window)
	e.physics = physics.NewWorld()

	// Initialize renderer
	if err := e.renderer.Init(); err != nil {
		return err
	}

	// Initialize input manager
	if err := e.input.Init(); err != nil {
		return err
	}

	// Set up window callbacks
	e.setupCallbacks()

	log.Println("Engine initialized successfully")
	return nil
}

// Run starts the main game loop
func (e *Engine) Run() {
	e.running = true
	e.lastTime = glfw.GetTime()

	for !e.window.ShouldClose() && e.running {
		currentTime := glfw.GetTime()
		deltaTime := currentTime - e.lastTime
		e.lastTime = currentTime

		// Update systems
		e.update(deltaTime)

		// Render
		e.render()

		// Poll events and swap buffers
		glfw.PollEvents()
		e.window.SwapBuffers()
	}
}

// update updates all engine systems
func (e *Engine) update(deltaTime float64) {
	// Update input
	e.input.Update()

	// Update physics
	e.physics.Update(deltaTime)

	// Update ECS world
	e.ecs.Update(deltaTime)
}

// render renders the current frame
func (e *Engine) render() {
	// Clear the screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	// Render the scene
	e.renderer.Render(e.ecs)
}

// Shutdown cleans up the engine and all its resources
func (e *Engine) Shutdown() {
	e.running = false

	// Shutdown systems
	if e.renderer != nil {
		e.renderer.Shutdown()
	}

	// Terminate GLFW
	if e.window != nil {
		e.window.Destroy()
	}
	glfw.Terminate()

	log.Println("Engine shutdown complete")
}

// setupCallbacks sets up window event callbacks
func (e *Engine) setupCallbacks() {
	e.window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		e.width = width
		e.height = height
	})

	e.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			e.window.SetShouldClose(true)
		}
	})
}

// GetWindow returns the GLFW window
func (e *Engine) GetWindow() *glfw.Window {
	return e.window
}

// GetECS returns the ECS world
func (e *Engine) GetECS() *ecs.World {
	return e.ecs
}

// GetRenderer returns the renderer
func (e *Engine) GetRenderer() *graphics.Renderer {
	return e.renderer
}

// GetInput returns the input manager
func (e *Engine) GetInput() *input.Manager {
	return e.input
}

// GetPhysics returns the physics world
func (e *Engine) GetPhysics() *physics.World {
	return e.physics
}
