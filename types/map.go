package types

import (
	"fmt"
	"strings"
)

type Map struct {
	Node
	KV []KV
}

// Returns a new empty Map<T>
func NewEmptyMap() *Map {
	return &Map{}
}

// Adds an operation id to the set pres(k). The set pres(k), is the set of all
//operations that have asserted the existence of the map.
func (m Map) AddOpPresence(opId string) {}

// Removes an operation id from the set pres(k). The set pres(k), is the set of
//all operations that have asserted the existence of the map.
func (m Map) RmOpPresence(opId string) {}

// Get a node value from input key. If the key does not exist in the map,
// creates a new KV in which key is the input passed and value is empty
func (m *Map) Get(key string) Node {
	for _, kv := range m.KV {
		if kv.Key == key {
			return kv.Value
		}
	}
	// If key does not exist, create one, add it to the map and returns new empty
	// value
	nKv := KV{
		Key: key,
	}
	return nKv.Value
}

// Gets all keys from Map<T>
func (m *Map) Keys() []string {
	keys := []string{}
	for _, kv := range m.KV {
		keys = append(keys, kv.Key)
	}
	return keys
}

// Gets all values from Map<T>
func (m *Map) Values() []Node {
	vals := []Node{}
	for _, kv := range m.KV {
		vals = append(vals, kv.Value)
	}
	return vals
}

// TODO: Deletes the whole Map<T>
func (m Map) Delete() {}

// TODO: Deletes a KV
func (m *Map) DeleteKey() {}

func (m Map) String() string {
	if len(m.KV) == 0 {
		return fmt.Sprintf("{}")
	}

	out := []string{}
	for _, kv := range m.KV {
		out = append(out, kv.String())
	}
	return fmt.Sprintf("[%v]", strings.Join(out, ","))
}

type KV struct {
	Key   string
	Value Node
}

func (kv KV) String() string {
	return fmt.Sprintf("{%v:%v}", kv.Key, kv.Value)
}

func (kv KV) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}
