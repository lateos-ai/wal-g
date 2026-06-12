package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files under testdata/")

func goldenPath(name string) string {
	return filepath.Join("testdata", name+".golden")
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = orig })
	ch := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		buf.ReadFrom(r)
		ch <- buf.String()
	}()
	f()
	w.Close()
	return <-ch
}

func TestCharacterization_VersionValue(t *testing.T) {
	const want = "0.14.2-lts.1"
	if Version != want {
		t.Fatalf("Version = %q, want %q", Version, want)
	}
}

func TestCharacterization_MainOutput(t *testing.T) {
	tests := []struct {
		name       string
		goldenFile string
	}{
		{name: "default", goldenFile: "main_output"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, main)
			gp := goldenPath(tt.goldenFile)
			if *updateGolden {
				if err := os.MkdirAll(filepath.Dir(gp), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(gp, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(gp)
			if err != nil {
				t.Fatalf("read golden file %s: %v (run with -update to create)", gp, err)
			}
			if string(want) != got {
				t.Fatalf("stdout mismatch:\ngot:  %q\nwant: %q", got, string(want))
			}
		})
	}
}

func TestCharacterization_NoStderr(t *testing.T) {
	origOut := os.Stdout
	origErr := os.Stderr
	outR, outW, _ := os.Pipe()
	errR, errW, _ := os.Pipe()
	os.Stdout = outW
	os.Stderr = errW
	t.Cleanup(func() {
		os.Stdout = origOut
		os.Stderr = origErr
	})
	done := make(chan struct{}, 2)
	var outBuf, errBuf bytes.Buffer
	go func() { outBuf.ReadFrom(outR); done <- struct{}{} }()
	go func() { errBuf.ReadFrom(errR); done <- struct{}{} }()
	main()
	outW.Close()
	errW.Close()
	<-done
	<-done

	wantStdout := "WAL-G version " + Version + "\n"
	if outBuf.String() != wantStdout {
		t.Fatalf("stdout:\ngot:  %q\nwant: %q", outBuf.String(), wantStdout)
	}
	if stderr := errBuf.String(); stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestCharacterization_Repeatable(t *testing.T) {
	first := captureStdout(t, main)
	second := captureStdout(t, main)
	if first != second {
		t.Fatalf("calls differ:\nrun1: %q\nrun2: %q", first, second)
	}
}

func TestCharacterization_LDFlags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ldflags integration test in -short mode")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not found in PATH")
	}
	cmd := exec.Command("go", "run", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\noutput: %s", err, out)
	}
	want := "WAL-G version " + Version + "\n"
	if string(out) != want {
		t.Fatalf("stdout:\ngot:  %q\nwant: %q", string(out), want)
	}
}
