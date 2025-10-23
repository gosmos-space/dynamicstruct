package dynamicstruct_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gosmos-space/dynamicstruct"
)

func TestAddField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		fieldType interface{}
		wantErr   error
	}{
		{
			name:      "add_string_field",
			fieldName: "Name",
			fieldType: "",
			wantErr:   nil,
		},
		{
			name:      "add_int_field",
			fieldName: "Age",
			fieldType: int(0),
			wantErr:   nil,
		},
		{
			name:      "add_bool_field",
			fieldName: "IsActive",
			fieldType: false,
			wantErr:   nil,
		},
		{
			name:      "add_slice_field",
			fieldName: "Tags",
			fieldType: []string{},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				builder := dynamicstruct.New()
				err := builder.AddField(tt.fieldName, tt.fieldType)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("AddField() error = %v, wantErr %v", err, tt.wantErr)
				}

				// Try adding same field again - should get already exists error
				err = builder.AddField(tt.fieldName, tt.fieldType)
				if !errors.Is(err, dynamicstruct.ErrFieldAlreadyExists) {
					t.Errorf(
						"AddField() with duplicate field error = %v, want %v",
						err,
						dynamicstruct.ErrFieldAlreadyExists,
					)
				}
			},
		)
	}

	// Test adding field after build
	t.Run(
		"add_field_after_build", func(t *testing.T) {
			builder := dynamicstruct.New()
			err := builder.AddField("Test", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			_, err = builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			err = builder.AddField("AnotherField", "")
			if !errors.Is(err, dynamicstruct.ErrInstanceAlreadyBuilt) {
				t.Errorf(
					"AddField() after build error = %v, want %v",
					err,
					dynamicstruct.ErrInstanceAlreadyBuilt,
				)
			}
		},
	)
}

func TestRemoveField(t *testing.T) {
	t.Run(
		"remove_existing_field", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add a field
			err := builder.AddField("Test", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Remove the field
			err = builder.RemoveField("Test")
			if err != nil {
				t.Errorf("RemoveField() error = %v, wantErr nil", err)
			}

			// Adding the same field again should work if removal was successful
			err = builder.AddField("Test", "")
			if err != nil {
				t.Errorf("AddField() after remove error = %v, wantErr nil", err)
			}
		},
	)

	t.Run(
		"remove_nonexistent_field", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Remove a field that doesn't exist
			err := builder.RemoveField("NonExistent")
			if err != nil {
				t.Errorf("RemoveField() error = %v, wantErr nil", err)
			}
		},
	)

	t.Run(
		"remove_field_after_build", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add a field
			err := builder.AddField("Test", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Build the struct
			_, err = builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Try to remove the field after build
			err = builder.RemoveField("Test")
			if !errors.Is(err, dynamicstruct.ErrInstanceAlreadyBuilt) {
				t.Errorf(
					"RemoveField() after build error = %v, want %v",
					err,
					dynamicstruct.ErrInstanceAlreadyBuilt,
				)
			}
		},
	)
}

func TestBuild(t *testing.T) {
	t.Run(
		"build_empty_struct", func(t *testing.T) {
			builder := dynamicstruct.New()
			instance, err := builder.Build()
			if err != nil {
				t.Errorf("Build() empty struct error = %v, wantErr nil", err)
			}
			if instance == nil {
				t.Error("Build() empty struct returned nil instance")
			}
		},
	)

	t.Run(
		"build_with_fields", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields
			_ = builder.AddField("Name", "")
			_ = builder.AddField("Age", int(0))
			_ = builder.AddField("IsActive", false)

			// Build the struct
			instance, err := builder.Build()
			if err != nil {
				t.Errorf("Build() with fields error = %v, wantErr nil", err)
			}
			if instance == nil {
				t.Error("Build() with fields returned nil instance")
			}
		},
	)

	t.Run(
		"build_twice", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Build once
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("First Build() error = %v", err)
			}

			// Build again
			_, err = builder.Build()
			if !errors.Is(err, dynamicstruct.ErrInstanceAlreadyBuilt) {
				t.Errorf(
					"Second Build() error = %v, want %v",
					err,
					dynamicstruct.ErrInstanceAlreadyBuilt,
				)
			}
		},
	)

	t.Run(
		"build_after_reset", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Build once
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("First Build() error = %v", err)
			}

			// Reset
			builder.Reset()

			// Build again
			_, err = builder.Build()
			if err != nil {
				t.Errorf("Build() after reset error = %v, wantErr nil", err)
			}
		},
	)
}

func TestGetFieldValue(t *testing.T) {
	t.Run(
		"get_before_build", func(t *testing.T) {
			builder := dynamicstruct.New()
			var val string
			err := builder.GetFieldValue("Test", &val)
			if !errors.Is(err, dynamicstruct.ErrInstanceNotBuilt) {
				t.Errorf(
					"GetFieldValue() before build error = %v, want %v",
					err,
					dynamicstruct.ErrInstanceNotBuilt,
				)
			}
		},
	)

	t.Run(
		"get_with_non_pointer", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			var val string
			err = builder.GetFieldValue("Name", val) // Note: not passing pointer
			if !errors.Is(err, dynamicstruct.ErrValueMustBePointer) {
				t.Errorf(
					"GetFieldValue() with non-pointer error = %v, want %v",
					err,
					dynamicstruct.ErrValueMustBePointer,
				)
			}
		},
	)

	t.Run(
		"get_with_nil_pointer", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Create a nil pointer of type *string
			var nilPtr *string
			err = builder.GetFieldValue("Name", nilPtr)
			if !errors.Is(err, dynamicstruct.ErrValueCannotBeNil) {
				t.Errorf(
					"GetFieldValue() with nil pointer error = %v, want %v",
					err,
					dynamicstruct.ErrValueCannotBeNil,
				)
			}
		},
	)

	t.Run(
		"get_nonexistent_field", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			var val string
			err = builder.GetFieldValue("NonExistent", &val)
			if !errors.Is(err, dynamicstruct.ErrFieldNotFound) {
				t.Errorf(
					"GetFieldValue() nonexistent field error = %v, want %v",
					err,
					dynamicstruct.ErrFieldNotFound,
				)
			}
		},
	)

	t.Run(
		"get_incompatible_type", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			var val int // Different type than field
			err = builder.GetFieldValue("Name", &val)
			if !errors.Is(err, dynamicstruct.ErrIncompatibleTypes) {
				t.Errorf(
					"GetFieldValue() incompatible type error = %v, want %v",
					err,
					dynamicstruct.ErrIncompatibleTypes,
				)
			}
		},
	)

	t.Run(
		"get_field_successfully", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")
			_ = builder.AddField("Age", int(0))
			_ = builder.AddField("IsActive", false)

			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			var name string
			err = builder.GetFieldValue("Name", &name)
			if err != nil {
				t.Errorf("GetFieldValue() for Name error = %v, wantErr nil", err)
			}

			var age int
			err = builder.GetFieldValue("Age", &age)
			if err != nil {
				t.Errorf("GetFieldValue() for Age error = %v, wantErr nil", err)
			}

			var isActive bool
			err = builder.GetFieldValue("IsActive", &isActive)
			if err != nil {
				t.Errorf("GetFieldValue() for IsActive error = %v, wantErr nil", err)
			}
		},
	)
}

