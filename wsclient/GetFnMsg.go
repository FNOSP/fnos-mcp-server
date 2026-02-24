package wsclient

import (
	"encoding/base64"
	"encoding/json"
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

type FileLsDto struct {
	DefaultDto
	Path string `json:"path,omitempty"`
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
	aes, err := f.aesCbcEncryptBase64(plaintext, f.Key, f.Iv)
	if err != nil {
		panic(err)
	}
	pubKeyByte := []byte(pubKey) // string -> []byte
	rsa, err := f.rsaEncrypt(f.Key, pubKeyByte)

	msg := LoginDto{
		Req: "encrypted",
		Iv:  base64.StdEncoding.EncodeToString(f.Iv),
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

func (f *FnOsWsBase) GetFileLsData(path *string) FnOsRequestBase[string] {
	reqId := f.GetReqId()
	_d := DefaultDto{
		Req:   "file.ls",
		ReqID: reqId,
	}
	var msgData FileLsDto
	if path == nil {
		msgData = FileLsDto{
			DefaultDto: _d,
		}
	} else {
		msgData = FileLsDto{
			DefaultDto: _d,
			Path:       *path,
		}
	}
	msg, err := f.GetSignReq(msgData)
	if err != nil {
		panic(err)
	}
	req := FnOsRequestBase[string]{
		Msg:   msg,
		ReqID: reqId,
	}
	return req
}

func (f *FnOsWsBase) GetMachineIdData() FnOsRequestBase[string] {
	reqId := f.GetReqId()
	_d := DefaultDto{
		Req:   "appcgi.sysinfo.getMachineId",
		ReqID: reqId,
	}
	msg, err := f.GetSignReq(_d)
	if err != nil {
		panic(err)
	}
	req := FnOsRequestBase[string]{
		Msg:   msg,
		ReqID: reqId,
	}
	return req
}
func (f *FnOsWsBase) GetNetworkNetListData() FnOsRequestBase[string] {
	reqId := f.GetReqId()
	_d := DefaultDto{
		Req:   "appcgi.network.net.list",
		ReqID: reqId,
	}
	msg, err := f.GetSignReq(_d)
	if err != nil {
		panic(err)
	}
	req := FnOsRequestBase[string]{
		Msg:   msg,
		ReqID: reqId,
	}
	return req
}
func (f *FnOsWsBase) GetHardwareInfoData() FnOsRequestBase[DefaultDto] {
	reqId := f.GetReqId()
	_d := DefaultDto{
		Req:   "appcgi.sysinfo.getHardwareInfo",
		ReqID: reqId,
	}
	req := FnOsRequestBase[DefaultDto]{
		Msg:   _d,
		ReqID: reqId,
	}
	return req
}
