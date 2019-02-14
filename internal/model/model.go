package model

// ModuleFile contains binary and the checksum
type ModuleFile struct {
	Path string
	Sum  string
}

// ModuleMeta is a metainformation of the module
type ModuleMeta struct {
	Name      string            `json:"package"`
	ID        string            `json:"id"`
	Version   int               `json:"version"`
	Extension map[string]string `json:"extension,omitempty"`
}

// Module contains package metadata, binary and a checksum
type Module struct {
	File ModuleFile
	Meta ModuleMeta
}
