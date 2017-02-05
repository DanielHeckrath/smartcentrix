package main

import (
	"crypto/rsa"
	"time"

	"github.com/101loops/clock"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/juju/errors"
	"github.com/satori/go.uuid"
)

const (
	tokenKeyCreationTimestamp = "iat"
	tokenKeyExpiry            = "exp"
	tokenKeyUserID            = "sid"
)

var tokenDuration = time.Hour * 24 * 30

// TODO: store public and private key in a way that allows certificate rotation
var publicKeyBytes = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7q0h8fLaeWAuvn0pZNn3
KVO+g6jd9TOLIyuou1P61qmlU8BDCTmcgI0+sk/EyAQKh6ZRIl349xOXq98wg9a4
BgouvyuM0RmP2weGKhamPrMK/D0hFYlN6xnUz8TCGgFiXvK4VX3onUKQGa+i/eLo
K8S1gpE03Ztevw2DrbdXvZrmrDFzYTUlW9VikEU+75gXyoZ0whvlMIyWvCfhWiFr
oq+/jgmaTPj26UVfxzB8nb/WTelDay4RYLgC5MHZqiwZm7f8l0GpTjJkNF5c6dJu
4rtk0BzDRXvBEGk4HH9QSOVUrp0Md/akSh4peeTQeaH66fFxkcuz4zws5lJTDR7i
TwIDAQAB
-----END PUBLIC KEY-----`)

var privateKeyBytes = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA7q0h8fLaeWAuvn0pZNn3KVO+g6jd9TOLIyuou1P61qmlU8BD
CTmcgI0+sk/EyAQKh6ZRIl349xOXq98wg9a4BgouvyuM0RmP2weGKhamPrMK/D0h
FYlN6xnUz8TCGgFiXvK4VX3onUKQGa+i/eLoK8S1gpE03Ztevw2DrbdXvZrmrDFz
YTUlW9VikEU+75gXyoZ0whvlMIyWvCfhWiFroq+/jgmaTPj26UVfxzB8nb/WTelD
ay4RYLgC5MHZqiwZm7f8l0GpTjJkNF5c6dJu4rtk0BzDRXvBEGk4HH9QSOVUrp0M
d/akSh4peeTQeaH66fFxkcuz4zws5lJTDR7iTwIDAQABAoIBAQCSpJk4mKeW73uI
2km2hx0OrT9ekUBeSR7xqv9uqThy76R+gqAtaNI5QY8F63DSG3mEwrES1n4DAGyt
0uFx/1jtjRAPsOhCCNyUDFloPqZB26uiMsTDAKt5CVPxm8hORg90mtia3lMvXBhB
T+Jq7yEK0z9aS3EZvz5FUD5ZW05zVSsilsjAV6UMG+Lz8Ao/bMnFm+uIWkRdn9Cw
pKm3sK1EWiJCuKO8zChGWUy2hmS1/cG3bA7yrh38pU4OV5siQBNjrlx1bGvngC+d
EAc3Bp50VbGB1cG+nprWuXVT3yM426ehfGEgkKlSWoxx1sqS6IKF2T7kyPuidqbd
RTSD/DqxAoGBAPqLaA6Gc+5O/q979RCc+QWGTz1wsvbtX/yLLCZ1BgruOp5t7pAt
FPooknh9REuxA2TVIi6I1LX/K1phd2+KEoJFv2uqj/XgrCOoeymZf6YGLP++zd96
51ND5P0xxkxPp+Gme9W4A0S8Kl5u3UctQHSTEbDUiUpqeWX2egjy+qdLAoGBAPPf
kdgp5+Y4OwSCmcOXvGX0/NuPx8iJT7urDkEQ6L0gfeLtwmw3Dww2JDD50i4f+fMf
qBQ34QoJF7W5+B+gZB5gGBL3fyYIWNnV+4+8URJDAOx9dSp8noQEXGEukcGJElDB
oqCzmvaoioeVz4WtJrtvz3wwraN+7Q9qWFKWpXqNAoGAfdrzfYBq0gYah9pTw3Gm
4fCS84EeVU3ujrT2i6bzTyBWj+kXEpOi2vrwgNgkK4WS997cmdWgTIAOrgsR7RTF
sW0J+DKouFGRBySGIeJ6rdKiXiHh1uYtN7V+XPXY79J/ualgwX37HlcLTX6RZ0TD
AQwzsclB4gDUVLTYnpA3+zsCgYBLxxB+ZqcUNizAfgRhbmiwFavsXYTqnyATZFeN
iD+JZOs49EAReBpI5Rnhzf6tLmpwTUng3mwiviiL4zliOmhht+JDInxzyOwy4/bC
9vUKA5/p3CHoDckDpIc/+0R3KqxyQ3jRDn38XuqMrtRI6UC7xUZnhIiv4OSwgY7o
Or84kQKBgQCNnPsS6WSM08Z/3riO8QafwY0OPF3UypaOzXIzspv4CgCZF0/1I5RH
fwALkBT7iclOrZ4e9dYK+wiOsleba0HRMQ4MBmn49HurkevTifpi6UuEHHHJgAjy
Zh1lBpInYxc9Hg5Qlcas244e3SAbwuVSlGjnarC7gUGRCAPxQdylYQ==
-----END RSA PRIVATE KEY-----`)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

var (
	errTokenEmpty         = errors.New("access token cannot be empty")
	errTokenInvalid       = errors.New("access token is invalid")
	errTokenExpired       = errors.New("access token is expired")
	errTokenClaimsInvalid = errors.New("access token claims are invalid")

	errUserIDEmpty  = errors.New("cannot create access token without user id")
	errSigninFailed = errors.New("failed to sign new access token")
)

func init() {
	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		panic(err)
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		panic(err)
	}

	publicKey = pub
	privateKey = priv
}

// claims is a custom jwt claims struct that adds a user id to the standard claims
type claims struct {
	UserID string `json:"uid"`

	jwt.StandardClaims
}

// keyFunc validates that a given token was signed with a RSA256 key
func keyFunc(token *jwt.Token) (interface{}, error) {
	// check that token was signed with RSA
	method, ok := token.Method.(*jwt.SigningMethodRSA)
	if !ok {
		return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	// check that token was signed with RSA256
	if method.Name != jwt.SigningMethodRS256.Name {
		return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return publicKey, nil
}

// validateToken parses a jwt token and validates its signature
func validateToken(token string) (*claims, error) {
	if token == "" {
		return nil, errTokenEmpty
	}

	t, err := jwt.ParseWithClaims(token, &claims{}, keyFunc)

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errTokenInvalid
	}

	claims, ok := t.Claims.(*claims)

	if !ok {
		return nil, errTokenClaimsInvalid
	}

	return claims, nil
}

// generateToken creates a new jwt access token
//
// see https://jwt.io/ for further details
func generateToken(userID string) (string, error) {
	return generateTokenWithID(uuid.NewV4().String(), userID)
}

func generateTokenWithID(tokenID, userID string) (string, error) {
	if userID == "" {
		return "", errUserIDEmpty
	}

	now := clock.Now()

	// create token claims
	claims := claims{
		userID,
		jwt.StandardClaims{
			Id:        tokenID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(tokenDuration).Unix(),
		},
	}

	// create new signed token string
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
