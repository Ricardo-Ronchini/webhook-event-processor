package common

import (
	"math"
	"os"
	"reflect"
	"strings"
)

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetEnvArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")

	items := strings.Split(value, ",")
	for i, item := range items {
		items[i] = strings.Trim(strings.TrimSpace(item), `"`)
	}

	return items
}

func RoundTo(n float32, places int) float32 {
	pow := math.Pow(10, float64(places))
	return float32(math.Round(float64(n)*pow) / pow)
}

// CombineObjectsWithNulls fills zero-value fields in source with the corresponding values from target.
// Both arguments must be pointers to structs of the same type.
func CombineObjectsWithNulls(source, target any) {
	srcVal := reflect.ValueOf(source)
	tgtVal := reflect.ValueOf(target)

	if srcVal.Kind() == reflect.Pointer {
		srcVal = srcVal.Elem()
	}
	if tgtVal.Kind() == reflect.Pointer {
		tgtVal = tgtVal.Elem()
	}

	if srcVal.Kind() != reflect.Struct || tgtVal.Kind() != reflect.Struct {
		return
	}

	if srcVal.Type() != tgtVal.Type() {
		return
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		tgtField := tgtVal.Field(i)

		if !srcField.CanSet() {
			continue
		}

		if srcField.IsZero() {
			srcField.Set(tgtField)
		}
	}
}