func TestGetField(t *testing.T) {
	t.Run(
		"get_before_build", func(t *testing.T) {
			builder := dynamicstruct.New()
			_, err := builder.GetField("Test")
			if !errors.Is(err, dynamicstruct.ErrInstanceNotBuilt) {
				t.Errorf(
					"GetField() before build error = %v, want %v",
					err,
					dynamicstruct.ErrInstanceNotBuilt,
				)
			}
		},
	)

	t.Run(
		"get_nonexistent_field", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			_, err = builder.GetField("NonExistent")
			if !errors.Is(err, dynamicstruct.ErrFieldNotFound) {
				t.Errorf(
					"GetField() nonexistent field error = %v, want %v",
					err,
					dynamicstruct.ErrFieldNotFound,
				)
			}
		},
	)

	t.Run(
		"get_field_successfully", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")        // AddField only uses type, creates zero value
			_ = builder.AddField("Age", int(0))     // Zero value int
			_ = builder.AddField("IsActive", false) // Zero value bool
			_ = builder.AddField("Score", 0.0)      // Zero value float64

			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Test string field (zero value)
			name, err := builder.GetField("Name")
			if err != nil {
				t.Errorf("GetField() for Name error = %v, wantErr nil", err)
			}
			if nameStr, ok := name.(string); !ok || nameStr != "" {
				t.Errorf("GetField() Name = %v, want \"\"", name)
			}

			// Test int field (zero value)
			age, err := builder.GetField("Age")
			if err != nil {
				t.Errorf("GetField() for Age error = %v, wantErr nil", err)
			}
			if ageInt, ok := age.(int); !ok || ageInt != 0 {
				t.Errorf("GetField() Age = %v, want 0", age)
			}

			// Test bool field (zero value)
			isActive, err := builder.GetField("IsActive")
			if err != nil {
				t.Errorf("GetField() for IsActive error = %v, wantErr nil", err)
			}
			if activeBool, ok := isActive.(bool); !ok || activeBool {
				t.Errorf("GetField() IsActive = %v, want false", isActive)
			}

			// Test float field (zero value)
			score, err := builder.GetField("Score")
			if err != nil {
				t.Errorf("GetField() for Score error = %v, wantErr nil", err)
			}
			if scoreFloat, ok := score.(float64); !ok || scoreFloat != 0.0 {
				t.Errorf("GetField() Score = %v, want 0.0", score)
			}
		},
	)

	t.Run(
		"get_complex_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add complex type fields - AddField only uses types, creates zero values
			_ = builder.AddField("Slice", []string{})     // Zero value slice
			_ = builder.AddField("Map", map[string]int{}) // Zero value map
			_ = builder.AddField("Struct", Person{})      // Zero value struct

			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Test slice field (zero value)
			slice, err := builder.GetField("Slice")
			if err != nil {
				t.Errorf("GetField() for Slice error = %v, wantErr nil", err)
			}
			if sliceVal, ok := slice.([]string); !ok || sliceVal != nil {
				t.Errorf("GetField() Slice = %v, want nil slice", slice)
			}

			// Test map field (zero value)
			mapVal, err := builder.GetField("Map")
			if err != nil {
				t.Errorf("GetField() for Map error = %v, wantErr nil", err)
			}
			if mapResult, ok := mapVal.(map[string]int); !ok || mapResult != nil {
				t.Errorf("GetField() Map = %v, want nil map", mapVal)
			}

			// Test struct field (zero value)
			structVal, err := builder.GetField("Struct")
			if err != nil {
				t.Errorf("GetField() for Struct error = %v, wantErr nil", err)
			}
			if structResult, ok := structVal.(Person); !ok || structResult.Name != "" || structResult.Age != 0 || structResult.Active != false {
				t.Errorf("GetField() Struct = %v, want Person{Name:\"\", Age:0, Active:false}", structVal)
			}
		},
	)

	t.Run(
		"get_zero_values", func(t *testing.T) {
			builder := dynamicstruct.New()
			_ = builder.AddField("Name", "")      // zero value string
			_ = builder.AddField("Age", int(0))   // zero value int
			_ = builder.AddField("Active", false) // zero value bool

			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Test zero value string
			name, err := builder.GetField("Name")
			if err != nil {
				t.Errorf("GetField() for Name error = %v, wantErr nil", err)
			}
			if nameStr, ok := name.(string); !ok || nameStr != "" {
				t.Errorf("GetField() Name = %v, want \"\"", name)
			}

			// Test zero value int
			age, err := builder.GetField("Age")
			if err != nil {
				t.Errorf("GetField() for Age error = %v, wantErr nil", err)
			}
			if ageInt, ok := age.(int); !ok || ageInt != 0 {
				t.Errorf("GetField() Age = %v, want 0", age)
			}

			// Test zero value bool
			active, err := builder.GetField("Active")
			if err != nil {
				t.Errorf("GetField() for Active error = %v, wantErr nil", err)
			}
			if activeBool, ok := active.(bool); !ok || activeBool {
				t.Errorf("GetField() Active = %v, want false", active)
			}
		},
	)
}

func TestReset(t *testing.T) {
	t.Run(
		"reset_after_build", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields and build
			_ = builder.AddField("Name", "")
			_, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Reset
			builder.Reset()

			// Should be able to add fields again after reset
			err = builder.AddField("Age", int(0))
			if err != nil {
				t.Errorf("AddField() after reset error = %v, wantErr nil", err)
			}

			// Should be able to build again
			_, err = builder.Build()
			if err != nil {
				t.Errorf("Build() after reset error = %v, wantErr nil", err)
			}
		},
	)

	t.Run(
		"reset_without_build", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields but don't build
			_ = builder.AddField("Name", "")

			// Reset
			builder.Reset()

			// Should be able to add fields
			err := builder.AddField("Age", int(0))
			if err != nil {
				t.Errorf("AddField() after reset without build error = %v, wantErr nil", err)
			}
		},
	)
}

