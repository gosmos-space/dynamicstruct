package dynamicstruct

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/structtag"
)

type Builder struct {
	fields          map[string]reflect.StructField
	anonymousFields []reflect.StructField
	instance        *reflect.Value
	m               sync.Mutex
}

func New() *Builder {
	return &Builder{
		fields: make(map[string]reflect.StructField),
	}
}

func (b *Builder) AddField(name string, kind any, tags ...string) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.instance != nil {
		return ErrInstanceAlreadyBuilt
	}

	if _, ok := b.fields[name]; ok {
		return ErrFieldAlreadyExists
	}

	// Build tag string from variadic tags
	var tag reflect.StructTag

	if len(tags) > 0 {
		tagString := strings.Join(tags, " ")

		// Validate tag format using structtag library, but only if not empty
		if tagString != "" {
			if _, err := structtag.Parse(tagString); err != nil {
				return ErrInvalidTag
			}
		}

		tag = reflect.StructTag(tagString)
	}

	b.fields[name] = reflect.StructField{
		Name: name,
		Type: reflect.TypeOf(kind),
		Tag:  tag,
	}

	return nil
}

func (b *Builder) AddAnonymousField(fieldType any, tags ...string) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.instance != nil {
		return ErrInstanceAlreadyBuilt
	}

	fieldTypeReflect := reflect.TypeOf(fieldType)

	// Check if anonymous field of this type already exists
	for _, field := range b.anonymousFields {
		if field.Type == fieldTypeReflect {
			return ErrAnonymousFieldAlreadyExists
		}
	}

	// Build tag string from variadic tags
	var tag reflect.StructTag

	if len(tags) > 0 {
		tagString := strings.Join(tags, " ")

		// Validate tag format using structtag library, but only if not empty
		if tagString != "" {
			if _, err := structtag.Parse(tagString); err != nil {
				return ErrInvalidTag
			}
		}

		tag = reflect.StructTag(tagString)
	}

	// Generate a unique name for the anonymous field
	fieldName := fieldTypeReflect.Name()
	if fieldName == "" {
		// For built-in types like string, int, etc., use the type kind
		fieldName = fieldTypeReflect.Kind().String()
	}
	// Ensure the name is exported (starts with uppercase)
	if len(fieldName) > 0 && fieldName[0] >= 'a' && fieldName[0] <= 'z' {
		fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]
	}

	b.anonymousFields = append(b.anonymousFields, reflect.StructField{
		Name:      fieldName,
		Type:      fieldTypeReflect,
		Tag:       tag,
		Anonymous: true,
	})

	return nil
}

func (b *Builder) RemoveField(name string) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.instance != nil {
		return ErrInstanceAlreadyBuilt
	}

	delete(b.fields, name)

	return nil
}

func (b *Builder) buildStructFields() []reflect.StructField {
	totalFields := len(b.anonymousFields) + len(b.fields)
	fields := make([]reflect.StructField, 0, totalFields)

	// Add anonymous fields first (as specified)
	fields = append(fields, b.anonymousFields...)

	// Add regular fields
	for _, field := range b.fields {
		fields = append(fields, field)
	}

	return fields
}

func (b *Builder) Build() (any, error) {
	b.m.Lock()
	defer b.m.Unlock()

	if b.instance != nil {
		return nil, ErrInstanceAlreadyBuilt
	}

	instance := reflect.New(
		reflect.StructOf(b.buildStructFields()),
	).Elem()

	b.instance = &instance

	return b.instance.Interface(), nil
}

func (b *Builder) Reset() {
	b.m.Lock()
	defer b.m.Unlock()

	b.instance = nil
	b.anonymousFields = nil
}

func (b *Builder) GetFieldValue(name string, value any) error {
	b.m.Lock()
	defer b.m.Unlock()

	// Check if instance is built
	if b.instance == nil {
		return ErrInstanceNotBuilt
	}

	valueReflect := reflect.ValueOf(value)

	// Check if value is a pointer and not nil
	if valueReflect.Kind() != reflect.Ptr {
		return ErrValueMustBePointer
	}

	// Check if value is not nil
	if valueReflect.IsNil() {
		return ErrValueCannotBeNil
	}

	// Get the field by name
	field := b.instance.FieldByName(name)

	if !field.IsValid() {
		return ErrFieldNotFound
	}

	// Check if the types are compatible
	if field.Type() != valueReflect.Elem().Type() {
		return fmt.Errorf(
			"%w: field type: %s, value type: %s",
			ErrIncompatibleTypes,
			field.Type().String(),
			valueReflect.Elem().Type().String(),
		)
	}

	// Set the value
	valueReflect.Elem().Set(field)

	return nil
}

func (b *Builder) GetAnonymousField(fieldType any) (any, error) {
	b.m.Lock()
	defer b.m.Unlock()

	// Check if instance is built
	if b.instance == nil {
		return nil, ErrInstanceNotBuilt
	}

	fieldTypeReflect := reflect.TypeOf(fieldType)

	// Find the anonymous field by type
	structType := b.instance.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if field.Anonymous && field.Type == fieldTypeReflect {
			return b.instance.Field(i).Interface(), nil
		}
	}

	return nil, ErrAnonymousFieldNotFound
}

func (b *Builder) GetAnonymousFieldValue(fieldType any, value any) error {
	b.m.Lock()
	defer b.m.Unlock()

	// Check if instance is built
	if b.instance == nil {
		return ErrInstanceNotBuilt
	}

	valueReflect := reflect.ValueOf(value)

	// Check if value is a pointer and not nil
	if valueReflect.Kind() != reflect.Ptr {
		return ErrValueMustBePointer
	}

	// Check if value is not nil
	if valueReflect.IsNil() {
		return ErrValueCannotBeNil
	}

	fieldTypeReflect := reflect.TypeOf(fieldType)

	// Find the anonymous field by type
	structType := b.instance.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if field.Anonymous && field.Type == fieldTypeReflect {
			fieldValue := b.instance.Field(i)

			// Check if the types are compatible
			if fieldValue.Type() != valueReflect.Elem().Type() {
				return fmt.Errorf(
					"%w: field type: %s, value type: %s",
					ErrIncompatibleTypes,
					fieldValue.Type().String(),
					valueReflect.Elem().Type().String(),
				)
			}

			// Set the value
			valueReflect.Elem().Set(fieldValue)

			return nil
		}
	}

	return ErrAnonymousFieldNotFound
}

func (b *Builder) GetField(name string) (any, error) {
	b.m.Lock()
	defer b.m.Unlock()

	// Check if instance is built
	if b.instance == nil {
		return nil, ErrInstanceNotBuilt
	}

	// Get the field by name
	field := b.instance.FieldByName(name)

	if !field.IsValid() {
		return nil, ErrFieldNotFound
	}

	// Return the field value as interface{}
	return field.Interface(), nil
}
