// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dl

import (
	"sort"
	"strings"
	"testing"
)

func TestParseVersion(t *testing.T) {
	for _, c := range []struct {
		in       string
		maj, min int
		tail     string
	}{
		{"go1.5", 5, 0, ""},
		{"go1.5beta1", 5, 0, "beta1"},
		{"go1.5.1", 5, 1, ""},
		{"go1.5.1rc1", 5, 1, "rc1"},
	} {
		maj, min, tail := parseVersion(c.in)
		if maj != c.maj || min != c.min || tail != c.tail {
			t.Errorf("parseVersion(%q) = %v, %v, %q; want %v, %v, %q",
				c.in, maj, min, tail, c.maj, c.min, c.tail)
		}
	}
}

func TestFileOrder(t *testing.T) {
	fs := []File{
		{Filename: "go1.16.src.tar.gz", Version: "go1.16", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.16.1.src.tar.gz", Version: "go1.16.1", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.16.linux-amd64.tar.gz", Version: "go1.16", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.1.linux-amd64.tar.gz", Version: "go1.16.1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.darwin-amd64.tar.gz", Version: "go1.16", OS: "darwin", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.darwin-amd64.pkg", Version: "go1.16", OS: "darwin", Arch: "amd64", Kind: "installer"},
		{Filename: "go1.16.darwin-arm64.tar.gz", Version: "go1.16", OS: "darwin", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.16.darwin-arm64.pkg", Version: "go1.16", OS: "darwin", Arch: "arm64", Kind: "installer"},
		{Filename: "go1.16beta1.linux-amd64.tar.gz", Version: "go1.16beta1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16beta2.linux-amd64.tar.gz", Version: "go1.16beta2", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16rc1.linux-amd64.tar.gz", Version: "go1.16rc1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.15.linux-amd64.tar.gz", Version: "go1.15", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.15.2.linux-amd64.tar.gz", Version: "go1.15.2", OS: "linux", Arch: "amd64", Kind: "archive"},
	}
	sort.Sort(fileOrder(fs))
	var s []string
	for _, f := range fs {
		s = append(s, f.Filename)
	}
	got := strings.Join(s, "\n")
	want := strings.Join([]string{
		"go1.16.1.src.tar.gz",
		"go1.16.1.linux-amd64.tar.gz",
		"go1.16.src.tar.gz",
		"go1.16.darwin-amd64.tar.gz",
		"go1.16.darwin-amd64.pkg",
		"go1.16.darwin-arm64.tar.gz",
		"go1.16.darwin-arm64.pkg",
		"go1.16.linux-amd64.tar.gz",
		"go1.15.2.linux-amd64.tar.gz",
		"go1.15.linux-amd64.tar.gz",
		"go1.16rc1.linux-amd64.tar.gz",
		"go1.16beta2.linux-amd64.tar.gz",
		"go1.16beta1.linux-amd64.tar.gz",
	}, "\n")
	if got != want {
		t.Errorf("sort order is\n%s\nwant:\n%s", got, want)
	}
}

func TestFilesToReleases(t *testing.T) {
	fs := []File{
		{Version: "go1.7.4", OS: "darwin"},
		{Version: "go1.7.4", OS: "windows"},
		{Version: "go1.7", OS: "darwin"},
		{Version: "go1.7", OS: "windows"},
		{Version: "go1.6.2", OS: "darwin"},
		{Version: "go1.6.2", OS: "windows"},
		{Version: "go1.6", OS: "darwin"},
		{Version: "go1.6", OS: "windows"},
		{Version: "go1.5.2", OS: "darwin"},
		{Version: "go1.5.2", OS: "windows"},
		{Version: "go1.5", OS: "darwin"},
		{Version: "go1.5", OS: "windows"},
		{Version: "go1.5beta1", OS: "windows"},
	}
	stable, unstable, archive := filesToReleases(fs)
	if got, want := list(stable), "go1.7.4, go1.6.2"; got != want {
		t.Errorf("stable = %q; want %q", got, want)
	}
	if got, want := list(unstable), ""; got != want {
		t.Errorf("unstable = %q; want %q", got, want)
	}
	if got, want := list(archive), "go1.7, go1.6, go1.5.2, go1.5, go1.5beta1"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

func TestHighlightedFiles(t *testing.T) {
	fs := []File{
		{Filename: "go1.16beta1.src.tar.gz", Version: "go1.16beta1", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.16beta1.linux-386.tar.gz", Version: "go1.16beta1", OS: "linux", Arch: "386", Kind: "archive"},
		{Filename: "go1.16beta1.linux-amd64.tar.gz", Version: "go1.16beta1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16beta1.darwin-amd64.tar.gz", Version: "go1.16beta1", OS: "darwin", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16beta1.darwin-amd64.pkg", Version: "go1.16beta1", OS: "darwin", Arch: "amd64", Kind: "installer"},
		{Filename: "go1.16beta1.darwin-arm64.tar.gz", Version: "go1.16beta1", OS: "darwin", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.16beta1.darwin-arm64.pkg", Version: "go1.16beta1", OS: "darwin", Arch: "arm64", Kind: "installer"},
		{Filename: "go1.16beta1.windows-386.zip", Version: "go1.16beta1", OS: "windows", Arch: "386", Kind: "archive"},
		{Filename: "go1.16beta1.windows-386.msi", Version: "go1.16beta1", OS: "windows", Arch: "386", Kind: "installer"},
		{Filename: "go1.16beta1.windows-amd64.zip", Version: "go1.16beta1", OS: "windows", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16beta1.windows-amd64.msi", Version: "go1.16beta1", OS: "windows", Arch: "amd64", Kind: "installer"},
	}
	sort.Sort(fileOrder(fs))
	var highlighted []string
	for _, f := range fs {
		if !f.Highlight() {
			continue
		}
		highlighted = append(highlighted, f.Filename)
	}
	got := strings.Join(highlighted, "\n")
	want := strings.Join([]string{
		"go1.16beta1.src.tar.gz",
		"go1.16beta1.darwin-amd64.pkg",
		"go1.16beta1.darwin-arm64.pkg",
		"go1.16beta1.linux-amd64.tar.gz",
		"go1.16beta1.windows-amd64.msi",
	}, "\n")
	if got != want {
		t.Errorf("highlighted files:\n%s\nwant:\n%s", got, want)
	}
}

func TestOldUnstableNotShown(t *testing.T) {
	fs := []File{
		{Version: "go1.7.4"},
		{Version: "go1.7"},
		{Version: "go1.7beta1"},
	}
	_, unstable, archive := filesToReleases(fs)
	if len(unstable) != 0 {
		t.Errorf("got unstable, want none")
	}
	if got, want := list(archive), "go1.7, go1.7beta1"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

// A new beta should show up under unstable, but not show up under archive. See golang.org/issue/29669.
func TestNewUnstableShownOnce(t *testing.T) {
	fs := []File{
		{Version: "go1.12beta2"},
		{Version: "go1.11.4"},
		{Version: "go1.11"},
		{Version: "go1.10.7"},
		{Version: "go1.10"},
		{Version: "go1.9"},
	}
	stable, unstable, archive := filesToReleases(fs)
	if got, want := list(stable), "go1.11.4, go1.10.7"; got != want {
		t.Errorf("stable = %q; want %q", got, want)
	}
	if got, want := list(unstable), "go1.12beta2"; got != want {
		t.Errorf("unstable = %q; want %q", got, want)
	}
	if got, want := list(archive), "go1.11, go1.10, go1.9"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

func TestUnstableShown(t *testing.T) {
	fs := []File{
		{Version: "go1.8beta2"},
		{Version: "go1.8rc1"},
		{Version: "go1.7.4"},
		{Version: "go1.7"},
		{Version: "go1.7beta1"},
	}
	_, unstable, archive := filesToReleases(fs)
	// Show RCs ahead of betas.
	if got, want := list(unstable), "go1.8rc1"; got != want {
		t.Errorf("unstable = %q; want %q", got, want)
	}
	if got, want := list(archive), "go1.7, go1.8beta2, go1.7beta1"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

// list returns a version list string for the given releases.
func list(rs []Release) string {
	var s string
	for i, r := range rs {
		if i > 0 {
			s += ", "
		}
		s += r.Version
	}
	return s
}
