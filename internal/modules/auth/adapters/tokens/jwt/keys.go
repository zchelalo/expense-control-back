package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Keys struct {
	AccessPrivate  *rsa.PrivateKey
	AccessPublic   *rsa.PublicKey
	RefreshPrivate *rsa.PrivateKey
	RefreshPublic  *rsa.PublicKey
}

type KeyPaths struct {
	AccessPrivatePath  string
	AccessPublicPath   string
	RefreshPrivatePath string
	RefreshPublicPath  string
}

func LoadKeys(paths KeyPaths) (*Keys, error) {
	read := func(p string) ([]byte, error) {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", p, err)
		}
		return b, nil
	}

	ap, err := read(paths.AccessPrivatePath)
	if err != nil { return nil, err }
	apu, err := read(paths.AccessPublicPath)
	if err != nil { return nil, err }
	rp, err := read(paths.RefreshPrivatePath)
	if err != nil { return nil, err }
	rpu, err := read(paths.RefreshPublicPath)
	if err != nil { return nil, err }

	accessPriv, err := jwtlib.ParseRSAPrivateKeyFromPEM(ap)
	if err != nil { return nil, fmt.Errorf("parse access private: %w", err) }
	accessPub, err := jwtlib.ParseRSAPublicKeyFromPEM(apu)
	if err != nil { return nil, fmt.Errorf("parse access public: %w", err) }

	refreshPriv, err := jwtlib.ParseRSAPrivateKeyFromPEM(rp)
	if err != nil { return nil, fmt.Errorf("parse refresh private: %w", err) }
	refreshPub, err := jwtlib.ParseRSAPublicKeyFromPEM(rpu)
	if err != nil { return nil, fmt.Errorf("parse refresh public: %w", err) }

	return &Keys{
		AccessPrivate:  accessPriv,
		AccessPublic:   accessPub,
		RefreshPrivate: refreshPriv,
		RefreshPublic:  refreshPub,
	}, nil
}