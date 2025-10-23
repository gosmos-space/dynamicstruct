package dynamicstruct

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/structtag"
)

type Builder struct {
	fields   map[string]reflect.StructField
	instance *reflect.Value
	m        sync.Mutex
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
	fields := make([]reflect.StructField, 0, len(b.fields))

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
