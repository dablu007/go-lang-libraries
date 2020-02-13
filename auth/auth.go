package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"flow/logger"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var jwtMiddleWare *jwtmiddleware.JWTMiddleware

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	// Note that the scope can be string or an array
	RawScope json.RawMessage `json:"scope"`
	// Scopes need to be unmarshalled post the initial unmarshalling as we can't be sure of the type
	Scopes []string `json:"-"`
	jwt.StandardClaims
}

func Init() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := viper.GetString("JwksAudience")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			fmt.Println(checkAud)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := viper.GetString("JwksIssuer")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	jwtMiddleWare = jwtMiddleware
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client secret key
		err := jwtMiddleWare.CheckJWT(c.Writer, c.Request)
		if err != nil {
			// Token not found
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(viper.GetString("JwksUrl"))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func ValidateScope(token string,validScopes string) bool {
	jsonTokens := strings.Split(token, ".")
	if len(jsonTokens) != 3 {
		logger.SugarLogger.Warnw("Token structure does not seem to be as expected. Token: %scope", token)
		return false
	}

	payloadToken := jsonTokens[1]
	decodedToken, decodeError := b64.StdEncoding.DecodeString(payloadToken + "==")
	if decodeError != nil {
		logger.SugarLogger.Warnw("Unable to decode token. Payload token:",map[string]string{
			"payloadToken":payloadToken,
		})
		return false
	}

	claims := CustomClaims{}
	marshallError := json.Unmarshal([]byte(decodedToken), &claims)
	if marshallError != nil {
		logger.SugarLogger.Warnw("Unable to unmarshal decoded claims. decodedToken:", map[string]string{
			"decodedToken": decodedToken,
		})
		return false
	}

	claims.Scopes, marshallError = claims.getScopes()
	if marshallError != nil {
		logger.SugarLogger.Warnw("Unable to get scopes. for token ")
		return false
	}

	for _, scope := range claims.Scopes {
		if strings.Contains(validScopes, scope){
			return true
		}
	}
	// Returning true if in-case it is not a customer token
	return false
}

func (claims *CustomClaims) getScopes() ([]string, error) {
	if len(claims.RawScope) == 0 {
		logger.SugarLogger.Warnf("Scope raw message is empty. Claims: %s", claims)
		return nil, errors.New("Scope raw message is empty")
	}

	switch claims.RawScope[0] {
	case '"':
		var scope string
		if err := json.Unmarshal(claims.RawScope, &scope); err != nil {
			logger.SugarLogger.Warnf("Unable to unmarshall stringified scope. RawScope: %s", claims.RawScope)
			return nil, errors.New("Unable to unmarshall string scope")
		}
		return []string{scope}, nil

	case '[':
		var scopes []string
		if err := json.Unmarshal(claims.RawScope, &scopes); err != nil {
			logger.SugarLogger.Warnf("Unable to unmarshall arrayed scopes. RawScopes: %s", claims.RawScope)
			return nil, errors.New("Unable to unmarshall arrayed scopes")
		}
		return scopes, nil
	}
	logger.SugarLogger.Warnf("Unable to unmarshal scopes. RawScopes: %s", claims.RawScope)
	return nil, errors.New("Unable to unmarshal scopes")
}


