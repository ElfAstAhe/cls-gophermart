package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AppClaims struct {
	jwt.RegisteredClaims
	Admin  bool   `json:"admin,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Roles  Roles  `json:"roles,omitempty"`
}

const secretKey = "secret-key"
const timeExpirationDuration = time.Minute * 120

var TestRoles Roles = []string{
	"test_role1",
	"test_role2",
}

func NewJWTStringFromUserInfo(userInfo *UserInfo) (string, error) {
	if userInfo == nil {
		return "", errors.New("userInfo is nil")
	}

	return NewJWTString(userInfo.Admin, userInfo.UserID, userInfo.User, userInfo.Roles...)
}

func NewJWTString(admin bool, userID string, user string, roles ...string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newAppClaims(admin, userID, user, roles...))

	res, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return res, nil
}

func newAppClaims(admin bool, userID string, user string, roles ...string) *AppClaims {
	return &AppClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        buildUniqueTokenID(),
			Issuer:    "url-shortener",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeExpirationDuration)),
			Subject:   user,
		},
		Admin:  admin,
		UserID: userID,
		Roles:  roles,
	}
}

func buildUniqueTokenID() string {
	const template = "shortener-token-%v"
	randID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Sprintf(template, time.Now().Nanosecond())
	}

	return fmt.Sprintf(template, randID.String())
}

func retrieveJWT(r *http.Request) (*jwt.Token, error) {
	cookie, _ := r.Cookie(CookieName)
	if cookie == nil {
		return nil, errs.NewAuthCookieAbsentError(CookieName, nil)
	}
	if err := cookie.Valid(); err != nil {
		return nil, errs.NewAuthCookieAbsentError("invalid cookie", err)
	}

	claims := &AppClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errs.NewAuthInfoAbsentError("error parsing token", err)
	}

	return token, nil
}

func UserInfoFromRequestJWT(r *http.Request) (*UserInfo, error) {
	jwtToken, err := retrieveJWT(r)
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errs.NewAuthInfoAbsentError("JWT is invalid", nil)
	}

	res, err := UserInfoFromJWT(jwtToken)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UserInfoFromJWT(jwt *jwt.Token) (*UserInfo, error) {
	claims, ok := jwt.Claims.(*AppClaims)
	if !ok {
		return nil, errs.NewAuthInfoAbsentError("error transform raw JWT claims into AppClaims struct", nil)
	}

	return NewUserInfo(claims.Admin, claims.UserID, claims.Subject, claims.Roles), nil
}
