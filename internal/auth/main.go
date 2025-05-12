package auth

import (
	"fmt"

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