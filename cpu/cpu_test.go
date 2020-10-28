// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu_test

import (
	. "github.com/qugust/nibiru-internal/cpu"
	"github.com/qugust/nibiru-internal/testenv"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestMinimalFeatures(t *testing.T) {
	// TODO: maybe do MustSupportFeatureDectection(t) ?
	if runtime.GOARCH == "arm64" {
		switch runtime.GOOS {
		case "linux", "android":
		default:
			t.Skipf("%s/%s is not supported", runtime.GOOS, runtime.GOARCH)
		}
	}

	for _, o := range Options {
		if o.Required && !*o.Feature {
			t.Errorf("%v expected true, got false", o.Name)
		}
	}
}

func MustHaveDebugOptionsSupport(t *testing.T) {
	if !DebugOptions {
		t.Skipf("skipping test: cpu feature options not supported by OS")
	}
}

func MustSupportFeatureDectection(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skipf("CPU feature detection is not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	// TODO: maybe there are other platforms?
}

func runDebugOptionsTest(t *testing.T, test string, options string) {
	MustHaveDebugOptionsSupport(t)

	testenv.MustHaveExec(t)

	env := "GODEBUG=" + options

	cmd := exec.Command(os.Args[0], "-test.run="+test)
	cmd.Env = append(cmd.Env, env)

	output, err := cmd.CombinedOutput()
	lines := strings.Fields(string(output))
	lastline := lines[len(lines)-1]

	got := strings.TrimSpace(lastline)
	want := "PASS"
	if err != nil || got != want {
		t.Fatalf("%s with %s: want %s, got %v", test, env, want, got)
	}
}

func TestDisableAllCapabilities(t *testing.T) {
	MustSupportFeatureDectection(t)
	runDebugOptionsTest(t, "TestAllCapabilitiesDisabled", "cpu.all=off")
}

func TestAllCapabilitiesDisabled(t *testing.T) {
	MustHaveDebugOptionsSupport(t)

	if os.Getenv("GODEBUG") != "cpu.all=off" {
		t.Skipf("skipping test: GODEBUG=cpu.all=off not set")
	}

	for _, o := range Options {
		want := o.Required
		if got := *o.Feature; got != want {
			t.Errorf("%v: expected %v, got %v", o.Name, want, got)
		}
	}
}
