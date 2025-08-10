package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Manager handles all input operations
type Manager struct {
	window *glfw.Window

	// Keyboard state
	keys     map[glfw.Key]bool
	prevKeys map[glfw.Key]bool

	// Mouse state
	mousePos struct {
		x, y float64
	}
	prevMousePos struct {
		x, y float64
	}
	mouseButtons     map[glfw.MouseButton]bool
	prevMouseButtons map[glfw.MouseButton]bool

	// Mouse scroll
	scrollX, scrollY float64
}

// NewManager creates a new input manager
func NewManager(window *glfw.Window) *Manager {
	return &Manager{
		window:           window,
		keys:             make(map[glfw.Key]bool),
		prevKeys:         make(map[glfw.Key]bool),
		mouseButtons:     make(map[glfw.MouseButton]bool),
		prevMouseButtons: make(map[glfw.MouseButton]bool),
	}
}

// Init initializes the input manager
func (m *Manager) Init() error {
	// Set up input callbacks
	m.window.SetKeyCallback(m.keyCallback)
	m.window.SetMouseButtonCallback(m.mouseButtonCallback)
	m.window.SetCursorPosCallback(m.cursorPosCallback)
	m.window.SetScrollCallback(m.scrollCallback)

	return nil
}

// Update updates the input state
func (m *Manager) Update() {
	// Update previous states
	for key := range m.keys {
		m.prevKeys[key] = m.keys[key]
	}

	for button := range m.mouseButtons {
		m.prevMouseButtons[button] = m.mouseButtons[button]
	}

	// Update mouse position
	m.prevMousePos.x = m.mousePos.x
	m.prevMousePos.y = m.mousePos.y

	// Reset scroll
	m.scrollX = 0
	m.scrollY = 0
}

// IsKeyPressed returns true if a key is currently pressed
func (m *Manager) IsKeyPressed(key glfw.Key) bool {
	return m.keys[key]
}

// IsKeyJustPressed returns true if a key was just pressed this frame
func (m *Manager) IsKeyJustPressed(key glfw.Key) bool {
	return m.keys[key] && !m.prevKeys[key]
}

// IsKeyJustReleased returns true if a key was just released this frame
func (m *Manager) IsKeyJustReleased(key glfw.Key) bool {
	return !m.keys[key] && m.prevKeys[key]
}

// IsMouseButtonPressed returns true if a mouse button is currently pressed
func (m *Manager) IsMouseButtonPressed(button glfw.MouseButton) bool {
	return m.mouseButtons[button]
}

// IsMouseButtonJustPressed returns true if a mouse button was just pressed this frame
func (m *Manager) IsMouseButtonJustPressed(button glfw.MouseButton) bool {
	return m.mouseButtons[button] && !m.prevMouseButtons[button]
}

// IsMouseButtonJustReleased returns true if a mouse button was just released this frame
func (m *Manager) IsMouseButtonJustReleased(button glfw.MouseButton) bool {
	return !m.mouseButtons[button] && m.prevMouseButtons[button]
}

// GetMousePosition returns the current mouse position
func (m *Manager) GetMousePosition() (float64, float64) {
	return m.mousePos.x, m.mousePos.y
}

// GetMouseDelta returns the mouse movement delta since last frame
func (m *Manager) GetMouseDelta() (float64, float64) {
	return m.mousePos.x - m.prevMousePos.x, m.mousePos.y - m.prevMousePos.y
}

// GetScroll returns the scroll delta since last frame
func (m *Manager) GetScroll() (float64, float64) {
	return m.scrollX, m.scrollY
}

// SetCursorMode sets the cursor mode
func (m *Manager) SetCursorMode(mode int) {
	m.window.SetInputMode(glfw.CursorMode, mode)
}

// Callbacks
func (m *Manager) keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		m.keys[key] = true
	} else if action == glfw.Release {
		m.keys[key] = false
	}
}

func (m *Manager) mouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		m.mouseButtons[button] = true
	} else if action == glfw.Release {
		m.mouseButtons[button] = false
	}
}

func (m *Manager) cursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	m.mousePos.x = xpos
	m.mousePos.y = ypos
}

func (m *Manager) scrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	m.scrollX = xoffset
	m.scrollY = yoffset
}
