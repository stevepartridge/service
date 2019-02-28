package service

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/stevepartridge/go/file"
)

// AppendCertsFromPEM support direct setting of rootCA
func (s *Service) AppendCertsFromPEM(rootCA []byte) error {
	ok := s.CertPool.AppendCertsFromPEM([]byte(rootCA))
	if !ok {
		return ErrAddingRootCA
	}
	return nil
}

// AddKeyPair allows for direct setting of cert/keys
func (s *Service) AddKeyPair(publicCert, privateKey []byte) error {
	if publicCert == nil {
		return ErrAddKeyPairPublicCertIsNil
	}

	if privateKey == nil {
		return ErrAddKeyPairPublicCertIsNil
	}

	s.PublicCert = publicCert
	s.PrivateKey = privateKey

	return nil
}

// AddKeyPairFromFiles allows for setting cert/key from files
func (s *Service) AddKeyPairFromFiles(publicCertFile, privateKeyFile string) error {

	if !file.Exists(publicCertFile) {
		return ErrReplacer(ErrAddKeyPairFromFilePublicCertNotFound, publicCertFile)
	}

	if !file.Exists(privateKeyFile) {
		return ErrReplacer(ErrAddKeyPairFromFilePrivateKeyNotFound, privateKeyFile)
	}

	publicCert, err := ioutil.ReadFile(publicCertFile)
	if err != nil {
		return err
	}

	privateKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}

	return s.AddKeyPair(publicCert, privateKey)
}

// GetCertificate is a helper for pulling tls.Certificate from provided cert/key
func (s *Service) GetCertificate() (tls.Certificate, error) {

	if s.PublicCert == nil {
		return tls.Certificate{}, ErrMissingPublicCert
	}
	if s.PrivateKey == nil {
		return tls.Certificate{}, ErrMissingPrivateKey
	}

	return tls.X509KeyPair(s.PublicCert, s.PrivateKey)

}

// EnableInsecure is a setter for enabling the insecure flag
func (s *Service) EnableInsecure() {
	s.enableInsecure = true
	fmt.Println(" ! ! * Enabled insecure * ! !")
}
