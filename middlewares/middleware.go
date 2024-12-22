package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "No authorization header"})
            c.Abort()
            return
        }
        parts := strings.Split(tokenString, " ")
		fmt.Println(parts);
        tokenString = parts[1];
		fmt.Println(tokenString);
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
            fmt.Printf("Token algorithm: %v\n", token.Header["alg"])
			secret := os.Getenv("Secret")
			fmt.Printf("Secret length: %d\n", len(os.Getenv("Secret")));
			
			return []byte(secret), nil
					})
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error":err.Error()})
            c.Abort()
            return
        }
		claims, ok := token.Claims.(jwt.MapClaims)
		metadata, ok := claims["user_metadata"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User metadata not found"})
			c.Abort()
			return
		}

		if role, exists := metadata["role"].(string); exists {
			fmt.Println(role);
			c.Set("user_role", role)
		}

        c.Next()
    }
}


	// return func(c *gin.Context) {
	// 	authHeader := c.GetHeader("Authorization")
	// 	if authHeader == "" {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
	// 		c.Abort()
	// 		return
	// 	}
	// 	obtainedtoken := strings.TrimPrefix(authHeader, "Bearer ")
	// 	token, err := jwt.Parse(obtainedtoken, func(token *jwt.Token) (interface{}, error) {
	// 		// Don't forget to validate the alg is what you expect:
	// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 		}
	// 		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("supabaseSecret")))
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to parse public key: %v", err)
	// 		}

	// 		return publicKey, nil
	// 	})
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
	// 		c.Abort()
	// 		return
	// 	}
	
	// 	c.Next()
	// }
// }

