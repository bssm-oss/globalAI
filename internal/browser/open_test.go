package browser

import "testing"

func TestCommandForCurrentPlatform(t *testing.T) {
	command, args, err := commandFor("http://127.0.0.1:4000")
	if err != nil {
		t.Fatalf("commandFor() error = %v", err)
	}
	if command == "" {
		t.Fatal("expected command to be populated")
	}
	if len(args) == 0 {
		t.Fatal("expected args to be populated")
	}
}
