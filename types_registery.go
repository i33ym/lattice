package lattice

import "fmt"

// TypeRegistry allows registering custom types at runtime.
type TypeRegistry struct {
	types map[string]Type
}

// NewTypeRegistry creates a new empty TypeRegistry.
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{types: make(map[string]Type)}
}

// Register adds a type to the registry, returning an error if the name is already taken.
func (r *TypeRegistry) Register(t Type) error {
	name := t.Name()
	if _, exists := r.types[name]; exists {
		return fmt.Errorf("registering type %q: %w", name, ErrTypeMismatch)
	}
	r.types[name] = t
	return nil
}

// Lookup returns the type with the given name, if it exists.
func (r *TypeRegistry) Lookup(name string) (Type, bool) {
	t, ok := r.types[name]
	return t, ok
}

// All returns all registered types.
func (r *TypeRegistry) All() []Type {
	result := make([]Type, 0, len(r.types))
	for _, t := range r.types {
		result = append(result, t)
	}
	return result
}