func TestIntegration(t *testing.T) {
	t.Run(
		"full_workflow", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add different types of fields
			err := builder.AddField("String", "test")
			if err != nil {
				t.Fatalf("AddField() string error = %v", err)
			}

			err = builder.AddField("Int", int(42))
			if err != nil {
				t.Fatalf("AddField() int error = %v", err)
			}

			err = builder.AddField("Bool", true)
			if err != nil {
				t.Fatalf("AddField() bool error = %v", err)
			}

			err = builder.AddField("Float", 3.14)
			if err != nil {
				t.Fatalf("AddField() float error = %v", err)
			}

			err = builder.AddField("Slice", []string{"one", "two"})
			if err != nil {
				t.Fatalf("AddField() slice error = %v", err)
			}

			// Remove a field
			err = builder.RemoveField("Slice")
			if err != nil {
				t.Fatalf("RemoveField() error = %v", err)
			}

			// Build
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Get field values
			var strVal string
			err = builder.GetFieldValue("String", &strVal)
			if err != nil {
				t.Errorf("GetFieldValue() string error = %v", err)
			}

			var intVal int
			err = builder.GetFieldValue("Int", &intVal)
			if err != nil {
				t.Errorf("GetFieldValue() int error = %v", err)
			}

			var boolVal bool
			err = builder.GetFieldValue("Bool", &boolVal)
			if err != nil {
				t.Errorf("GetFieldValue() bool error = %v", err)
			}

			var floatVal float64
			err = builder.GetFieldValue("Float", &floatVal)
			if err != nil {
				t.Errorf("GetFieldValue() float error = %v", err)
			}

			// Verify the removed field is not found
			err = builder.GetFieldValue("Slice", &[]string{})
			if !errors.Is(err, dynamicstruct.ErrFieldNotFound) {
				t.Errorf(
					"GetFieldValue() removed field error = %v, want %v",
					err,
					dynamicstruct.ErrFieldNotFound,
				)
			}

			// Reset and rebuild
			builder.Reset()

			// Should be able to add fields again
			err = builder.AddField("NewField", "new value")
			if err != nil {
				t.Errorf("AddField() after reset error = %v", err)
			}

			// Should be able to build again
			instance, err = builder.Build()
			if err != nil {
				t.Errorf("Build() after reset error = %v", err)
			}
			if instance == nil {
				t.Error("Build() after reset returned nil instance")
			}
		},
	)
}

// Custom struct for testing
type Person struct {
	Name   string
	Age    int
	Active bool
}

