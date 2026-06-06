package main

import (
"testing"
)

func TestVersion(t *testing.T) {
if Version == "" {
t.Fatal("Version not set")
}
}
