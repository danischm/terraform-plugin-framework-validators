package listvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.List
		max         int
		expectError bool
	}
	tests := map[string]testCase{
		"List unknown": {
			val: types.ListUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"List null": {
			val: types.ListNull(
				types.StringType,
			),
			expectError: false,
		},
		"List size less than max": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
				},
			),
			max:         2,
			expectError: false,
		},
		"List size equal to max": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			max:         2,
			expectError: false,
		},
		"List size greater than max": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
					types.StringValue("third"),
				},
			),
			max:         2,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ListResponse{}
			SizeAtMost(test.max).ValidateList(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
