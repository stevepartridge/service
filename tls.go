package service

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/stevepartridge/go/file"
)

func (s *Service) AppendCertsFromPEM(rootCA []byte) error {
	ok := s.CertPool.AppendCertsFromPEM([]byte(rootCA))
	if !ok {
		return ErrAddingRootCA
	}
	return nil
}

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

func (s *Service) GetCertificate() (tls.Certificate, error) {

	if s.PublicCert == nil {
		return tls.Certificate{}, ErrMissingPublicCert
	}
	if s.PrivateKey == nil {
		return tls.Certificate{}, ErrMissingPrivateKey
	}

	return tls.X509KeyPair(s.PublicCert, s.PrivateKey)

}

func (s *Service) EnableInsecure() {
	s.enableInsecure = true
	fmt.Println(" ! ! * Enabled insecure * ! !")
}
