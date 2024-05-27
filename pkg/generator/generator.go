package generator

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

func GeneratePrivateKey(start string) (string, error) {
	// 生成随机的32字节私钥
	privateKeyBytes := make([]byte, 32)
	_, err := rand.Read(privateKeyBytes)
	if err != nil {
		return "", err
	}
	privateKeyHex := start + hex.EncodeToString(privateKeyBytes)[1:]
	return privateKeyHex, nil
}

func PrivateKeyToAddress(privateKeyHex string) string {
	// 将16进制私钥转换为ECDSA私钥
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	// 获取以太坊地址
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return address
}
