package nullable

import "reflect"

func IsNullable(value any) any {
	if reflect.ValueOf(value).IsZero() {
		return nil
	}
	return value
}
