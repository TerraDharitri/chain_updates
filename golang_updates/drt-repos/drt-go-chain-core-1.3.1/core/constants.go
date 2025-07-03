package core

// HeaderType defines the type to be used for the header that is sent
type HeaderType string

const (
	// MetaHeader defines the type of *block.MetaBlock
	MetaHeader HeaderType = "MetaBlock"
	// ShardHeaderV1 defines the type of *block.Header
	ShardHeaderV1 HeaderType = "Header"
	// ShardHeaderV2 defines the type of *block.HeaderV2
	ShardHeaderV2 HeaderType = "HeaderV2"
)

// NodeType represents the node's role in the network
type NodeType string

// NodeTypeObserver signals that a node is running as observer node
const NodeTypeObserver NodeType = "observer"

// NodeTypeValidator signals that a node is running as validator node
const NodeTypeValidator NodeType = "validator"

// pkPrefixSize specifies the max numbers of chars to be displayed from one publc key
const pkPrefixSize = 12

// FileModeUserReadWrite represents the permission for a file which allows the user for reading and writing
const FileModeUserReadWrite = 0600

// FileModeReadWrite represents the permission for a file which allows reading and writing for user and group and read
// for others
const FileModeReadWrite = 0664

// MetachainShardId will be used to identify a shard ID as metachain
const MetachainShardId = uint32(0xFFFFFFFF)

// AllShardId will be used to identify that a message is for all shards
const AllShardId = uint32(0xFFFFFFF0)

// MegabyteSize represents the size in bytes of a megabyte
const MegabyteSize = 1024 * 1024

// MaxMachineIDLen is the maximum machine ID length
const MaxMachineIDLen = 10

// BuiltInFunctionClaimDeveloperRewards is the key for the claim developer rewards built-in function
const BuiltInFunctionClaimDeveloperRewards = "ClaimDeveloperRewards"

// BuiltInFunctionChangeOwnerAddress is the key for the change owner built in function built-in function
const BuiltInFunctionChangeOwnerAddress = "ChangeOwnerAddress"

// BuiltInFunctionSetUserName is the key for the set user name built-in function
const BuiltInFunctionSetUserName = "SetUserName"

// BuiltInFunctionSaveKeyValue is the key for the save key value built-in function
const BuiltInFunctionSaveKeyValue = "SaveKeyValue"

// BuiltInFunctionDCDTTransfer is the key for the electronic standard digital token transfer built-in function
const BuiltInFunctionDCDTTransfer = "DCDTTransfer"

// BuiltInFunctionDCDTBurn is the key for the electronic standard digital token burn built-in function
const BuiltInFunctionDCDTBurn = "DCDTBurn"

// BuiltInFunctionDCDTFreeze is the key for the electronic standard digital token freeze built-in function
const BuiltInFunctionDCDTFreeze = "DCDTFreeze"

// BuiltInFunctionDCDTUnFreeze is the key for the electronic standard digital token unfreeze built-in function
const BuiltInFunctionDCDTUnFreeze = "DCDTUnFreeze"

// BuiltInFunctionDCDTWipe is the key for the electronic standard digital token wipe built-in function
const BuiltInFunctionDCDTWipe = "DCDTWipe"

// BuiltInFunctionDCDTPause is the key for the electronic standard digital token pause built-in function
const BuiltInFunctionDCDTPause = "DCDTPause"

// BuiltInFunctionDCDTUnPause is the key for the electronic standard digital token unpause built-in function
const BuiltInFunctionDCDTUnPause = "DCDTUnPause"

// BuiltInFunctionSetDCDTRole is the key for the electronic standard digital token set built-in function
const BuiltInFunctionSetDCDTRole = "DCDTSetRole"

// BuiltInFunctionUnSetDCDTRole is the key for the electronic standard digital token unset built-in function
const BuiltInFunctionUnSetDCDTRole = "DCDTUnSetRole"

// BuiltInFunctionDCDTSetLimitedTransfer is the key for the electronic standard digital token built-in function which sets the property
// for the token to be transferable only through accounts with transfer roles
const BuiltInFunctionDCDTSetLimitedTransfer = "DCDTSetLimitedTransfer"

// BuiltInFunctionDCDTUnSetLimitedTransfer is the key for the electronic standard digital token built-in function which unsets the property
// for the token to be transferable only through accounts with transfer roles
const BuiltInFunctionDCDTUnSetLimitedTransfer = "DCDTUnSetLimitedTransfer"

