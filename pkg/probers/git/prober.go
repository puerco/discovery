// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"fmt"

	"github.com/openvex/discovery/pkg/discovery/options"
	"github.com/openvex/go-vex/pkg/vex"
	purl "github.com/package-url/packageurl-go"
)

const Kind = "git"

type Prober struct {
	Options options.Options
	impl    gitImplementation
}

type Options struct {
	// RepositoryURL is the URL of the repo.
	RepositoryURL string
	// Clone path is a directory where the repository will be clones
	ClonePath string
}

func New() *Prober {
	p := &Prober{
		impl:    &defaultImplementation{},
		Options: options.Default,
	}
	p.Options.ProberOptions[purl.TypeOCI] = Options{}
	return p
}

func (p *Prober) FindDocumentsFromPurl(options.Options, purl.PackageURL) ([]*vex.VEX, error) {
	// Clone the repository
	if err := p.impl.CloneOrUpdateRepository(&p.Options); err != nil {
		return nil, fmt.Errorf("cloning repository: %w", err)
	}

	// Read the clone repo and find any applicable documents
	docs, err := p.impl.FindVexDocuments(&p.Options)
	if err != nil {
		return nil, fmt.Errorf("finding VEX documents: %w", err)
	}

	filteredDocs, err := p.impl.FilterVexDocuments(&p.Options, docs)
	if err != nil {
		return nil, fmt.Errorf("filtering documents: %w", err)
	}

	return filteredDocs, nil
}

func (p *Prober) SetOptions(options.Options) {}
