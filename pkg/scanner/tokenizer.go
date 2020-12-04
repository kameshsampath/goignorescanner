package scanner

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"

	"k8s.io/apimachinery/pkg/util/sets"
)

// tokenize is used to perform normalization on the pattern like
// - clean up the path to be good file path
// - make sure paths start with /
// - check if the patterns has inversions i.e !foo kind of things

func tokenize(pattern string) (*IgnorePattern, error) {

	// clean the pattern to be well formed Go path
	filepath.Clean(pattern)

	// make sure the path starts with /
	filepath.FromSlash(pattern)

	ignorePattern := &IgnorePattern{}
	ignorePattern.Pattern = pattern
	paths := sets.NewString()

	// check if it has inverts and remove them before creating paths
	if strings.HasPrefix(pattern, "!") {
		pattern = strings.TrimPrefix(pattern, "!")
		ignorePattern.Invert = true
	}

	// remove the root slash and split the patterns as paths
	if strings.HasPrefix(pattern, string(os.PathSeparator)) {
		pattern = strings.TrimPrefix(pattern, string(os.PathSeparator))
		t := strings.Split(pattern, string(os.PathSeparator))
		paths.Insert(string(os.PathSeparator))
		paths.Insert(t...)
	} else {
		t := strings.Split(pattern, string(os.PathSeparator))
		paths.Insert(t...)
	}

	ignorePattern.Paths = paths
	ignorePattern.IsDir = isDir(pattern)

	expr, err := asRegExp(pattern)

	expr = strings.TrimPrefix(expr, "^")

	// Since the patterns are relative to the root, compile regex with dir root
	// prepended to it
	expr = filepath.Join(directory, expr)

	expr = "^" + expr

	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(expr)

	if err != nil {
		return nil, err
	}

	ignorePattern.RegexPattern = re

	return ignorePattern, nil
}

// asRegExp  builds a regular expression of the pattern
func asRegExp(pattern string) (string, error) {
	pathSep := string(os.PathSeparator)
	escPath := pathSep

	// make sure the unix paths are escaped with \\
	if pathSep == `\` {
		escPath += `\`
	}

	//start
	rexPat := "^"

	var s scanner.Scanner
	s.Init(strings.NewReader(pattern))

	for s.Peek() != scanner.EOF {
		ch := s.Next()

		//handle *
		if '*' == ch {
			if '*' == s.Peek() {
				//check if next char is also *, typically like **
				s.Next()
				//Treat **/ as **
				if string(s.Peek()) == pathSep {
					s.Next()
				}

				//If pattern EOF
				if s.Peek() == scanner.EOF {
					rexPat += ".*"
				} else {
					//make sure we escape  path seperator after **
					rexPat += "(.*" + escPath + ")?"
				}
			} else {
				rexPat += ".*"
			}
		} else if '?' == ch {
			// make sure ? escapes any character than path seperator
			rexPat += "[^" + pathSep + "]"
		} else if '.' == ch || ch == '$' {
			rexPat += `\` + string(ch)
		} else if ch == '\\' {
			//handle windows path
			if pathSep == `\` {
				rexPat += escPath
				continue
			}
			if s.Peek() == scanner.EOF {
				rexPat += `\` + string(s.Next())
			} else {
				rexPat += `\`
			}
		} else {
			rexPat += string(ch)
		}
	}

	//end
	rexPat += "$"

	//compile regular expression
	//regx, err := regexp.Compile(rexPat)
	//if err != nil {
	//	return nil, err
	//}

	return rexPat, nil
}

// checks is the pattern is directory or not
// e.g
// pattern !README.md  will return false
// pattern target/     will return true
// pattern lib/*   will return true
// pattern target/one/two/one.txt  will return false
func isDir(pattern string) bool {
	re, err := regexp.Compile("/\\*?$")

	if err != nil {
		return false
	}

	return re.MatchString(pattern)
}
