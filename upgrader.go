// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package pubspecmgr

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type Upgrader struct {
	ch chan *Package
	wg *sync.WaitGroup // Add WaitGroup to wait for tasks to complete
	mu sync.Mutex      // Add mutex to protect AST modification operations

	parsed *Parsed

	constraint bool
	workers    int
}

type UpgraderOption func(*Upgrader)

// WithConstraint Sets constraint mode
func WithConstraint(constraint bool) UpgraderOption {
	return func(p *Upgrader) {
		p.constraint = constraint
	}
}

// WithWorkers Sets number of workers
func WithWorkers(workers int) UpgraderOption {
	return func(p *Upgrader) {
		p.workers = workers
	}
}

// NewUpgrader Create a new Upgrader
func NewUpgrader(parsed *Parsed, opts ...UpgraderOption) *Upgrader {
	u := &Upgrader{
		wg: &sync.WaitGroup{},
		mu: sync.Mutex{},

		parsed: parsed,

		constraint: false,
		workers:    20,
	}

	for _, opt := range opts {
		opt(u)
	}

	u.ch = make(chan *Package, u.workers)

	u.wg.Add(u.workers)
	for i := 0; i < u.workers; i++ {
		go u.worker()
	}

	return u
}

// Run Starts processing the packages
func (u *Upgrader) Run() {
	for _, pkg := range u.parsed.Packages {
		u.ch <- pkg
	}

	u.wait()
}

// wait Waits for all tasks to complete
func (u *Upgrader) wait() {
	// Close the channel to signal no more tasks will be sent
	close(u.ch)

	// Wait for all goroutines to finish
	u.wg.Wait()
}

// Worker goroutine to process nodes from the channel
func (u *Upgrader) worker() {
	defer func() {
		// When doUpdate ends, mark the WaitGroup as done
		u.wg.Done()
	}()

	for pkg := range u.ch {
		latest := pkg.GetLatest(u.constraint)
		if latest != "" {
			// Using mutex to protect AST modification to avoid concurrent modification conflicts
			u.mu.Lock()
			tk := pkg.node.GetToken()
			tk.Value = latest
			p, _ := yaml.PathString(pkg.node.GetPath())
			err := p.ReplaceWithNode(u.parsed.File, ast.String(tk))
			u.mu.Unlock()

			if err != nil {
				log.Printf("failed to replace package %s version %s with %s", pkg.Name, pkg.Version, latest)
			}
		}
	}
}
