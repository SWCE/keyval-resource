package models

type EmptyVersion struct {
}

type Version map[string]string

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version  `json:"version"`
}

type OutParams struct {
	File string `json:"file"`
}

type OutRequest struct {
	Source Source `json:"source"`
	Params  OutParams `json:"params"`
}

type OutResponse struct {
	Version  Version  `json:"version"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version  Version  `json:"version"`
}

type CheckResponse []Version

type Source struct {}

