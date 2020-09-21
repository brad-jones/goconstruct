# goconstruct

[![PkgGoDev](https://pkg.go.dev/badge/github.com/brad-jones/goconstruct)](https://pkg.go.dev/github.com/brad-jones/goconstruct)
[![GoReport](https://goreportcard.com/badge/github.com/brad-jones/goconstruct)](https://goreportcard.com/report/github.com/brad-jones/goconstruct)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.15.1-lightblue.svg)](https://golang.org)
![.github/workflows/main.yml](https://github.com/brad-jones/goconstruct/workflows/.github/workflows/main.yml/badge.svg?branch=master)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![KeepAChangelog](https://img.shields.io/badge/Keep%20A%20Changelog-1.0.0-%23E05735)](https://keepachangelog.com/)
[![License](https://img.shields.io/github/license/brad-jones/goconstruct.svg)](https://github.com/brad-jones/goconstruct/blob/master/LICENSE)

A _"go"_ version of [github.com/aws/constructs](https://github.com/aws/constructs),
based on functions instead of classes.

## Quick Start

`go get -u github.com/brad-jones/goconstruct`

```go
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
```

Outputs the following JSON:

```json
{
    "id": "/",
    "type": "Root",
    "props": null,
    "children": [
        {
            "id": "/MyFileData",
            "type": "Local/FileData",
            "props": {
                "Filename": "./go.mod"
            },
            "children": []
        },
        {
            "id": "/MyFile",
            "type": "Local/File",
            "props": {
                "SensitiveContent": "",
                "ContentBase64": "",
                "Filename": "./foo.txt",
                "FilePermission": "",
                "DirectoryPermission": "",
                "Source": "",
                "Content": "/MyFileData.Content"
            },
            "children": []
        }
    ]
}
```

The idea being that this JSON would then be passed on to a parser/generator
which understands how to convert it into Terraform JSON/HCL or Cloudformation
or say Github Actions YAML as the case may be.

_Also see further working examples under: <https://github.com/brad-jones/goconstruct/tree/master/examples>_

## Motivation

This project started life because I wanted the AWS CDK experience but with Terraform.
I had played with <https://www.pulumi.com/> but found that it had strayed too far
from what terraform is, ie: I couldn't `pulumi synth` and see some terraform.

I later discovered <https://github.com/hashicorp/terraform-cdk> however decided
to continue with this mainly because _"go"_ still isn't supported by <https://github.com/aws/jsii>.

This RFC PR <https://github.com/aws/aws-cdk-rfcs/pull/206> details the specifics
of _"go"_ support in the CDK and to quote the RFC:

> The programming model for Go differs significantly from that of Typescript.
> Imposing an object-oriented philosophy on a procedural language may result in
> non-idiomatic constructs and APIs in the target language. However, the tradeoff
> for having CDK constructs available in more languages outweighs this disadvantage.

I then asked myself, _"Well what would it look like in go if I started from scratch?"_
And this is the result. As it is largely based on functions instead of classes it
should translate well to almost any language.

It is in use _(or will be)_ at <https://github.com/brad-jones/tdk>

My ultimate goal is to be able to define a stack using `go` and compile it all
down to a single binary. eg: `mystack-v1.0.0 plan && mystack-v1.0.0 apply`.

My immediate itch that I am scratching with all of this is to be able to download
that single binary on a new development machine, be it Linux/MacOS or Windows and
have my dev environment automatically configured. I have tried out various home
directory _"dot-file"_ managers (latest being <https://www.chezmoi.io>) but they
have always left something to be desired.

At work _([Xero](https://www.xero.com))_ I can also see other advantages of being
able to create a binary artifact in a CI/CD pipeline that represents a particular
version of a stack.
