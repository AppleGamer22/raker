package types

import (
	"database/sql/driver"
	"fmt"
)

const format = "(%f,%f)"

type Point struct {
	X float64
	Y float64
}

func (p *Point) String() string {
	return fmt.Sprintf(format, p.X, p.Y)
}

type NullPoint struct {
	Point Point
	Valid bool
}

// Scan implements the sql.Scanner interface.
func (np *NullPoint) Scan(value interface{}) error {
	if value == nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}

	var err error
	np.Valid = true

	// database/sql with pq driver returns point as string format "(x,y)"
	switch v := value.(type) {
	case string:
		_, err = fmt.Sscanf(v, format, &np.Point.X, &np.Point.Y)
	case []byte:
		_, err = fmt.Sscanf(string(v), format, &np.Point.X, &np.Point.Y)
	default:
		np.Valid = false
		return fmt.Errorf("cannot scan type %T into NullPoint", value)
	}

	if err != nil {
		np.Valid = false
	}
	return err
}

// Value implements the driver.Valuer interface.
func (np NullPoint) Value() (driver.Value, error) {
	if !np.Valid {
		return nil, nil
	}
	return np.Point.String(), nil
}
