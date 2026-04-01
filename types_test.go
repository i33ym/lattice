package lattice

import (
	"testing"
)

func TestTypeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		typ      Type
		wantName string
	}{
		{"StringType/returns/string", StringType, "string"},
		{"IntType/returns/int", IntType, "int"},
		{"FloatType/returns/float", FloatType, "float"},
		{"BoolType/returns/bool", BoolType, "bool"},
		{"ImageType/returns/image", ImageType, "image"},
		{"AudioType/returns/audio", AudioType, "audio"},
		{"VideoType/returns/video", VideoType, "video"},
		{"JSONType/returns/json", JSONType, "json"},
		{"TimestampType/returns/timestamp", TimestampType, "timestamp"},
		{"VectorType/returns/vector(128)", VectorType(128), "vector(128)"},
		{"ListOf/returns/list(string)", ListOf(StringType), "list(string)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.typ.Name(); got != tt.wantName {
				t.Fatalf("Name() = %q, want %q", got, tt.wantName)
			}
		})
	}
}

func TestTypeEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    Type
		b    Type
		want bool
	}{
		{"SameType/StringType/equal", StringType, StringType, true},
		{"SameType/IntType/equal", IntType, IntType, true},
		{"DifferentType/StringVsInt/notEqual", StringType, IntType, false},
		{"DifferentType/BoolVsFloat/notEqual", BoolType, FloatType, false},
		{"VectorType/sameDim/equal", VectorType(128), VectorType(128), true},
		{"VectorType/differentDim/notEqual", VectorType(128), VectorType(256), false},
		{"VectorType/vsScalar/notEqual", VectorType(128), IntType, false},
		{"ListOf/sameElem/equal", ListOf(StringType), ListOf(StringType), true},
		{"ListOf/differentElem/notEqual", ListOf(StringType), ListOf(IntType), false},
		{"ListOf/vsScalar/notEqual", ListOf(StringType), StringType, false},
		{"ListOf/nested/equal", ListOf(ListOf(IntType)), ListOf(ListOf(IntType)), true},
		{"ListOf/nestedDifferent/notEqual", ListOf(ListOf(IntType)), ListOf(ListOf(FloatType)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Fatalf("%s.Equal(%s) = %v, want %v", tt.a.Name(), tt.b.Name(), got, tt.want)
			}
		})
	}
}

func TestVectorTypeDim(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dim  int
		want int
	}{
		{"dim128/returns/128", 128, 128},
		{"dim256/returns/256", 256, 256},
		{"dim1/returns/1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vt := VectorType(tt.dim)
			typed, ok := vt.(vectorType)
			if !ok {
				t.Fatalf("VectorType(%d) did not return vectorType", tt.dim)
			}
			if got := typed.Dim(); got != tt.want {
				t.Fatalf("Dim() = %d, want %d", got, tt.want)
			}
		})
	}
}
