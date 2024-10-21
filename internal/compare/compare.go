package compare

import (
	version "github.com/knqyf263/go-rpm-version"
)

func Compare(version1 string, version2 string) (res int, err error) {
	v1 := version.NewVersion(version1)
	v2 := version.NewVersion(version2)

	switch {
	case v1.LessThan(v2):
		res = -1
	case v1.Equal(v2):
		res = 0
	case v1.GreaterThan(v2):
		res = 1
	}

	return
}