// TestComplexTypes tests creating and retrieving more complex data types
func TestComplexTypes(t *testing.T) {
	t.Run(
		"map_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Test with different map types
			mapTypes := []interface{}{
				map[string]string{},
				map[string]int{},
				map[int]string{},
				map[string]interface{}{},
				map[string][]int{},
			}

			for i, mapType := range mapTypes {
				fieldName := "Map" + string(rune('A'+i))
				err := builder.AddField(fieldName, mapType)
				if err != nil {
					t.Fatalf("AddField() for map type %T error = %v", mapType, err)
				}
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Test retrieving map fields
			var strStrMap map[string]string
			err = builder.GetFieldValue("MapA", &strStrMap)
			if err != nil {
				t.Errorf("GetFieldValue() for map[string]string error = %v", err)
			}

			var strIntMap map[string]int
			err = builder.GetFieldValue("MapB", &strIntMap)
			if err != nil {
				t.Errorf("GetFieldValue() for map[string]int error = %v", err)
			}

			var intStrMap map[int]string
			err = builder.GetFieldValue("MapC", &intStrMap)
			if err != nil {
				t.Errorf("GetFieldValue() for map[int]string error = %v", err)
			}

			var strInterfaceMap map[string]interface{}
			err = builder.GetFieldValue("MapD", &strInterfaceMap)
			if err != nil {
				t.Errorf("GetFieldValue() for map[string]interface{} error = %v", err)
			}

			var strIntSliceMap map[string][]int
			err = builder.GetFieldValue("MapE", &strIntSliceMap)
			if err != nil {
				t.Errorf("GetFieldValue() for map[string][]int error = %v", err)
			}
		},
	)

	t.Run(
		"struct_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields with struct types
			err := builder.AddField("Person", Person{})
			if err != nil {
				t.Fatalf("AddField() for Person struct error = %v", err)
			}

			err = builder.AddField("Time", time.Time{})
			if err != nil {
				t.Fatalf("AddField() for time.Time struct error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Test retrieving struct fields
			var person Person
			err = builder.GetFieldValue("Person", &person)
			if err != nil {
				t.Errorf("GetFieldValue() for Person struct error = %v", err)
			}

			var timeVal time.Time
			err = builder.GetFieldValue("Time", &timeVal)
			if err != nil {
				t.Errorf("GetFieldValue() for time.Time struct error = %v", err)
			}
		},
	)

	t.Run(
		"pointer_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields with pointer types
			strPtr := new(string)
			intPtr := new(int)
			boolPtr := new(bool)
			personPtr := new(Person)
			mapPtr := new(map[string]int)
			slicePtr := new([]string)

			err := builder.AddField("StringPtr", strPtr)
			if err != nil {
				t.Fatalf("AddField() for *string error = %v", err)
			}

			err = builder.AddField("IntPtr", intPtr)
			if err != nil {
				t.Fatalf("AddField() for *int error = %v", err)
			}

			err = builder.AddField("BoolPtr", boolPtr)
			if err != nil {
				t.Fatalf("AddField() for *bool error = %v", err)
			}

			err = builder.AddField("PersonPtr", personPtr)
			if err != nil {
				t.Fatalf("AddField() for *Person error = %v", err)
			}

			err = builder.AddField("MapPtr", mapPtr)
			if err != nil {
				t.Fatalf("AddField() for *map[string]int error = %v", err)
			}

			err = builder.AddField("SlicePtr", slicePtr)
			if err != nil {
				t.Fatalf("AddField() for *[]string error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Test retrieving pointer fields
			var strPtrOut *string
			err = builder.GetFieldValue("StringPtr", &strPtrOut)
			if err != nil {
				t.Errorf("GetFieldValue() for *string error = %v", err)
			}

			var intPtrOut *int
			err = builder.GetFieldValue("IntPtr", &intPtrOut)
			if err != nil {
				t.Errorf("GetFieldValue() for *int error = %v", err)
			}

			var boolPtrOut *bool
			err = builder.GetFieldValue("BoolPtr", &boolPtrOut)
			if err != nil {
				t.Errorf("GetFieldValue() for *bool error = %v", err)
			}

			var personPtrOut *Person
			err = builder.GetFieldValue("PersonPtr", &personPtrOut)
			if err != nil {
				t.Errorf("GetFieldValue() for *Person error = %v", err)
			}

			var mapPtrOut *map[string]int
			err = builder.GetFieldValue("MapPtr", &mapPtrOut)
			if err != nil {
				t.Errorf("GetFieldValue() for *map[string]int error = %v", err)
			}

			var slicePtrOut *[]string
			err = builder.GetFieldValue("SlicePtr", &slicePtrOut)
			if err != nil {
				t.Errorf("GetFieldValue() for *[]string error = %v", err)
			}
		},
	)

	t.Run(
		"nested_complex_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields with nested complex types
			nestedMap := map[string]map[int]string{}
			nestedSlice := [][]int{}
			nestedStruct := struct {
				Person Person
				Count  int
				Tags   []string
			}{}
			ptrToMap := &map[string]Person{}
			mapOfPtrs := map[string]*Person{}
			sliceOfPtrs := []*Person{}

			err := builder.AddField("NestedMap", nestedMap)
			if err != nil {
				t.Fatalf("AddField() for nested map error = %v", err)
			}

			err = builder.AddField("NestedSlice", nestedSlice)
			if err != nil {
				t.Fatalf("AddField() for nested slice error = %v", err)
			}

			err = builder.AddField("NestedStruct", nestedStruct)
			if err != nil {
				t.Fatalf("AddField() for nested struct error = %v", err)
			}

			err = builder.AddField("PtrToMap", ptrToMap)
			if err != nil {
				t.Fatalf("AddField() for pointer to map error = %v", err)
			}

			err = builder.AddField("MapOfPtrs", mapOfPtrs)
			if err != nil {
				t.Fatalf("AddField() for map of pointers error = %v", err)
			}

			err = builder.AddField("SliceOfPtrs", sliceOfPtrs)
			if err != nil {
				t.Fatalf("AddField() for slice of pointers error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Test retrieving nested complex type fields
			var nestedMapOut map[string]map[int]string
			err = builder.GetFieldValue("NestedMap", &nestedMapOut)
			if err != nil {
				t.Errorf("GetFieldValue() for nested map error = %v", err)
			}

			var nestedSliceOut [][]int
			err = builder.GetFieldValue("NestedSlice", &nestedSliceOut)
			if err != nil {
				t.Errorf("GetFieldValue() for nested slice error = %v", err)
			}

			var nestedStructOut struct {
				Person Person
				Count  int
				Tags   []string
			}
			err = builder.GetFieldValue("NestedStruct", &nestedStructOut)
			if err != nil {
				t.Errorf("GetFieldValue() for nested struct error = %v", err)
			}

			var ptrToMapOut *map[string]Person
			err = builder.GetFieldValue("PtrToMap", &ptrToMapOut)
			if err != nil {
				t.Errorf("GetFieldValue() for pointer to map error = %v", err)
			}

			var mapOfPtrsOut map[string]*Person
			err = builder.GetFieldValue("MapOfPtrs", &mapOfPtrsOut)
			if err != nil {
				t.Errorf("GetFieldValue() for map of pointers error = %v", err)
			}

			var sliceOfPtrsOut []*Person
			err = builder.GetFieldValue("SliceOfPtrs", &sliceOfPtrsOut)
			if err != nil {
				t.Errorf("GetFieldValue() for slice of pointers error = %v", err)
			}
		},
	)

	t.Run(
		"channel_and_function_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add channel types
			strChan := make(chan string)
			intChan := make(chan int)
			bidirChan := make(chan interface{})

			// Add function types
			func1 := func() {}
			func2 := func(s string) int { return len(s) }

			err := builder.AddField("StringChan", strChan)
			if err != nil {
				t.Fatalf("AddField() for chan string error = %v", err)
			}

			err = builder.AddField("IntChan", intChan)
			if err != nil {
				t.Fatalf("AddField() for chan int error = %v", err)
			}

			err = builder.AddField("BidirChan", bidirChan)
			if err != nil {
				t.Fatalf("AddField() for chan interface{} error = %v", err)
			}

			err = builder.AddField("SimpleFunc", func1)
			if err != nil {
				t.Fatalf("AddField() for simple function error = %v", err)
			}

			err = builder.AddField("ParamFunc", func2)
			if err != nil {
				t.Fatalf("AddField() for function with params error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Test retrieving channel fields
			var strChanOut chan string
			err = builder.GetFieldValue("StringChan", &strChanOut)
			if err != nil {
				t.Errorf("GetFieldValue() for chan string error = %v", err)
			}

			var intChanOut chan int
			err = builder.GetFieldValue("IntChan", &intChanOut)
			if err != nil {
				t.Errorf("GetFieldValue() for chan int error = %v", err)
			}

			var bidirChanOut chan interface{}
			err = builder.GetFieldValue("BidirChan", &bidirChanOut)
			if err != nil {
				t.Errorf("GetFieldValue() for chan interface{} error = %v", err)
			}

			// Test retrieving function fields
			var func1Out func()
			err = builder.GetFieldValue("SimpleFunc", &func1Out)
			if err != nil {
				t.Errorf("GetFieldValue() for simple function error = %v", err)
			}

			var func2Out func(string) int
			err = builder.GetFieldValue("ParamFunc", &func2Out)
			if err != nil {
				t.Errorf("GetFieldValue() for function with params error = %v", err)
			}
		},
	)

	t.Run(
		"incompatible_complex_types", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add fields with complex types
			err := builder.AddField("Map", map[string]int{})
			if err != nil {
				t.Fatalf("AddField() for map error = %v", err)
			}

			err = builder.AddField("Struct", Person{})
			if err != nil {
				t.Fatalf("AddField() for struct error = %v", err)
			}

			err = builder.AddField("Slice", []string{})
			if err != nil {
				t.Fatalf("AddField() for slice error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if instance == nil {
				t.Fatal("Build() returned nil instance")
			}

			// Test with incompatible types
			var wrongMapType map[int]string // Different key type
			err = builder.GetFieldValue("Map", &wrongMapType)
			if !errors.Is(err, dynamicstruct.ErrIncompatibleTypes) {
				t.Errorf(
					"GetFieldValue() with wrong map type error = %v, want %v",
					err,
					dynamicstruct.ErrIncompatibleTypes,
				)
			}

			type OtherStruct struct {
				Field string
			}
			var wrongStructType OtherStruct // Different struct type
			err = builder.GetFieldValue("Struct", &wrongStructType)
			if !errors.Is(err, dynamicstruct.ErrIncompatibleTypes) {
				t.Errorf(
					"GetFieldValue() with wrong struct type error = %v, want %v",
					err,
					dynamicstruct.ErrIncompatibleTypes,
				)
			}

			var wrongSliceType []int // Different element type
			err = builder.GetFieldValue("Slice", &wrongSliceType)
			if !errors.Is(err, dynamicstruct.ErrIncompatibleTypes) {
				t.Errorf(
					"GetFieldValue() with wrong slice type error = %v, want %v",
					err,
					dynamicstruct.ErrIncompatibleTypes,
				)
			}
		},
	)
}

func TestDynamicStructWithJSON(t *testing.T) {
	t.Run(
		"marshal_dynamic_struct", func(t *testing.T) {
			// Create a dynamic struct with typical JSON fields
			builder := dynamicstruct.New()
			err := builder.AddField("ID", int(0))
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Name", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Email", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Active", false)
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Build the struct
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Marshal the empty struct to JSON
			emptyJSON, err := json.Marshal(instance)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}

			// Verify the JSON structure
			expectedEmptyJSON := `{"ID":0,"Name":"","Email":"","Active":false}`
			expected := map[string]interface{}{}

			err = json.Unmarshal([]byte(expectedEmptyJSON), &expected)

			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}

			actual := map[string]interface{}{}

			err = json.Unmarshal(emptyJSON, &actual)

			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("json.Marshal() = %v, want %v", actual, expected)
			}
		},
	)

	t.Run(
		"unmarshal_into_dynamic_struct_with_pointers", func(t *testing.T) {
			// Create a builder for a struct that will have json tags
			builder := dynamicstruct.New()

			// Add fields with JSON tags using StructField directly
			err := builder.AddField("ID", int(0))
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Name", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Email", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Build the struct
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Create a new instance as a pointer for unmarshaling
			instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()

			// JSON data to unmarshal
			jsonData := []byte(`{"ID":123,"Name":"John Doe","Email":"john@example.com"}`)

			// Unmarshal into the pointer
			err = json.Unmarshal(jsonData, instancePtr)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}

			// Extract values using reflection
			instanceValue := reflect.ValueOf(instancePtr).Elem()

			// Check ID field
			idField := instanceValue.FieldByName("ID")
			if !idField.IsValid() || idField.Int() != 123 {
				t.Errorf("Unmarshaled ID = %v, want %v", idField.Int(), 123)
			}

			// Check Name field
			nameField := instanceValue.FieldByName("Name")
			if !nameField.IsValid() || nameField.String() != "John Doe" {
				t.Errorf("Unmarshaled Name = %v, want %v", nameField.String(), "John Doe")
			}

			// Check Email field
			emailField := instanceValue.FieldByName("Email")
			if !emailField.IsValid() || emailField.String() != "john@example.com" {
				t.Errorf("Unmarshaled Email = %v, want %v", emailField.String(), "john@example.com")
			}
		},
	)

	t.Run(
		"marshal_after_field_update", func(t *testing.T) {
			// Create a dynamic struct
			builder := dynamicstruct.New()
			err := builder.AddField("ID", int(0))
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Name", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Build the struct
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Create a pointer to the instance for modification
			instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()

			// Get the element that the pointer refers to
			instanceElem := reflect.ValueOf(instancePtr).Elem()

			// Set field values
			idField := instanceElem.FieldByName("ID")
			if idField.IsValid() && idField.CanSet() {
				idField.SetInt(42)
			}

			nameField := instanceElem.FieldByName("Name")
			if nameField.IsValid() && nameField.CanSet() {
				nameField.SetString("Alice")
			}

			// Marshal to JSON
			jsonData, err := json.Marshal(instancePtr)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}

			// Verify JSON output
			expectedJSON := `{"ID":42,"Name":"Alice"}`

			var expected map[string]interface{}

			err = json.Unmarshal([]byte(expectedJSON), &expected)

			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actual map[string]interface{}

			err = json.Unmarshal(jsonData, &actual)

			if err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("json.Marshal() = %v, want %v", actual, expected)
			}
		},
	)

	t.Run(
		"json_struct_with_nested_struct", func(t *testing.T) {
			// Create a struct for Address
			addressBuilder := dynamicstruct.New()
			err := addressBuilder.AddField("Street", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = addressBuilder.AddField("City", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = addressBuilder.AddField("ZIP", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			addressInstance, err := addressBuilder.Build()
			if err != nil {
				t.Fatalf("Build() address error = %v", err)
			}

			// Create the main struct with the address field
			userBuilder := dynamicstruct.New()
			err = userBuilder.AddField("Name", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Add the address as a field with the same type as addressInstance
			err = userBuilder.AddField("Address", addressInstance)
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			userInstance, err := userBuilder.Build()
			if err != nil {
				t.Fatalf("Build() user error = %v", err)
			}

			// Create a pointer to the instance for modification
			userPtr := reflect.New(reflect.TypeOf(userInstance)).Interface()

			// Get the element that the pointer refers to
			userElem := reflect.ValueOf(userPtr).Elem()

			// Set name
			nameField := userElem.FieldByName("Name")
			if nameField.IsValid() && nameField.CanSet() {
				nameField.SetString("Bob")
			}

			// Get address field
			addressField := userElem.FieldByName("Address")
			if !addressField.IsValid() {
				t.Fatalf("Address field not found or not valid")
			}

			// Set address fields
			streetField := addressField.FieldByName("Street")
			if streetField.IsValid() && streetField.CanSet() {
				streetField.SetString("123 Main St")
			}

			cityField := addressField.FieldByName("City")
			if cityField.IsValid() && cityField.CanSet() {
				cityField.SetString("Anytown")
			}

			zipField := addressField.FieldByName("ZIP")
			if zipField.IsValid() && zipField.CanSet() {
				zipField.SetString("12345")
			}

			// Marshal to JSON
			jsonData, err := json.Marshal(userPtr)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}

			// Verify JSON output
			expectedJSON := `{"Name":"Bob","Address":{"Street":"123 Main St","City":"Anytown","ZIP":"12345"}}`

			var expected map[string]interface{}

			err = json.Unmarshal([]byte(expectedJSON), &expected)

			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actual map[string]interface{}

			err = json.Unmarshal(jsonData, &actual)

			if err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("json.Marshal() = %v, want %v", actual, expected)
			}
		},
	)

	t.Run(
		"json_struct_with_slices", func(t *testing.T) {
			// Create a dynamic struct with a slice field
			builder := dynamicstruct.New()
			err := builder.AddField("Name", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// Add a slice field
			err = builder.AddField("Tags", []string{})
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Create a pointer to the instance for modification
			instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()

			// Get the element that the pointer refers to
			instanceElem := reflect.ValueOf(instancePtr).Elem()

			// Set name
			nameField := instanceElem.FieldByName("Name")
			if nameField.IsValid() && nameField.CanSet() {
				nameField.SetString("Product")
			}

			// Set tags
			tagsField := instanceElem.FieldByName("Tags")
			if tagsField.IsValid() && tagsField.CanSet() {
				tagsField.Set(reflect.ValueOf([]string{"tag1", "tag2", "tag3"}))
			}

			// Marshal to JSON
			jsonData, err := json.Marshal(instancePtr)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}

			// Verify JSON output
			var expected map[string]interface{}
			err = json.Unmarshal(
				[]byte(`{"Name":"Product","Tags":["tag1","tag2","tag3"]}`),
				&expected,
			)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actual map[string]interface{}
			err = json.Unmarshal(jsonData, &actual)
			if err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("json.Marshal() = %v, want %v", actual, expected)
			}
		},
	)
}

