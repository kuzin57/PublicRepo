package tools

import (
	"path/filepath"
)

type Validator struct {
	restrictTo []string
	exclude    []string
}

func NewValidator(restrictTo, exclude []string) *Validator {
	return &Validator{
		restrictTo: restrictTo,
		exclude:    exclude,
	}
}

func findMatch(name string, restrictions []string) (bool, error) {
	var (
		matched bool
		err     error
	)
	for _, restr := range restrictions {
		matched, err = filepath.Match(restr, name)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func (v *Validator) Validate(name, path string) (bool, error) {
	var (
		matched bool
		err     error
	)

	if len(v.restrictTo) > 0 {
		matchedName, e := findMatch(name, v.restrictTo)
		if e != nil {
			return true, e
		}

		matchedPath, e := findMatch(path, v.restrictTo)
		if e != nil {
			return true, e
		}

		if !matchedName && !matchedPath {
			return true, nil
		}
	}

	matched, err = findMatch(path, v.exclude)
	if err != nil || matched {
		return matched, err
	}

	matched, err = findMatch(name, v.exclude)
	if err != nil || matched {
		return matched, err
	}
	return false, nil
}
