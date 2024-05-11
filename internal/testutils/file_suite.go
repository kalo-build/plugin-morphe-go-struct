package testutils

import "github.com/stretchr/testify/suite"

type FileSuite struct {
	suite.Suite
}

// FileEquals checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (f *FileSuite) FileEquals(actualPath string, expectedPath string, msgAndArgs ...any) bool {
	f.T().Helper()

	return AssertFileEquals(f.T(), actualPath, expectedPath, msgAndArgs...)
}
