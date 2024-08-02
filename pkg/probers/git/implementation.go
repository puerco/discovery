// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openvex/discovery/pkg/discovery/options"
	"github.com/openvex/go-vex/pkg/vex"
	purl "github.com/package-url/packageurl-go"
	"sigs.k8s.io/release-sdk/git"
	"sigs.k8s.io/release-utils/util"
)

type gitImplementation interface {
	CloneOrUpdateRepository(*options.Options) error
	FindVexDocuments(*options.Options) ([]*vex.VEX, error)
	FilterVexDocuments(*options.Options, []*vex.VEX) ([]*vex.VEX, error)
}

type defaultImplementation struct{}

func (di *defaultImplementation) LocalOptions(opts *options.Options) Options {
	if _, ok := opts.ProberOptions[Kind]; !ok {
		return Options{}
	}
	if lo, ok := opts.ProberOptions[Kind].(Options); ok {
		return lo
	}
	return Options{}
}

// CloneOrUpdateRepository performs a shallow clone of the repository
// specified in the prober options. It will be cloned to the specified path
// or to a temporary directory f none is set.
func (di *defaultImplementation) CloneOrUpdateRepository(opts *options.Options) error {
	if opts.ProberOptions[Kind].(Options).ClonePath == "" {
		lo := opts.ProberOptions[Kind].(Options)
		path, err := os.MkdirTemp("", "vex-discovery-tmp-*")
		if err != nil {
			return fmt.Errorf("creating tem dir: %w", err)
		}
		lo.ClonePath = path
		opts.ProberOptions[Kind] = lo
	}

	if _, err := git.ShallowCloneOrOpenRepo(
		opts.ProberOptions[Kind].(Options).ClonePath,
		opts.ProberOptions[Kind].(Options).RepositoryURL,
		strings.HasPrefix("https://", opts.ProberOptions[Kind].(Options).RepositoryURL),
	); err != nil {
		return fmt.Errorf("cloning repository: %w", err)
	}
	return nil
}

// FindVexDocuments will check the well known VEX paths for VEX documents.
func (di *defaultImplementation) FindVexDocuments(opts *options.Options, p purl.PackageURL) ([]*vex.VEX, error) {
	filePaths := []string{}
	docs := []*vex.VEX{}
	path := opts.ProberOptions[Kind].(Options).ClonePath

	tryPaths := []string{
		filepath.Join(path, "openvex.json"),
		filepath.Join(path, ".openvex.json"),
	}

	// If there is a main VEX document, we always return i
	for _, p := range tryPaths {
		if util.Exists(filepath.Join(path, "openvex.json")) {
			filePaths = append(filePaths, p)
		}
	}

	if len(filePaths) == 0 {
		return nil, nil
	}

	for _, file := range filePaths {
		doc, err := vex.Open(file)
		if err != nil {
			return nil, fmt.Errorf("error parsing file %q: %w", strings.TrimPrefix(file, path), err)
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

func (di *defaultImplementation) FilterVexDocuments(opts *options.Options, docs []*vex.VEX, p purl.PackageURL) ([]*vex.VEX, error) {
	filteredDocs := []*vex.VEX{}
	for _, doc := range docs {
		for _, s := range doc.Statements {
			if s.MatchesProduct(p.String()) {
				filteredDocs = append(filteredDocs, doc)
			}
		}
	}
	return docs, nil
}
