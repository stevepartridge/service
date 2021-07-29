package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stevepartridge/service/insecure"
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

	if !fileExists(publicCertFile) {
		return fmt.Errorf(ErrAddKeyPairFromFilePublicCertNotFound.Error(), publicCertFile)
	}

	if !fileExists(privateKeyFile) {
		return fmt.Errorf(ErrAddKeyPairFromFilePrivateKeyNotFound.Error(), privateKeyFile)
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

	if s.PublicCert == nil && s.PrivateKey == nil {
		fmt.Println("TLS Certificate is not set, using self-signed")
		s.AddKeyPair([]byte(insecure.Cert), []byte(insecure.Key))
	}

	if s.PublicCert == nil {
		return tls.Certificate{}, ErrMissingPublicCert
	}
	if s.PrivateKey == nil {
		return tls.Certificate{}, ErrMissingPrivateKey
	}

	return tls.X509KeyPair(s.PublicCert, s.PrivateKey)

}

func (s *Service) insecureCert() *tls.Certificate {

	cert, err := s.GetCertificate()
	if err != nil {
		fmt.Println("unable to get insecure certificates: ", err.Error())
		return nil
	}

	if !s.CertPool.AppendCertsFromPEM([]byte(insecure.RootCA)) {
		fmt.Printf("failed to append insecure certificate to certificate pool: %s\n", err)
		return nil
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		fmt.Printf("failed to parse certificate: %s\n", err)
		return nil
	}

	s.CertPool = x509.NewCertPool()
	s.CertPool.AddCert(cert.Leaf)

	return &cert

}

func WithInsecureSkipVerify() func(*Service) error {
	return func(s *Service) error {
		s.insecureSkipVerify = true
		return nil
	}
}

func WithTLSDisabled() func(*Service) error {
	return func(s *Service) error {
		s.disableTLSCerts = true
		return nil
	}
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
