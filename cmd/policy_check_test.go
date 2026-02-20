package cmd

import (
	"testing"
)

func TestPolicyCheckAttrsFields(t *testing.T) {
	a := policyCheckAttrs{
		Status: "soft_failed",
		Scope:  "organization",
		Actions: map[string]bool{
			"is-overridable": true,
		},
	}
	if a.Status != "soft_failed" {
		t.Errorf("expected status soft_failed, got %s", a.Status)
	}
	if a.Scope != "organization" {
		t.Errorf("expected scope organization, got %s", a.Scope)
	}
	if !a.Actions["is-overridable"] {
		t.Error("expected is-overridable to be true")
	}
}

func TestPolicyCheckCmdRegistered(t *testing.T) {
	// Verify subcommands are registered
	found := map[string]bool{}
	for _, sub := range policyCheckCmd.Commands() {
		found[sub.Name()] = true
	}
	for _, name := range []string{"list", "show", "override"} {
		if !found[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestPolicyCheckListRequiresRunFlag(t *testing.T) {
	// Reset flag
	flagPCRunID = ""
	// runPolicyCheckList should fail without --run
	err := runPolicyCheckList(policyCheckListCmd, nil)
	if err == nil {
		t.Fatal("expected error when --run not set")
	}
}

func TestPolicyCheckOverrideRequiresArg(t *testing.T) {
	// cobra.ExactArgs(1) should reject 0 args
	if policyCheckOverrideCmd.Args == nil {
		t.Fatal("expected Args validator on override command")
	}
}

func TestPolicyCheckShowRequiresArg(t *testing.T) {
	if policyCheckShowCmd.Args == nil {
		t.Fatal("expected Args validator on show command")
	}
}
