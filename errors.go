package service

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidHost = errors.New("Host wasn't provided or is invalid")
	ErrInvalidPort = errors.New("Invalid port %d")

	ErrAddingRootCA      = errors.New("Unable to add Root CA to cert pool, invalid cert.")
	ErrMissingPublicCert = errors.New("Missing public cert. Fix by using: service.AddKeyPair or server.AddKeyPairFromFiles")
	ErrMissingPrivateKey = errors.New("Missing private key. Fix by using: service.AddKeyPair or server.AddKeyPairFromFiles")

	ErrAddKeyPairPublicCertIsNil = errors.New("Unable to add key pair, public cert is nil")
	ErrAddKeyPairPrivateKeyIsNil = errors.New("Unable to add key pair, private key is nil")

	ErrAddKeyPairFromFilePublicCertEmpty    = errors.New("Unable to add key pair from files, public cert is empty")
	ErrAddKeyPairFromFilePublicCertNotFound = errors.New("Unable to add key pair, public cert not found at %s")
	ErrAddKeyPairPrivateKeyEmpty            = errors.New("Unable to add key pair, private key is nil")
	ErrAddKeyPairFromFilePrivateKeyNotFound = errors.New("Unable to add key pair, private key not found at %s")

	ErrGatewayHandlerIsNil = errors.New("Gateway entrypoint handler nil")
)

func ErrReplacer(err error, replacers ...interface{}) error {
	return errors.New(fmt.Errorf(err.Error(), replacers...))
}
