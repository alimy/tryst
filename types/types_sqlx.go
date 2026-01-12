// Copyright 2023 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package types

import (
	"database/sql/driver"
	"errors"
)

// Serializable data marshal/unmarshal constraint for Binary type.
type Serializable[T any] interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) (T, error)
}

// Binary[T] is a []byte which transparently Binary[T] data being submitted to
// a database and unmarshal data being Scanned from a database.
type Binary[T Serializable[T]] struct {
	Data T
}

// NullBinary[T] represents a Binary that may be null.
// NullBinary[T] implements the scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBinary[T Serializable[T]] struct {
	Data  T
	Valid bool // Valid is true if Binary is not NULL
}

// Value implements the driver.Valuer interface, marshal the raw value of
// this Binary[T].
func (b *Binary[T]) Value() (driver.Value, error) {
	return b.Data.MarshalBinary()
}

// Scan implements the sql.Scanner interface, unmashal the value coming off
// the wire and storing the raw result in the Binary[T].
func (b *Binary[T]) Scan(src any) (err error) {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
	default:
		return errors.New("incompatible type for Binary")
	}
	b.Data, err = b.Data.UnmarshalBinary(source)
	return
}

// Value implements the driver.Valuer interface, marshal the raw value of
// this Binary[T].
func (b *NullBinary[T]) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Data.MarshalBinary()
}

// Scan implements the sql.Scanner interface, unmashal the value coming off
// the wire and storing the raw result in the Binary[T].
func (b *NullBinary[T]) Scan(src any) (err error) {
	if b.Valid = (src != nil); !b.Valid {
		return nil
	}
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
	default:
		return errors.New("incompatible type for Binary")
	}
	b.Data, err = b.Data.UnmarshalBinary(source)
	return
}