// Helper function to demonstrate how to build a dynamic struct with a map
func TestDynamicStructWithMap(t *testing.T) {
	t.Run(
		"dynamic_struct_with_map", func(t *testing.T) {
			builder := dynamicstruct.New()

			// Add a map field
			err := builder.AddField("Properties", map[string]interface{}{})
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Create a pointer to work with the instance
			instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()

			// Get the element
			instanceElem := reflect.ValueOf(instancePtr).Elem()

			// Get the Properties field
			propsField := instanceElem.FieldByName("Properties")
			if !propsField.IsValid() {
				t.Fatal("Properties field not valid")
			}

			// Create a map and set it
			props := map[string]interface{}{
				"color": "red",
				"size":  "large",
				"count": 42,
			}
			propsField.Set(reflect.ValueOf(props))

			// Marshal to JSON
			jsonData, err := json.Marshal(instancePtr)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}

			// The order of map keys in JSON is not guaranteed, so we'll unmarshal back
			// and compare the structures instead of the raw JSON strings
			var unmarshaledMap map[string]interface{}
			err = json.Unmarshal(jsonData, &unmarshaledMap)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}

			// Check if Properties exists and has the right keys/values
			propertiesMap, ok := unmarshaledMap["Properties"].(map[string]interface{})
			if !ok {
				t.Errorf("Properties not found or not a map in unmarshaled JSON")
			}

			// Check if all expected keys and values are present
			expectedKeys := []string{"color", "size", "count"}
			for _, key := range expectedKeys {
				if _, exists := propertiesMap[key]; !exists {
					t.Errorf("Key %s not found in Properties map", key)
				}
			}

			if propertiesMap["color"] != "red" {
				t.Errorf("Properties[\"color\"] = %v, want \"red\"", propertiesMap["color"])
			}

			if propertiesMap["size"] != "large" {
				t.Errorf("Properties[\"size\"] = %v, want \"large\"", propertiesMap["size"])
			}

			// Note: JSON numbers are unmarshaled as float64 by default
			if propertiesMap["count"].(float64) != float64(42) {
				t.Errorf("Properties[\"count\"] = %v, want 42", propertiesMap["count"])
			}
		},
	)
}

