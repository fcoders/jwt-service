package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fcoders/jwt-service/common"
	"github.com/fcoders/jwt-service/core/cache"
	"github.com/fcoders/jwt-service/settings"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	expireOffset = 60 // indicates how is the expiration time represented (60=minutes, 3600=hours, etc)
)

var authBackendInstance *JWTAuthenticationBackendKeys
var tokenCache cache.Connector

// JWTAuthenticationBackendKeys hols the keys used for signing the token
type JWTAuthenticationBackendKeys struct {
	Store map[string]*KeyStore
}

// GetStore returns the store instance identified by id
func (backend *JWTAuthenticationBackendKeys) GetStore(id string) (ks *KeyStore, exists bool) {
	if backend.Store != nil {
		ks, exists = backend.Store[id]
	}
	return
}

// ParseToken parse the string with the token and returns a jwt.Token
func (backend *JWTAuthenticationBackendKeys) ParseToken(token string, id string) (*jwt.Token, error) {
	return jwt.Parse(token, backend.KeyFunc(id))
}

// KeyFunc is a function used to validate the algorithm and extract the token from the request
func (backend *JWTAuthenticationBackendKeys) KeyFunc(id string) jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, err error) {

		// validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		if ks, exists := backend.GetStore(id); exists {
			return ks.PublicKey, nil
		}

		err = fmt.Errorf("No keys defined for client ID %s", id)
		return
	}
}

// KeyStore represents an in memory store for private/public keys
type KeyStore struct {
	ID         string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// IsLoaded returns true if private and public keys are loaded correcrtly
func (ks *KeyStore) IsLoaded() bool {
	return ks.PrivateKey != nil && ks.PublicKey != nil
}

// LoadPrivateKey loads the private key from the file
func (ks *KeyStore) LoadPrivateKey(path string) (err error) {

	var pk *rsa.PrivateKey

	privateKeyFile, errOpenFile := os.Open(path)
	if errOpenFile != nil {
		err = fmt.Errorf("Error opening private key file: %s", errOpenFile.Error())
		return
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	pk, err = x509.ParsePKCS1PrivateKey(data.Bytes)
	ks.PrivateKey = pk
	return
}

// LoadPublicKey loads the public key from file
func (ks *KeyStore) LoadPublicKey(path string) (err error) {

	publicKeyFile, errOpenFile := os.Open(path)
	if errOpenFile != nil {
		err = fmt.Errorf("Error opening private key file: %s", errOpenFile.Error())
		return
	}
	defer publicKeyFile.Close()

	pemfileinfo, _ := publicKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))

	publicKeyImported, errParse := x509.ParsePKIXPublicKey(data.Bytes)
	if errParse != nil {
		log.Fatalf("Error parsing public key file: %s", err.Error())
	}

	pk, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("Error converting public key to RSA")
	}

	ks.PublicKey = pk
	return
}

// InitJWTAuthenticationBackend initializes the JWT auth system with the keys
func InitJWTAuthenticationBackend(cacheConnector cache.Connector) (bk *JWTAuthenticationBackendKeys, err error) {
	if authBackendInstance == nil {
		authBackendInstance = new(JWTAuthenticationBackendKeys)

		// load store
		store, errStore := loadKeyStores()
		if errStore != nil {
			err = errStore
			return
		}

		authBackendInstance.Store = store

		// Start connections with Redis server
		conf := settings.Get().Redis
		tokenCache = cacheConnector
		tokenCache.Init(conf.Address, conf.Password)
	}

	return authBackendInstance, nil
}

// GenerateToken generates a new token for the user.
// Parameter id represents the user ID or subject and grant represent the token type generated
// (currently only 'access_token' supported)
func (backend *JWTAuthenticationBackendKeys) GenerateToken(requestClaims map[string]interface{}, id string) (tokenString string, expiresIn int, err error) {

	token := jwt.New(jwt.SigningMethodRS512)
	claims := token.Claims.(jwt.MapClaims)
	now := time.Now()

	claims["exp"] = now.Add(time.Minute * time.Duration(settings.Get().JWT.TokenExpiration)).Unix()
	claims["iat"] = now.Unix()

	for k, v := range requestClaims {

		switch k {
		case "grant":
			claims["typ"] = v
		// case "id":
		// 	claims["sub"] = v

		default:
			claims[k] = v
		}
	}

	token.Claims = claims

	if store, exists := backend.GetStore(id); exists {

		tokenString, err = token.SignedString(store.PrivateKey)
		if err != nil {
			log.Fatalf("Error signing the token: %s", err.Error())
		}

		expiresIn = settings.Get().JWT.TokenExpiration * 60
	} else {
		err = fmt.Errorf("No keys defined for client ID %s", id)
	}

	return
}

// Destroy invalidates the token by saving it in our cache, until it expires
func (backend *JWTAuthenticationBackendKeys) Destroy(token *jwt.Token) error {
	claims := token.Claims.(jwt.MapClaims)
	return tokenCache.SetValue(token.Raw, token.Raw, backend.GetTokenRemainingValidity(claims["exp"]))
}

// IsInBlacklist checks if the token has been marked as invalid
func (backend *JWTAuthenticationBackendKeys) IsInBlacklist(token string) bool {
	if redisToken, _ := tokenCache.GetValue(token); redisToken == nil {
		return false
	}

	return true
}

// CloseCacheConnections closes all the current connections with the current cache system
func CloseCacheConnections() {
	if tokenCache != nil {
		tokenCache.Close()
	}
}

// GetTokenRemainingValidity returns the remaining time for the token expiration,
// expressed in seconds.
func (backend *JWTAuthenticationBackendKeys) GetTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}

	return expireOffset
}

func loadKeyStores() (store map[string]*KeyStore, err error) {

	keysPath := path.Join(common.GetAppPath(), "keys")
	files, errRead := ioutil.ReadDir(keysPath)
	if errRead != nil {
		err = fmt.Errorf("Failed to read keys folder: %s", errRead)
		return
	}

	store = make(map[string]*KeyStore)

	for i := range files {
		if files[i].IsDir() {

			// look for key files
			dirPath := path.Join(keysPath, files[i].Name())
			filesDir, errReadDir := ioutil.ReadDir(dirPath)
			if errReadDir != nil {
				err = fmt.Errorf("Cannot read content from '%s': %s", dirPath, errReadDir)
				return
			}

			ks := new(KeyStore)

			for k := range filesDir {

				switch strings.ToLower(filesDir[k].Name()) {

				case "key":
					// load private key
					if errPK := ks.LoadPrivateKey(path.Join(dirPath, filesDir[k].Name())); errPK != nil {
						err = fmt.Errorf("Cannot load private key from %s/%s: %s", dirPath, filesDir[k].Name(), errPK)
						return
					}

				case "key.pub":
					// load public key
					if errPK := ks.LoadPublicKey(path.Join(dirPath, filesDir[k].Name())); errPK != nil {
						err = fmt.Errorf("Cannot load public key from %s/%s: %s", dirPath, filesDir[k].Name(), errPK)
						return
					}

				default:

				}

			}

			if ks.IsLoaded() {
				store[files[i].Name()] = ks
			}

		}
	}

	return
}
