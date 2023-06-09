package util

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// GenerateCode 生成六位数验证码
func GenerateCode() int {
	var code string
	for i := 0; i < 6; i++ {
		t := strconv.Itoa(rand.Intn(10))
		code = code + t
	}
	s1, _ := strconv.Atoi(code)
	return s1
}

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)

}

// CheckIdCard 检验身份证
func CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// GetIdFromToken 从token中获取id
func GetIdFromToken(c *gin.Context) int {
	// get token
	tokenString := c.GetHeader("Authorization")
	//token is nil or token is valid
	//Determine if the token starts with Bearer
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		c.Abort()
	}

	// delete Bearer from token
	tokenString = tokenString[6:]

	//parse token
	token, claims, err := ParseToken(tokenString)
	// error or token is not valid
	if err != nil || !token.Valid {
		panic(err)
	}
	return claims.UserId
}
