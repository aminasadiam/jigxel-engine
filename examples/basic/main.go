package main

import (
	"log"
	"time"

	"github.com/aminasadiam/jigxer-engine/pkg/ecs"
	"github.com/aminasadiam/jigxer-engine/pkg/engine"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	// Create engine
	gameEngine := engine.NewEngine("Jigxer Engine - Basic Example", 800, 600)

	// Initialize
	if err := gameEngine.Init(); err != nil {
		log.Fatal("Failed to initialize engine:", err)
	}
	defer gameEngine.Shutdown()

	// Get ECS world
	world := gameEngine.GetECS()

	// Create a rotating triangle entity
	entityID := world.CreateEntity()

	// Add transform component
	transform := ecs.NewTransformComponent(
		mgl32.Vec3{0, 0, 0}, // Position
		mgl32.Vec3{0, 0, 0}, // Rotation
		mgl32.Vec3{1, 1, 1}, // Scale
	)
	world.AddComponent(entityID, transform)

	// Add mesh component
	mesh := ecs.NewMeshComponent("default")
	world.AddComponent(entityID, mesh)

	// Add tag component
	tag := ecs.NewTagComponent("rotating", "triangle")
	world.AddComponent(entityID, tag)

	// Create a simple rotation system
	rotationSystem := &RotationSystem{
		startTime: time.Now(),
	}
	world.AddSystem(rotationSystem)

	// Run the engine
	gameEngine.Run()
}

// RotationSystem rotates entities with the "rotating" tag
type RotationSystem struct {
	startTime time.Time
}

func (rs *RotationSystem) Update(deltaTime float64, world *ecs.World) {
	// Get all entities with the "rotating" tag
	entities := world.GetEntitiesWithComponent("tag")

	for _, entityID := range entities {
		tagComponent := world.GetComponent(entityID, "tag")
		if tagComponent == nil {
			continue
		}

		tag := tagComponent.(*ecs.TagComponent)
		if !tag.HasTag("rotating") {
			continue
		}

		// Get transform component
		transformComponent := world.GetComponent(entityID, "transform")
		if transformComponent == nil {
			continue
		}

		transform := transformComponent.(*ecs.TransformComponent)

		// Calculate rotation based on time
		elapsed := time.Since(rs.startTime).Seconds()
		rotationSpeed := 1.0 // radians per second
		rotation := elapsed * rotationSpeed

		// Apply rotation around Z-axis
		transform.Rotation = mgl32.Vec3{0, 0, float32(rotation)}
	}
}

func (rs *RotationSystem) GetName() string {
	return "RotationSystem"
}
