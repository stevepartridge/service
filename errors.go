package service

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidHost error message
	ErrInvalidHost = errors.New("Host wasn't provided or is invalid")

	// ErrInvalidPort error message
	ErrInvalidPort = errors.New("Invalid port %d")

	// ErrAddingRootCA error message
	ErrAddingRootCA = errors.New("Unable to add Root CA to cert pool, invalid cert.")

	// ErrMissingPublicCert error message
	ErrMissingPublicCert = errors.New("Missing public cert. Fix by using: service.AddKeyPair or server.AddKeyPairFromFiles")

	// ErrMissingPrivateKey error message
	ErrMissingPrivateKey = errors.New("Missing private key. Fix by using: service.AddKeyPair or server.AddKeyPairFromFiles")

	// ErrAddKeyPairPublicCertIsNil error message
	ErrAddKeyPairPublicCertIsNil = errors.New("Unable to add key pair, public cert is nil")

	// ErrAddKeyPairPrivateKeyIsNil error message
	ErrAddKeyPairPrivateKeyIsNil = errors.New("Unable to add key pair, private key is nil")

	// ErrAddKeyPairFromFilePublicCertEmpty error message
	ErrAddKeyPairFromFilePublicCertEmpty = errors.New("Unable to add key pair from files, public cert is empty")

	// ErrAddKeyPairFromFilePublicCertNotFound error message
	ErrAddKeyPairFromFilePublicCertNotFound = errors.New("Unable to add key pair, public cert not found at %s")

	// ErrAddKeyPairPrivateKeyEmpty error message
	ErrAddKeyPairPrivateKeyEmpty = errors.New("Unable to add key pair, private key is nil")

	// ErrAddKeyPairFromFilePrivateKeyNotFound error message
	ErrAddKeyPairFromFilePrivateKeyNotFound = errors.New("Unable to add key pair, private key not found at %s")

	// ErrGatewayHandlerIsNil error message
	ErrGatewayHandlerIsNil = errors.New("Gateway entrypoint handler nil")
)

// ErrReplacer allows for static errors to have dynamic values
func ErrReplacer(err error, replacers ...interface{}) error {
	return errors.New(fmt.Errorf(err.Error(), replacers...))
}
