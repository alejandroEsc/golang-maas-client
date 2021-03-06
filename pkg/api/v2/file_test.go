// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE File for details.

package v2

import (
	"net/http"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestReadFilesBadSchema(t *testing.T) {
	var f File
	err = json.Unmarshal([]byte("wat?"), &f)
	assert.Error(t, err)
}

func TestReadFiles(t *testing.T) {
	var files []File
	err = json.Unmarshal([]byte(filesResponse), &files)
	assert.Nil(t, err)
	assert.Len(t, files, 2)
	file := files[0]
	assert.Equal(t, file.Filename, "test")
}

func TestReadAllFromGetFile(t *testing.T) {
	// When Get File is used, the response includes the body of the File
	// base64 encoded, so ReadAll just decodes it.
	server, controller := createTestServerController(t)
	server.AddGetResponse("/api/2.0/files/testing/", http.StatusOK, fileResponse)
	defer server.Close()

	file, err := controller.GetFile("testing")
	assert.Nil(t, err)
	assert.Equal(t, file.Content, "this is a test\n")
}

func TestReadAllFromFiles(t *testing.T) {
	// When Get File is used, the response includes the body of the File
	// base64 encoded, so ReadAll just decodes it.
	server, controller := createTestServerController(t)
	defer server.Close()
	server.AddGetResponse("/api/2.0/files/", http.StatusOK, filesResponse)
	server.AddGetResponse("/api/2.0/files/?Filename=test&op=Get", http.StatusOK, "some Content\n")

	files, err := controller.getFiles("")
	assert.Nil(t, err)
	file := files[0]
	content, err := controller.ReadFileContent(&file)
	assert.Nil(t, err)
	assert.Equal(t, "some Content\n", string(content))
}

const (
	fileResponse = `
{
    "resource_uri": "/MAAS/api/2.0/files/testing/",
    "Content": "dGhpcyBpcyBhIHRlc3QK",
    "anon_resource_uri": "/MAAS/api/2.0/files/?op=get_by_key&key=88e64b76-fb82-11e5-932f-52540051bf22",
    "Filename": "testing"
}
`
	filesResponse = `
[
    {
        "resource_uri": "/MAAS/api/2.0/files/test/",
        "anon_resource_uri": "/MAAS/api/2.0/files/?op=get_by_key&key=3afba564-fb7d-11e5-932f-52540051bf22",
        "Filename": "test"
    },
    {
        "resource_uri": "/MAAS/api/2.0/files/test-File.txt/",
        "anon_resource_uri": "/MAAS/api/2.0/files/?op=get_by_key&key=69913e62-fad2-11e5-932f-52540051bf22",
        "Filename": "test-File.txt"
    }
]
`
)
