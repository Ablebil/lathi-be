package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB json.RawMessage

func (j *JSONB) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*j = JSONB(bytes)
	return nil
}

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
