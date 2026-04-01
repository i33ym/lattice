package lattice

// UDF represents a reference to a user-defined function.
type UDF struct {
	Name     string
	Endpoint string
	Version  string
}

// ComputedSpec defines a computed column's UDF and input dependencies.
type ComputedSpec struct {
	UDF    UDF
	Inputs []string
}

// Computed creates a ComputedSpec from a UDF and input column names.
func Computed(udf UDF, inputs ...string) *ComputedSpec {
	return &ComputedSpec{UDF: udf, Inputs: inputs}
}
