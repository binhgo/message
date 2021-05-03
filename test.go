package main

import (
	"crypto"
	"crypto/rsa"
	"fmt"

	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/action"
	"gitlab.ghn.vn/internal-tools/message/api"
	"gitlab.ghn.vn/internal-tools/message/cip"
	"gitlab.ghn.vn/internal-tools/message/model"
)

func TEST() {
	// test_firebase_send_msg()
	// test01()
	// testEncrypt3()
	// testEncrypt2()
	// testEncrypt()
	// test0001()
}

// func test01() {
//
// 	str := "fuckkk yoouuu hahah"
//
// 	// The GenerateKey method takes in a reader that returns random bits, and
// 	// the number of bits
// 	priKey := cip.GenerateRsaKeyPair()
// 	pubKey := priKey.PublicKey
//
// 	encryptedBytes := api.EncryptMsg(str, &pubKey)
//
// 	fmt.Println("encrypted bytes: ", encryptedBytes)
//
// 	decryptedStr := action.DecryptMsg(encryptedBytes, priKey)
// 	fmt.Println(decryptedStr)
//
// 	decryptedBytes, err := priKey.Decrypt(nil, []byte(encryptedBytes), &rsa.OAEPOptions{Hash: crypto.SHA256})
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// We get back the original information in the form of bytes, which we
// 	// the cast to a string and print
// 	fmt.Println("decrypted message: ", string(decryptedBytes))
// }

// func testEncrypt3() {
//
// 	var messageInput model.MsgInput
// 	msgId := bson.ObjectIdHex("5f947dabe2d51c70c5ebaadc")
// 	msgFilter := model.Message{ID: msgId}
// 	queryRs := model.MessageDB.QueryOne(msgFilter)
// 	if queryRs.Status == sdk.APIStatus.Ok {
// 		msg := queryRs.Data.([]*model.Message)[0]
// 		fmt.Println(msg.Content)
//
// 		contentByte, _ := json.Marshal(msg.Content)
// 		json.Unmarshal(contentByte, &messageInput)
// 	}
//
// 	userId := "1005"
//
// 	filter := model.User{
// 		UserID: userId,
// 	}
//
// 	queryRs = model.UserDB.QueryOne(filter)
// 	if queryRs.Status == sdk.APIStatus.Ok {
// 		user1004 := queryRs.Data.([]*model.User)[0]
//
// 		priKey, err := cip.PemToPrivateKey(user1004.Key)
// 		fmt.Println(err)
//
// 		decryptedBytes, err := priKey.Decrypt(nil, messageInput.Msg, &rsa.OAEPOptions{Hash: crypto.SHA256})
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		fmt.Println(string(decryptedBytes))
// 	}
// }

func testEncrypt2() {

	str := "hi how are you there"

	userId := "1004"

	filter := model.User{
		UserID: userId,
	}

	queryRs := model.UserDB.QueryOne(filter)
	if queryRs.Status == sdk.APIStatus.Ok {
		user1004 := queryRs.Data.([]*model.User)[0]

		priKey, err := cip.PemToPrivateKey(user1004.Key)
		fmt.Println(err)

		pubKey, err := cip.PemToPublicKey(user1004.PubKey)
		// pubKey := priKey.PublicKey
		fmt.Println(err)

		encryptedStr := api.EncryptMsg(str, pubKey)
		fmt.Println(encryptedStr)

		// decryptedBytes, err := priKey.Decrypt(nil, []byte(encryptedStr), &rsa.OAEPOptions{Hash: crypto.SHA256})
		// if err != nil {
		// 	panic(err)
		// }

		decryptedStr := action.DecryptMsg([]byte(encryptedStr), priKey)
		fmt.Println(string(decryptedStr))
	}
}

func testEncrypt() {

	str := "hi how are you there"
	//
	pripem, pubpem, err := cip.GenerateRsaKeyPairPem()
	fmt.Println(err)

	prikey, err := cip.PemToPrivateKey(pripem)
	fmt.Println(err)

	pubkey, err := cip.PemToPublicKey(pubpem)
	fmt.Println(err)

	// pubKey := priKey.PublicKey

	encryptedStr := api.EncryptMsg(str, pubkey)
	fmt.Println(encryptedStr)

	decryptedStr := action.DecryptMsg([]byte(encryptedStr), prikey)
	fmt.Println(decryptedStr)

	userId := "1004"

	filter := model.User{
		UserID: userId,
	}

	queryRs := model.UserDB.QueryOne(filter)
	if queryRs.Status == sdk.APIStatus.Ok {
		user1004 := queryRs.Data.([]*model.User)[0]

		// pubKey, err := cip.PemToPublicKey(user1004.PubKey)
		// fmt.Println(err)
		//
		// encryptedStr := api.EncryptMsg(str, pubKey)
		// fmt.Println(encryptedStr)

		priKey, err := cip.PemToPrivateKey(user1004.Key)
		fmt.Println(err)

		pubKey := priKey.PublicKey
		fmt.Println(err)

		encryptedStr := api.EncryptMsg(str, &pubKey)
		fmt.Println(encryptedStr)

		decryptedBytes, err := priKey.Decrypt(nil, []byte(encryptedStr), &rsa.OAEPOptions{Hash: crypto.SHA256})
		if err != nil {
			panic(err)
		}

		// decryptedStr := action.DecryptMsg(str, priKey)
		fmt.Println(decryptedBytes)
	}
}
