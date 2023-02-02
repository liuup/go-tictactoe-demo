package token

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TOKEN_HOUR_LIFESPAN = "1"
	API_SECRET          = "yoursecrectstring"
)

func GenerateToken(user_id uint) (string, error) {
	token_lifespan, err := strconv.Atoi(TOKEN_HOUR_LIFESPAN) // 参数
	if err != nil {
		return "", err
	}

	// 除了加上userid，还可以加上其他身份内容
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(API_SECRET)) // 参数
}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	// log.Println(tokenString)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	// 要么来自于query参数，要么来自于request header
	// 需要把token提取出来
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")

	// log.Println(bearerToken)

	// if len(strings.Split(bearerToken, " ")) == 2 {
	// 	return strings.Split(bearerToken, " ")[1]
	// }
	// return ""

	return bearerToken
}

func ExtractTokenID(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}
