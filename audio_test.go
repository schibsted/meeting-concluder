package concluder

import (
	"fmt"
	"testing"
)

func TestInputDevices(t *testing.T) {
	a := NewAudioRecorder()
	defer a.Done()

	inputDevices, err := a.InputDevices()
	if err != nil {
		t.Fatalf("InputDevices() returned an error: %v", err)
	}
	if len(inputDevices) > 0 {
		fmt.Println("Found these input devices:")
	}
	for i, dev := range inputDevices {
		fmt.Printf("%d: %s\n", i+1, dev.Name)
	}
}