// BuiltInFunctionDCDTLocalMint is the key for the electronic standard digital token local mint built-in function
const BuiltInFunctionDCDTLocalMint = "DCDTLocalMint"

// BuiltInFunctionDCDTLocalBurn is the key for the electronic standard digital token local burn built-in function
const BuiltInFunctionDCDTLocalBurn = "DCDTLocalBurn"

// BuiltInFunctionDCDTNFTTransfer is the key for the electronic standard digital token NFT transfer built-in function
const BuiltInFunctionDCDTNFTTransfer = "DCDTNFTTransfer"

// BuiltInFunctionDCDTNFTCreate is the key for the electronic standard digital token NFT create built-in function
const BuiltInFunctionDCDTNFTCreate = "DCDTNFTCreate"

// BuiltInFunctionDCDTNFTAddQuantity is the key for the electronic standard digital token NFT add quantity built-in function
const BuiltInFunctionDCDTNFTAddQuantity = "DCDTNFTAddQuantity"

// BuiltInFunctionDCDTNFTCreateRoleTransfer is the key for the electronic standard digital token create role transfer function
const BuiltInFunctionDCDTNFTCreateRoleTransfer = "DCDTNFTCreateRoleTransfer"

// BuiltInFunctionDCDTNFTBurn is the key for the electronic standard digital token NFT burn built-in function
const BuiltInFunctionDCDTNFTBurn = "DCDTNFTBurn"

// BuiltInFunctionDCDTNFTAddURI is the key for the electronic standard digital token NFT add URI built-in function
const BuiltInFunctionDCDTNFTAddURI = "DCDTNFTAddURI"

// BuiltInFunctionDCDTNFTUpdateAttributes is the key for the electronic standard digital token NFT update attributes built-in function
const BuiltInFunctionDCDTNFTUpdateAttributes = "DCDTNFTUpdateAttributes"

// BuiltInFunctionMultiDCDTNFTTransfer is the key for the electronic standard digital token multi transfer built-in function
const BuiltInFunctionMultiDCDTNFTTransfer = "MultiDCDTNFTTransfer"

// BuiltInFunctionSetGuardian is the key for setting a guardian built-in function
const BuiltInFunctionSetGuardian = "SetGuardian"

// BuiltInFunctionGuardAccount is the built-in function key for guarding an account
const BuiltInFunctionGuardAccount = "GuardAccount"

// BuiltInFunctionUnGuardAccount is the built-in function key for un-guarding an account
const BuiltInFunctionUnGuardAccount = "UnGuardAccount"

// BuiltInFunctionMigrateDataTrie is the built-in function key for migrating the data trie
const BuiltInFunctionMigrateDataTrie = "MigrateDataTrie"

// DCDTSetTokenType represents the builtin function name to set token type
const DCDTSetTokenType = "DCDTSetTokenType"

// DCDTModifyRoyalties represents the builtin function name to modify royalties
const DCDTModifyRoyalties = "DCDTModifyRoyalties"

// DCDTSetNewURIs represents the builtin function name to set new URIs for NFTs
const DCDTSetNewURIs = "DCDTSetNewURIs"

// DCDTModifyCreator represents the builtin function name to modify creator for NFTs
const DCDTModifyCreator = "DCDTModifyCreator"

// DCDTMetaDataRecreate represents the builtin function name to recreate the metadata for DCDT tokens
const DCDTMetaDataRecreate = "DCDTMetaDataRecreate"

// DCDTMetaDataUpdate represents the builtin function name to update the metadata for DCDT tokens
const DCDTMetaDataUpdate = "DCDTMetaDataUpdate"

// DCDTRoleLocalMint is the constant string for the local role of mint for DCDT tokens
const DCDTRoleLocalMint = "DCDTRoleLocalMint"

// DCDTRoleLocalBurn is the constant string for the local role of burn for DCDT tokens
const DCDTRoleLocalBurn = "DCDTRoleLocalBurn"

// DCDTRoleNFTCreate is the constant string for the local role of create for DCDT NFT tokens
const DCDTRoleNFTCreate = "DCDTRoleNFTCreate"

// DCDTRoleNFTCreateMultiShard is the constant string for the local role of create for DCDT NFT tokens multishard
const DCDTRoleNFTCreateMultiShard = "DCDTRoleNFTCreateMultiShard"

