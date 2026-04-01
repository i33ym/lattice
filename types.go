package lattice

import "fmt"

// Type describes a column data type in Lattice.
type Type interface {
	Name() string
	Equal(Type) bool
}

type stringType struct{}
type intType struct{}
type floatType struct{}
type boolType struct{}
type imageType struct{}
type audioType struct{}
type videoType struct{}
type jsonType struct{}
type timestampType struct{}

var (
	// StringType is the type for string columns.
	StringType Type = stringType{}
	// IntType is the type for integer columns.
	IntType Type = intType{}
	// FloatType is the type for float columns.
	FloatType Type = floatType{}
	// BoolType is the type for boolean columns.
	BoolType Type = boolType{}
	// ImageType is the type for image columns.
	ImageType Type = imageType{}
	// AudioType is the type for audio columns.
	AudioType Type = audioType{}
	// VideoType is the type for video columns.
	VideoType Type = videoType{}
	// JSONType is the type for JSON columns.
	JSONType Type = jsonType{}
	// TimestampType is the type for timestamp columns.
	TimestampType Type = timestampType{}
)

func (stringType) Name() string         { return "string" }
func (stringType) Equal(t Type) bool    { _, ok := t.(stringType); return ok }
func (intType) Name() string            { return "int" }
func (intType) Equal(t Type) bool       { _, ok := t.(intType); return ok }
func (floatType) Name() string          { return "float" }
func (floatType) Equal(t Type) bool     { _, ok := t.(floatType); return ok }
func (boolType) Name() string           { return "bool" }
func (boolType) Equal(t Type) bool      { _, ok := t.(boolType); return ok }
func (imageType) Name() string          { return "image" }
func (imageType) Equal(t Type) bool     { _, ok := t.(imageType); return ok }
func (audioType) Name() string          { return "audio" }
func (audioType) Equal(t Type) bool     { _, ok := t.(audioType); return ok }
func (videoType) Name() string          { return "video" }
func (videoType) Equal(t Type) bool     { _, ok := t.(videoType); return ok }
func (jsonType) Name() string           { return "json" }
func (jsonType) Equal(t Type) bool      { _, ok := t.(jsonType); return ok }
func (timestampType) Name() string      { return "timestamp" }
func (timestampType) Equal(t Type) bool { _, ok := t.(timestampType); return ok }

type vectorType struct{ dim int }

// VectorType returns a vector type with the given dimension.
func VectorType(dim int) Type { return vectorType{dim: dim} }

func (v vectorType) Name() string   { return fmt.Sprintf("vector(%d)", v.dim) }
func (v vectorType) Equal(t Type) bool {
	other, ok := t.(vectorType)
	return ok && other.dim == v.dim
}

// Dim returns the dimension of the vector type.
func (v vectorType) Dim() int { return v.dim }

type listType struct{ elem Type }

// ListOf returns a list type wrapping an element type.
func ListOf(elem Type) Type { return listType{elem: elem} }

func (l listType) Name() string { return fmt.Sprintf("list(%s)", l.elem.Name()) }
func (l listType) Equal(t Type) bool {
	other, ok := t.(listType)
	return ok && other.elem.Equal(l.elem)
}
