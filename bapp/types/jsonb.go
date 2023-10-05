package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONB is a type that provides helper methods for storing/fetching JSON in sql
type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	bytes, err := json.Marshal(a)
	return string(bytes), err
}

func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		// Sqlite stores these as a string, which is a kind of byte array right?
		s, ok := value.(string)
		if !ok {
			return errors.New("type assertion to []byte or string failed")
		}
		b = []byte(s)
	}
	return json.Unmarshal(b, &a)
}

func (a *JSONB) FromStruct(obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("JSONB FromStruct: marshaling obj: (%w)", err)
	}

	var newMap map[string]interface{}
	if err = json.Unmarshal(data, &newMap); err != nil {
		return fmt.Errorf("JSONB FromStruct: Unmarshaling to new map: (%w)", err)
	}

	*a = JSONB(newMap)

	return nil
}
