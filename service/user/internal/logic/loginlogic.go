package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/room/room"
	"github.com/Qsgs-Fans/freekill-server/service/user/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/user/model"
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

func (l *LoginLogic) checkDeviceIdBanned (deviceId string) error {
	if len(deviceId) > 64 {
		return fmt.Errorf("Invalid device id %v: too long", deviceId)
	}

	_, err := l.svcCtx.BannedDevicesModel.FindOneByDeviceUuid(l.ctx, deviceId)
	if err == nil {
		return fmt.Errorf("Refused banned device id %v", deviceId)
	} else if err != model.ErrNotFound {
		return fmt.Errorf("Query error: %v", err)
	}

	return nil
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

func (l *LoginLogic) register(in *user.LoginRequest, password string) (*model.Users, error) {
	var salt int32
	_ = binary.Read(rand.Reader, binary.LittleEndian, &salt)
	saltHex := fmt.Sprintf("%x", salt)
	pwAndSalt := password + saltHex
	pwHash := sha256.Sum256([]byte(pwAndSalt))
	pwHashStr := hex.EncodeToString(pwHash[:])

	newUser := &model.Users{
		Username: in.Username,
		Password: pwHashStr,
		Salt: saltHex,
		Avatar: "liubei",
	}
	res, err := l.svcCtx.UsersModel.Insert(l.ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.Id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// TODO 添加注册信息 这样才能防止同device id以及ip反复注册

	return newUser, nil
}

func (l *LoginLogic) checkPassword(in *user.LoginRequest, password string) (*model.Users, error) {
	res, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, in.Username)

	var isReg bool
	if err != nil {
		if err == model.ErrNotFound {
			res2, err2 := l.register(in, password)
			if err2 != nil {
				return nil, err2
			}
			res = res2
			isReg = true
		} else {
			return nil, err
		}
	}

	var passed bool
	if isReg {
		passed = true
	} else {
		salt := res.Salt
		pwAndSalt := password + salt
		pwHash := sha256.Sum256([]byte(pwAndSalt))
		pwHashStr := hex.EncodeToString(pwHash[:])

		passed = pwHashStr == res.Password
	}

	if passed {
		return res, nil
	} else {
		return nil, fmt.Errorf("password error: user=%v", in.Username)
	}
}

func (l *LoginLogic) setupPlayer(connId string, info *model.Users) {
	sender := l.svcCtx.Sender
	jsonArr := []any{ info.Id, info.Username, info.Avatar }
	sender.Notify(l.ctx, "Setup", jsonArr, connId)
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginReply, error) {
	// TODO checkVersion 服务端能支持多个客户端版本吗？
	// TODO checkMd5 MD5保存在哪？能支持多个Md5登录吗？还是必须强制更新到最新的？
	// TODO 黑名单白名单 此为CRUD 先todo得了
	err := l.checkDeviceIdBanned(in.Deviceid)
	if err != nil {
		errmsg := "you have been banned!"
		err2 := l.svcCtx.Sender.NotifyRaw(l.ctx, "ErrorDlg", errmsg, in.ConnId)
		if err2 != nil {
			return &user.LoginReply{}, fmt.Errorf("Error when sending banned message: %v", err2)
		}
		return &user.LoginReply{}, err
	}

	aesKey, password, err := l.decrypt(in.Password)
	if err != nil {
		return &user.LoginReply{}, err
	}

	userInfo, err := l.checkPassword(in, password)
	if err != nil {
		return &user.LoginReply{}, err
	}

	// TODO 在数据库中插入登录信息 此为CRUD

	// TODO 断线重连相关

	// 进入server->createNewPlayer()
	l.svcCtx.AddNewLoginUser(userInfo, in)

	l.setupPlayer(in.ConnId, userInfo)
	// TODO notify("SetServerSettings")
	// TODO notify("AddTotalGametime")
	rrpc := l.svcCtx.RoomRpc
	rrpc.EnterRoom(l.ctx, &room.UidAndRidRequest{
		UserId: userInfo.Id,
		RoomId: 0,
	})

	// 最后 aesKey交给网关层调用者
	return &user.LoginReply{
		AesKey: aesKey,
		UserId: userInfo.Id,
	}, nil
}
