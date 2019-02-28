package service

import (
	"testing"

	"github.com/stevepartridge/service/insecure"
)

func Test_Unit_ServiceAddRootCA_Invalid(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}
}

func Test_Unit_ServiceAddKeyPairFromFiles_Success(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}

	err = svc.AddKeyPairFromFiles(
		"certificates/out/host.local.crt",
		"certificates/out/host.local.key",
	)

	svc.EnableInsecure()
	if !svc.enableInsecure {
		t.Error("Enable insecure should be true")
	}

	if err != nil {
		t.Errorf("Error adding key pair from files %s", err.Error())
	}
}

func Test_Unit_ServiceAddKeyPairFromFiles_MissingCert(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}

	err = svc.AddKeyPairFromFiles(
		"bogus.crt",
		"certificates/out/host.local.key",
	)

	if err == nil {
		t.Error("Expected error, but didn't receive one.")
	}

	expectedErrMsg := ErrReplacer(ErrAddKeyPairFromFilePublicCertNotFound, "bogus.crt")
	if err.Error() != expectedErrMsg.Error() {
		t.Errorf("Error adding key pair from files %s", err.Error())
	}
}

func Test_Unit_ServiceAddKeyPairFromFiles_MissingKey(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}

	err = svc.AddKeyPairFromFiles(
		"certificates/out/host.local.crt",
		"bogus.key",
	)

	if err == nil {
		t.Error("Expected error, but didn't receive one.")
	}

	expectedErrMsg := ErrReplacer(ErrAddKeyPairFromFilePrivateKeyNotFound, "bogus.key")
	if err.Error() != expectedErrMsg.Error() {
		t.Errorf("Expected error message %s but saw %s",
			expectedErrMsg.Error(),
			err.Error(),
		)
	}
}

func Test_Unit_ServiceGetCertificates_CertNil(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}

	err = svc.AddKeyPair(nil, []byte(insecure.Key))
	if err == nil {
		t.Error("Expected error but didn't receive one.")
	}

	if err.Error() != ErrAddKeyPairPublicCertIsNil.Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrAddKeyPairPublicCertIsNil.Error(),
			err.Error())
	}

}

func Test_Unit_ServiceGetCertificates_IsNil(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc == nil {
		t.Error("Service should not be nil")
	}

	_, err = svc.GetCertificate()

	if err.Error() != ErrMissingPublicCert.Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrMissingPublicCert.Error(),
			err.Error())
	}

}
