package lattice

import (
	"errors"
	"testing"
)

func TestTableSchemaValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		schema    TableSchema
		wantErr   error
		wantNoErr bool
	}{
		{
			"Validate/validSchema/succeeds",
			TableSchema{
				Name: "users",
				Columns: []ColumnSpec{
					{Name: "id", Type: StringType},
					{Name: "name", Type: StringType},
				},
			},
			nil,
			true,
		},
		{
			"Validate/emptyTableName/returnsErrInvalidSchema",
			TableSchema{
				Name:    "",
				Columns: []ColumnSpec{{Name: "id", Type: StringType}},
			},
			ErrInvalidSchema,
			false,
		},
		{
			"Validate/noColumns/returnsErrInvalidSchema",
			TableSchema{
				Name:    "empty",
				Columns: nil,
			},
			ErrInvalidSchema,
			false,
		},
		{
			"Validate/duplicateColumnNames/returnsErrDuplicateColumn",
			TableSchema{
				Name: "dupes",
				Columns: []ColumnSpec{
					{Name: "id", Type: StringType},
					{Name: "id", Type: IntType},
				},
			},
			ErrDuplicateColumn,
			false,
		},
		{
			"Validate/computedReferencesNonExistentColumn/returnsErrColumnNotFound",
			TableSchema{
				Name: "computed",
				Columns: []ColumnSpec{
					{Name: "id", Type: StringType},
					{Name: "derived", Computed: Computed(UDF{Name: "fn"}, "missing_col")},
				},
			},
			ErrColumnNotFound,
			false,
		},
		{
			"Validate/computedReferencesExistingColumn/succeeds",
			TableSchema{
				Name: "computed_ok",
				Columns: []ColumnSpec{
					{Name: "id", Type: StringType},
					{Name: "audio", Type: AudioType},
					{Name: "transcript", Computed: Computed(UDF{Name: "transcribe"}, "audio")},
				},
			},
			nil,
			true,
		},
		{
			"Validate/nilTypeOnNonComputedColumn/returnsErrInvalidSchema",
			TableSchema{
				Name: "niltype",
				Columns: []ColumnSpec{
					{Name: "id", Type: nil},
				},
			},
			ErrInvalidSchema,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.schema.Validate()
			if tt.wantNoErr {
				if err != nil {
					t.Fatalf("Validate() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() expected error wrapping %v, got nil", tt.wantErr)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want errors.Is %v", err, tt.wantErr)
			}
		})
	}
}
