// Package maker provides CLI tools to manage and use GNU make-compliant snippet
// files, which contain rules to deal with specific domains such as programming
// language toolchains or helper tools as docker, linters, dependency injectors
// (google wire, uber fx, etc) that are needed during any stage of a repository
// management.
//
// Those snippet files rely on a convention-over-configuration style
// that uses variables that can be either set before the include happens or
// overridden later on. They also should not be modified directly on the client
// repositories, as they are meant to be updated with the upstream version.
//
// Copyright (c) William Artero. MIT License
package maker // import "github.com/wwmoraes/maker"
