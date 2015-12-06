// Package interfaces
package interfaces

// BUG(rjeczalik): Does not work with recursive types.

// BUG(rjeczalik): Does not work with with more than one level of indirection
// (pointer to pointers).

// BUG(rjeczalik): Does not and will not work with struct literals.

// BUG(rjeczalik): May incorrectly generate dependencies for a map types which
// key and value are named types imported from different packages.
// As a workaround run goimports over the output file.
