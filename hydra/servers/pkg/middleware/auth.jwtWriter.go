package middleware

import (
	"fmt"

	xjwt "github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/security/jwt"
)

//JwtWriter 将jwt信息写入到请求中
func JwtWriter() Handler {
	return func(ctx IMiddleContext) {
		ctx.Next()
		conf, err := ctx.APPConf().GetJWTConf()
		if err != nil {
			ctx.Response().Abort(xjwt.JWTStatusConfError, err)
			return
		}

		if conf.Disable {
			return
		}
		setJwtResponse(ctx, conf, ctx.User().Auth().Response())
	}
}

func setJwtResponse(ctx context.IContext, jwtAuth *xjwt.JWTAuth, data interface{}) {
	var jwtToken = ""
	var err error
	if data != nil {
		jwtToken, err = jwt.Encrypt(jwtAuth.Secret, jwtAuth.Mode, data, jwtAuth.ExpireAt)
		if err != nil {
			ctx.Response().Abort(xjwt.JWTStatusConfDataError, fmt.Errorf("jwt配置出错：%v", err))
			return
		}
	}
	//检查是否需要跳过请求
	if ok, _ := jwtAuth.Match(ctx.Request().Path().GetRequestPath()); ok && jwtToken == "" {
		return
	}

	setToken(ctx, jwtAuth, jwtToken)
}

//setToken 设置jwt到响应头或cookie中
func setToken(ctx context.IContext, jwt *xjwt.JWTAuth, jwtToken string) {
	k, v := jwt.GetJWTForRspns(jwtToken, jwtToken == "")
	ctx.Response().Header(k, v)
}
