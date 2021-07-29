package service

import "errors"

var (
	// ErrInvalidPort error message
	ErrInvalidPort                   = errors.New("invalid port %d")
	ErrInvalidRecieveSize            = errors.New("invalid max recieve size %d")
	ErrInvalidSendSize               = errors.New("invalid max send size %d")
	ErrMissingGrpcServerOptions      = errors.New("invalid grpc server options: missing")
	ErrMissingGrpcUnaryInterceptors  = errors.New("invalid grpc unary interceptors: missing")
	ErrMissingGrpcStreamInterceptors = errors.New("invalid grpc stream interceptors: missing")
	ErrMissingGatewayHandlers        = errors.New("invalid grpc gateway handlers: missing")
	ErrWithGRPCServerIsNil           = errors.New("invalid grpc server: is nil ")
	ErrServeGRPCNotYetDefined        = errors.New("service serve gprc not yet defined (has pb.RegisterYourServer(service.GRPC(), ...) been called?)")

	ErrDisableTLSCertsMissingGRPCPort = errors.New("invalid grpc port: must set GRPC port when disabling TLS (use environment " + EnvGRPCPort + " or WithGRPCPort(port)")

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
)
