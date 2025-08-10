package graphics

import (
	"fmt"
	"log"

	"github.com/aminasadiam/jigxer-engine/pkg/ecs"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Renderer handles all rendering operations
type Renderer struct {
	shaders map[string]*Shader
	meshes  map[string]*Mesh
}

// Shader represents an OpenGL shader program
type Shader struct {
	ID uint32
}

// Mesh represents a 3D mesh
type Mesh struct {
	VAO         uint32
	VBO         uint32
	EBO         uint32
	VertexCount int32
}

// NewRenderer creates a new renderer
func NewRenderer() *Renderer {
	return &Renderer{
		shaders: make(map[string]*Shader),
		meshes:  make(map[string]*Mesh),
	}
}

// Init initializes the renderer
func (r *Renderer) Init() error {
	// Enable depth testing
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Enable blending
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Create default shaders
	if err := r.createDefaultShaders(); err != nil {
		return err
	}

	log.Println("Renderer initialized successfully")
	return nil
}

// Render renders the current scene
func (r *Renderer) Render(world *ecs.World) {
	// Get default shader
	shader, exists := r.shaders["default"]
	if !exists {
		return
	}

	shader.Use()

	// Set up projection matrix
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), 800.0/600.0, 0.1, 100.0)
	shader.SetMat4("projection", projection)

	// Set up view matrix
	view := mgl32.LookAtV(
		mgl32.Vec3{0, 0, 3},
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 1, 0},
	)
	shader.SetMat4("view", view)

	// Render entities with transform and mesh components
	// This is a simplified version - in a real engine you'd have proper component types
	entities := world.GetEntitiesWithComponent("transform")
	for range entities {
		// Set model matrix
		model := mgl32.Ident4()
		shader.SetMat4("model", model)

		// Render mesh if available
		if mesh, exists := r.meshes["default"]; exists {
			r.renderMesh(mesh)
		}
	}
}

// Shutdown cleans up the renderer
func (r *Renderer) Shutdown() {
	// Clean up shaders
	for _, shader := range r.shaders {
		gl.DeleteProgram(shader.ID)
	}

	// Clean up meshes
	for _, mesh := range r.meshes {
		gl.DeleteVertexArrays(1, &mesh.VAO)
		gl.DeleteBuffers(1, &mesh.VBO)
		gl.DeleteBuffers(1, &mesh.EBO)
	}
}

// createDefaultShaders creates the default shaders
func (r *Renderer) createDefaultShaders() error {
	vertexShaderSource := `
		#version 410 core
		layout (location = 0) in vec3 aPos;
		layout (location = 1) in vec3 aColor;
		
		out vec3 ourColor;
		
		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;
		
		void main()
		{
			gl_Position = projection * view * model * vec4(aPos, 1.0);
			ourColor = aColor;
		}
	` + "\x00"

	fragmentShaderSource := `
		#version 410 core
		out vec4 FragColor;
		in vec3 ourColor;
		
		void main()
		{
			FragColor = vec4(ourColor, 1.0);
		}
	` + "\x00"

	shader, err := NewShader(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return err
	}

	r.shaders["default"] = shader

	// Create a simple triangle mesh
	r.createDefaultMesh()

	return nil
}

// createDefaultMesh creates a simple triangle mesh
func (r *Renderer) createDefaultMesh() {
	vertices := []float32{
		// positions        // colors
		-0.5, -0.5, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0,
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	r.meshes["default"] = &Mesh{
		VAO:         VAO,
		VBO:         VBO,
		EBO:         0,
		VertexCount: 3,
	}
}

// renderMesh renders a mesh
func (r *Renderer) renderMesh(mesh *Mesh) {
	gl.BindVertexArray(mesh.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, mesh.VertexCount)
	gl.BindVertexArray(0)
}

// NewShader creates a new shader program
func NewShader(vertexSource, fragmentSource string) (*Shader, error) {
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var success int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var infoLog [512]byte
		gl.GetProgramInfoLog(program, 512, nil, &infoLog[0])
		return nil, fmt.Errorf("shader linking failed: %s", string(infoLog[:]))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{ID: program}, nil
}

// Use activates the shader
func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

// SetMat4 sets a mat4 uniform
func (s *Shader) SetMat4(name string, value mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &value[0])
}

// compileShader compiles a shader
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var infoLog [512]byte
		gl.GetShaderInfoLog(shader, 512, nil, &infoLog[0])
		return 0, fmt.Errorf("shader compilation failed: %s", string(infoLog[:]))
	}

	return shader, nil
}
