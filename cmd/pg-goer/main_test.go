package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("version flag", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping version test in short mode")
		}
	})

	t.Run("help flag", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping help test in short mode")
		}
	})
}