// Example of using the dynamic struct with a common JSON workflow
func TestAddFieldWithTags(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		fieldType interface{}
		tags      []string
		wantErr   error
	}{
		{
			name:      "add_field_with_json_tag",
			fieldName: "Name",
			fieldType: "",
			tags:      []string{`json:"name"`},
			wantErr:   nil,
		},
		{
			name:      "add_field_with_multiple_tags",
			fieldName: "Email",
			fieldType: "",
			tags:      []string{`json:"email"`, `validate:"required"`},
			wantErr:   nil,
		},
		{
			name:      "add_field_with_complex_json_tag",
			fieldName: "Price",
			fieldType: float64(0),
			tags:      []string{`json:"price,omitempty"`},
			wantErr:   nil,
		},
		{
			name:      "add_field_without_tags",
			fieldName: "SimpleField",
			fieldType: int(0),
			tags:      []string{},
			wantErr:   nil,
		},
		{
			name:      "add_field_with_invalid_tag",
			fieldName: "InvalidField",
			fieldType: "",
			tags:      []string{`json:"name" invalid_tag_format`},
			wantErr:   dynamicstruct.ErrInvalidTag,
		},
		{
			name:      "add_field_with_empty_tag",
			fieldName: "EmptyTag",
			fieldType: "",
			tags:      []string{``},
			wantErr:   nil,
		},
		{
			name:      "add_field_with_xml_tag",
			fieldName: "XMLField",
			fieldType: "",
			tags:      []string{`xml:"xmlfield,attr"`},
			wantErr:   nil,
		},
		{
			name:      "add_field_with_db_tag",
			fieldName: "DBField",
			fieldType: "",
			tags:      []string{`db:"db_field"`},
			wantErr:   nil,
		},
		{
			name:      "add_field_with_multiple_complex_tags",
			fieldName: "ComplexField",
			fieldType: "",
			tags:      []string{`json:"complex_field,omitempty"`, `xml:"ComplexField"`, `validate:"required,min=1"`},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := dynamicstruct.New()
			err := builder.AddField(tt.fieldName, tt.fieldType, tt.tags...)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("AddField() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("AddField() unexpected error = %v", err)
				return
			}

			// Build and verify the struct has the correct tag
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Use reflection to check if the tag was set correctly
			structType := reflect.TypeOf(instance)
			field, found := structType.FieldByName(tt.fieldName)
			if !found {
				t.Fatalf("Field %s not found in struct", tt.fieldName)
			}

			if len(tt.tags) > 0 {
				expectedTag := strings.Join(tt.tags, " ")
				if string(field.Tag) != expectedTag {
					t.Errorf("Field tag = %q, want %q", field.Tag, expectedTag)
				}
			} else {
				if field.Tag != "" {
					t.Errorf("Expected empty tag, got %q", field.Tag)
				}
			}
		})
	}
}

