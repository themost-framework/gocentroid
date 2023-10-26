package query

import (
	"encoding/json"
	"errors"

	"golang.org/x/exp/maps"
)

type QueryEntity struct {
	Name  string
	Alias *string
}

func (entity *QueryEntity) MarshalJSON() ([]byte, error) {
	if entity.Alias == nil {
		return json.Marshal(map[string]interface{}{
			entity.Name: 1,
		})
	}
	return json.Marshal(map[string]interface{}{
		*entity.Alias: FormatNameReference(entity.Name),
	})
}

func (entity *QueryEntity) UnmarshalJSON(str []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(str, &data); err != nil {
		return err
	}
	keys := maps.Keys(data)
	if len(keys) > 0 {
		// get first key
		key := keys[0]
		if data[key] == 1 {
			entity.Name = key
		} else {
			entity.Alias = &key
			name, ok := (data[key]).(string)
			if !ok {
				return errors.New("expected query collection name to be a valid string")
			}
			entity.Name = TrimNameReference(name)
		}
	} else {
		return errors.New("expected query collection name or alias")
	}
	return nil
}
