package crypto_test

import (
	"fmt"
	"github.com/kaydxh/golang/go/crypto"
	"testing"
)

func TestAesCbcEncryptDecrypt(t *testing.T) {
	plainText := []byte("Hello World")
	fmt.Println("plainText: ", string(plainText))

	key := []byte("daW3eDgPEa9TjknE")
	cryptText, err := crypto.AesCbcEncrypt(plainText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("cryptText: ", string(cryptText))

	newplainText, err := crypto.AesCbcDecrypt(cryptText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("newplainText: ", string(newplainText))
	if string(newplainText) != string(plainText) {
		t.Errorf("newplainText not equal plainText")
	}
}