func TestTagValidation(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		wantErr bool
	}{
		{
			name:    "valid_json_tag",
			tags:    []string{`json:"name"`},
			wantErr: false,
		},
		{
			name:    "valid_multiple_tags",
			tags:    []string{`json:"name"`, `xml:"Name"`},
			wantErr: false,
		},
		{
			name:    "invalid_tag_format",
			tags:    []string{`json:"name" invalid`},
			wantErr: true,
		},
		{
			name:    "valid_complex_tag",
			tags:    []string{`json:"name,omitempty"`, `validate:"required,min=1,max=100"`},
			wantErr: false,
		},
		{
			name:    "empty_tag",
			tags:    []string{``},
			wantErr: false,
		},
		{
			name:    "tag_with_spaces",
			tags:    []string{`json:"field name"`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := dynamicstruct.New()
			err := builder.AddField("TestField", "", tt.tags...)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	t.Run("old_addfield_syntax", func(t *testing.T) {
		builder := dynamicstruct.New()

		// Test that old syntax still works (no tags parameter)
		err := builder.AddField("Name", "")
		if err != nil {
			t.Errorf("AddField() without tags error = %v", err)
		}

		err = builder.AddField("Age", int(0))
		if err != nil {
			t.Errorf("AddField() without tags error = %v", err)
		}

		instance, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		// Verify fields exist and have no tags
		structType := reflect.TypeOf(instance)

		nameField, found := structType.FieldByName("Name")
		if !found {
			t.Fatal("Name field not found")
		}
		if nameField.Tag != "" {
			t.Errorf("Name field should have empty tag, got %q", nameField.Tag)
		}

		ageField, found := structType.FieldByName("Age")
		if !found {
			t.Fatal("Age field not found")
		}
		if ageField.Tag != "" {
			t.Errorf("Age field should have empty tag, got %q", ageField.Tag)
		}
	})
}

func TestJSONWithTags(t *testing.T) {
	t.Run("json_marshaling_with_tags", func(t *testing.T) {
		builder := dynamicstruct.New()

		// Add fields with JSON tags
		err := builder.AddField("UserName", "", `json:"username"`)
		if err != nil {
			t.Fatalf("AddField() error = %v", err)
		}

		err = builder.AddField("EmailAddress", "", `json:"email"`)
		if err != nil {
			t.Fatalf("AddField() error = %v", err)
		}

		err = builder.AddField("IsActive", false, `json:"active"`)
		if err != nil {
			t.Fatalf("AddField() error = %v", err)
		}

		err = builder.AddField("SecretField", "", `json:"-"`)
		if err != nil {
			t.Fatalf("AddField() error = %v", err)
		}

		err = builder.AddField("OptionalField", "", `json:"optional,omitempty"`)
		if err != nil {
			t.Fatalf("AddField() error = %v", err)
		}

		instance, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		// Create pointer and set values
		instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()
		instanceElem := reflect.ValueOf(instancePtr).Elem()

		// Set values
		instanceElem.FieldByName("UserName").SetString("testuser")
		instanceElem.FieldByName("EmailAddress").SetString("test@example.com")
		instanceElem.FieldByName("IsActive").SetBool(true)
		instanceElem.FieldByName("SecretField").SetString("secret")
		// Leave OptionalField empty to test omitempty

		// Marshal to JSON
		jsonData, err := json.Marshal(instancePtr)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		// Verify JSON uses the tag names
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		// Check that JSON keys match the tag names
		if result["username"] != "testuser" {
			t.Errorf("Expected username in JSON, got %v", result)
		}
		if result["email"] != "test@example.com" {
			t.Errorf("Expected email in JSON, got %v", result)
		}
		if result["active"] != true {
			t.Errorf("Expected active in JSON, got %v", result)
		}

		// Secret field should not appear in JSON (json:"-")
		if _, exists := result["SecretField"]; exists {
			t.Error("SecretField should not appear in JSON")
		}

		// Optional field should not appear due to omitempty and empty value
		if _, exists := result["optional"]; exists {
			t.Error("OptionalField should not appear in JSON due to omitempty")
		}
	})
}

func TestJSONWorkflow(t *testing.T) {
	t.Run(
		"complete_json_workflow", func(t *testing.T) {
			// 1. Create a dynamic struct based on expected JSON structure
			builder := dynamicstruct.New()

			err := builder.AddField("ID", int(0))
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Name", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Email", "")
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Active", false)
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			err = builder.AddField("Metadata", map[string]interface{}{})
			if err != nil {
				t.Fatalf("AddField() error = %v", err)
			}

			// 2. Build the struct
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// 3. Create a pointer to work with the instance
			instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()

			// 4. Simulate receiving JSON data (e.g., from an API)
			jsonData := []byte(`{
			"ID": 12345,
			"Name": "Test User",
			"Email": "test@example.com",
			"Active": true,
			"Metadata": {
				"lastLogin": "2025-03-24T10:00:00Z",
				"preferences": {
					"theme": "dark",
					"notifications": true
				}
			}
		}`)

			// 5. Unmarshal JSON into our dynamic struct
			err = json.Unmarshal(jsonData, instancePtr)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}

			// 6. Access and verify data using reflection
			instanceElem := reflect.ValueOf(instancePtr).Elem()

			// Check ID
			idField := instanceElem.FieldByName("ID")
			if !idField.IsValid() || idField.Int() != 12345 {
				t.Errorf("ID = %v, want 12345", idField.Int())
			}

			// Check Name
			nameField := instanceElem.FieldByName("Name")
			if !nameField.IsValid() || nameField.String() != "Test User" {
				t.Errorf("Name = %v, want \"Test User\"", nameField.String())
			}

			// Check Active
			activeField := instanceElem.FieldByName("Active")
			if !activeField.IsValid() || !activeField.Bool() {
				t.Errorf("Active = %v, want true", activeField.Bool())
			}

			// 7. Modify the struct
			nameField.SetString("Updated User")

			// 8. Marshal back to JSON
			updatedJSON, err := json.Marshal(instancePtr)
			if err != nil {
				t.Errorf("json.Marshal() updated struct error = %v", err)
			}

			// 9. Verify the updated JSON
			var updatedData map[string]interface{}
			err = json.Unmarshal(updatedJSON, &updatedData)
			if err != nil {
				t.Errorf("json.Unmarshal() updated JSON error = %v", err)
			}

			if updatedData["Name"] != "Updated User" {
				t.Errorf("Updated JSON Name = %v, want \"Updated User\"", updatedData["Name"])
			}
		},
	)
}

// Test types for anonymous field testing
type PersonTest struct {
	Name string
	Age  int
}

type AddressTest struct {
	Street string
	City   string
}

type ContactTest struct {
	Email string
	Phone string
}

func TestAddAnonymousField(t *testing.T) {
	tests := []struct {
		name      string
		fieldType interface{}
		tags      []string
		wantErr   error
	}{
		{
			name:      "add_person_anonymous_field",
			fieldType: PersonTest{},
			tags:      []string{},
			wantErr:   nil,
		},
		{
			name:      "add_address_anonymous_field_with_tags",
			fieldType: AddressTest{},
			tags:      []string{`json:",inline"`},
			wantErr:   nil,
		},
		{
			name:      "add_string_anonymous_field",
			fieldType: "",
			tags:      []string{},
			wantErr:   nil,
		},
		{
			name:      "add_int_anonymous_field",
			fieldType: int(0),
			tags:      []string{},
			wantErr:   nil,
		},
		{
			name:      "add_duplicate_anonymous_field",
			fieldType: PersonTest{},
			tags:      []string{},
			wantErr:   dynamicstruct.ErrAnonymousFieldAlreadyExists,
		},
		{
			name:      "add_anonymous_field_with_invalid_tag",
			fieldType: ContactTest{},
			tags:      []string{`json:"name" invalid_format`},
			wantErr:   dynamicstruct.ErrInvalidTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := dynamicstruct.New()

			if tt.name == "add_duplicate_anonymous_field" {
				// First add a PersonTest field
				err := builder.AddAnonymousField(PersonTest{})
				if err != nil {
					t.Fatalf("First AddAnonymousField() failed: %v", err)
				}
			}

			err := builder.AddAnonymousField(tt.fieldType, tt.tags...)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("AddAnonymousField() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("AddAnonymousField() unexpected error = %v", err)
				return
			}

			// Build and verify the struct has the correct anonymous field
			instance, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			// Use reflection to check if the anonymous field was added correctly
			structType := reflect.TypeOf(instance)
			found := false
			expectedType := reflect.TypeOf(tt.fieldType)

			for i := 0; i < structType.NumField(); i++ {
				field := structType.Field(i)
				if field.Anonymous && field.Type == expectedType {
					found = true
					if len(tt.tags) > 0 {
						expectedTag := strings.Join(tt.tags, " ")
						if string(field.Tag) != expectedTag {
							t.Errorf("Anonymous field tag = %q, want %q", field.Tag, expectedTag)
						}
					}
					break
				}
			}

			if !found {
				t.Errorf("Anonymous field of type %s not found in struct", expectedType.Name())
			}
		})
	}
}

func TestAnonymousFieldAfterBuild(t *testing.T) {
	t.Run("add_anonymous_field_after_build", func(t *testing.T) {
		builder := dynamicstruct.New()
		builder.AddField("Name", "")
		_, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		err = builder.AddAnonymousField(PersonTest{})
		if !errors.Is(err, dynamicstruct.ErrInstanceAlreadyBuilt) {
			t.Errorf("AddAnonymousField() after build error = %v, want %v", err, dynamicstruct.ErrInstanceAlreadyBuilt)
		}
	})
}

