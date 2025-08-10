package physics

import (
	"math"
	"sync"
)

// World represents the physics world
type World struct {
	bodies    map[uint64]*RigidBody
	gravity   Vector2
	timeStep  float64
	mutex     sync.RWMutex
}

// Vector2 represents a 2D vector
type Vector2 struct {
	X, Y float64
}

// RigidBody represents a physics body
type RigidBody struct {
	ID       uint64
	Position Vector2
	Velocity Vector2
	Force    Vector2
	Mass     float64
	InverseMass float64
	Width    float64
	Height   float64
	Active   bool
}

// NewWorld creates a new physics world
func NewWorld() *World {
	return &World{
		bodies:   make(map[uint64]*RigidBody),
		gravity:  Vector2{0, -9.81},
		timeStep: 1.0 / 60.0,
	}
}

// AddBody adds a rigid body to the physics world
func (w *World) AddBody(body *RigidBody) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.bodies[body.ID] = body
}

// RemoveBody removes a rigid body from the physics world
func (w *World) RemoveBody(id uint64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	delete(w.bodies, id)
}

// GetBody returns a rigid body by ID
func (w *World) GetBody(id uint64) *RigidBody {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	return w.bodies[id]
}

// SetGravity sets the gravity vector
func (w *World) SetGravity(gravity Vector2) {
	w.gravity = gravity
}

// Update updates the physics simulation
func (w *World) Update(deltaTime float64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	// Update all bodies
	for _, body := range w.bodies {
		if !body.Active {
			continue
		}
		
		// Apply gravity
		body.Force = body.Force.Add(w.gravity.Mul(body.Mass))
		
		// Update velocity
		body.Velocity = body.Velocity.Add(body.Force.Mul(deltaTime).Mul(body.InverseMass))
		
		// Update position
		body.Position = body.Position.Add(body.Velocity.Mul(deltaTime))
		
		// Reset force
		body.Force = Vector2{0, 0}
	}
	
	// Check collisions
	w.checkCollisions()
}

// checkCollisions checks for collisions between all bodies
func (w *World) checkCollisions() {
	bodies := make([]*RigidBody, 0, len(w.bodies))
	for _, body := range w.bodies {
		if body.Active {
			bodies = append(bodies, body)
		}
	}
	
	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			if w.checkCollision(bodies[i], bodies[j]) {
				w.resolveCollision(bodies[i], bodies[j])
			}
		}
	}
}

// checkCollision checks if two bodies are colliding
func (w *World) checkCollision(body1, body2 *RigidBody) bool {
	// Simple AABB collision detection
	left1 := body1.Position.X - body1.Width/2
	right1 := body1.Position.X + body1.Width/2
	top1 := body1.Position.Y + body1.Height/2
	bottom1 := body1.Position.Y - body1.Height/2
	
	left2 := body2.Position.X - body2.Width/2
	right2 := body2.Position.X + body2.Width/2
	top2 := body2.Position.Y + body2.Height/2
	bottom2 := body2.Position.Y - body2.Height/2
	
	return !(right1 < left2 || left1 > right2 || bottom1 > top2 || top1 < bottom2)
}

// resolveCollision resolves a collision between two bodies
func (w *World) resolveCollision(body1, body2 *RigidBody) {
	// Simple collision resolution - just separate the bodies
	// In a real implementation, you'd use proper impulse-based collision response
	
	// Calculate separation vector
	separation := body2.Position.Sub(body1.Position)
	distance := separation.Length()
	
	if distance == 0 {
		separation = Vector2{1, 0}
		distance = 1
	}
	
	// Normalize separation vector
	normal := separation.Div(distance)
	
	// Calculate overlap
	overlap := (body1.Width + body2.Width) / 2 - distance
	
	if overlap > 0 {
		// Move bodies apart
		separationVector := normal.Mul(overlap / 2)
		
		if body1.InverseMass > 0 {
			body1.Position = body1.Position.Sub(separationVector)
		}
		if body2.InverseMass > 0 {
			body2.Position = body2.Position.Add(separationVector)
		}
	}
}

// Vector2 methods
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{v.X + other.X, v.Y + other.Y}
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{v.X - other.X, v.Y - other.Y}
}

func (v Vector2) Mul(scalar float64) Vector2 {
	return Vector2{v.X * scalar, v.Y * scalar}
}

func (v Vector2) Div(scalar float64) Vector2 {
	return Vector2{v.X / scalar, v.Y / scalar}
}

func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// NewRigidBody creates a new rigid body
func NewRigidBody(id uint64, position Vector2, width, height, mass float64) *RigidBody {
	inverseMass := 0.0
	if mass > 0 {
		inverseMass = 1.0 / mass
	}
	
	return &RigidBody{
		ID:          id,
		Position:    position,
		Velocity:    Vector2{0, 0},
		Force:       Vector2{0, 0},
		Mass:        mass,
		InverseMass: inverseMass,
		Width:       width,
		Height:      height,
		Active:      true,
	}
}
