// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"io"
	"unsafe"
)

var (
	// API the json codec in use.
	API Core

	// Marshal returns the JSON encoding of v.
	Marshal func(v any) ([]byte, error)

	// Unmarshal parses the JSON-encoded data and stores the result
	// in the value pointed to by v. If v is nil or not a pointer,
	// Unmarshal returns an [InvalidUnmarshalError].
	Unmarshal func(data []byte, v any) error

	// MarshalIndent is like [Marshal] but applies [Indent] to format the output.
	// Each JSON element in the output will begin on a new line beginning with prefix
	// followed by one or more copies of indent according to the indentation nesting.
	MarshalIndent func(v any, prefix, indent string) ([]byte, error)

	// NewEncoder returns a new encoder that writes to writer.
	NewEncoder func(writer io.Writer) Encoder

	// NewDecoder returns a new decoder that reads from reader.
	NewDecoder func(reader io.Reader) Decoder
)

// SetAPI setup glocal API object and help functions.
// Notice: doesnt check the param of api whether effect like nil.
func SetAPI(api Core) {
	API, MarshalIndent = api, api.MarshalIndent
	Marshal, Unmarshal = api.Marshal, api.Unmarshal
	NewEncoder, NewDecoder = api.NewEncoder, api.NewDecoder
}

// Core the api for json codec.
type Core interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
	NewEncoder(writer io.Writer) Encoder
	NewDecoder(reader io.Reader) Decoder
}

// Encoder an interface writes JSON values to an output stream.
type Encoder interface {
	// SetEscapeHTML specifies whether problematic HTML characters
	// should be escaped inside JSON quoted strings.
	// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
	// to avoid certain safety problems that can arise when embedding JSON in HTML.
	//
	// In non-HTML settings where the escaping interferes with the readability
	// of the output, SetEscapeHTML(false) disables this behavior.
	SetEscapeHTML(on bool)

	// Encode writes the JSON encoding of v to the stream,
	// followed by a newline character.
	//
	// See the documentation for Marshal for details about the
	// conversion of Go values to JSON.
	Encode(v any) error
}

// Decoder an interface reads and decodes JSON values from an input stream.
type Decoder interface {
	// UseNumber causes the Decoder to unmarshal a number into an any as a
	// Number instead of as a float64.
	UseNumber()

	// DisallowUnknownFields causes the Decoder to return an error when the destination
	// is a struct and the input contains object keys which do not match any
	// non-ignored, exported fields in the destination.
	DisallowUnknownFields()

	// Decode reads the next JSON-encoded value from its
	// input and stores it in the value pointed to by v.
	//
	// See the documentation for Unmarshal for details about
	// the conversion of JSON into a Go value.
	Decode(v any) error
}

// MarshalToString is like [Marshal] but returns the JSON string encoding of v..
func MarshalToString(v any) (string, error) {
	out, err := API.Marshal(v)
	if err != nil {
		return "", err
	}
	return unsafe.String(unsafe.SliceData(out), len(out)), nil
}

// UnmarshalFromString is like [Unmarshal] but from data string.
func UnmarshalFromString(data string, v any) error {
	return API.Unmarshal(unsafe.Slice(unsafe.StringData(data), len(data)), v)
}
