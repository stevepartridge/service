package service

import (
	"fmt"

	"testing"
)

var (
	testHost1 = "service.local"
	testPort1 = 1234
)

func Test_Unit_ServiceNew_ValidAddress(t *testing.T) {
	svc, err := New(testHost1, testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	expected := fmt.Sprintf("%s:%d", testHost1, testPort1)

	if svc.Addr() != expected {
		t.Errorf("Addr should be %s but saw %s", expected, svc.Addr())
	}

}

func Test_Unit_ServiceNew_InvalidHost(t *testing.T) {
	_, err := New("", testPort1)
	if err == nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if err != ErrInvalidHost {
		t.Errorf("Expected error %s but saw %s",
			ErrInvalidHost.Error(),
			err.Error(),
		)
	}

}

func Test_Unit_ServiceNew_InvalidPort(t *testing.T) {
	_, err := New(testHost1, 0)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	if err.Error() != ErrReplacer(ErrInvalidPort, 0).Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrReplacer(ErrInvalidPort, 0).Error(),
			err.Error(),
		)
	}

}

func Test_Unit_ServiceNew_InvalidNegativePort(t *testing.T) {
	_, err := New(testHost1, -123)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	if err.Error() != ErrReplacer(ErrInvalidPort, -123).Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrReplacer(ErrInvalidPort, -123).Error(),
			err.Error(),
		)
	}

}