func TestGetAnonymousField(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(builder *dynamicstruct.Builder)
		fieldType interface{}
		wantErr   error
	}{
		{
			name: "get_person_anonymous_field",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: PersonTest{},
			wantErr:   nil,
		},
		{
			name: "get_address_anonymous_field",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(AddressTest{})
				builder.Build()
			},
			fieldType: AddressTest{},
			wantErr:   nil,
		},
		{
			name: "get_nonexistent_anonymous_field",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: ContactTest{},
			wantErr:   dynamicstruct.ErrAnonymousFieldNotFound,
		},
		{
			name: "get_anonymous_field_before_build",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				// Don't build
			},
			fieldType: PersonTest{},
			wantErr:   dynamicstruct.ErrInstanceNotBuilt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := dynamicstruct.New()
			tt.setupFunc(builder)

			value, err := builder.GetAnonymousField(tt.fieldType)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetAnonymousField() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetAnonymousField() unexpected error = %v", err)
				return
			}

			if value == nil {
				t.Error("GetAnonymousField() returned nil value")
				return
			}

			// Verify the type matches
			expectedType := reflect.TypeOf(tt.fieldType)
			actualType := reflect.TypeOf(value)
			if actualType != expectedType {
				t.Errorf("GetAnonymousField() returned type %v, want %v", actualType, expectedType)
			}
		})
	}
}

func TestGetAnonymousFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(builder *dynamicstruct.Builder)
		fieldType interface{}
		valuePtr  interface{}
		wantErr   error
	}{
		{
			name: "get_person_anonymous_field_value",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: PersonTest{},
			valuePtr:  &PersonTest{},
			wantErr:   nil,
		},
		{
			name: "get_string_anonymous_field_value",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField("")
				builder.Build()
			},
			fieldType: "",
			valuePtr:  new(string),
			wantErr:   nil,
		},
		{
			name: "get_anonymous_field_value_wrong_type",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: PersonTest{},
			valuePtr:  &AddressTest{},
			wantErr:   dynamicstruct.ErrIncompatibleTypes,
		},
		{
			name: "get_anonymous_field_value_non_pointer",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: PersonTest{},
			valuePtr:  PersonTest{},
			wantErr:   dynamicstruct.ErrValueMustBePointer,
		},
		{
			name: "get_anonymous_field_value_nil_pointer",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: PersonTest{},
			valuePtr:  (*PersonTest)(nil),
			wantErr:   dynamicstruct.ErrValueCannotBeNil,
		},
		{
			name: "get_nonexistent_anonymous_field_value",
			setupFunc: func(builder *dynamicstruct.Builder) {
				builder.AddAnonymousField(PersonTest{})
				builder.Build()
			},
			fieldType: ContactTest{},
			valuePtr:  &ContactTest{},
			wantErr:   dynamicstruct.ErrAnonymousFieldNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := dynamicstruct.New()
			tt.setupFunc(builder)

			err := builder.GetAnonymousFieldValue(tt.fieldType, tt.valuePtr)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetAnonymousFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetAnonymousFieldValue() unexpected error = %v", err)
			}
		})
	}
}

func TestAnonymousFieldsOrder(t *testing.T) {
	t.Run("anonymous_fields_come_first", func(t *testing.T) {
		builder := dynamicstruct.New()

		// Add regular field first
		builder.AddField("RegularField", "")

		// Add anonymous fields
		builder.AddAnonymousField(PersonTest{})
		builder.AddAnonymousField(AddressTest{})

		instance, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		structType := reflect.TypeOf(instance)

		// Check that first two fields are anonymous
		if structType.NumField() < 3 {
			t.Fatalf("Expected at least 3 fields, got %d", structType.NumField())
		}

		// First field should be PersonTest (anonymous)
		firstField := structType.Field(0)
		if !firstField.Anonymous || firstField.Type != reflect.TypeOf(PersonTest{}) {
			t.Errorf("First field should be anonymous PersonTest, got %+v", firstField)
		}

		// Second field should be AddressTest (anonymous)
		secondField := structType.Field(1)
		if !secondField.Anonymous || secondField.Type != reflect.TypeOf(AddressTest{}) {
			t.Errorf("Second field should be anonymous AddressTest, got %+v", secondField)
		}

		// Third field should be RegularField (not anonymous)
		thirdField := structType.Field(2)
		if thirdField.Anonymous || thirdField.Name != "RegularField" {
			t.Errorf("Third field should be regular RegularField, got %+v", thirdField)
		}
	})
}

func TestAnonymousFieldsReset(t *testing.T) {
	t.Run("reset_clears_anonymous_fields", func(t *testing.T) {
		builder := dynamicstruct.New()

		// Add anonymous fields
		builder.AddAnonymousField(PersonTest{})
		builder.AddAnonymousField(AddressTest{})

		// Build
		_, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		// Reset
		builder.Reset()

		// Try to add the same anonymous field again - should work
		err = builder.AddAnonymousField(PersonTest{})
		if err != nil {
			t.Errorf("AddAnonymousField() after reset error = %v", err)
		}

		// Build again
		instance, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() after reset error = %v", err)
		}

		// Should only have one anonymous field now
		structType := reflect.TypeOf(instance)
		anonymousCount := 0
		for i := 0; i < structType.NumField(); i++ {
			if structType.Field(i).Anonymous {
				anonymousCount++
			}
		}

		if anonymousCount != 1 {
			t.Errorf("Expected 1 anonymous field after reset, got %d", anonymousCount)
		}
	})
}

func TestAnonymousFieldsIntegration(t *testing.T) {
	t.Run("mixed_fields_workflow", func(t *testing.T) {
		builder := dynamicstruct.New()

		// Add mix of regular and anonymous fields
		builder.AddField("ID", int(0))
		builder.AddAnonymousField(PersonTest{}, `json:",inline"`)
		builder.AddField("Active", false)
		builder.AddAnonymousField(AddressTest{})

		_, err := builder.Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		// Test regular field access
		id, err := builder.GetField("ID")
		if err != nil {
			t.Errorf("GetField(ID) error = %v", err)
		}
		if _, ok := id.(int); !ok {
			t.Errorf("ID field should be int, got %T", id)
		}

		// Test anonymous field access
		person, err := builder.GetAnonymousField(PersonTest{})
		if err != nil {
			t.Errorf("GetAnonymousField(PersonTest) error = %v", err)
		}
		if _, ok := person.(PersonTest); !ok {
			t.Errorf("PersonTest field should be PersonTest, got %T", person)
		}

		address, err := builder.GetAnonymousField(AddressTest{})
		if err != nil {
			t.Errorf("GetAnonymousField(AddressTest) error = %v", err)
		}
		if _, ok := address.(AddressTest); !ok {
			t.Errorf("AddressTest field should be AddressTest, got %T", address)
		}

		// Test type safety with GetAnonymousFieldValue
		var personValue PersonTest
		err = builder.GetAnonymousFieldValue(PersonTest{}, &personValue)
		if err != nil {
			t.Errorf("GetAnonymousFieldValue(PersonTest) error = %v", err)
		}

		var addressValue AddressTest
		err = builder.GetAnonymousFieldValue(AddressTest{}, &addressValue)
		if err != nil {
			t.Errorf("GetAnonymousFieldValue(AddressTest) error = %v", err)
		}
	})
}
