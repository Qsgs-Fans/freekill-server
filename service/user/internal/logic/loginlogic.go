package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/user/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) decrypt(password string) ([]byte, string, error) {
	var empty []byte
	// base64 string -> bytes
	ciphertext, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return empty, "", err
	}

	// bytes -> plainText
	decryptedPassword, err := rsa.DecryptPKCS1v15(
		rand.Reader,
		l.svcCtx.RsaKeyPair.PrivateKey,
		ciphertext,
	)
	if err != nil {
		return empty, "", err
	}

	// plainText -> AES key + password
	if len(decryptedPassword) < 32 {
		return empty, "", fmt.Errorf("Invalid password data")
	}
	hexAESKey := string(decryptedPassword[:32])
	userPassword := string(decryptedPassword[32:])
	aesKey, err := hex.DecodeString(hexAESKey)
	if err != nil || len(aesKey) != 16 { // 128-bit = 16字节
		return empty, "", fmt.Errorf("Invalid AES key")
	}

	return aesKey, userPassword, nil
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginReply, error) {
	// TODO checkVersion
	// TODO checkMd5
	// TODO checkUUidBanned
	aesKey, password, err := l.decrypt(in.Password)
	if err != nil {
		return &user.LoginReply{}, err
	}
	// TODO checkPassword
	fmt.Println(password)
	// TODO 将用户添加到用户列表（Redis）中之类的 aesKey交给网关层调用者
	return &user.LoginReply{
		AesKey: aesKey,
	}, nil
}
