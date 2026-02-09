package wsclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type FnOsRequestBase[T any] struct {
	Msg   T      `json:"msg"`
	ReqID string `json:"req_id"`
}

type DefaultDto struct {
	Req   string `json:"req"`
	ReqID string `json:"reqid"`
}
type LoginDto struct {
	Req string `json:"req"`
	Iv  string `json:"iv"`
	Rsa string `json:"rsa"`
	Aes string `json:"aes"`
}

func (f *FnOsWsBase) RsaPubData() FnOsRequestBase[DefaultDto] {
	reqId := f.GetReqId()
	msg := DefaultDto{
		Req:   "util.crypto.getRSAPub",
		ReqID: reqId,
	}
	req := FnOsRequestBase[DefaultDto]{
		Msg:   msg,
		ReqID: reqId,
	}
	return req
}
func (f *FnOsWsBase) GetLoginData(userName, passWord, si, pubKey string) FnOsRequestBase[LoginDto] {
	reqId := f.GetReqId()
	aesStr := map[string]interface{}{
		"reqid":      reqId,
		"user":       userName,
		"password":   passWord,
		"deviceType": "Browser",
		"deviceName": "Windows-Google Chrome",
		"stay":       true,
		"req":        "user.login",
		"si":         si,
	}
	plaintext, err := json.Marshal(aesStr)
	if err != nil {
		panic(err)
	}
	aes, err := f.aesCbcEncryptBase64(plaintext, f.key, f.iv)
	if err != nil {
		panic(err)
	}
	pubKeyByte := []byte(pubKey) // string -> []byte
	fmt.Printf("%v\n", pubKeyByte)
	rsa, err := f.rsaEncrypt(f.key, pubKeyByte)

	msg := LoginDto{
		Req: "encrypted",
		Iv:  base64.StdEncoding.EncodeToString(f.iv),
		Rsa: rsa,
		Aes: aes,
	}
	req := FnOsRequestBase[LoginDto]{
		Msg:   msg,
		ReqID: reqId,
	}
	return req
}

func (f *FnOsWsBase) GetHostNameData() FnOsRequestBase[DefaultDto] {
	reqId := f.GetReqId()
	msg := DefaultDto{
		Req:   "appcgi.sysinfo.getHostName",
		ReqID: reqId,
	}
	req := FnOsRequestBase[DefaultDto]{
		Msg:   msg,
		ReqID: reqId,
	}
	return req
}
