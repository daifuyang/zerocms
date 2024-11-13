package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
	"zerocms/api/model"
	"zerocms/utils"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	username := req.Username
	password := req.Password

	userModel, err := l.svcCtx.UserModel.FindOneByLoginName(l.ctx, sql.NullString{String: username, Valid: true})
	if err != nil {
		return nil, err
	}

	userPwd := userModel.Password.String
	salt := userModel.Salt.String

	if !utils.ComparePassword(userPwd+salt, password+salt) {
		return nil, errors.New("账号密码不正确")
	}

	data, err := generateTokens(strconv.FormatInt(userModel.Id, 10), l.svcCtx.Config.TokenSecretKey)
	if err != nil {
		return nil, err
	}

	// 存入sql中
	_, err = l.svcCtx.UserTokenModel.Insert(l.ctx, &model.SysUserToken{
		UserId:      userModel.Id,
		AccessToken: data.AccessToken,
		ExpiresAt: sql.NullTime{
			Time:  time.Unix(data.ExpiresAt, 0),
			Valid: true,
		},
		RefreshToken: data.RefreshToken,
		ReExpiresAt: sql.NullTime{
			Time:  time.Unix(data.ReExpiresAt, 0),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.LoginResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: data,
	}
	return
}

// 生成 JWT 令牌
func generateTokens(userID string, secretKey string) (types.TokenData, error) {
	// 计算 token 的过期时间
	accessTokenExpiry := time.Now().Add(time.Hour * 24).Unix()       // 1天有效期
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 30).Unix() // 30天有效期

	// 生成 access_token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     accessTokenExpiry,
	})

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return types.TokenData{}, fmt.Errorf("could not create access token: %v", err)
	}

	// 生成 refresh_token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshTokenExpiry,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return types.TokenData{}, fmt.Errorf("could not create refresh token: %v", err)
	}

	// 返回令牌和过期时间
	return types.TokenData{
		AccessToken:  accessTokenString,
		ExpiresAt:    accessTokenExpiry,
		RefreshToken: refreshTokenString,
		ReExpiresAt:  refreshTokenExpiry,
	}, nil
}
