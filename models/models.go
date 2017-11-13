package models

type EmptyVersion struct {
	Dummy string  `json:"dummy"`
}

type Version map[string]string

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type OutRequest struct {
	Source Source `json:"source"`
	Params  OutParams `json:"params"`
}

type OutParams struct {
	File string `json:"file"`
}

type OutResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version  EmptyVersion  `json:"version"`
}

type CheckResponse []EmptyVersion

type Source struct {}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
