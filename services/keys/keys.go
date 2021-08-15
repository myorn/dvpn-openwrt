package keys

type Key struct {
	Name     string `json:"Name"`
	Operator string `json:"Operator"`
	Address  string `json:"Address"`
}

type Keys struct {
	Keys []Key `json:"Keys"`
}

type AddRecoverRequest struct {
	Name     string `json:"Name"`
	Mnemonic string `json:"Mnemonic"`
}
