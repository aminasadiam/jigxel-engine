package ecs

import (
	"sync"
)

// EntityID represents a unique entity identifier
type EntityID uint64

// Component represents a component interface
type Component interface {
	GetType() string
}

// System represents a system interface
type System interface {
	Update(deltaTime float64, world *World)
	GetName() string
}

// World represents the ECS world
type World struct {
	entities     map[EntityID]*Entity
	components   map[string][]Component
	systems      []System
	nextEntityID EntityID
	mutex        sync.RWMutex
}

// Entity represents a game entity
type Entity struct {
	ID         EntityID
	Components map[string]Component
	Active     bool
}

// NewWorld creates a new ECS world
func NewWorld() *World {
	return &World{
		entities:   make(map[EntityID]*Entity),
		components: make(map[string][]Component),
		systems:    make([]System, 0),
	}
}

// CreateEntity creates a new entity
func (w *World) CreateEntity() EntityID {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	entityID := w.nextEntityID
	w.nextEntityID++

	entity := &Entity{
		ID:         entityID,
		Components: make(map[string]Component),
		Active:     true,
	}

	w.entities[entityID] = entity
	return entityID
}

// DestroyEntity destroys an entity
func (w *World) DestroyEntity(entityID EntityID) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if entity, exists := w.entities[entityID]; exists {
		// Remove all components
		for componentType := range entity.Components {
			w.removeComponentFromList(entityID, componentType)
		}

		// Remove entity
		delete(w.entities, entityID)
	}
}

// AddComponent adds a component to an entity
func (w *World) AddComponent(entityID EntityID, component Component) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if entity, exists := w.entities[entityID]; exists {
		componentType := component.GetType()
		entity.Components[componentType] = component

		// Add to component list
		if w.components[componentType] == nil {
			w.components[componentType] = make([]Component, 0)
		}
		w.components[componentType] = append(w.components[componentType], component)
	}
}

// RemoveComponent removes a component from an entity
func (w *World) RemoveComponent(entityID EntityID, componentType string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if entity, exists := w.entities[entityID]; exists {
		if _, hasComponent := entity.Components[componentType]; hasComponent {
			delete(entity.Components, componentType)
			w.removeComponentFromList(entityID, componentType)
		}
	}
}

// GetComponent gets a component from an entity
func (w *World) GetComponent(entityID EntityID, componentType string) Component {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if entity, exists := w.entities[entityID]; exists {
		return entity.Components[componentType]
	}
	return nil
}

// GetEntitiesWithComponent gets all entities that have a specific component
func (w *World) GetEntitiesWithComponent(componentType string) []EntityID {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	var entities []EntityID
	for entityID, entity := range w.entities {
		if _, hasComponent := entity.Components[componentType]; hasComponent {
			entities = append(entities, entityID)
		}
	}
	return entities
}

// AddSystem adds a system to the world
func (w *World) AddSystem(system System) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.systems = append(w.systems, system)
}

// RemoveSystem removes a system from the world
func (w *World) RemoveSystem(systemName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for i, system := range w.systems {
		if system.GetName() == systemName {
			w.systems = append(w.systems[:i], w.systems[i+1:]...)
			break
		}
	}
}

// Update updates all systems
func (w *World) Update(deltaTime float64) {
	w.mutex.RLock()
	systems := make([]System, len(w.systems))
	copy(systems, w.systems)
	w.mutex.RUnlock()

	for _, system := range systems {
		system.Update(deltaTime, w)
	}
}

// GetEntityCount returns the number of entities
func (w *World) GetEntityCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return len(w.entities)
}

// GetSystemCount returns the number of systems
func (w *World) GetSystemCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return len(w.systems)
}

// removeComponentFromList removes a component from the component list
func (w *World) removeComponentFromList(entityID EntityID, componentType string) {
	if components, exists := w.components[componentType]; exists {
		for i, component := range components {
			// This is a simplified removal - in a real implementation,
			// you'd want to store entity references with components
			if component != nil {
				w.components[componentType] = append(components[:i], components[i+1:]...)
				break
			}
		}
	}
}
