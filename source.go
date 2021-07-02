package main

import (
	"io/ioutil"
	"path/filepath"
)

type Source struct {
	Name string
	Code string
	Body []Stmt
}

type SourceResolver interface {
	Resolve(context *LoxContext, name string) (*Source, error)
}

type FileSourceResolver struct {
	Directory string
}

func MakeFileSourceResolver(directory string) *FileSourceResolver {
	return &FileSourceResolver{directory}
}

func (r *FileSourceResolver) Resolve(context *LoxContext, name string) (*Source, error) {
	filename := filepath.Join(r.Directory, name)
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	source := &Source{Name: name}
	source.Code = string(contents)

	scanner := MakeScanner(context, source.Code)
	tokens := scanner.scanTokens()

	if context.hadError {
		return source, nil
	}

	parser := MakeParser(context, tokens)
	statements, _ := parser.parse()

	if context.hadError {
		return source, nil
	}

	source.Body = statements
	return source, nil
}
