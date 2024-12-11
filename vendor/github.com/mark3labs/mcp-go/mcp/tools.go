package mcp

// ToolOption is a function that configures a Tool
type ToolOption func(*Tool)

// PropertyOption is a function that configures a property
type PropertyOption func(map[string]interface{})

//
// Core Tool Functions
//

// NewTool creates a new Tool with the given name and options
func NewTool(name string, opts ...ToolOption) Tool {
	tool := Tool{
		Name: name,
		InputSchema: ToolInputSchema{
			Type:       "object",
			Properties: make(map[string]interface{}),
			Required:   nil, // Will be omitted from JSON if empty
		},
	}

	for _, opt := range opts {
		opt(&tool)
	}

	return tool
}

// WithDescription adds a description to the Tool
func WithDescription(description string) ToolOption {
	return func(t *Tool) {
		t.Description = description
	}
}

//
// Common Property Options
//

// Description adds a description to a property
func Description(desc string) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["description"] = desc
	}
}

// Required marks a property as required
func Required() PropertyOption {
	return func(schema map[string]interface{}) {
		schema["required"] = true
	}
}

// Title adds a display-friendly title to a property
func Title(title string) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["title"] = title
	}
}

//
// String Property Options
//

// DefaultString sets the default value for a string property
func DefaultString(value string) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["default"] = value
	}
}

// Enum specifies a list of allowed values for a string property
func Enum(values ...string) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["enum"] = values
	}
}

// MaxLength sets the maximum length for a string property
func MaxLength(max int) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["maxLength"] = max
	}
}

// MinLength sets the minimum length for a string property
func MinLength(min int) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["minLength"] = min
	}
}

// Pattern sets a regex pattern that a string property must match
func Pattern(pattern string) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["pattern"] = pattern
	}
}

//
// Number Property Options
//

// DefaultNumber sets the default value for a number property
func DefaultNumber(value float64) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["default"] = value
	}
}

// Max sets the maximum value for a number property
func Max(max float64) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["maximum"] = max
	}
}

// Min sets the minimum value for a number property
func Min(min float64) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["minimum"] = min
	}
}

// MultipleOf specifies that a number must be a multiple of the given value
func MultipleOf(value float64) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["multipleOf"] = value
	}
}

//
// Boolean Property Options
//

// DefaultBool sets the default value for a boolean property
func DefaultBool(value bool) PropertyOption {
	return func(schema map[string]interface{}) {
		schema["default"] = value
	}
}

//
// Property Type Helpers
//

// WithBoolean adds a boolean property to the tool schema
func WithBoolean(name string, opts ...PropertyOption) ToolOption {
	return func(t *Tool) {
		schema := map[string]interface{}{
			"type": "boolean",
		}

		for _, opt := range opts {
			opt(schema)
		}

		// Remove required from property schema and add to InputSchema.required
		if required, ok := schema["required"].(bool); ok && required {
			delete(schema, "required")
			if t.InputSchema.Required == nil {
				t.InputSchema.Required = []string{name}
			} else {
				t.InputSchema.Required = append(t.InputSchema.Required, name)
			}
		}

		t.InputSchema.Properties[name] = schema
	}
}

// WithNumber adds a number property to the tool schema
func WithNumber(name string, opts ...PropertyOption) ToolOption {
	return func(t *Tool) {
		schema := map[string]interface{}{
			"type": "number",
		}

		for _, opt := range opts {
			opt(schema)
		}

		// Remove required from property schema and add to InputSchema.required
		if required, ok := schema["required"].(bool); ok && required {
			delete(schema, "required")
			if t.InputSchema.Required == nil {
				t.InputSchema.Required = []string{name}
			} else {
				t.InputSchema.Required = append(t.InputSchema.Required, name)
			}
		}

		t.InputSchema.Properties[name] = schema
	}
}

// WithString adds a string property to the tool schema
func WithString(name string, opts ...PropertyOption) ToolOption {
	return func(t *Tool) {
		schema := map[string]interface{}{
			"type": "string",
		}

		for _, opt := range opts {
			opt(schema)
		}

		// Remove required from property schema and add to InputSchema.required
		if required, ok := schema["required"].(bool); ok && required {
			delete(schema, "required")
			if t.InputSchema.Required == nil {
				t.InputSchema.Required = []string{name}
			} else {
				t.InputSchema.Required = append(t.InputSchema.Required, name)
			}
		}

		t.InputSchema.Properties[name] = schema
	}
}
