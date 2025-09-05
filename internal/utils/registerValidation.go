package utils

import (
	"errors"
	"log"
	"regexp"

	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

func RegisterValidation(body models.Auth) error {
	// cek format email
	// ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	regexEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexEmail.Match([]byte(body.Email)) {
		return errors.New("format penulisa email salah")
	}

	// cek format password
	// harus : huruf, angka, simbol, 8 karakter
	log.Println(body.Password)
	islengEight := len(body.Password) >= 8
	isNotHvSymbl := regexp.MustCompile(`[!@#$%^&*/><]`).MatchString(body.Password)
	isNotHvChar := regexp.MustCompile(`[a-zA-Z]`).MatchString(body.Password)
	isNotHvDigit := regexp.MustCompile(`\d`).MatchString(body.Password)

	log.Println(isNotHvSymbl, isNotHvChar, isNotHvDigit, islengEight)
	if !isNotHvChar || !isNotHvSymbl || !isNotHvDigit || !islengEight {
		return errors.New("email harus mengandung : huruf, angka, simbol, 8 karakter")
	}
	return nil
}
