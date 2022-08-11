package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const ENCRYPTKEY string = "tcdts-encrypting"

type EncryptionController struct {
	//16, 24, or 32 bytes
	encrypting_key string
}

func NewEncryptionController() *EncryptionController {
	k := len(ENCRYPTKEY)
	if k == 16 || k == 24 || k == 32 {
		return &EncryptionController{encrypting_key: ENCRYPTKEY}
	} else {
		return nil
	}
}

//解密
func (e *EncryptionController) AesDecrypt(cryted string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(e.encrypting_key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

func (e *EncryptionController) AesEncrypt(orig string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(e.encrypting_key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

//补码
func (e *EncryptionController) PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func (e *EncryptionController) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func main() {
	ec := NewEncryptionController()
	passWord := "42wdInpJouggEiGfSnnmtpmV"
	fmt.Println("原文：", passWord)

	encryptCode := ec.AesEncrypt(passWord)
	fmt.Println("密文：", encryptCode)

	dbconn, _ := sql.Open("mysql", "root:password#dbr@tcp(10.100.156.210:3306)/shop_lyp")
	dbconn.SetMaxOpenConns(10)
	dbconn.SetMaxIdleConns(5)

	mSql := fmt.Sprintf("select password from role where id=%d", 5)
	rows, _ := dbconn.Query(mSql)
	defer rows.Close() //这里如果不释放连接到池里，执行5次后其他并发就会阻塞
	for rows.Next() {
		var password string
		err := rows.Scan(&password)
		if nil != err {
			fmt.Println(err)
		}
	}
	decryptCode := ec.AesDecrypt(encryptCode)
	fmt.Println("解密结果：", decryptCode)
}

// func main() {
// 	passWord := "42wdInpJouggEiGfSnnmtpmV"
// 	key := "123456781234567812345678"
// 	fmt.Println("原文：", passWord)

// 	encryptCode := AesEncrypt(passWord, key)
// 	fmt.Println("密文：", encryptCode)

// 	dbconn, _ := sql.Open("mysql", "root:password#dbr@tcp(10.100.156.210:3306)/shop_lyp")
// 	dbconn.SetMaxOpenConns(10)
// 	dbconn.SetMaxIdleConns(5)
// 	// result, err := dbconn.Exec("INSERT INTO role (password) VALUES (?)", encryptCode)
// 	// if nil != err {
// 	// 	fmt.Println(err)
// 	// }

// 	// id, err := result.LastInsertId()
// 	// if nil != err {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(id)

// 	mSql := fmt.Sprintf("select password from role where id=%d", 5)
// 	rows, _ := dbconn.Query(mSql)
// 	defer rows.Close() //这里如果不释放连接到池里，执行5次后其他并发就会阻塞
// 	for rows.Next() {
// 		var password string
// 		err := rows.Scan(&password)
// 		if nil != err {
// 			fmt.Println(err)
// 		}
// 	}

// 	decryptCode := AesDecrypt(encryptCode, key)
// 	fmt.Println("解密结果：", decryptCode)
// }

func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// func main() {
// 	db, _ := sql.Open("mysql", "root:password#dbr@tcp(10.100.156.210:3306)/shop_lyp")
// 	db.SetMaxOpenConns(10)
// 	db.SetMaxIdleConns(5)
// 	//连接数据库查询
// 	// for i := 0; i < 100; i++ {
// 	// 	go func(i int) {
// 	// 		mSql := "select * from user"
// 	// 		rows, _ := db.Query(mSql)
// 	// 		rows.Close() //这里如果不释放连接到池里，执行5次后其他并发就会阻塞
// 	// 		fmt.Println("第 ", i)
// 	// 	}(i)
// 	// }

// 	mSql := "select password from role where id=10 "
// 	rows, _ := db.Query(mSql)

// 	rows.Close() //这里如果不释放连接到池里，执行5次后其他并发就会阻塞

// 	for {
// 		time.Sleep(time.Second)
// 	}
// }
