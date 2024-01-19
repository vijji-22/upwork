package database

import (
	"database/sql/driver"
	"encoding/json"
)

type DbJson[V any] struct {
	v *V
}

func NewDbJson[V any]() *DbJson[V] {
	return &DbJson[V]{}
}

func (d *DbJson[V]) Scan(value any) error {
	if s, ok := value.([]uint8); ok {
		return d.UnmarshalJSON([]byte(string(s)))
	}

	return nil
}

func (d *DbJson[V]) Value() (driver.Value, error) {
	v, err := json.Marshal(d.v)
	return string(v), err
}

func (d *DbJson[V]) Json() *V {
	return d.v
}

// MarshalJSON implements json.Marshaler.
func (d *DbJson[V]) MarshalJSON() ([]byte, error) {
	if d.v == nil {
		return nil, nil
	}
	return json.Marshal(d.v)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DbJson[V]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	v := new(V)
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	d.v = v
	return nil
}

type DbMap[K comparable, V any] struct {
	m map[K]V
}

func NewDbMap[K comparable, V any]() *DbMap[K, V] {
	return &DbMap[K, V]{
		m: make(map[K]V),
	}
}

func (d *DbMap[K, V]) Scan(value any) error {

	if s, ok := value.([]uint8); ok {
		return json.Unmarshal([]byte(string(s)), &d.m)
	}

	return nil
}

func (d *DbMap[K, V]) Value() (driver.Value, error) {
	v, err := json.Marshal(d.m)
	return string(v), err
}

func (d *DbMap[K, V]) Map() map[K]V {
	return d.m
}

// MarshalJSON implements json.Marshaler.
func (d *DbMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.m)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DbMap[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.m)
}

type DbSlice[V any] struct {
	s []V
}

func NewDbSlice[V any]() *DbSlice[V] {
	return &DbSlice[V]{
		s: make([]V, 0),
	}
}

func (d *DbSlice[V]) Scan(value any) error {
	if s, ok := value.([]uint8); ok {
		return json.Unmarshal([]byte(string(s)), &d.s)
	}

	return nil
}

func (d *DbSlice[V]) Value() (driver.Value, error) {
	v, err := json.Marshal(d.s)
	return string(v), err
}

func (d *DbSlice[V]) Slice() []V {
	return d.s
}

// MarshalJSON implements json.Marshaler.
func (d *DbSlice[V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.s)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DbSlice[V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.s)
}
