package huff

type Contract struct {
	Code       []byte
	DeployCode []byte
}

// defaultEVMVersion is the latest EVM version supported by Huff.
var defaultEVMVersion = EVMVersionShanghai

// EVMVersion represents the EVM version to compile for.
type EVMVersion string

const (
	EVMVersionShanghai   EVMVersion = "shanghai"
	EVMVersionParis      EVMVersion = "paris"
	EVMVersionLondon     EVMVersion = "london"
	EVMVersionBerlin     EVMVersion = "berlin"
	EVMVersionIstanbul   EVMVersion = "istanbul"
	EVMVersionPetersburg EVMVersion = "petersburg"
	EVMVersionByzantium  EVMVersion = "byzantium"
)
