package routes

import (
	"fmt"
	"model/user"
	"net/http"
)

//signin
// 逻辑请参考 jtw/src/wechat/ndcwx/routes/signin.go
//里面部分信息的查询和更改（比如用户，抵用金流量的增加）是用的rcp。rpc的 client在golibs/src/git.ot24.net/ndp/jtm ,server在jts
//
const (
	// signin 主页加载
	signinHtmlPath = "/signin/index"
	//TODO 抽奖
	doSigninPath = "/signin/api"
)

func init() {
	http.Handle(signinHtmlPath, ReqURLPrt(http.HandlerFunc(signinHTML)))
	http.Handle(doSigninPath, ReqURLPrt(http.HandlerFunc(doSignin)))
}

func signinHTML(w http.ResponseWriter, r *http.Request) {
	u, ok := auth(w, r)
	if !ok {
		logl.Println("0o")
		return
	}

	name := u.UserName

	userDB, err := user.NewUserDB()
	if err != nil {
		logl.Error(err)
		String(w, 500, err.Error())
		return
	}

	user, err := userDB.GetUserByName(name)
	if err != nil {
		logl.Error(err)
		String(w, 500, err.Error())
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logl.Error(err.(error))
		}
	}()
	fmt.Println(user)
	//TODO:签到记录表信息,查看用户是否签过到，给前端传一个状态

	Render(w, 200, "public/user/info.html",
		H{
			"user": user,
			//是否签过到。前端根据这个参数决定签到按钮是否失效
			//"status": status,
			//TODO:活动等信息
		})
}

//TODO 抽奖
func doSignin(w http.ResponseWriter, r *http.Request) {
	//TODO:向签到记录表记录签到信息。
	//TODO:抽奖,将奖品信息传到前端
	prize := "签到成功！假装你中奖了：iphone7s "
	JSON(w, 200, H{
		"mess": prize,
	})
}
