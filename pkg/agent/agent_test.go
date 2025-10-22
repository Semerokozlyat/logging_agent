package agent

import (
	"testing"
)

func TestNew(t *testing.T) {
	a := New()
	if a == nil {
		t.Fatal("Expected non-nil agent")
	}
	
	if a.config.LogLevel != "info" {
		t.Errorf("Expected LogLevel 'info', got '%s'", a.config.LogLevel)
	}
}

func TestRun(t *testing.T) {
	a := New()
	err := a.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
