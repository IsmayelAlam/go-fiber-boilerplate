package utils

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var tokenType = "access"

type JWTConfig struct {
	Issuer              string
	Audience            string
	Secret              string
	TokenExpiry         int
	RefreshExpiry       int
	RefreshCookieDomain string
	RefreshCookiePath   string
	RefreshCookieName   string
}

type TokenPair struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (j *JWTConfig) GenerateToken(id uuid.UUID) (TokenPair, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = fmt.Sprint(id)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = tokenType
	claims["exp"] = time.Now().UTC().Add(time.Duration(j.TokenExpiry) * time.Hour).Unix()

	signedAccessToken, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["sub"] = fmt.Sprint(id)
	refreshClaims["iat"] = time.Now().UTC().Unix()
	refreshClaims["exp"] = time.Now().UTC().Add(time.Duration(j.RefreshExpiry) * 24 * time.Hour).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	TokenPair := TokenPair{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return TokenPair, nil
}

func (j *JWTConfig) SetRefreshCookie(c *fiber.Ctx, token string) {
	cookie := new(fiber.Cookie)
	cookie.Name = j.RefreshCookieName
	cookie.Path = j.RefreshCookiePath
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.MaxAge = int((time.Duration(j.RefreshExpiry) * 24 * time.Hour).Seconds())
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.Domain = j.RefreshCookieDomain
	cookie.Secure = true
	cookie.HTTPOnly = true

	// Set the cookie in the response
	c.Cookie(cookie)
}

func (j *JWTConfig) GetExpiredRefreshCookie(c *fiber.Ctx) {

	cookie := new(fiber.Cookie)
	cookie.Name = j.RefreshCookieName
	cookie.Path = j.RefreshCookiePath
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.Domain = j.RefreshCookieDomain
	cookie.Secure = true
	cookie.HTTPOnly = true

	// Set the cookie in the response
	c.Cookie(cookie)
}

func (j *JWTConfig) AccessTokenValidate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect (e.g., HMAC with SHA256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		// jwt.Parse returns errors for invalid signatures, malformed tokens, etc.
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims format")
	}

	// --- Manual validation of standard claims (in case jwt lib skips some) ---

	now := time.Now().UTC().Unix()

	// 1. Validate "exp" (expiration time)
	if exp, ok := claims["exp"].(float64); !ok || now >= int64(exp) {
		return "", fmt.Errorf("token is expired")
	}

	// 2. Validate "nbf" (not before) – optional but good practice
	if nbf, ok := claims["nbf"].(float64); ok && now < int64(nbf) {
		return "", fmt.Errorf("token not yet valid")
	}

	// 3. Validate "iat" (issued at) – optional sanity check
	if iat, ok := claims["iat"].(float64); ok && now < int64(iat) {
		return "", fmt.Errorf("token issued in the future")
	}

	// 4. Validate "iss" (issuer)
	if iss, ok := claims["iss"].(string); !ok || iss != j.Issuer {
		return "", fmt.Errorf("invalid token issuer")
	}

	// 5. Validate "aud" (audience)
	if aud, ok := claims["aud"].(string); !ok || aud != j.Audience {
		return "", fmt.Errorf("invalid token audience")
	}

	// 6. Validate custom "typ" claim (your token type, e.g., "access")
	if typ, ok := claims["typ"].(string); !ok || typ != tokenType {
		return "", fmt.Errorf("invalid token type")
	}

	// 7. Validate "sub" (subject/user ID)
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", fmt.Errorf("invalid or missing subject (sub) claim")
	}

	return sub, nil
}

func (j *JWTConfig) RefreshTokenValidate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect (e.g., HMAC with SHA256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		// jwt.Parse returns errors for invalid signatures, malformed tokens, etc.
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims format")
	}

	// --- Manual validation of standard claims (in case jwt lib skips some) ---

	now := time.Now().UTC().Unix()

	// 1. Validate "exp" (expiration time)
	if exp, ok := claims["exp"].(float64); !ok || now >= int64(exp) {
		return "", fmt.Errorf("token is expired")
	}

	// 2. Validate "nbf" (not before) – optional but good practice
	if nbf, ok := claims["nbf"].(float64); ok && now < int64(nbf) {
		return "", fmt.Errorf("token not yet valid")
	}

	// 3. Validate "iat" (issued at) – optional sanity check
	if iat, ok := claims["iat"].(float64); ok && now < int64(iat) {
		return "", fmt.Errorf("token issued in the future")
	}

	// 4. Validate "sub" (subject/user ID)
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", fmt.Errorf("invalid or missing subject (sub) claim")
	}

	return sub, nil
}
