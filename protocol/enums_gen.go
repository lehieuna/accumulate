package protocol

// GENERATED BY go run ./tools/cmd/gen-enum. DO NOT EDIT.

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AccumulateNetwork/accumulate/internal/encoding"
)

// ChainTypeUnknown is used when the chain type is not known.
const ChainTypeUnknown ChainType = 0

// ChainTypeTransaction holds transaction hashes.
const ChainTypeTransaction ChainType = 1

// ChainTypeAnchor holds chain anchors.
const ChainTypeAnchor ChainType = 2

// ChainTypeData holds data entry hashes.
const ChainTypeData ChainType = 3

// KeyPageOperationUnknown is used when the key page operation is not known.
const KeyPageOperationUnknown KeyPageOperation = 0

// KeyPageOperationUpdate replaces a key in the page with a new key.
const KeyPageOperationUpdate KeyPageOperation = 1

// KeyPageOperationRemove removes a key from the page.
const KeyPageOperationRemove KeyPageOperation = 2

// KeyPageOperationAdd adds a key to the page.
const KeyPageOperationAdd KeyPageOperation = 3

// KeyPageOperationSetThreshold sets the signing threshold (the M of "M of N" signatures required).
const KeyPageOperationSetThreshold KeyPageOperation = 4

// ObjectTypeUnknown is used when the object type is not known.
const ObjectTypeUnknown ObjectType = 0

// ObjectTypeAccount represents an account object.
const ObjectTypeAccount ObjectType = 1

// ObjectTypeTransaction represents a transaction object.
const ObjectTypeTransaction ObjectType = 2

// ID returns the ID of the Chain Type
func (v ChainType) ID() uint64 { return uint64(v) }

// String returns the name of the Chain Type
func (v ChainType) String() string {
	switch v {
	case ChainTypeUnknown:
		return "unknown"
	case ChainTypeTransaction:
		return "transaction"
	case ChainTypeAnchor:
		return "anchor"
	case ChainTypeData:
		return "data"
	default:
		return fmt.Sprintf("ChainType:%d", v)
	}
}

// ChainTypeByName returns the named Chain Type.
func ChainTypeByName(name string) (ChainType, bool) {
	switch name {
	case "unknown":
		return ChainTypeUnknown, true
	case "transaction":
		return ChainTypeTransaction, true
	case "anchor":
		return ChainTypeAnchor, true
	case "data":
		return ChainTypeData, true
	default:
		return 0, false
	}
}

// MarshalJSON marshals the Chain Type to JSON as a string.
func (v ChainType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON unmarshals the Chain Type from JSON as a string.
func (v *ChainType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	var ok bool
	*v, ok = ChainTypeByName(s)
	if !ok || strings.ContainsRune(v.String(), ':') {
		return fmt.Errorf("invalid Chain Type %q", s)
	}
	return nil
}

// BinarySize returns the number of bytes required to binary marshal the Chain Type.
func (v ChainType) BinarySize() int {
	return encoding.UvarintBinarySize(v.ID())
}

// MarshalBinary marshals the Chain Type to bytes as a unsigned varint.
func (v ChainType) MarshalBinary() ([]byte, error) {
	return encoding.UvarintMarshalBinary(v.ID()), nil
}

// UnmarshalBinary unmarshals the Chain Type from bytes as a unsigned varint.
func (v *ChainType) UnmarshalBinary(data []byte) error {
	u, err := encoding.UvarintUnmarshalBinary(data)
	if err != nil {
		return err
	}

	*v = ChainType(u)
	return nil
}

// ID returns the ID of the Key PageOpe ration
func (v KeyPageOperation) ID() uint64 { return uint64(v) }

// String returns the name of the Key PageOpe ration
func (v KeyPageOperation) String() string {
	switch v {
	case KeyPageOperationUnknown:
		return "unknown"
	case KeyPageOperationUpdate:
		return "update"
	case KeyPageOperationRemove:
		return "remove"
	case KeyPageOperationAdd:
		return "add"
	case KeyPageOperationSetThreshold:
		return "setThreshold"
	default:
		return fmt.Sprintf("KeyPageOperation:%d", v)
	}
}

// KeyPageOperationByName returns the named Key PageOpe ration.
func KeyPageOperationByName(name string) (KeyPageOperation, bool) {
	switch name {
	case "unknown":
		return KeyPageOperationUnknown, true
	case "update":
		return KeyPageOperationUpdate, true
	case "remove":
		return KeyPageOperationRemove, true
	case "add":
		return KeyPageOperationAdd, true
	case "setThreshold":
		return KeyPageOperationSetThreshold, true
	default:
		return 0, false
	}
}

// MarshalJSON marshals the Key PageOpe ration to JSON as a string.
func (v KeyPageOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON unmarshals the Key PageOpe ration from JSON as a string.
func (v *KeyPageOperation) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	var ok bool
	*v, ok = KeyPageOperationByName(s)
	if !ok || strings.ContainsRune(v.String(), ':') {
		return fmt.Errorf("invalid Key PageOpe ration %q", s)
	}
	return nil
}

// BinarySize returns the number of bytes required to binary marshal the Key PageOpe ration.
func (v KeyPageOperation) BinarySize() int {
	return encoding.UvarintBinarySize(v.ID())
}

// MarshalBinary marshals the Key PageOpe ration to bytes as a unsigned varint.
func (v KeyPageOperation) MarshalBinary() ([]byte, error) {
	return encoding.UvarintMarshalBinary(v.ID()), nil
}

// UnmarshalBinary unmarshals the Key PageOpe ration from bytes as a unsigned varint.
func (v *KeyPageOperation) UnmarshalBinary(data []byte) error {
	u, err := encoding.UvarintUnmarshalBinary(data)
	if err != nil {
		return err
	}

	*v = KeyPageOperation(u)
	return nil
}

// ID returns the ID of the Object Type
func (v ObjectType) ID() uint64 { return uint64(v) }

// String returns the name of the Object Type
func (v ObjectType) String() string {
	switch v {
	case ObjectTypeUnknown:
		return "unknown"
	case ObjectTypeAccount:
		return "account"
	case ObjectTypeTransaction:
		return "transaction"
	default:
		return fmt.Sprintf("ObjectType:%d", v)
	}
}

// ObjectTypeByName returns the named Object Type.
func ObjectTypeByName(name string) (ObjectType, bool) {
	switch name {
	case "unknown":
		return ObjectTypeUnknown, true
	case "account":
		return ObjectTypeAccount, true
	case "transaction":
		return ObjectTypeTransaction, true
	default:
		return 0, false
	}
}

// MarshalJSON marshals the Object Type to JSON as a string.
func (v ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON unmarshals the Object Type from JSON as a string.
func (v *ObjectType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	var ok bool
	*v, ok = ObjectTypeByName(s)
	if !ok || strings.ContainsRune(v.String(), ':') {
		return fmt.Errorf("invalid Object Type %q", s)
	}
	return nil
}

// BinarySize returns the number of bytes required to binary marshal the Object Type.
func (v ObjectType) BinarySize() int {
	return encoding.UvarintBinarySize(v.ID())
}

// MarshalBinary marshals the Object Type to bytes as a unsigned varint.
func (v ObjectType) MarshalBinary() ([]byte, error) {
	return encoding.UvarintMarshalBinary(v.ID()), nil
}

// UnmarshalBinary unmarshals the Object Type from bytes as a unsigned varint.
func (v *ObjectType) UnmarshalBinary(data []byte) error {
	u, err := encoding.UvarintUnmarshalBinary(data)
	if err != nil {
		return err
	}

	*v = ObjectType(u)
	return nil
}
