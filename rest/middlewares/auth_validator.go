package middlewares

import (
	"fmt"
	"strings"

	"statsv0/models/security"
	"statsv0/tools/custerror"

	"github.com/gin-gonic/gin"
)

/**
 * @apiDefine AuthHeader
 *
 * @apiExample {String} Header Autorización
 *    Authorization=bearer {token}
 *
 * @apiErrorExample 401 Unauthorized
 *    HTTP/1.1 401 Unauthorized
 */

// ValidateAuthentication validate gets and check variable body to create new variable
// puts model.Variable in context as body if everything is correct
func ValidateAuthentication(c *gin.Context) {
	fmt.Println("ValidateAuthentication...")
	if err := validateToken(c); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}

//ejecuta la funcion validate de security que busca en cache o en la api de auth el usuario por el token dentro del header del context de la request

var securityValidate func(token string) (*security.User, error) = security.Validate

func validateToken(c *gin.Context) error {
	tokenString, err := GetHeaderToken(c)
	if err != nil {
		return custerror.Unauthorized
	}

	if _, err = securityValidate(tokenString); err != nil {
		return custerror.Unauthorized
	}

	return nil
}

// get token from Authorization header
func GetHeaderToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", custerror.Unauthorized
	}
	return tokenString[7:], nil
}
