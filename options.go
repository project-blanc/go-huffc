package huffc

// Options for the [Compiler]. A zero Options consists entirely of default values.
type Options struct {
	EVMVersion EVMVersion // EVM version to compile for (default: shanghai)
}

func (opts *Options) setDefaults() {
	if opts.EVMVersion == "" {
		opts.EVMVersion = defaultEVMVersion
	}
}

// defaultEVMVersion is the latest EVM version supported by Huff.
var defaultEVMVersion = EVMVersionShanghai

// EVMVersion represents the EVM version to compile for.
type EVMVersion string

const (
	EVMVersionShanghai EVMVersion = "shanghai"
	EVMVersionParis    EVMVersion = "paris"
)
