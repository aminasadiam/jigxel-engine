# Jigxer Engine

A modern, modular game engine written in Go with OpenGL rendering, physics simulation, and an Entity-Component-System (ECS) architecture.

## Features

### Core Engine
- **Main Engine Loop**: Efficient game loop with fixed time step
- **Window Management**: GLFW-based window creation and management
- **Scene Management**: Entity and component management system
- **Resource Management**: Centralized resource loading and caching

### Graphics System
- **OpenGL 4.1 Rendering**: Modern OpenGL with shader support
- **Shader Management**: GLSL shader compilation and management
- **Mesh Rendering**: 3D mesh rendering with vertex buffers
- **Camera System**: Perspective and orthographic camera support

### Entity-Component-System (ECS)
- **Entity Management**: Efficient entity creation and destruction
- **Component System**: Flexible component-based architecture
- **System Processing**: Parallel system execution
- **Built-in Components**: Transform, Mesh, Physics, Audio, and Tag components

### Input System
- **Keyboard Input**: Full keyboard support with key state tracking
- **Mouse Input**: Mouse position, buttons, and scroll wheel
- **Input Events**: Press, release, and hold detection
- **Cursor Management**: Cursor mode control (normal, hidden, disabled)

### Physics System
- **Rigid Body Physics**: Mass-based physics simulation
- **Collision Detection**: AABB collision detection
- **Gravity System**: Configurable gravity vector
- **Collision Resolution**: Basic collision response

### Audio System
- **Sound Management**: Audio file loading and playback
- **Volume Control**: Per-sound volume adjustment
- **Looping Support**: Sound looping capabilities
- **Audio Context**: Centralized audio management

## Project Structure

```
jigxer-engine/
├── cmd/
│   └── main.go              # Main entry point
├── internal/
│   ├── engine/
│   │   └── engine.go        # Core engine implementation
│   ├── ecs/
│   │   ├── world.go         # ECS world management
│   │   └── components.go    # Built-in components
│   ├── graphics/
│   │   └── renderer.go      # OpenGL renderer
│   ├── input/
│   │   └── manager.go       # Input management
│   ├── audio/
│   │   └── manager.go       # Audio management
│   └── physics/
│       └── world.go         # Physics simulation
├── go.mod                   # Go module file
└── README.md               # This file
```

## Getting Started

### Prerequisites

- Go 1.24.5 or later
- OpenGL 4.1 compatible graphics card
- GLFW development libraries

### Installation

1. Clone the repository:
```bash
git clone https://github.com/aminasadiam/jigxer-engine.git
cd jigxer-engine
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the example:
```bash
go run cmd/main.go
```

## Usage Example

```go
package main

import (
    "log"
    "github.com/aminasadiam/jigxer-engine/internal/engine"
    "github.com/aminasadiam/jigxer-engine/internal/ecs"
    "github.com/go-gl/mathgl/mgl32"
)

func main() {
    // Create engine
    gameEngine := engine.NewEngine("My Game", 800, 600)
    
    // Initialize
    if err := gameEngine.Init(); err != nil {
        log.Fatal("Failed to initialize engine:", err)
    }
    defer gameEngine.Shutdown()
    
    // Get ECS world
    world := gameEngine.GetECS()
    
    // Create an entity
    entityID := world.CreateEntity()
    
    // Add components
    transform := ecs.NewTransformComponent(
        mgl32.Vec3{0, 0, 0},    // Position
        mgl32.Vec3{0, 0, 0},    // Rotation
        mgl32.Vec3{1, 1, 1},    // Scale
    )
    world.AddComponent(entityID, transform)
    
    // Add mesh component
    mesh := ecs.NewMeshComponent("default")
    world.AddComponent(entityID, mesh)
    
    // Run the engine
    gameEngine.Run()
}
```

## Architecture

### Engine Core
The engine core manages the main game loop, window creation, and system initialization. It coordinates all subsystems and provides a clean API for game development.

### ECS System
The Entity-Component-System provides a flexible architecture for game objects:
- **Entities**: Unique identifiers for game objects
- **Components**: Data containers (Transform, Mesh, Physics, etc.)
- **Systems**: Logic processors that operate on components

### Rendering Pipeline
The graphics system uses modern OpenGL with:
- Vertex Buffer Objects (VBO) for efficient data transfer
- Vertex Array Objects (VAO) for vertex attribute configuration
- Shader programs for programmable rendering pipeline
- Matrix transformations for 3D rendering

### Physics Simulation
The physics system provides:
- Rigid body dynamics with mass and forces
- AABB collision detection
- Basic collision response
- Configurable gravity

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [ ] Advanced physics (constraints, joints)
- [ ] Particle system
- [ ] UI system
- [ ] Networking support
- [ ] Asset pipeline
- [ ] Scene serialization
- [ ] Performance profiling
- [ ] Multi-platform support

## Dependencies

- `github.com/go-gl/gl/v4.1-core/gl` - OpenGL bindings
- `github.com/go-gl/glfw/v3.3/glfw` - Window management
- `github.com/go-gl/mathgl` - Mathematics library
- `github.com/hajimehoshi/oto/v2` - Audio playback
- `github.com/faiface/beep` - Audio format support
