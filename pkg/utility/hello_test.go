package utility

import (
	"testing"
)

// TestHelloWorld
func TestHelloWorld(t *testing.T) {
	HelloWorld()
}

func TestHello(t *testing.T) {
	message := Hello("Test")

	if message != "Hello Test" {
		t.Errorf("The message is wrong, got: %s, want: %s.", message, "Hello Test")
	}
}
