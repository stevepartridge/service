package service

import (
	"os"
	"strconv"
	"testing"
)

func TestNewService(t *testing.T) {

	t.Run("With default settings", func(t *testing.T) {

		s, err := New()
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if s.Port != defaultPort {
			t.Errorf("Expected Port to be %d but saw %d", defaultPort, s.Port)
		}

	})

	t.Run("With env port", func(t *testing.T) {

		port := 1234
		os.Setenv("PORT", strconv.Itoa(port))

		s, err := New()
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if s.Port != port {
			t.Errorf("Expected Port to be %d but saw %d", port, s.Port)
		}

	})
}
