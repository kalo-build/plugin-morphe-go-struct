package testutils

import "github.com/stretchr/testify/suite"

type FileSuite struct {
	suite.Suite
}

// FileEquals checks whether the file contents match. It also fails if
// the paths point to a directory or there is an error when trying to check a file.
func (f *FileSuite) FileEquals(actualPath string, expectedPath string, msgAndArgs ...any) bool {
	f.T().Helper()

	return AssertFileEquals(f.T(), actualPath, expectedPath, msgAndArgs...)
}

// FileContentsEquals checks whether a file matches the passed contents. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (f *FileSuite) FileContentsEquals(actualPath string, expectedContents []byte, msgAndArgs ...any) bool {
	f.T().Helper()

	return AssertFileContentsEquals(f.T(), actualPath, expectedContents, msgAndArgs...)
}
