package auth

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type KeyCloakAuth struct {
	Keys Jwts
	Role string
}

type Jwts struct {
	Keys []JWT `json:"keys"`
}

type JWT struct {
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	N   string `json:"n,omitempty"`
	E   string `json:"e,omitempty"`
	Use string `json:"use,omitempty"`
}

func StartKeyCloakAuth(keyCloakUrl string, role string) (*KeyCloakAuth, error) {
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(fiber.MethodGet, keyCloakUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("KeyCloak not working")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	keys := make([]JWT, 2)
	j := Jwts{Keys: keys}

	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	kCA := new(KeyCloakAuth)
	kCA.Keys = j
	kCA.Role = role

	return kCA, nil
}

func (k *KeyCloakAuth) Auth(c *fiber.Ctx) error {

	if c.Path() == "/healthz" || c.Path() == "/metrics" || c.Path() == "/swagger/index.html" {
		return c.Next()
	}

	headers := c.GetReqHeaders()
	authHeader, ok := headers["Authorization"]
	if !ok {
		return c.Status(401).JSON("Empty Token")
	}

	if !k.verify(authHeader[0]) {
		return c.Status(401).JSON("Wrong Token")
	}

	return c.Next()
}

func (k *KeyCloakAuth) ExtractBearerToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		return token[len("Bearer "):]
	}
	return token
}

func (k *KeyCloakAuth) verify(ah string) bool {
	token, err := jwt.Parse(k.ExtractBearerToken(ah), func(token *jwt.Token) (any, error) {
		keyID, _ := token.Header["kid"].(string)

		var selectedJWK *JWT

		for _, jwk := range k.Keys.Keys {
			if jwk.Kid == keyID {
				selectedJWK = &jwk
				break
			}
		}

		if selectedJWK == nil {
			return errors.New("not found kid"), nil
		}

		n, _ := base64.RawURLEncoding.DecodeString(selectedJWK.N)
		e, _ := base64.RawURLEncoding.DecodeString(selectedJWK.E)
		z := new(big.Int)
		z.SetBytes(n)
		var buffer bytes.Buffer
		buffer.WriteByte(0)
		buffer.Write(e)
		exponent := binary.BigEndian.Uint32(buffer.Bytes())

		publicKey := &rsa.PublicKey{N: z, E: int(exponent)}

		return publicKey, nil
	})
	if err != nil {
		log.Println("Error from Keycloak: ", err)
		return false
	} else if !token.Valid {
		return false
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	roles := claims["realm_access"].(map[string]any)

	for _, r := range roles["roles"].([]any) {
		if r.(string) == k.Role {
			return true
		}
	}

	return false
}
