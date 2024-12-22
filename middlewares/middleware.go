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
func CheckRole(c *gin.Context, requiredRole string) bool {
    if userRole, exists := c.Get("user_role"); exists {
        if role, ok := userRole.(string); ok && role == requiredRole {
            return true
        }
    }
    return false
}
