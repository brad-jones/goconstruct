package main_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTfFile(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		assert.Equal(t,
			[]string{
				"{",
				"    \"id\": \"/\",",
				"    \"type\": \"Root\",",
				"    \"props\": null,",
				"    \"children\": [",
				"        {",
				"            \"id\": \"/MyFileData\",",
				"            \"type\": \"Local/FileData\",",
				"            \"props\": {",
				"                \"Filename\": \"./go.mod\"",
				"            },",
				"            \"children\": []",
				"        },",
				"        {",
				"            \"id\": \"/MyFile\",",
				"            \"type\": \"Local/File\",",
				"            \"props\": {",
				"                \"SensitiveContent\": \"\",",
				"                \"ContentBase64\": \"\",",
				"                \"Filename\": \"./foo.txt\",",
				"                \"FilePermission\": \"\",",
				"                \"DirectoryPermission\": \"\",",
				"                \"Source\": \"\",",
				"                \"Content\": \"/MyFileData.Content\"",
				"            },",
				"            \"children\": []",
				"        }",
				"    ]",
				"}",
				"",
			},
			normaliseCmdOutput(out),
		)
	}
}

func normaliseCmdOutput(in []byte) []string {
	out := string(in)
	return strings.Split(out, "\n")
}
