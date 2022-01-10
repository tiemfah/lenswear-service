package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"
)

func GenerateRSA(publicKeyPath, privateKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		logrus.Errorf("err read private : ", err)
	}
	pub, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		logrus.Errorf("err read public : ", err)
	}

	privateKey := bytesToPrivateKey(priv)
	publicKey := bytesToPublicKey(pub)

	return privateKey, publicKey
}

func bytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Panic(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Panic(err)
	}
	return key
}

func bytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Panic(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Panic(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Panic("bytesToPublicKey not ok")
	}
	return key
}
