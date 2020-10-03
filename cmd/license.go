package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman/constants/legal"
	"github.com/tkw1536/ggman/program"
)

// License is the 'ggman license' command.
//
// The license command prints to standard output legal notices about the ggman program.
var License program.Command = license{}

type license struct{}

func (license) Name() string {
	return "license"
}

func (license) Options(flagset *pflag.FlagSet) program.Options {
	return program.Options{}
}

func (license) AfterParse() error {
	return nil
}

func (license) Run(context program.Context) error {
	context.Printf(stringLicenseInfo, stringLicenseText, legal.Notices)
	return nil
}

const stringLicenseInfo = `
ggman -- A golang script that can manage multiple git repositories locally
https://github.com/tkw1536/ggman

================================================================================
ggman is licensed under the terms of the MIT License:

%s
================================================================================

Furthermore, this executable may include code from the following projects:
%s
`

const stringLicenseText = `MIT License

Copyright (c) 2018-20 Tom Wiesing

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`
