package main

import (
	"fmt"
	"database/sql"
	"encoding/json"	
	"net/http"
	"time"
	"strings"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Authentication methods
type accessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var hmacAccessSecret = []byte("abcdefgh")
var hmacRefreshSecret = []byte("ijklmnopq")

// Middlewares
func APImiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
    })
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("Authorization")
		stringToken := strings.TrimSpace(strings.Replace(authorizationHeader, "Bearer", "", 1))
		JWTToken, err := jwt.ParseWithClaims(stringToken, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
			return hmacAccessSecret, nil
		})
		if err != nil {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			fmt.Println(err)
			return
		}
		

		if claims, ok := JWTToken.Claims.(*accessClaims); ok && JWTToken.Valid {
			fmt.Printf("%v %v", claims.Email, claims.RegisteredClaims.Issuer)
		} else {
			fmt.Println(err)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
    })
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/authentication/login server login")
		var userLoginData UserLoginData
		err := json.NewDecoder(r.Body).Decode(&userLoginData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("UserLoginData values",userLoginData)
		userLoginDataDB := getUserByEmail(userLoginData.Email, db)
		fmt.Println("userLoginDataDB values",userLoginDataDB)
		
		if userLoginDataDB.Password == userLoginData.Password {
			accessToken, refreshToken, err := loginSignToken(userLoginData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{
				"accessToken": accessToken,
				"refreshToken": refreshToken})
			return
		}
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
		
	}
}

// Sign Access Token
func signAccessToken(user UserLoginData) (string, error) {	
	fmt.Println("signAccessToken")
	// Create claims with multiple fields populated
	accessClaims := accessClaims{
		user.Email,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Issue at time
			NotBefore: jwt.NewNumericDate(time.Now()), // Work since
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)		
	accessTokenString, err := accessToken.SignedString(hmacAccessSecret)

	return accessTokenString, err
}

// Sign Refresh Token
func signRefreshToken(user UserLoginData) (string, error) {
	fmt.Println("signRefreshToken")
	// Create claims with multiple fields populated	
	refreshClaims := refreshClaims{
		user.Email,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 14)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Issue at time
			NotBefore: jwt.NewNumericDate(time.Now()), // Work since
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)		
	refreshTokenString, err := refreshToken.SignedString(hmacRefreshSecret)

	return refreshTokenString, err
}

// Login Sign Token
func loginSignToken(user UserLoginData) (string, string, error) {
	accessTokenString, accessErr := signAccessToken(user)
	refreshTokenString, refreshErr := signRefreshToken(user)
	err := errors.Join(accessErr, refreshErr)
	fmt.Println(accessTokenString, accessErr, refreshTokenString, refreshErr)
	return accessTokenString, refreshTokenString, err
}

// Check if Refresh Token is Valid 
func validRefreshToken(stringToken string) (*refreshClaims, error) {
    
	JWTToken, err := jwt.ParseWithClaims(stringToken, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacRefreshSecret, nil
	})
	if err != nil {
		return nil, err
	}	

	if claims, ok := JWTToken.Claims.(*refreshClaims); ok && JWTToken.Valid {
		return claims, nil
	} else {		
		return nil, err
	}
}

// Refresh Sign Token
func refreshSignToken(user UserLoginData) (string, error) {
	accessTokenString, accessErr := signAccessToken(user)
	fmt.Println(accessTokenString, accessErr)
	return accessTokenString, accessErr
}