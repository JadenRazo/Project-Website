package validator

import (
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
    "strings"

    "github.com/xeipuuv/gojsonschema"
)

// Validator handles request validation
type Validator struct {
    schemas map[string]*gojsonschema.Schema
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
    return &Validator{
        schemas: make(map[string]*gojsonschema.Schema),
    }
}

// RegisterSchema registers a JSON schema for validation
func (v *Validator) RegisterSchema(name string, schema []byte) error {
    loader := gojsonschema.NewBytesLoader(schema)
    schema, err := gojsonschema.NewSchema(loader)
    if err != nil {
        return fmt.Errorf("invalid schema %s: %v", name, err)
    }
    v.schemas[name] = schema
    return nil
}

// ValidateRequest validates an HTTP request against a registered schema
func (v *Validator) ValidateRequest(r *http.Request, schemaName string) error {
    schema, exists := v.schemas[schemaName]
    if !exists {
        return fmt.Errorf("schema %s not found", schemaName)
    }

    // Read request body
    var body interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return fmt.Errorf("invalid JSON: %v", err)
    }

    // Validate against schema
    result, err := schema.Validate(gojsonschema.NewGoLoader(body))
    if err != nil {
        return fmt.Errorf("validation error: %v", err)
    }

    if !result.Valid() {
        var errors []string
        for _, err := range result.Errors() {
            errors = append(errors, err.String())
        }
        return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
    }

    return nil
}

// ValidateStruct validates a struct against validation tags
func (v *Validator) ValidateStruct(s interface{}) error {
    val := reflect.ValueOf(s)
    if val.Kind() != reflect.Ptr || val.IsNil() {
        return fmt.Errorf("invalid struct: must be a non-nil pointer")
    }

    val = val.Elem()
    typ := val.Type()
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        tag := fieldType.Tag.Get("validate")

        if tag == "" {
            continue
        }

        rules := strings.Split(tag, ",")
        for _, rule := range rules {
            if err := v.validateField(field, rule); err != nil {
                return fmt.Errorf("%s: %v", fieldType.Name, err)
            }
        }
    }

    return nil
}

// validateField validates a single field against a validation rule
func (v *Validator) validateField(field reflect.Value, rule string) error {
    parts := strings.Split(rule, "=")
    if len(parts) != 2 {
        return fmt.Errorf("invalid validation rule: %s", rule)
    }

    switch parts[0] {
    case "required":
        if field.IsZero() {
            return fmt.Errorf("field is required")
        }
    case "min":
        if field.Kind() == reflect.String {
            if field.Len() < parseInt(parts[1]) {
                return fmt.Errorf("minimum length is %s", parts[1])
            }
        } else if field.Kind() == reflect.Int || field.Kind() == reflect.Int64 {
            if field.Int() < parseInt(parts[1]) {
                return fmt.Errorf("minimum value is %s", parts[1])
            }
        }
    case "max":
        if field.Kind() == reflect.String {
            if field.Len() > parseInt(parts[1]) {
                return fmt.Errorf("maximum length is %s", parts[1])
            }
        } else if field.Kind() == reflect.Int || field.Kind() == reflect.Int64 {
            if field.Int() > parseInt(parts[1]) {
                return fmt.Errorf("maximum value is %s", parts[1])
            }
        }
    case "email":
        if field.Kind() == reflect.String {
            if !isValidEmail(field.String()) {
                return fmt.Errorf("invalid email format")
            }
        }
    case "url":
        if field.Kind() == reflect.String {
            if !isValidURL(field.String()) {
                return fmt.Errorf("invalid URL format")
            }
        }
    }

    return nil
}

// ValidationMiddleware creates a middleware that validates requests
func ValidationMiddleware(validator *Validator, schemaName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if err := validator.ValidateRequest(r, schemaName); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Helper functions
func parseInt(s string) int64 {
    var i int64
    fmt.Sscanf(s, "%d", &i)
    return i
}

func isValidEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidURL(url string) bool {
    return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
