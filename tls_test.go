package service

import (
	"testing"
)

func Test_Unit_ServiceAddRootCA_Invalid(t *testing.T) {
	svc, err := New(testHost1, testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}
}