// DCDTRoleNFTAddQuantity is the constant string for the local role of adding quantity for existing DCDT NFT tokens
const DCDTRoleNFTAddQuantity = "DCDTRoleNFTAddQuantity"

// DCDTRoleNFTBurn is the constant string for the local role of burn for DCDT NFT tokens
const DCDTRoleNFTBurn = "DCDTRoleNFTBurn"

// DCDTRoleNFTAddURI is the constant string for the local role of adding a URI for DCDT NFT tokens
const DCDTRoleNFTAddURI = "DCDTRoleNFTAddURI"

// DCDTRoleNFTUpdateAttributes is the constant string for the local role of updating attributes for DCDT NFT tokens
const DCDTRoleNFTUpdateAttributes = "DCDTRoleNFTUpdateAttributes"

// DCDTRoleTransfer is the constant string for the local role to transfer DCDT, only for special tokens
const DCDTRoleTransfer = "DCDTTransferRole"

// DCDTRoleSetNewURI represents the role which can rewrite the URI in the token metadata
const DCDTRoleSetNewURI = "DCDTRoleSetNewURI"

// DCDTRoleModifyRoyalties represents the role which can rewrite the royalties of a token
const DCDTRoleModifyRoyalties = "DCDTRoleModifyRoyalties"

// DCDTRoleModifyCreator represents the role which can rewrite the creator in the token metadata
const DCDTRoleModifyCreator = "DCDTRoleModifyCreator"

// DCDTRoleNFTRecreate represents the role which can recreate the token metadata
const DCDTRoleNFTRecreate = "DCDTRoleNFTRecreate"

// DCDTRoleNFTUpdate represents the role which can update the token metadata
const DCDTRoleNFTUpdate = "DCDTRoleNFTUpdate"

// DCDTType defines the possible types in case of DCDT tokens
type DCDTType uint32

const (
	// Fungible defines the token type for DCDT fungible tokens
	Fungible DCDTType = iota
	// NonFungible defines the token type for DCDT non fungible tokens
	NonFungible
	// NonFungibleV2 defines the token type for DCDT non fungible tokens
	NonFungibleV2
	// SemiFungible defines the token type for DCDT semi fungible tokens
	SemiFungible
	// MetaFungible defines the token type for DCDT meta fungible tokens
	MetaFungible
	// DynamicNFT defines the token type for DCDT dynamic NFT tokens
	DynamicNFT
	// DynamicSFT defines the token type for DCDT dynamic SFT tokens
	DynamicSFT
	// DynamicMeta defines the token type for DCDT dynamic meta tokens
	DynamicMeta
)

// String will convert number type in string
func (t DCDTType) String() string {
	switch t {
	case Fungible:
		return FungibleDCDT
	case NonFungible:
		return NonFungibleDCDT
	case NonFungibleV2:
		return NonFungibleDCDTv2
	case SemiFungible:
		return SemiFungibleDCDT
	case MetaFungible:
		return MetaDCDT
	case DynamicNFT:
		return DynamicNFTDCDT
	case DynamicSFT:
		return DynamicSFTDCDT
	case DynamicMeta:
		return DynamicMetaDCDT
	default:
		return "Unknown"
	}
}

// FungibleDCDT defines the string for the token type of fungible DCDT
const FungibleDCDT = "FungibleDCDT"

// NonFungibleDCDT defines the string for the token type of non fungible DCDT
const NonFungibleDCDT = "NonFungibleDCDT"

// NonFungibleDCDTv2 defines the string for the token type of non fungible DCDT
const NonFungibleDCDTv2 = "NonFungibleDCDTv2"

// MetaDCDT defines the string for the token type of meta DCDT
const MetaDCDT = "MetaDCDT"

// SemiFungibleDCDT defines the string for the token type of semi fungible DCDT
const SemiFungibleDCDT = "SemiFungibleDCDT"

// Dynamic is the prefix used for dynamic DCDT tokens
const Dynamic = "Dynamic"

// DynamicNFTDCDT defines the string for the token type of dynamic NFT DCDT
const DynamicNFTDCDT = Dynamic + NonFungibleDCDT

// DynamicSFTDCDT defines the string for the token type of dynamic SFT DCDT
const DynamicSFTDCDT = Dynamic + SemiFungibleDCDT

// DynamicMetaDCDT defines the string for the token type of dynamic meta DCDT
const DynamicMetaDCDT = Dynamic + MetaDCDT

