# DynamicStruct

DynamicStruct is a Go package that enables runtime creation and manipulation of struct types. It provides a clean API to dynamically create, modify, and access struct fields, which is particularly useful when dealing with dynamic data structures or when the structure of your data is not known at compile time.

## Features

- Create struct types dynamically at runtime
- Add and remove fields with type safety
- Thread-safe operations with mutex protection
- Access field values with type checking
- Works seamlessly with Go's standard library, including JSON encoding/decoding

## Installation

```bash
go get github.com/gosmos-space/dynamicstruct
```

## Usage

### Creating a Dynamic Struct

```go
// Create a new builder
builder := dynamicstruct.New()

// Add fields with their types
_ = builder.AddField("Name", "")      // String field
_ = builder.AddField("Age", int(0))   // Integer field
_ = builder.AddField("Active", false) // Boolean field

// Build the struct instance
instance, err := builder.Build()
if err != nil {
    panic(err)
}
```

### Removing Fields

```go
// Remove a field
err := builder.RemoveField("Age")
if err != nil {
    // Handle error
    // Possible errors: ErrInstanceAlreadyBuilt
}
```

### Accessing Field Values

```go
// Get field values
var name string
err := builder.GetFieldValue("Name", &name)
if err != nil {
    // Handle error
    // Possible errors: 
    // - ErrInstanceNotBuilt
    // - ErrValueMustBePointer
    // - ErrValueCannotBeNil
    // - ErrFieldNotFound
    // - ErrIncompatibleTypes
}
```

### Resetting the Builder

```go
// Reset the builder to reuse it
builder.Reset()

// Now you can add new fields and build again
_ = builder.AddField("ID", int(0))
```

### Working with JSON

DynamicStruct works well with Go's standard JSON encoding/decoding:

```go
// Create a struct dynamically
builder := dynamicstruct.New()
_ = builder.AddField("ID", int(0))
_ = builder.AddField("Name", "")
_ = builder.AddField("Email", "")

// Build it
instance, _ := builder.Build()

// Convert to a pointer for JSON operation
instancePtr := reflect.New(reflect.TypeOf(instance)).Interface()
instanceElem := reflect.ValueOf(instancePtr).Elem()

// Set some values
idField := instanceElem.FieldByName("ID")
if idField.IsValid() && idField.CanSet() {
    idField.SetInt(42)
}

nameField := instanceElem.FieldByName("Name")
if nameField.IsValid() && nameField.CanSet() {
    nameField.SetString("Alice")
}

// Marshal to JSON
jsonData, _ := json.Marshal(instancePtr)
fmt.Println(string(jsonData))
// Output: {"ID":42,"Name":"Alice","Email":""}

// Unmarshal JSON back into a dynamic struct
newData := []byte(`{"ID":123,"Name":"Bob","Email":"bob@example.com"}`)
json.Unmarshal(newData, instancePtr)
```

## Error Handling

The package provides specific error types:

- `ErrFieldAlreadyExists`: When trying to add a field with a name that already exists
- `ErrInstanceAlreadyBuilt`: When trying to modify an already built instance
- `ErrInstanceNotBuilt`: When trying to access an instance that hasn't been built
- `ErrValueMustBePointer`: When passing a non-pointer to GetFieldValue
- `ErrValueCannotBeNil`: When passing a nil pointer to GetFieldValue
- `ErrFieldNotFound`: When the requested field doesn't exist
- `ErrIncompatibleTypes`: When the field type doesn't match the pointer type

Use `errors.Is()` to check for these specific errors:

```go
err := builder.AddField("Name", "")
if errors.Is(err, dynamicstruct.ErrFieldAlreadyExists) {
    // Handle duplicate field
}
```

## Thread Safety

All operations in DynamicStruct are protected by a mutex, making it safe to use from multiple goroutines.

## Limitations

- Currently, there's no direct support for JSON tags or other struct tags
- Field visibility is limited (all fields are exported)

## Cautions and Best Practices

⚠️ **Use with caution**: Dynamic struct creation should be used only when necessary, not as a default approach.

### Potential Risks and Drawbacks

1. **Type Safety**: Using dynamic structs bypasses Go's compile-time type checking, which can lead to runtime errors instead of compile-time errors.

2. **Performance Impact**:
    - Dynamic struct creation has overhead compared to using static structs
    - Each reflection operation is significantly slower than direct field access
    - Heavy use in performance-critical paths can cause noticeable slowdowns

3. **Code Readability**: Code using dynamic structs is typically harder to understand, maintain, and debug than code using static types.

4. **IDE Support**: You lose IDE features like autocomplete and refactoring support when working with dynamic fields.

5. **Error Prone**: It's easy to make mistakes when accessing fields by string names, which won't be caught until runtime.

### Appropriate Use Cases

Use DynamicStruct when:

- Working with truly dynamic data where the structure isn't known until runtime
- Implementing generic data processing pipelines
- Adapting to external systems with changing schemas
- Building tools for data conversion or mapping
- Creating mock objects for testing

**Do not use** when:

- The struct structure is known at compile time
- Performance is critical
- Type safety is a primary concern
- Code clarity and maintenance are priorities
- Simple maps or existing structs would suffice

### Alternative Approaches

Consider these alternatives before reaching for dynamic structs:

- Use `map[string]interface{}` for simple key-value storage
- Define interfaces for polymorphic behavior
- Use code generation for known schema variations
- Define a union type that covers all possible fields

## License

[MIT](LICENSE)
