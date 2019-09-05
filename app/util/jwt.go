package util

import (
	"Goo/app/db"
	"Goo/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"time"
)

var (
	JWTSecretKey       = []byte("you-will-never-guess")
	JWTSecretKeyGetter = func(token *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	}
)

func GenerateJWToken(user model.User) (string, error) {
	claim := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		//Token签发者，格式是区分大小写的字符串或者uri，用于唯一标识签发token的一方。
		"iss": "Datagrand",
		//Token的主体，即它的所有人，格式是区分大小写的字符串或者uri。
		"sub": "Anyone who is a datagrand user",
		//指定Token在nbf时间之前不能使用，即token开始生效的时间，格式为时间戳。
		"nbf": time.Now().Unix(),
		//Token的签发时间，格式为时间戳。
		"iat": time.Now().Unix(),
		//Token的过期时间，格式为时间戳。
		"exp": time.Now().Add(time.Hour * time.Duration(24)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(JWTSecretKey)
	return tokenStr, err
}

func ParseJWToken(tokenStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, JWTSecretKeyGetter)
	if err != nil {
		err = errors.New("Cannot parse token")
		return nil, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("Cannot convert claim")
		return nil, err
	}
	if !token.Valid {
		err = errors.New("Token is invalid")
		return nil, err
	}
	return claim, err
}

func GetCurrentClaims(ctx iris.Context) (jwt.Claims, error) {
	token := ctx.Values().Get("jwt")
	if token == nil {
		return nil, errors.New("Cannot get jwt claims")
	}
	return token.(*jwt.Token).Claims, nil
}

func GetCurrentUser(ctx iris.Context) (user model.User, err error) {
	claims, err := GetCurrentClaims(ctx)
	if err != nil {
		return
	}
	userID := claims.(jwt.MapClaims)["id"]
	result := db.Session.First(&user, "id = ?", userID)
	return user, result.Error
}
