package util

//aes 加密解密

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
)

const (
	KEY = "A clever brown fox jumps over a lazy dog"
)

// padding
func padding(src []byte,blocksize int) []byte {
	padnum 	:= blocksize - len(src) % blocksize
	pad 	:= bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src,pad...)
}

// unpadding
func unpadding(src []byte) []byte {
	n	:= len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

// EncryptAES 数据加密
func EncryptAES(src []byte) []byte {
	aesKey := md5.Sum([]byte(KEY))
	block,_:=aes.NewCipher(aesKey[:])
	src=padding(src, block.BlockSize())
	blockmode:=cipher.NewCBCEncrypter(block, aesKey[:])
	blockmode.CryptBlocks(src, src)
	return src
}

// DecryptAES 数据解密
func DecryptAES(src []byte,key []byte) []byte {
	aesKey := md5.Sum([]byte(KEY))
	block,_:=aes.NewCipher(aesKey[:])
	blockmode:=cipher.NewCBCDecrypter(block, aesKey[:])
	blockmode.CryptBlocks(src,src)
	src=unpadding(src)
	return src
}



