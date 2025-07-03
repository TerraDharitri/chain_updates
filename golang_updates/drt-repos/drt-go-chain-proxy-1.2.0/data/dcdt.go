package data

const (
	FungibleTokens     = "fungible-tokens"
	SemiFungibleTokens = "semi-fungible-tokens"
	NonFungibleTokens  = "non-fungible-tokens"
)

// ValidTokenTypes holds a slice containing the valid dcdt token types
var ValidTokenTypes = []string{FungibleTokens, SemiFungibleTokens, NonFungibleTokens}

// DCDTSupplyResponse is a response holding dcdt supply
type DCDTSupplyResponse struct {
	Data  DCDTSupply `json:"data"`
	Error string     `json:"error"`
	Code  ReturnCode `json:"code"`
}

// DCDTSupply is a DTO holding dcdt supply
type DCDTSupply struct {
	Supply           string `json:"supply"`
	Minted           string `json:"minted"`
	Burned           string `json:"burned"`
	InitialMinted    string `json:"initialMinted"`
	RecomputedSupply bool   `json:"recomputedSupply"`
}

// IsValidDcdtPath returns true if the provided path is a valid dcdt token type
func IsValidDcdtPath(path string) bool {
	for _, tokenType := range ValidTokenTypes {
		if tokenType == path {
			return true
		}
	}

	return false
}
