// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package mage

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func testPackageSpec() PackageSpec {
	return PackageSpec{
		Name:     "brewbeat",
		Version:  "7.0.0",
		Snapshot: true,
		OS:       "windows",
		Arch:     "x86_64",
		ExtraTags: []string{
			"git-{{ substring commit 0 12 }}",
		},
		Files: map[string]PackageFile{
			"brewbeat.yml": PackageFile{
				Source: "./testdata/config.yml",
				Mode:   0644,
			},
			"README.txt": PackageFile{
				Content: "Hello! {{.Version}}\n",
				Mode:    0644,
			},
		},
	}
}

func TestPackageZip(t *testing.T) {
	testPackage(t, PackageZip)
}

func TestPackageTarGz(t *testing.T) {
	testPackage(t, PackageTarGz)
}

func TestPackageRPM(t *testing.T) {
	t.Skip("Flaky test")
	if err := HaveDocker(); err != nil {
		t.Skip("docker is required")
	}

	testPackage(t, PackageRPM)
}

func TestPackageDeb(t *testing.T) {
	t.Skip("Flaky test")
	if err := HaveDocker(); err != nil {
		t.Skip("docker is required")
	}

	testPackage(t, PackageDeb)
}

func testPackage(t testing.TB, pack func(PackageSpec) error) {
	spec := testPackageSpec().Evaluate()

	readme := spec.Files["README.txt"]
	readmePath := filepath.ToSlash(filepath.Clean(readme.Source))
	assert.True(t, strings.HasPrefix(readmePath, packageStagingDir))

	commit := spec.ExtraTags[0]
	expected := "git-" + commitHash[:12]
	assert.Equal(t, expected, commit)

	if err := pack(spec); err != nil {
		t.Fatal(err)
	}
}

func TestRepoRoot(t *testing.T) {
	repo, err := GetProjectRepoInfo()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "github.com/elastic/elastic-agent", repo.RootImportPath)
	assert.True(t, filepath.IsAbs(repo.RootDir))
	cwd := filepath.Join(repo.RootDir, repo.SubDir)
	assert.Equal(t, CWD(), cwd)
}

func TestDumpVariables(t *testing.T) {
	out, err := dumpVariables()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
}

func TestLoadSpecs(t *testing.T) {
	pkgs, err := LoadSpecs("../packaging/packages.yml")
	if err != nil {
		t.Fatal(err)
	}

	for flavor, s := range pkgs {
		out, err := yaml.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		if testing.Verbose() {
			t.Log("Packaging flavor:", flavor, "\n", string(out))
		}
	}
}
