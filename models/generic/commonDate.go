package generic

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type CommonDate struct {
	time.Time
}

func (t *CommonDate) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (t *CommonDate) MarshalJSON() ([]byte, error) {
	formatted := t.Time.Format("2006-01-02")
	return []byte(`"` + formatted + `"`), nil
}

func (t *CommonDate) String() string {
	return t.Format("2006-01-02")
}

func (t *CommonDate) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.Time, nil
}

func (t *CommonDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(time.Time); ok {
		t.Time = v
		return nil
	}
	return fmt.Errorf("failed to scan CommonDate: %v", value)
}

type CommonDuration struct {
	time.Duration
}

func (d CommonDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *CommonDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
