package auth

import (
	"errors"
	"fmt"
	
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string)(string,error)  {
	hashByte,err := bcrypt.GenerateFromPassword([]byte(password),10);
	if err != nil{
	fmt.Printf("Error occured generating HashedPAssword: %v",err);
		return "",err
	}
	return string(hashByte),nil

}

func CheckPasswordHash(hash,password string)error  {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password));
return err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	
	expiresAt := time.Now().Add(expiresIn * time.Second)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject: userID.String(),
	});

	signedToken,err := token.SignedString([]byte(tokenSecret));
	if err != nil{
		return "",err
	}
	return signedToken,nil
}

func ValidateJWT(tokenString,tokenSecret string)(uuid.UUID,error)  {
	
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}
	
	token,err := jwt.ParseWithClaims(tokenString,&MyCustomClaims{},func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret),nil
	});


	if err != nil{
		fmt.Printf("Error: %v\n",err)
		return uuid.Nil,err
	}else if claims,ok := token.Claims.(*MyCustomClaims);ok{
		uid,err := uuid.Parse(claims.RegisteredClaims.Subject)
		if err != nil{
			//fmt.Printf("Error occurred converting string to uuid %v",err)
			return uuid.Nil, fmt.Errorf("invalid UUID in token: %w", err)
		}
		return uid,nil
	}else{
		//log.Fatal("Unknown claims type,cannot proceed")
		return uuid.Nil,errors.New("unkown Claims")
	}
}

func GetBearerToken(headers http.Header)(string,error)  {
	authHeader := headers.Get("Authorization");
	reqToken := strings.TrimPrefix(authHeader,"Bearer ");
	//log.Printf("token: %v\n",reqToken)
	if authHeader == "" || reqToken == authHeader{
		return "",errors.New("authentication header not present")
	}
	return reqToken,nil
}