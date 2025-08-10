package main

import (
	"fmt"

	"github.com/aminasadiam/jigxel-engine/pkg/ecs"
	"github.com/aminasadiam/jigxel-engine/pkg/physics"
)

func main() {
	fmt.Println("jigxel Engine - Basic Demo")

	// Create ECS world
	world := ecs.NewWorld()

	// Create a simple entity
	entityID := world.CreateEntity()

	// Add a tag component
	tag := ecs.NewTagComponent("player", "game_object")
	world.AddComponent(entityID, tag)

	// Create physics world
	physicsWorld := physics.NewWorld()

	// Create a physics body
	body := physics.NewRigidBody(1, physics.Vector2{0, 0}, 1.0, 1.0, 1.0)
	physicsWorld.AddBody(body)

	// Add physics component to entity
	physicsComp := ecs.NewPhysicsComponent(1, 1.0)
	world.AddComponent(entityID, physicsComp)

	fmt.Printf("Created entity %d with %d components\n", entityID, world.GetEntityCount())
	fmt.Printf("Physics world has body at position (%.2f, %.2f)\n", physicsWorld.GetBody(1).Position.X, physicsWorld.GetBody(1).Position.Y)

	// Simulate a few physics steps
	for i := 0; i < 10; i++ {
		physicsWorld.Update(1.0 / 60.0)
		body := physicsWorld.GetBody(1)
		fmt.Printf("Step %d: Position (%.2f, %.2f)\n", i, body.Position.X, body.Position.Y)
	}

	fmt.Println("Engine demo completed successfully!")
}
