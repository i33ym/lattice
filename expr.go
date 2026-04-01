package lattice

// Expr represents a filter expression for queries.
type Expr interface {
	expr()
}

// And combines expressions with logical AND.
type And struct{ Exprs []Expr }

// Or combines expressions with logical OR.
type Or struct{ Exprs []Expr }

// Eq checks column equality.
type Eq struct {
	Column string
	Value  any
}

// GT checks if a column value is greater than the given value.
type GT struct {
	Column string
	Value  any
}

// LT checks if a column value is less than the given value.
type LT struct {
	Column string
	Value  any
}

// GTE checks if a column value is greater than or equal to the given value.
type GTE struct {
	Column string
	Value  any
}

// LTE checks if a column value is less than or equal to the given value.
type LTE struct {
	Column string
	Value  any
}

// JSONPath filters on a JSON column.
type JSONPath struct {
	Column string
	Path   string
	Value  any
}

// VectorSearch performs similarity search.
type VectorSearch struct {
	Column string
	Vector []float32
	K      int
}

func (And) expr()          {}
func (Or) expr()           {}
func (Eq) expr()           {}
func (GT) expr()           {}
func (LT) expr()           {}
func (GTE) expr()          {}
func (LTE) expr()          {}
func (JSONPath) expr()     {}
func (VectorSearch) expr() {}

// NewAnd creates an And expression from the given sub-expressions.
func NewAnd(exprs ...Expr) And { return And{Exprs: exprs} }

// NewOr creates an Or expression from the given sub-expressions.
func NewOr(exprs ...Expr) Or { return Or{Exprs: exprs} }

// NewEq creates an equality expression.
func NewEq(column string, value any) Eq { return Eq{Column: column, Value: value} }

// NewGT creates a greater-than expression.
func NewGT(column string, value any) GT { return GT{Column: column, Value: value} }

// NewLT creates a less-than expression.
func NewLT(column string, value any) LT { return LT{Column: column, Value: value} }

// NewGTE creates a greater-than-or-equal expression.
func NewGTE(column string, value any) GTE { return GTE{Column: column, Value: value} }

// NewLTE creates a less-than-or-equal expression.
func NewLTE(column string, value any) LTE { return LTE{Column: column, Value: value} }

// NewJSONPath creates a JSON path filter expression.
func NewJSONPath(column, path string, value any) JSONPath {
	return JSONPath{Column: column, Path: path, Value: value}
}

// NewVectorSearch creates a vector similarity search expression.
func NewVectorSearch(column string, vector []float32, k int) VectorSearch {
	return VectorSearch{Column: column, Vector: vector, K: k}
}