// MaxRoyalty defines 100% as uint32
const MaxRoyalty = uint32(10000)

// RelayedTransaction is the key for the electronic meta/gassless/relayed transaction standard
const RelayedTransaction = "relayedTx"

// RelayedTransactionV2 is the key for the optimized electronic meta/gassless/relayed transaction standard
const RelayedTransactionV2 = "relayedTxV2"

// SCDeployInitFunctionName is the key for the function which is called at smart contract deploy time
const SCDeployInitFunctionName = "_init"

// ProtectedKeyPrefix is the key prefix which is protected from writing in the trie - only for special builtin functions
const ProtectedKeyPrefix = "N" + "U" + "M" + "B" + "A" + "T"

// DelegationSystemSCKey is the key under which there is data in case of system delegation smart contracts
const DelegationSystemSCKey = "delegation"

// DCDTKeyIdentifier is the key prefix for dcdt tokens
const DCDTKeyIdentifier = "dcdt"

// DCDTRoleIdentifier is the key prefix for dcdt role identifier
const DCDTRoleIdentifier = "role"

// DCDTNFTLatestNonceIdentifier is the key prefix for dcdt latest nonce identifier
const DCDTNFTLatestNonceIdentifier = "nonce"

// GuardiansKeyIdentifier is the key prefix for guardians
const GuardiansKeyIdentifier = "guardians"

// MaxNumShards represents the maximum number of shards possible in the system
const MaxNumShards = 256

// MinMetaTxExtraGasCost is the constant defined for minimum gas value to be sent in meta transaction
const MinMetaTxExtraGasCost = uint64(1_000_000)

// MaxLeafSize represents maximum amount of data which can be saved under one leaf
const MaxLeafSize = uint64(1 << 26) //64MB

// MaxBufferSizeToSendTrieNodes represents max buffer size to send in bytes used when resolving trie nodes
// Every trie node that has a greater size than this constant is considered a large trie node and should be split in
// smaller chunks
const MaxBufferSizeToSendTrieNodes = 1 << 18 //256KB

// MaxUserNameLength represents the maximum number of bytes a UserName can have
const MaxUserNameLength = 32

// MinLenArgumentsDCDTTransfer defines the min length of arguments for the DCDT transfer
const MinLenArgumentsDCDTTransfer = 2

// MinLenArgumentsDCDTNFTTransfer defines the minimum length for dcdt nft transfer
const MinLenArgumentsDCDTNFTTransfer = 4

// MaxLenForDCDTIssueMint defines the maximum length in bytes for the issued/minted balance
const MaxLenForDCDTIssueMint = 100

// BaseOperationCostString represents the field name for base operation costs
const BaseOperationCostString = "BaseOperationCost"

// BuiltInCostString represents the field name for built in operation costs
const BuiltInCostString = "BuiltInCost"

// DCDTSCAddress is the hard-coded address for dcdt issuing smart contract
var DCDTSCAddress = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 255, 255}

// SCDeployIdentifier is the identifier for a smart contract deploy
const SCDeployIdentifier = "SCDeploy"

// SCUpgradeIdentifier is the identifier for a smart contract upgrade
const SCUpgradeIdentifier = "SCUpgrade"

// WriteLogIdentifier is the identifier for the information log that is generated by a smart contract call/dcdt transfer
const WriteLogIdentifier = "writeLog"

// SignalErrorOperation is the identifier for the log that is generated when a smart contract is executed but return an error
const SignalErrorOperation = "signalError"

// CompletedTxEventIdentifier is the identifier for the log that is generated when the execution of a smart contract call is done
const CompletedTxEventIdentifier = "completedTxEvent"

// InternalVMErrorsOperation is the identifier for the log that is generated when the execution of a smart contract generates an interval vm error
const InternalVMErrorsOperation = "internalVMErrors"

// GasRefundForRelayerMessage is the return message for to the smart contract result with refund for the relayer
const GasRefundForRelayerMessage = "gas refund for relayer"

// InitialVersionOfTransaction defines the initial version for a transaction
const InitialVersionOfTransaction = uint32(1)

// DefaultAddressPrefix is the default hrp of Dharitri/Numbat
const DefaultAddressPrefix = "drt"

// TopicRequestSuffix represents the topic name suffix for requests
const TopicRequestSuffix = "_REQUEST"
