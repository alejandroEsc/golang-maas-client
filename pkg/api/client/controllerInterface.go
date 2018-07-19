package client

import (
	"net/url"

	"github.com/juju/utils/set"
)

type ControllerInterface interface {
	Put(path string, params url.Values) ([]byte, error)

	Post(path string, op string, params url.Values) ([]byte, error)

	PostFile(path string, op string, params url.Values, fileContent []byte) ([]byte, error)

	Delete(path string) error

	Get(path string, op string, params url.Values) ([]byte, error)

	GetAPIVersionInfo() (set.Strings, error)
}



type ApiHelper interface {

	// Files
	GetFile(filename string) (*[]byte, error)
	ReadFileContent(filename string) ([]byte, error)

	// Fabrics
	Fabrics() ([]byte, error)

	//Spaces
	Spaces() ([]byte, error)
}