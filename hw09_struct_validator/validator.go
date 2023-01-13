package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := make([]string, len(v))
	for idx, validationError := range v {
		result[idx] = fmt.Sprintf("invalid %s: %s", validationError.Field, validationError.Err.Error())
	}

	return strings.Join(result, ", ")
}

const (
	typeString      = "string"
	typeStringSlice = "[]string"
	typeInt         = "int"
	typeIntSlice    = "[]int"

	lenStr = "len"
	inStr  = "in"
	reStr  = "regexp"
	maxNum = "max"
	minNum = "min"
	inNum  = "in"
)

var typeValidators = map[string][]string{
	typeString: {lenStr, inStr, reStr},
	typeInt:    {maxNum, minNum, inNum},
}

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind().String() != "struct" {
		return errors.New("variable is not a struct")
	}

	objectFields := reflect.VisibleFields(reflect.TypeOf(v))
	objectValue := reflect.Indirect(reflect.ValueOf(v))

	validationErrors := ValidationErrors{}

	// проходимся по полям объекта
	for _, field := range objectFields {
		// пробуем найти тэг validate
		validateValue := field.Tag.Get("validate")
		if validateValue == "" {
			continue
		}

		// разделяем валидаторы
		validations := strings.Split(validateValue, "|")
		if len(validations) == 0 {
			return fmt.Errorf("validator parsing error (%s)", validateValue)
		}

		// проходимся по валидаторам
		for _, validation := range validations {
			// у каждого валидатора должно быть название и параметры
			validationParams := strings.Split(validation, ":")
			if len(validationParams) != 2 {
				return fmt.Errorf("validator parsing error (%s)", validation)
			}

			// валидируем значение поля
			validator, condition := validationParams[0], validationParams[1]
			fieldValue := objectValue.FieldByName(field.Name)

			var err error
			validated := true

			switch field.Type.String() {
			case typeString:
				validated, err = validateStringSlice([]string{fieldValue.String()}, validator, condition)
			case typeInt:
				validated, err = validateIntSlice([]int{int(fieldValue.Int())}, validator, condition)
			case typeStringSlice:
				if values, converted := fieldValue.Interface().([]string); converted {
					validated, err = validateStringSlice(values, validator, condition)
				}
			case typeIntSlice:
				if values, converted := fieldValue.Interface().([]int); converted {
					validated, err = validateIntSlice(values, validator, condition)
				}
			}

			if err != nil {
				return err
			}

			if !validated {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   fmt.Errorf("validation failed on '%s'", validation),
				})
			}
		}
	}

	return validationErrors
}

func checkValidator(validator string, valueType string) bool {
	// у каждого типа должны быть свои валидаторы
	return slices.Contains(typeValidators[valueType], validator)
}

func validateStringSlice(values []string, validator, condition string) (bool, error) {
	if !checkValidator(validator, typeString) {
		return false, errors.New("invalid validator")
	}

	var re *regexp.Regexp

	for _, value := range values {
		switch validator {
		case lenStr:
			conditionValue, err := strconv.Atoi(condition)
			if err != nil {
				return false, errors.New("invalid validator value")
			}

			if len(value) > conditionValue {
				return false, nil
			}
		case inStr:
			conditionValue := strings.Split(condition, ",")
			if len(conditionValue) == 0 {
				return false, errors.New("invalid validator value")
			}

			if !slices.Contains(conditionValue, value) {
				return false, nil
			}
		case reStr:
			if re == nil {
				newRe, err := regexp.Compile(condition)
				if err != nil {
					return false, err
				}

				re = newRe
			}

			if !re.MatchString(value) {
				return false, nil
			}
		}
	}

	return true, nil
}

func validateIntSlice(values []int, validator, condition string) (bool, error) {
	if !checkValidator(validator, typeInt) {
		return false, errors.New("invalid validator")
	}

	for _, value := range values {
		switch validator {
		case maxNum:
			conditionValue, err := strconv.Atoi(condition)
			if err != nil {
				return false, errors.New("invalid validator value")
			}

			if value > conditionValue {
				return false, nil
			}
		case minNum:
			conditionValue, err := strconv.Atoi(condition)
			if err != nil {
				return false, errors.New("invalid validator value")
			}

			if value < conditionValue {
				return false, nil
			}
		case inNum:
			conditionValue := strings.Split(condition, ",")
			if len(conditionValue) == 0 {
				return false, errors.New("invalid validator value")
			}

			if !slices.Contains(conditionValue, strconv.Itoa(value)) {
				return false, nil
			}
		}
	}

	return true, nil
}
