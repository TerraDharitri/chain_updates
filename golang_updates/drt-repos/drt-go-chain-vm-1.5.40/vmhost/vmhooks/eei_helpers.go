package vmhooks

const dcdtRoleLocalMint = "DCDTRoleLocalMint"
const dcdtRoleLocalBurn = "DCDTRoleLocalBurn"
const dcdtRoleNFTCreate = "DCDTRoleNFTCreate"
const dcdtRoleNFTAddQuantity = "DCDTRoleNFTAddQuantity"
const dcdtRoleNFTBurn = "DCDTRoleNFTBurn"
const dcdtRoleNFTUpdateAttributes = "DCDTRoleNFTUpdateAttributes"
const dcdtRoleNFTAddURI = "DCDTRoleNFTAddURI"
const dcdtRoleNFTRecreate = "DCDTRoleNFTRecreate"
const dcdtRoleModifyCreator = "DCDTRoleModifyCreator"
const dcdtRoleModifyRoyalties = "DCDTRoleModifyRoyalties"
const dcdtRoleSetNewURI = "DCDTRoleSetNewURI"

const tickerMinLength = 3
const tickerMaxLength = 10
const additionalRandomCharsLength = 6
const identifierMinLength = tickerMinLength + additionalRandomCharsLength + 1
const identifierMaxLength = tickerMaxLength + additionalRandomCharsLength + 1

// constants defining roles values
const (
	RoleMint = 1 << iota
	RoleBurn
	RoleNFTCreate
	RoleNFTAddQuantity
	RoleNFTBurn
	RoleNFTUpdateAttributes
	RoleNFTAddURI
	RoleNFTRecreate
	RoleModifyCreator
	RoleModifyRoyalties
	RoleSetNewURI
)

func roleFromByteArray(bytes []byte) int64 {
	stringValue := string(bytes)
	switch stringValue {
	case dcdtRoleLocalMint:
		return RoleMint
	case dcdtRoleLocalBurn:
		return RoleBurn
	case dcdtRoleNFTCreate:
		return RoleNFTCreate
	case dcdtRoleNFTAddQuantity:
		return RoleNFTAddQuantity
	case dcdtRoleNFTBurn:
		return RoleNFTBurn
	default:
		return 0
	}
}

func roleFromByteArrayV2(bytes []byte) int64 {
	stringValue := string(bytes)
	switch stringValue {
	case dcdtRoleLocalMint:
		return RoleMint
	case dcdtRoleLocalBurn:
		return RoleBurn
	case dcdtRoleNFTCreate:
		return RoleNFTCreate
	case dcdtRoleNFTAddQuantity:
		return RoleNFTAddQuantity
	case dcdtRoleNFTBurn:
		return RoleNFTBurn
	case dcdtRoleNFTUpdateAttributes:
		return RoleNFTUpdateAttributes
	case dcdtRoleNFTAddURI:
		return RoleNFTAddURI
	case dcdtRoleNFTRecreate:
		return RoleNFTRecreate
	case dcdtRoleModifyCreator:
		return RoleModifyCreator
	case dcdtRoleModifyRoyalties:
		return RoleModifyRoyalties
	case dcdtRoleSetNewURI:
		return RoleSetNewURI
	default:
		return 0
	}
}

func getDCDTRoles(dataBuffer []byte, cryptoOpcodesV2Enabled bool) int64 {
	result := int64(0)
	currentIndex := 0
	valueLen := len(dataBuffer)

	for currentIndex < valueLen {
		// first character before each role is a \n, so we skip it
		currentIndex += 1

		// next is the length of the role as string
		roleLen := int(dataBuffer[currentIndex])
		currentIndex += 1

		// next is role's ASCII string representation
		endIndex := currentIndex + roleLen
		roleName := dataBuffer[currentIndex:endIndex]
		currentIndex = endIndex

		if cryptoOpcodesV2Enabled {
			result |= roleFromByteArrayV2(roleName)
		} else {
			result |= roleFromByteArray(roleName)
		}
	}
	return result
}

// ValidateToken - validates the token ID
func ValidateToken(tokenID []byte) bool {
	tokenIDLen := len(tokenID)
	if tokenIDLen < identifierMinLength || tokenIDLen > identifierMaxLength {
		return false
	}

	tickerLen := tokenIDLen - additionalRandomCharsLength

	if !isTickerValid(tokenID[0 : tickerLen-1]) {
		return false
	}

	// dash char between the random chars and the ticker
	if tokenID[tickerLen-1] != '-' {
		return false
	}

	if !randomCharsAreValid(tokenID[tickerLen:tokenIDLen]) {
		return false
	}

	return true
}

// ticker must be all uppercase alphanumeric
func isTickerValid(tickerName []byte) bool {
	if len(tickerName) < tickerMinLength || len(tickerName) > tickerMaxLength {
		return false
	}
	for _, ch := range tickerName {
		isBigCharacter := ch >= 'A' && ch <= 'Z'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isBigCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}

// random chars are alphanumeric lowercase
func randomCharsAreValid(chars []byte) bool {
	if len(chars) != additionalRandomCharsLength {
		return false
	}
	for _, ch := range chars {
		isSmallCharacter := ch >= 'a' && ch <= 'f'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isSmallCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}
