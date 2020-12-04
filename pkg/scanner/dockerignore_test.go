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
	"os"
	"path/filepath"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
)

func TestDirectoryHasDockerIgnore(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd() =", err)
	}

	dir := filepath.Join(wd, "testdata", "dir2")

	igScanner, err := NewOrDefault(dir)

	if err != nil {
		t.Error("hasDockerIgnore() = ", err)
	}

	_, err = igScanner.Scan()

	if err != nil {
		t.Error("igScanner.Scan()= ", err)
	}

	if len(IgnorePatterns) == 0 {
		t.Errorf("The directory %s has '.dockerignore', but got it does not", dir)
	}

	if got, want := len(IgnorePatterns), int(9); got != want {
		t.Errorf("Patterns() = %d, wanted %d", got, want)
	}
}

func TestDirectoryHasNoDockerIgnore(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd() =", err)
	}
	dir := filepath.Join(wd, "testdata", "dir1")

	igScanner, err := NewOrDefault(dir)

	if err != nil {
		t.Error("hasDockerIgnore() = ", err)
	}

	_, err = igScanner.Scan()

	if err != nil {
		t.Error("igScanner.Scan()= ", err)
	}

	if got, want := len(IgnorePatterns), int(3); got != want {
		t.Errorf("Patterns() = %d, wanted %d", got, want)
	}

}

func TestDockerIgnoredPatterns(t *testing.T) {
	expected := sets.NewString("lib", "*.md", "!README.md", "temp?", "target", "!target/*-runner.jar")
	expected.Insert(defaultPatterns...)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd() =", err)
	}

	dir := filepath.Join(wd, "testdata", "dir2")

	igScanner, err := NewOrDefault(dir)

	if err != nil {
		t.Error("ignorablePatterns() = ", err)
	}

	_, err = igScanner.Scan()

	if err != nil {
		t.Error("igScanner.Scan()= ", err)
	}

	if got, want := len(IgnorePatterns), int(9); got != want {
		t.Errorf("Patterns() = %d, wanted %d", got, want)
	}

}

func TestEmptyDockerIgnoredPatterns(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd() =", err)
	}

	dir := filepath.Join(wd, "testdata", "empty")

	igScanner, err := NewOrDefault(dir)

	if err != nil {
		t.Error("ignorablePatterns() = ", err)
	}

	_, err = igScanner.Scan()

	if err != nil {
		t.Error("igScanner.Scan()= ", err)
	}

	if got, want := len(IgnorePatterns), int(3); got != want {
		t.Errorf("Patterns() = %d, wanted %d", got, want)
	}

}

func TestIgnoreables(t *testing.T) {

	eIncludes := sets.NewString([]string{
		"/Users/kameshs/git/kameshsampath/goignorescanner/pkg/scanner/testdata/starignore/README.md",
		"/Users/kameshs/git/kameshsampath/goignorescanner/pkg/scanner/testdata/starignore/target/foo-runner.jar",
		"/Users/kameshs/git/kameshsampath/goignorescanner/pkg/scanner/testdata/starignore/target/lib/one.jar",
		"/Users/kameshs/git/kameshsampath/goignorescanner/pkg/scanner/testdata/starignore/target/quarkus-app/one.txt",
	}...)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd() =", err)
	}

	dir := filepath.Join(wd, "testdata", "starignore")

	ignoreScanner, err := NewOrDefault(dir)

	if err != nil {
		t.Error("isIgnorable() = ", err)
	}

	incls, err := ignoreScanner.Scan()

	sIncludes := sets.NewString(incls...)

	if err != nil {
		t.Error("ignoreScanner.Scan() = ", err)
	}

	if got, want := len(incls), int(4); got != want {
		t.Errorf("Patterns() = %d, wanted %d", got, want)
	}

	if !eIncludes.Equal(sIncludes) {
		t.Errorf("Includes() = %s, wanted %s", sIncludes, eIncludes)
	}

}
