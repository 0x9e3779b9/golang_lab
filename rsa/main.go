/*********************************************
*  Last modified: 2015-12-07 20:40
*  Filename: main.go
*  Description:
*********************************************/
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func genKey() (private, public []byte) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}
	priDer := x509.MarshalPKCS1PrivateKey(priKey)
	private = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: priDer,
		},
	)

	pub := priKey.PublicKey
	pubDer, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		fmt.Println(err)
		return
	}
	public = pem.EncodeToMemory(
		&pem.Block{
			Type:    "PUBLIC KEY",
			Headers: nil,
			Bytes:   pubDer,
		},
	)
	ioutil.WriteFile("public.pem", public, 0644)
	ioutil.WriteFile("private.pem", private, 0644)
	return
}

func hash(data []byte) []byte {
	hashFunc := crypto.MD5
	h := hashFunc.New()
	h.Write(data)
	return h.Sum(nil)
}

func sign(src []byte, priPem []byte) ([]byte, error) {
	block, _ := pem.Decode(priPem)
	priKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	hashFunc := crypto.MD5
	return rsa.SignPKCS1v15(rand.Reader, priKey, hashFunc, hash(src))
}

func verify(src, sign []byte, pubPem []byte) error {
	block, _ := pem.Decode(pubPem)
	pubKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	return rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.MD5, hash(src), sign)
}

func Encrypt(data []byte, pubPem []byte) ([]byte, error) {
	block, _ := pem.Decode(pubPem)
	pubKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), data)
}

func Decrypt(ciphertext []byte, priPem []byte) ([]byte, error) {
	block, _ := pem.Decode(priPem)
	priKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return rsa.DecryptPKCS1v15(rand.Reader, priKey, ciphertext)
}

func test() {
	var data = []byte(`{Hello}`)
	pri, pub := genKey()
	sign, err := sign(data, pri)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = verify(data, sign, pub)
	if err != nil {
		fmt.Println(1, err)
		return
	}
}

func test1() {
	var data = []byte(`{Hello}`)
	pri, pub := genKey()
	ciphertext, _ := Encrypt(data, pub)
	dd, _ := Decrypt(ciphertext, pri)
	println(string(dd))
}

func main() {
	test()
	test1()
}
