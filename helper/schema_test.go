package helper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// -----------------------------------------------------------------------------
// Test Helpers
// -----------------------------------------------------------------------------

func resourceFoo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// -----------------------------------------------------------------------------
// DataSourceSchemaFromResourceSchema
// -----------------------------------------------------------------------------

// Ensures the copied schema map has the same number of attributes as the
// input schema map
func TestDataSourceSchemaFromResourceSchema_SameLength(t *testing.T) {
	obj := resourceFoo()
	expectedLen := len(obj.Schema)
	actualLen := len(DataSourceSchemaFromResourceSchema(obj.Schema))
	if expectedLen != actualLen {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Expected length [%d], got [%d].",
			expectedLen,
			actualLen,
		)
	}
}

// Ensures the copied schema map has the correct value. Each of the
// schema entries in the map should have the Type and Description
// preserved. The schema should be forced to a computed property.
//
// For sets, the hash function should be preserved. For types that have
// the Elem attribute, if the elem is a schema the type must be the same.
// Otherwise if the elem is a resource, the resource's schema map is
// compared to the expected resource schema mamp.
func TestDataSourceSchemaFromResourceSchema_Value(t *testing.T) {
	obj := resourceFoo()
	actualSchemaMap := DataSourceSchemaFromResourceSchema(obj.Schema)
	assertSchemaMapValue(
		t,
		obj.Schema,
		actualSchemaMap,
	)
}

// Ensures the actual schema map is the expected schema map.
func assertSchemaMapValue(t *testing.T, expectedSchemaMap, actualSchemaMap map[string]*schema.Schema) {
	for key, expectedSchema := range expectedSchemaMap {
		actualSchema, ok := actualSchemaMap[key]
		if !ok {
			t.Fatalf(
				"DataSourceSchemaFromResourceSchema did not return the correct "+
					"output. Expected key [%s] not found.",
				key,
			)
		}
		assertSchemaValue(
			t,
			expectedSchema,
			actualSchema,
			key,
		)
	}
}

// Ensures the actual schema is the expected schema.
func assertSchemaValue(t *testing.T, expectedSchema, actualSchema *schema.Schema, name string) {
	// nil check
	if (expectedSchema == nil && actualSchema != nil) ||
		(expectedSchema != nil && actualSchema == nil) {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Expected schema [%v], actual schema [%v] for [%s]",
			expectedSchema,
			actualSchema,
			name,
		)
	}
	// preserve type, description
	if actualSchema.Type != expectedSchema.Type {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Type [%s] for [%s], got [%s].",
			expectedSchema.Type,
			name,
			actualSchema.Type,
		)
	}
	if actualSchema.Description != expectedSchema.Description {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Description [%s] for [%s], got [%s].",
			expectedSchema.Description,
			name,
			actualSchema.Description,
		)
	}
	// force to computed
	if !actualSchema.Computed {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Computed [%t] for [%s], got [%t].",
			true,
			name,
			actualSchema.Computed,
		)
	}
	if actualSchema.Optional {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Optional [%t] for [%s], got [%t].",
			false,
			name,
			actualSchema.Optional,
		)
	}
	if actualSchema.Required {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Required [%t] for [%s], got [%t].",
			false,
			name,
			actualSchema.Required,
		)
	}
	// preserve set function
	//
	// NOTE(ALL): In Golang, you cannot test for func equality. You can only
	//   check for func equality with nil.
	if (actualSchema.Set == nil && expectedSchema.Set != nil) ||
		(actualSchema.Set != nil && expectedSchema.Set == nil) {
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Set [%v] for [%s], got [%v].",
			expectedSchema.Set,
			name,
			actualSchema.Set,
		)
	}

	// Check Elem for complex & simple types

	if (actualSchema.Elem == nil && expectedSchema.Elem != nil) ||
		(actualSchema.Elem != nil && expectedSchema.Elem == nil) {
		// elem was not set and is the zero value
		t.Fatalf(
			"DataSourceSchemaFromResourceSchema did not return the correct "+
				"output. Exepcted Elem [%v] for [%s], got [%v].",
			expectedSchema.Elem,
			name,
			actualSchema.Elem,
		)
	} else if expectedElem, ok := expectedSchema.Elem.(*schema.Resource); ok {
		// elem is a complex type (ie: a resource), recurse the check
		// on the elem definition

		// Differt assert values
		var actualElem *schema.Resource
		var ok bool
		if actualElem, ok = actualSchema.Elem.(*schema.Resource); !ok {
			t.Fatalf(
				"DataSourceSchemaFromResourceSchema did not return the correct "+
					"output. Expected Elem to be [*schema.Resource], got [%T] for "+
					"[%s].",
				actualSchema.Elem,
				name,
			)
		}
		// recurse
		assertSchemaMapValue(
			t,
			expectedElem.Schema,
			actualElem.Schema,
		)
	} else if expectedElem, ok := expectedSchema.Elem.(*schema.Schema); ok {
		// elem is a simple type (ie: a schema), check the type definition

		// Differt assert values
		var actualElem *schema.Schema
		var ok bool
		if actualElem, ok = actualSchema.Elem.(*schema.Schema); !ok {
			t.Fatalf(
				"DataSourceSchemaFromResourceSchema did not return the correct "+
					"output. Expected Elem to be [*schema.Schema], got [%T] for "+
					"[%s].",
				actualSchema.Elem,
				name,
			)
		}
		// Wrong type
		if expectedElem.Type != actualElem.Type {
			t.Fatalf(
				"DataSourceSchemaFromResourceSchema did not return the correct "+
					"output. Expected Elem to be [*schema.Schema] of type [%s], got [%s] "+
					"for [%s].",
				expectedElem.Type,
				actualElem.Type,
				name,
			)
		}
	}
}
