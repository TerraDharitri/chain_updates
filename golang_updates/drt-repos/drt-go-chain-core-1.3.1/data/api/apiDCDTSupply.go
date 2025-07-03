package api

// DCDTSupply represents the structure for dcdt supply that is returned by api routes
type DCDTSupply struct {
	InitialMinted    string `json:"initialMinted"`
	Supply           string `json:"supply"`
	Burned           string `json:"burned"`
	Minted           string `json:"minted"`
	RecomputedSupply bool   `json:"recomputedSupply"`
}
