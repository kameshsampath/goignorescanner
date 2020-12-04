/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scanner

import (
	"path/filepath"
	"regexp"

	"github.com/karrick/godirwalk"
	"k8s.io/apimachinery/pkg/util/sets"
)

// IgnorePattern holds the ignoreable patterns
type IgnorePattern struct {
	Pattern      string
	Paths        sets.String
	RegexPattern *regexp.Regexp
	Invert       bool
	IsDir        bool
}

// DirectoryScanner  helps identifying if a BundleFile needs to be ignored
type DirectoryScanner interface {
	// Scan checks file has to be ignored or not, returns true if it needs to be ignored
	Scan() ([]string, error)
}

var (
	_               DirectoryScanner = (*defaultIgnorer)(nil)
	defaultPatterns                  = []string{".git", "vendor", "node_modules"}
)

var dirOpts = &godirwalk.Options{
	Callback: func(osPathname string, dirEntry *godirwalk.Dirent) error {

		// if osPathName is one among the default excludes skip walking into them
		df := sets.NewString(defaultPatterns...)
		if df.Has(filepath.Base(osPathname)) {
			return godirwalk.SkipThis
		}

		for _, igp := range IgnorePatterns {

			re := igp.RegexPattern

			regxMatches := re.FindAllStringSubmatch(osPathname, -1)
			// since we ignore, just add to the includes if and only if we know its inverted pattern
			if len(regxMatches) > 0 && igp.Invert {
				for _, tuple := range regxMatches {
					includes = append(includes, tuple[0])
				}
			}

		}

		return nil
	},
}

// defaultIgnorer is the default DirectoryScanner which is returned when no .dockerignore file is present
// or error processing .dockerignore
type defaultIgnorer struct{}

// Ignore implements DirectoryScanner, for no dockerignore cases, where only .git is ignored
func (i *defaultIgnorer) Scan() ([]string, error) {
	err := godirwalk.Walk(directory, dirOpts)
	return includes, err
}
