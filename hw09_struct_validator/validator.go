package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrOnlyStructSupport     = errors.New("only structure type is supported")
	ErrValueNotSupported     = errors.New("value not supported")
	ErrInvalidValidationRule = errors.New("invalid validation rule")
	ErrValidateMin           = errors.New("less than min")
	ErrValidateMax           = errors.New("bigger than max")
	ErrValidateLen           = errors.New("not equal to len")
	ErrValidateIn            = errors.New("not in range")
	ErrValidateRegexp        = errors.New("does not match pattern")
)

type validationRules map[string]string

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	msg := make([]string, 0, len(v))
	for _, e := range v {
		msg = append(msg, fmt.Sprintf("%s: %v", e.Field, e.Err))
	}

	return strings.Join(msg, "\n")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	refValue := reflect.ValueOf(v)

	if refValue.Kind() != reflect.Struct {
		return ErrOnlyStructSupport
	}

	refType := refValue.Type()

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		if !field.IsExported() {
			continue
		}

		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		rules, err := parseValidationRules(validateTag)
		if err != nil {
			return fmt.Errorf("invalid validation rules for field %s: %w", field.Name, err)
		}

		value := refValue.Field(i)
		err = validateField(value, field.Name, rules)
		if err != nil {
			var verr ValidationErrors
			if errors.As(err, &verr) {
				validationErrors = append(validationErrors, verr...)
			} else {
				return fmt.Errorf("validation error for field %s: %w", field.Name, err)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(value reflect.Value, fieldName string, rules validationRules) error {
	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String:
		return validateString(value.String(), fieldName, rules)
	case reflect.Int:
		return validateInt(value.Int(), fieldName, rules)
	case reflect.Slice:
		return validateSlice(value, fieldName, rules)
	}

	return nil
}

func validateString(value string, fieldName string, rules validationRules) error {
	var errs ValidationErrors

	for key, rule := range rules {
		var err error

		switch key {
		case "in":
			if !slices.Contains(strings.Split(rule, ","), value) {
				err = ErrValidateIn
			}
		case "len":
			ln, e := strconv.Atoi(rule)
			if e != nil {
				err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, rule)
			} else if len([]rune(value)) != ln {
				err = ErrValidateLen
			}
		case "regexp":
			if _, e := regexp.Compile(rule); e != nil {
				err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, rule)
			} else if !regexp.MustCompile(rule).MatchString(value) {
				err = ErrValidateRegexp
			}
		default:
			err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, key)
		}

		if err != nil {
			errs = append(errs, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateInt(value int64, fieldName string, rules validationRules) error {
	var errs ValidationErrors

	for key, rule := range rules {
		var err error

		switch key {
		case "in":
			collection := strings.Split(rule, ",")
			found := false
			for _, item := range collection {
				num, e := strconv.ParseInt(item, 10, 64)
				if e != nil {
					err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, rule)
					break
				}
				if num == value {
					found = true
					break
				}
			}
			if err == nil && !found {
				err = ErrValidateIn
			}
		case "min":
			minVal, e := strconv.ParseInt(rule, 10, 64)
			if e != nil {
				err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, rule)
			} else if value < minVal {
				err = ErrValidateMin
			}
		case "max":
			maxVal, e := strconv.ParseInt(rule, 10, 64)
			if e != nil {
				err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, rule)
			} else if value > maxVal {
				err = ErrValidateMax
			}
		default:
			err = fmt.Errorf("%w: %s", ErrInvalidValidationRule, key)
		}

		if err != nil {
			errs = append(errs, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateSlice(value reflect.Value, fieldName string, rules validationRules) error {
	if value.Len() == 0 {
		return nil
	}

	var sliceErr error

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		//nolint:exhaustive
		switch item.Kind() {
		case reflect.String:
			if err := validateString(item.String(), fieldName, rules); err != nil {
				var verr ValidationErrors
				if errors.As(err, &verr) {
					if sliceErr == nil {
						sliceErr = verr[0].Err
					}
				}
			}
		case reflect.Int:
			if err := validateInt(item.Int(), fieldName, rules); err != nil {
				var verr ValidationErrors
				if errors.As(err, &verr) {
					if sliceErr == nil {
						sliceErr = verr[0].Err
					}
				}
			}
		default:
			return ErrValueNotSupported
		}
	}

	if sliceErr != nil {
		return ValidationErrors{
			ValidationError{
				Field: fieldName,
				Err:   sliceErr,
			},
		}
	}
	return nil
}

func parseValidationRules(str string) (validationRules, error) {
	rules := make(validationRules)
	if str == "" {
		return rules, nil
	}

	for _, item := range strings.Split(str, "|") {
		parts := strings.SplitN(item, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("%w: %s", ErrInvalidValidationRule, item)
		}
		rules[parts[0]] = parts[1]
	}

	return rules, nil
}
