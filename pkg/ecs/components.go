package ecs

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TransformComponent represents position, rotation, and scale
type TransformComponent struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}

func (t *TransformComponent) GetType() string {
	return "transform"
}

// NewTransformComponent creates a new transform component
func NewTransformComponent(position, rotation, scale mgl32.Vec3) *TransformComponent {
	return &TransformComponent{
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	}
}

// MeshComponent represents a 3D mesh
type MeshComponent struct {
	MeshID string
	Visible bool
}

func (m *MeshComponent) GetType() string {
	return "mesh"
}

// NewMeshComponent creates a new mesh component
func NewMeshComponent(meshID string) *MeshComponent {
	return &MeshComponent{
		MeshID:  meshID,
		Visible: true,
	}
}

// PhysicsComponent represents physics properties
type PhysicsComponent struct {
	BodyID uint64
	Mass   float64
	Active bool
}

func (p *PhysicsComponent) GetType() string {
	return "physics"
}

// NewPhysicsComponent creates a new physics component
func NewPhysicsComponent(bodyID uint64, mass float64) *PhysicsComponent {
	return &PhysicsComponent{
		BodyID: bodyID,
		Mass:   mass,
		Active: true,
	}
}

// AudioComponent represents audio properties
type AudioComponent struct {
	SoundID string
	Volume  float64
	Loop    bool
}

func (a *AudioComponent) GetType() string {
	return "audio"
}

// NewAudioComponent creates a new audio component
func NewAudioComponent(soundID string, volume float64, loop bool) *AudioComponent {
	return &AudioComponent{
		SoundID: soundID,
		Volume:  volume,
		Loop:    loop,
	}
}

// TagComponent represents entity tags
type TagComponent struct {
	Tags []string
}

func (t *TagComponent) GetType() string {
	return "tag"
}

// NewTagComponent creates a new tag component
func NewTagComponent(tags ...string) *TagComponent {
	return &TagComponent{
		Tags: tags,
	}
}

// AddTag adds a tag to the component
func (t *TagComponent) AddTag(tag string) {
	t.Tags = append(t.Tags, tag)
}

// HasTag checks if the component has a specific tag
func (t *TagComponent) HasTag(tag string) bool {
	for _, t := range t.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
