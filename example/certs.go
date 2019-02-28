package main

import (
	"io/ioutil"

	"crypto/tls"

	"github.com/spf13/viper"
	"github.com/stevepartridge/go/file"
	"github.com/stevepartridge/service/insecure"
)

const (
	envTLSCert   = "host.local.crt"
	envTLSKey    = "host.local.key"
	envTLSRootCA = "Local_Development_Corp._LLC._Inc._CA.crt"
)

func getCert() ([]byte, error) {

	if viper.GetString(envTLSCert) != "" {

		if file.Exists(viper.GetString(envTLSCert)) {

			publicCert, err := ioutil.ReadFile(viper.GetString(envTLSCert))
			if err != nil {
				return nil, err
			}
			return []byte(publicCert), nil
		}

		return []byte(viper.GetString(envTLSCert)), nil
	}

	return []byte(insecure.Cert), nil
}

func getKey() ([]byte, error) {

	if viper.GetString(envTLSKey) != "" {
		if file.Exists(viper.GetString(envTLSKey)) {
			privateKey, err := ioutil.ReadFile(viper.GetString(envTLSKey))
			if err != nil {
				return nil, err
			}
			return []byte(privateKey), nil
		}

		return []byte(viper.GetString(envTLSKey)), nil

	}

	return []byte(insecure.Key), nil
}

func getCertificate() (tls.Certificate, error) {

	publicCert, err := getCert()
	if err != nil {
		return tls.Certificate{}, err
	}

	if publicCert == nil {
		publicCert = []byte{}
	}

	privateKey, err := getKey()
	if err != nil {
		return tls.Certificate{}, err
	}
	if privateKey == nil {
		privateKey = []byte{}
	}

	return tls.X509KeyPair(publicCert, privateKey)
}

func getRootCA() []byte {

	if viper.GetString(envTLSRootCA) != "" {
		if file.Exists(viper.GetString(envTLSRootCA)) {
			rootCA, err := ioutil.ReadFile(viper.GetString(envTLSRootCA))
			ifError(err)
			if rootCA != nil {
				return rootCA
			}
		}

		return []byte(viper.GetString(envTLSRootCA))
	}

	return []byte(insecure.RootCA)
}
