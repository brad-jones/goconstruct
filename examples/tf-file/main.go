package main

import (
	"encoding/json"
	"fmt"

	"github.com/brad-jones/goconstruct"
)

// -- START: File Construct
// based on: https://registry.terraform.io/providers/hashicorp/local/latest/docs/resources/file

type FileProps struct {
	SensitiveContent    string
	ContentBase64       string
	Filename            string
	FilePermission      string
	DirectoryPermission string
	Source              string
	Content             string
}

type FileOutputs struct {
	*FileProps
}

func File(parent *goconstruct.Construct, id string, props *FileProps) *FileOutputs {
	return goconstruct.New(id,
		goconstruct.Scope(parent),
		goconstruct.Type("Local/File"),
		goconstruct.Props(props),
		goconstruct.Outputs(&FileOutputs{}),
	).Outputs.(*FileOutputs)
}

// -- START: FileData Construct
// based on: https://registry.terraform.io/providers/hashicorp/local/latest/docs/data-sources/file

type FileDataProps struct {
	Filename string
}

type FileDataOutputs struct {
	*FileDataProps
	Content       string
	ContentBase64 string
}

func FileData(parent *goconstruct.Construct, id string, props *FileDataProps) *FileDataOutputs {
	return goconstruct.New(id,
		goconstruct.Scope(parent),
		goconstruct.Type("Local/FileData"),
		goconstruct.Props(props),
		goconstruct.Outputs(&FileDataOutputs{}),
	).Outputs.(*FileDataOutputs)
}

// -- START: Construct Usage

func main() {
	root := goconstruct.New("/", goconstruct.Type("Root"),
		goconstruct.Constructor(func(c *goconstruct.Construct) {
			src := FileData(c, "MyFileData", &FileDataProps{
				Filename: "./go.mod",
			})
			File(c, "MyFile", &FileProps{
				Filename: "./foo.txt",
				Content:  src.Content,
			})
		}),
	)

	dat, _ := json.MarshalIndent(root, "", "    ")
	fmt.Println(string(dat))
}
