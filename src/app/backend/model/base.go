package model

type Model struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type ListModel struct {
	Kind  string        `json:"kind"`
	Items []interface{} `json:"items"`
}

const (
	KIND_STAUTS = "Status"
)
