package keycloak

import (
	"bitbucket.org/pinang-och-go/pkg/v1/utils/errors"
	"context"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"strings"
)

// isValidToken validates the status token
func IsValidToken(ctx context.Context) error {
	// path public key
	var publickeyPath = "./env/public_key.pem"

	// read header from incoming request
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.FormatError(errors.Authorization, codes.Unauthenticated, false, "401", "Unauthorization")
	}
	if len(headers["authorization"]) != 1 {
		return errors.FormatError(errors.Authorization, codes.Unauthenticated, false, "401", "Unauthorization")
	}

	token := headers["authorization"][0]

	if token == "" {
		return errors.FormatError(errors.InvalidFormat, codes.InvalidArgument, false, "401", "Missing authentication")
	}

	// read public key from path
	publikeyData, err := ioutil.ReadFile(publickeyPath)
	if err != nil {
		return err
	}

	// parse the public key
	publickey, err := jwt.ParseRSAPublicKeyFromPEM(publikeyData)
	if err != nil {
		return err
	}

	// validate token with public key
	parts := strings.Split(token, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], publickey)
	if err != nil {
		return errors.FormatError(errors.InvalidFormat, codes.InvalidArgument, false, "401", "Invalid token")
	}

	// validate token expired
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return publickey, nil
	})
	// log
	//if err != nil {
	//	return err
	//}
	if err != nil && err.Error() == "Token is expired" {
		return errors.FormatError(errors.InvalidFormat, codes.InvalidArgument, false, "401", "Expired token")
	}

	return nil
}
