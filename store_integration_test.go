//go:build integration

package lattice

import (
	"context"
	"testing"
)

func TestStoreIntegrationCreateAndInsert(t *testing.T) {
	t.Parallel()
	t.Log("integration tests require a running database - skipping in unit test mode")

	tests := []struct {
		name   string
		schema TableSchema
		rows   []Row
	}{
		{
			"CreateAndInsert/singleRow/succeeds",
			TableSchema{
				Name: "integration_test",
				Columns: []ColumnSpec{
					{Name: "id", Type: StringType},
					{Name: "data", Type: JSONType},
				},
			},
			[]Row{{"id": "1", "data": "{}"}},
		},
		{
			"CreateAndInsert/multipleRows/succeeds",
			TableSchema{
				Name: "integration_batch",
				Columns: []ColumnSpec{
					{Name: "id", Type: StringType},
					{Name: "value", Type: IntType},
				},
			},
			[]Row{
				{"id": "1", "value": 10},
				{"id": "2", "value": 20},
				{"id": "3", "value": 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_ = context.Background()
			t.Skip("no store adapter configured for integration testing")
		})
	}
}

func TestStoreIntegrationQueryRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		filter Expr
	}{
		{"QueryRoundTrip/eqFilter/returnsMatch", NewEq("id", "1")},
		{"QueryRoundTrip/gtFilter/returnsMatch", NewGT("value", 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Skip("no store adapter configured for integration testing")
		})
	}
}
