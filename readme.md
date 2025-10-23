# DynamicStruct

DynamicStruct is a Go package that enables runtime creation and manipulation of struct types. It provides a clean API to dynamically create, modify, and access struct fields, which is particularly useful when dealing with dynamic data structures or when the structure of your data is not known at compile time.

## Features

- Create struct types dynamically at runtime
- Add and remove fields with type safety
- Support for struct tags (JSON, XML, validation, etc.)
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

// Add fields with struct tags
_ = builder.AddField("Email", "", `json:"email"`)
_ = builder.AddField("UserID", int(0), `json:"user_id"`, `validate:"required"`)

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

### Getting Field Values Directly

For convenience, you can also get field values directly without providing a pointer:

```go
// Get field values directly as interface{}
name, err := builder.GetField("Name")
if err != nil {
    // Handle error
    // Possible errors: 
    // - ErrInstanceNotBuilt
    // - ErrFieldNotFound
}

// Type assertion is needed when using GetField
nameStr, ok := name.(string)
if !ok {
    // Handle type assertion failure
}

// Example with different types
age, err := builder.GetField("Age")
if err == nil {
    if ageInt, ok := age.(int); ok {
        fmt.Printf("Age: %d\n", ageInt)
    }
}

active, err := builder.GetField("Active")
if err == nil {
    if activeBool, ok := active.(bool); ok {
        fmt.Printf("Active: %t\n", activeBool)
    }
}
```

Note: `GetField` returns the field value as `interface{}`, so you need to perform type assertion to get the actual typed value. Use `GetFieldValue` if you prefer compile-time type safety.

### Working with Struct Tags

You can add struct tags to fields for use with JSON, XML, validation libraries, and more:

```go
builder := dynamicstruct.New()

// Add fields with JSON tags
_ = builder.AddField("UserName", "", `json:"username"`)
_ = builder.AddField("EmailAddress", "", `json:"email"`)
_ = builder.AddField("IsActive", false, `json:"active"`)

// Add fields with multiple tags
_ = builder.AddField("Price", float64(0), `json:"price,omitempty"`, `validate:"min=0"`)

// Add fields with complex tags
_ = builder.AddField("XMLField", "", `xml:"xmlfield,attr"`, `json:"-"`)

// Build the struct
instance, _ := builder.Build()

// JSON marshaling will respect the tags
jsonData, _ := json.Marshal(instance)
// Output uses tag names: {"username":"","email":"","active":false,"price":0}
```

**Tag Validation:**
The library validates struct tag format using `github.com/fatih/structtag`. Invalid tag formats will return `ErrInvalidTag`.

```go
// This will return an error due to invalid tag format
err := builder.AddField("Invalid", "", `json:"name" invalid_format`)
if errors.Is(err, dynamicstruct.ErrInvalidTag) {
    // Handle invalid tag error
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
// Create a struct dynamically with JSON tags
builder := dynamicstruct.New()
_ = builder.AddField("ID", int(0), `json:"id"`)
_ = builder.AddField("Name", "", `json:"name"`)
_ = builder.AddField("Email", "", `json:"email"`)

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
// Output: {"id":42,"name":"Alice","email":""}

// Unmarshal JSON back into a dynamic struct
newData := []byte(`{"id":123,"name":"Bob","email":"bob@example.com"}`)
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
- `ErrInvalidTag`: When providing an invalid struct tag format

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

- Field visibility is limited (all fields are exported)
- Struct tag validation requires the `github.com/fatih/structtag` dependency

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
