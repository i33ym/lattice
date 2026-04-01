package query

// Op represents a comparison operator for query filters.
type Op int

const (
	OpEq Op = iota
	OpGT
	OpLT
	OpGTE
	OpLTE
)

// Filter represents a parsed query filter on a single column.
type Filter struct {
	Column string
	Op     Op
	Value  any
}

// Query is the built query ready for execution.
type Query struct {
	Table   string
	Filters []Filter
	Limit   int
	Offset  int
}

// Builder constructs queries using a chainable API.
type Builder struct {
	table   string
	filters []Filter
	limit   int
	offset  int
}

// NewBuilder creates a Builder targeting the given table.
func NewBuilder(table string) *Builder {
	return &Builder{table: table}
}

// Where adds a filter condition and returns the builder for chaining.
func (b *Builder) Where(column string, op Op, value any) *Builder {
	b.filters = append(b.filters, Filter{Column: column, Op: op, Value: value})
	return b
}

// Limit sets the maximum number of results and returns the builder for chaining.
func (b *Builder) Limit(n int) *Builder {
	b.limit = n
	return b
}

// Offset sets the number of results to skip and returns the builder for chaining.
func (b *Builder) Offset(n int) *Builder {
	b.offset = n
	return b
}

// Build produces the final Query from the builder state.
func (b *Builder) Build() Query {
	filters := make([]Filter, len(b.filters))
	copy(filters, b.filters)
	return Query{
		Table:   b.table,
		Filters: filters,
		Limit:   b.limit,
		Offset:  b.offset,
	}
}
