# Validation

By default, data checking is supported for structs by `ValidateStruct(structData interface{}, msgForTag MsgForTagFunc) (map[string][]string, error)`

### Usage


Quick usage
```go
import "app/core/validation"

if errorData, err := c.ValidateStruct(loginDto, validation.MsgForTag); err != nil {
    return c.BadRequest(errorData)
}
```

Customize error message by yourself
```go
if errorData, err := c.ValidateStruct(loginDto, func(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "This field is required"
    case "gte":
        return fmt.Sprintf("This field is gte %s characters", fe.Param())
    case "email":
        return "Invalid email"
    }

    return fe.Error()
}); err != nil {
    return c.BadRequest(errorData)
}
```

### Message for Tag

Review and add more code in file `message_for_tag.go`.

### Add custom Validator

Add more functions in file `rules.go`. Don't forget declare in `validator.go` file

```go
// Custom validation for uuid.UUID fields. Use `validate:"uuid"`
_ = instance.RegisterValidation("uuid", uuidValidator)
```