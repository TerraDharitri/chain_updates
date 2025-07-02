package abi

// MultiValue is a multi-value
type MultiValue struct {
	Items []any
}

// VariadicValues holds variadic values
type VariadicValues struct {
	Items       []any
	ItemCreator func() any
}

// OptionalValue holds an optional value
type OptionalValue struct {
	Value any
}
