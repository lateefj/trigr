package main

import (
	"testing"

	"github.com/lateefj/trigr"
)

func TestProjectManager(t *testing.T) {
	pm := NewProjectManager()
	p := NewProject("test")
	p.AssignLocalSource(trigr.Git, "/tmp/test_trigr")
	err := pm.Add(p)
	if err != nil {
		t.Fatal(err)
	}
	err = pm.Add(p)
	if err == nil {
		t.Fatal("Duplicate project add should have failed")
	}

	np, err := pm.Get(p.Id)
	if err != nil {
		t.Fatal(err)
	}
	if np == nil {
		t.Fatal("Failed to Get a project based on the id")
	}

	err = pm.Remove(p.Id)
	if err != nil {
		t.Fatal(err)
	}

	tp, err := pm.Get(p.Id)
	if tp != nil {
		t.Fatal("Expected nil for removed project")
	}
	if err == nil {
		t.Fatal("Expected error trying to get project that already exists")
	}
}
