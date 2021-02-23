package logic

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"go-zero-book/common/errorx"
	"go-zero-book/service/user/model"
	"strings"
	"time"

	"go-zero-book/service/user/cmd/api/internal/svc"
	"go-zero-book/service/user/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginReply, error) {
	// todo: add your logic here and delete this line
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}

	userInfo, err := l.svcCtx.UserModel.FindOneByNumber(req.Username)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errorx.NewDefaultError("用户名不存在")
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		return nil, errorx.NewDefaultError("用户密码不正确")
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}

	return &types.LoginReply{
		Id: userInfo.Id,
		Name: userInfo.Name,
		Gender: userInfo.Gender,
		AccessToken: jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *LoginLogic) getJwtToken(secret string, now int64, expire int64, id int64) (string, error) {
	//claims := make(jwt.MapClaims)
	//claims["exp"] = now + expire
	//claims["iat"] = now
	//claims["userId"] = id
	//token := jwt.New(jwt.SigningMethodES256)
	//token.Claims = claims


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": now + expire,
		"iat": now,
		"userId": id,
	})

	//tokenString, err := token.SignedString([]byte("key"))

	return token.SignedString([]byte(secret))
}
