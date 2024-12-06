package utils_test

import (
	"mmddvg/chapar/pkg/services/utils"
	"testing"
)

func TestEncryption(t *testing.T) {
	u, p := "username", "password"
	ep, err := utils.Encrypt(u, p)
	if err != nil {
		t.Error(err)
	}

	isMatch, err := utils.CheckPassword(ep, p, u)
	if err != nil {
		t.Error(err.Error())
	}

	if !isMatch {
		t.Error("isMatch is false")
	}
}
