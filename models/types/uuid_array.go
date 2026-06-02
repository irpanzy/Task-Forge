package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("unsupported type for UUIDArray")
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	parts := strings.Split(str, ",")

	*a = make(UUIDArray, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(strings.Trim(part, `"`))
		if part == "" {
			continue
		}
		id, err := uuid.Parse(part)
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %v", err)
		}
		*a = append(*a, id)
	}

	return nil
}

func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	postgreFormat := make([]string, 0, len(a))
	for _, id := range a {
		postgreFormat = append(postgreFormat, fmt.Sprintf(`"%s"`, id.String()))
	}
	return "{" + strings.Join(postgreFormat, ",") + "}", nil
}

func (_ UUIDArray) GormDataType() string {
	return "uuid[]"
}
