package main

import (
	"context"
	"encoding/json"
	"ff14/common"
	"ff14/config"
	"ff14/respdata"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var Env = config.Env

func main() {

	var user string
	var pass string
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		user = os.Getenv("USER")
		pass = os.Getenv("PASS")
	} else {
		user = Env.GetString("user.username")
		if user == "" {
			fmt.Println("请输入账号：")
			fmt.Scanln(&user)
		}
		pass = Env.GetString("user.password")
		if pass == "" {
			fmt.Println("请输入密码：")
			passByte, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println("\n读取密码时发生错误：", err)
				return
			}
			pass = string(passByte)
		}
	}

	if user == "" || pass == "" {
		panic("账号或密码为空，无法运行")
	}

	Env.Set("user.username", user)
	Env.Set("user.password", pass)
	err := Env.WriteConfig()
	if err != nil {
		panic(err)
	}

	// 签到
	response := common.Post(Env.GetString("api.sign_in"), nil, "")
	err = common.FFXIVIsError(response)
	if err != nil {
		if strings.Contains(err.Error(), "登录") {
			fmt.Println("密钥失效，尝试重新登录获取密钥")
			u := Env.GetString("user.username")
			p := Env.GetString("user.password")
			var networkCookies []*network.Cookie
			var opts []chromedp.ExecAllocatorOption

			flags := []chromedp.ExecAllocatorOption{
				chromedp.Flag("headless", true), // 是否隐藏浏览器窗口，如果是win10要测试改为false
				chromedp.Flag("enable-automation", false),
			}
			opts = chromedp.DefaultExecAllocatorOptions[:]
			opts = append(opts, flags...)

			allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
			ctx, _ := chromedp.NewContext(
				allocCtx,
				chromedp.WithLogf(log.Printf),
			)

			defer cancel()

			if err := chromedp.Run(ctx,
				network.ClearBrowserCache(),
				chromedp.Navigate("https://ff14risingstones.web.sdo.com/pc/index.html#/post"),
				chromedp.ActionFunc(func(ctx context.Context) error {
					fmt.Println("开始登录")
					return nil
				}),
				chromedp.Sleep(5*time.Second),
				chromedp.WaitVisible(`#aside > div.el-card.is-always-shadow > div > div.mb20 > div.flex.h120.alcenter.jccenter > button`, chromedp.ByQuery),
				chromedp.Click("#aside > div.el-card.is-always-shadow > div > div.mb20 > div.flex.h120.alcenter.jccenter > button", chromedp.ByQuery),
				chromedp.WaitVisible(`#isAgreementAccept`, chromedp.ByQuery),
				chromedp.Click(`#isAgreementAccept`, chromedp.ByQuery),
				chromedp.Sleep(2*time.Second),
				chromedp.WaitVisible(`#wegame_btn`, chromedp.ByQuery),
				chromedp.Click(`#wegame_btn`, chromedp.ByQuery),
				chromedp.Sleep(5*time.Second),
			); err != nil {
				fmt.Println(err)
			}

			var iframes []*cdp.Node
			if err := chromedp.Run(ctx,
				chromedp.Nodes(`iframe`, &iframes, chromedp.ByQuery),
				chromedp.ActionFunc(func(ctx context.Context) error {
					fmt.Println("却换到wegame登录页面")
					return nil
				}),
			); err != nil {
				fmt.Println(err)
			}

			sectionNode := iframes[0]

			if err := chromedp.Run(ctx,
				chromedp.WaitVisible("#switcher_plogin"),
				chromedp.ActionFunc(func(ctx context.Context) error {
					fmt.Println("输入账号密码登录阶段")
					return nil
				}),
				chromedp.Click("#switcher_plogin", chromedp.ByID, chromedp.FromNode(sectionNode)),
				chromedp.Sleep(2*time.Second),
				chromedp.SetValue(`#u`, u, chromedp.ByID, chromedp.FromNode(sectionNode)),
				chromedp.Sleep(2*time.Second),
				chromedp.SetValue(`#p`, p, chromedp.ByID, chromedp.FromNode(sectionNode)),
				chromedp.Sleep(2*time.Second),
				chromedp.Click(`#login_button`, chromedp.ByID, chromedp.FromNode(sectionNode)),
				chromedp.Sleep(6*time.Second),
			); err != nil {
				panic(err)
			}

			if err := chromedp.Run(ctx,
				chromedp.WaitVisible(`#me`, chromedp.ByQuery),
				chromedp.Navigate("https://ff14risingstones.web.sdo.com/pc/index.html#/post"),
				chromedp.Sleep(3*time.Second),
				chromedp.ActionFunc(func(ctx context.Context) error {
					fmt.Println("登录成功，获取cookie")
					var err error
					networkCookies, err = network.GetCookies().Do(ctx)
					return err
				}),
			); err != nil {
				fmt.Println(err)
			}

			for _, cookie := range networkCookies {
				Env.Set("cookie."+cookie.Name, cookie.Value)
				fmt.Println("设置cookie：", cookie.Name, cookie.Value)
				err := Env.WriteConfig()
				if err != nil {
					panic(err)
				}
			}

			fmt.Println("重新请求，当前密钥：", Env.GetString("cookie.ff14risingstones"))
			response = common.Post(Env.GetString("api.sign_in"), nil, "")
			err = common.FFXIVIsError(response)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("签到成功")

	// 获取列表
	fmt.Println("开始获取最新发布列表")
	query := url.Values{}
	query.Add("type", "1")
	query.Add("order", "latest")
	query.Add("page", "2") // 因为第一页有置顶
	query.Add("limit", "15")
	response = common.Get(Env.GetString("api.posts_list"), query)
	err = common.FFXIVIsError(response)
	if err != nil {
		panic(err)
	} else {
		// 点赞
		var posts respdata.PostsList
		time.Sleep(5 * time.Second)
		json.Unmarshal([]byte(response), &posts)
		for k, v := range posts.Data.Rows {
			if k >= 5 {
				break
			}
			postLike(v.PostsId)
			time.Sleep(2 * time.Second)
		}

		// 评论一次
		postComment()

		// 盖章
		for i := 1; i < 4; i++ {
			doSeal(strconv.Itoa(i))
			time.Sleep(2 * time.Second)
		}
	}
}

func postLike(id string) {
	params := url.Values{}
	params.Add("id", id)
	params.Add("type", "1")
	response := common.Post(Env.GetString("api.post_like"), params, "")
	err := common.FFXIVIsError(response)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("点赞成功")
	}
}

func postComment() {
	params := url.Values{}
	params.Add("content", "<p>每日一水<span class=\"at-emo\">[emo2]</span>&nbsp;</p>")
	params.Add("posts_id", Env.GetString("official.post_id"))
	params.Add("parent_id", "0")
	params.Add("root_parent", "0")
	params.Add("comment_pic", "")
	response := common.Post(Env.GetString("api.post_comment"), params, "")
	err := common.FFXIVIsError(response)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("评论成功")
	}
}

func doSeal(typeString string) {
	params := url.Values{}
	params.Add("type", typeString)
	response := common.Post(Env.GetString("api.do_seal"), params, "")
	err := common.FFXIVIsError(response)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("盖章成功")
	}
}
