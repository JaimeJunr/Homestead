package types

import "errors"

// Domain errors
var (
	// ErrNotFound indicates that a resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists indicates that a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidInput indicates that the input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrPermissionDenied indicates that permission was denied
	ErrPermissionDenied = errors.New("permission denied")

	// ErrExecutionFailed indicates that execution failed
	ErrExecutionFailed = errors.New("execution failed")

	// ErrDependencyNotMet indicates that a dependency is not met
	ErrDependencyNotMet = errors.New("dependency not met")
)
