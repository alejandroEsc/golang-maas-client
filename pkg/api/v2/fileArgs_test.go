package v2

import (
	"bytes"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestAddFileArgsValidate(t *testing.T) {
	reader := bytes.NewBufferString("test")
	for _, test := range []struct {
		args    AddFileArgs
		errText string
	}{{
		errText: "missing Filename not valid",
	}, {
		args:    AddFileArgs{Filename: "/foo"},
		errText: `paths in Filename "/foo" not valid`,
	}, {
		args:    AddFileArgs{Filename: "a/foo"},
		errText: `paths in Filename "a/foo" not valid`,
	}, {
		args:    AddFileArgs{Filename: "foo.txt"},
		errText: `missing Content or Reader not valid`,
	}, {
		args: AddFileArgs{
			Filename: "foo.txt",
			Reader:   reader,
		},
		errText: `missing Length not valid`,
	}, {
		args: AddFileArgs{
			Filename: "foo.txt",
			Reader:   reader,
			Length:   4,
		},
	}, {
		args: AddFileArgs{
			Filename: "foo.txt",
			Content:  []byte("foo"),
			Reader:   reader,
		},
		errText: `specifying Content and Reader not valid`,
	}, {
		args: AddFileArgs{
			Filename: "foo.txt",
			Content:  []byte("foo"),
			Length:   20,
		},
		errText: `specifying Length and Content not valid`,
	}, {
		args: AddFileArgs{
			Filename: "foo.txt",
			Content:  []byte("foo"),
		},
	}} {
		err := test.args.Validate()
		if test.errText == "" {
			assert.Nil(t, err)
		} else {
			assert.True(t, errors.IsNotValid(err))
			assert.EqualError(t, err, test.errText)
		}
	}
}
