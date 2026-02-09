package wsclient

import (
	"fmt"
	"time"
)

// Login 登录方法
// userName 飞牛系统登录账号
// passWord	飞牛系统登录密码
func (f *FnOsWsBase) Login(userName, passWord string) error {
	rsa_pub := f.GetRsaPub()
	loginData := f.GetLoginData(userName, passWord, rsa_pub.Si, rsa_pub.Pub)
	data, err := f.Send(loginData.Msg, loginData.ReqID, 10*time.Second)
	if err != nil {
		panic(err)
	}
	ret, err := ConvertTo[LoginRetDto](data)
	if err != nil {
		panic(err)
	}
	f.LoginRetDto = ret
	return nil
}

func (f *FnOsWsBase) GetRsaPub() GetRSAPubRetDto {
	rsapubData := f.RsaPubData()
	req_id := rsapubData.ReqID
	data, err := f.Send(rsapubData.Msg, req_id, 10*time.Second)
	if err != nil {
		_ = fmt.Errorf("获取RSAPUB失败: %w", err)
		panic(err)
	}
	ret, err := ConvertTo[GetRSAPubRetDto](data)
	if err != nil {
		panic(err)
	}

	return ret
}

func (f *FnOsWsBase) GetHostName() GetHostNameRetDto {
	reqData := f.GetHostNameData()
	data, err := f.Send(reqData.Msg, reqData.ReqID, 10*time.Second)
	ret, err := ConvertTo[GetHostNameRetDto](data)
	if err != nil {
		panic(err)
	}
	return ret
}
