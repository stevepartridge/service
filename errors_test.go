package service

import (
	"testing"

	"fmt"
)

func Test_Unit_ErrReplacerSuccess(t *testing.T) {

	expectedStr := fmt.Sprintf(ErrInvalidPort.Error(), 1234)
	err := ErrReplacer(ErrInvalidPort, 1234)
	if err.Error() != expectedStr {
		t.Errorf("Expected error to be %s but saw %s",
			expectedStr,
			err.Error(),
		)
	}
}
