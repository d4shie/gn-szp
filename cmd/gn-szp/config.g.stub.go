// This is a stub file that is only used when the editor build constraint is
// present in the environment.
//
// This allows us to fool the editor into showing us definitions before
// generation, as well as prevent warnings. If introducing a new
// generated variable, please introduce it here first so that LSP
// is able to pick up on it.
//
// Guidelines:
// As the generator will create a fixed slice, whose size we cannot know
// during development (defined as constant at compile time), all definitions
// should be created as variable-sized slices.
//go:build editor || CODEGEN

package main

var (
	gSfvDirs    []string
	gGroupDirs  []string
	gIgnoreDirs []string
)
