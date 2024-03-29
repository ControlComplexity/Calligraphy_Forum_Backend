package api

import (
	"calligraphy-forum/model/constants"
	"calligraphy-forum/pkg/errs"
	"calligraphy-forum/pkg/msg"
	"calligraphy-forum/pkg/validate"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"calligraphy-forum/cache"
	"calligraphy-forum/controllers/render"
	"calligraphy-forum/model"
	"calligraphy-forum/services"
)

type UserController struct {
	Ctx iris.Context
}

// 获取当前登录用户
func (c *UserController) GetCurrent() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		return web.JsonData(render.BuildUserProfile(user))
	}
	return web.JsonSuccess()
}

// 用户详情
func (c *UserController) GetBy(userId int64) *web.JsonResult {
	user := cache.UserCache.Get(userId)
	if user != nil && user.Status != constants.StatusDeleted {
		return web.JsonData(render.BuildUserDetail(user))
	}
	return web.JsonErrorMsg("用户不存在")
}

// 修改用户资料
func (c *UserController) PostEditBy(userId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if user.Id != userId {
		return web.JsonErrorMsg("无权限")
	}
	var (
		nickname    = strings.TrimSpace(params.FormValue(c.Ctx, "nickname"))
		homePage    = params.FormValue(c.Ctx, "homePage")
		description = params.FormValue(c.Ctx, "description")
		gender      = strings.TrimSpace(params.FormValue(c.Ctx, "gender"))
		birthdayStr = strings.TrimSpace(params.FormValue(c.Ctx, "birthday"))
		birthday    *time.Time
		err         error
	)

	if len(nickname) == 0 {
		return web.JsonErrorMsg("昵称不能为空")
	}

	if strs.IsNotBlank(gender) {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			return web.JsonErrorMsg("性别数据错误")
		}
	}
	if strs.IsNotBlank(birthdayStr) {
		*birthday, err = dates.Parse(birthdayStr, dates.FmtDate)
		if err != nil {
			return web.JsonError(err)
		}
	}

	if len(homePage) > 0 && validate.IsURL(homePage) != nil {
		return web.JsonErrorMsg("个人主页地址错误")
	}

	columns := map[string]interface{}{
		"nickname":    nickname,
		"home_page":   homePage,
		"description": description,
		"gender":      gender,
	}
	if birthday != nil {
		columns["birthday"] = birthday
	}
	err = services.UserService.Updates(user.Id, columns)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 修改头像
func (c *UserController) PostUpdateAvatar() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	avatar := strings.TrimSpace(params.FormValue(c.Ctx, "avatar"))
	if len(avatar) == 0 {
		return web.JsonErrorMsg("头像不能为空")
	}
	err := services.UserService.UpdateAvatar(user.Id, avatar)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *UserController) PostUpdateNickname() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	nickname := strings.TrimSpace(params.FormValue(c.Ctx, "nickname"))
	if len(nickname) == 0 {
		return web.JsonErrorMsg("Nickname cannot be empty")
	}
	err := services.UserService.UpdateNickname(user.Id, nickname)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *UserController) PostUpdateDescription() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	description := strings.TrimSpace(params.FormValue(c.Ctx, "description"))
	err := services.UserService.UpdateDescription(user.Id, description)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *UserController) PostUpdateGender() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	gender := strings.TrimSpace(params.FormValue(c.Ctx, "gender"))
	err := services.UserService.UpdateGender(user.Id, gender)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *UserController) PostUpdateBirthday() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	birthday := strings.TrimSpace(params.FormValue(c.Ctx, "birthday"))
	err := services.UserService.UpdateBirthday(user.Id, birthday)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

// 设置用户名
func (c *UserController) PostSetUsername() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	username := strings.TrimSpace(params.FormValue(c.Ctx, "username"))
	err := services.UserService.SetUsername(user.Id, username)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 设置邮箱
func (c *UserController) PostSetEmail() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	email := strings.TrimSpace(params.FormValue(c.Ctx, "email"))
	err := services.UserService.SetEmail(user.Id, email)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 设置密码
func (c *UserController) PostSetPassword() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	password := params.FormValue(c.Ctx, "password")
	rePassword := params.FormValue(c.Ctx, "rePassword")
	err := services.UserService.SetPassword(user.Id, password, rePassword)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 修改密码
func (c *UserController) PostUpdatePassword() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	var (
		oldPassword = params.FormValue(c.Ctx, "oldPassword")
		password    = params.FormValue(c.Ctx, "password")
		rePassword  = params.FormValue(c.Ctx, "rePassword")
	)
	if err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 设置背景图
func (c *UserController) PostSetBackgroundImage() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	backgroundImage := params.FormValue(c.Ctx, "backgroundImage")
	if strs.IsBlank(backgroundImage) {
		return web.JsonErrorMsg("请上传图片")
	}
	if err := services.UserService.UpdateBackgroundImage(user.Id, backgroundImage); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 用户收藏
func (c *UserController) GetFavorites() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	// 用户必须登录
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	// 查询列表
	limit := 20
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = services.FavoriteService.Find(sqls.NewCnd().Where("user_id = ? and id < ?",
			user.Id, cursor).Desc("id").Limit(20))
	} else {
		favorites = services.FavoriteService.Find(sqls.NewCnd().Where("user_id = ?", user.Id).Desc("id").Limit(limit))
	}

	hasMore := false
	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].Id
		hasMore = len(favorites) >= limit
	}

	return web.JsonCursorData(render.BuildFavorites(favorites), strconv.FormatInt(cursor, 10), hasMore)
}

// 获取最近3条未读消息
func (c *UserController) GetMsgrecent() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	var count int64 = 0
	var messages []model.Message
	if user != nil {
		count = services.MessageService.GetUnReadCount(user.Id)
		messages = services.MessageService.Find(sqls.NewCnd().Eq("user_id", user.Id).
			Eq("status", msg.StatusUnread).Limit(3).Desc("id"))
	}
	return web.NewEmptyRspBuilder().Put("count", count).Put("messages", render.BuildMessages(messages)).JsonResult()
}

// 用户消息
func (c *UserController) GetMessages() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	// 用户必须登录
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	messages, paging := services.MessageService.FindPageByCnd(sqls.NewCnd().
		Eq("user_id", user.Id).
		Page(page, 20).Desc("id"))

	// 全部标记为已读
	services.MessageService.MarkRead(user.Id)

	return web.JsonPageData(render.BuildMessages(messages), paging)
}

// 用户积分记录
func (c *UserController) GetScorelogs() *web.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)
	user := services.UserTokenService.GetCurrent(c.Ctx)
	// 用户必须登录
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	logs, paging := services.UserScoreLogService.FindPageByCnd(sqls.NewCnd().
		Eq("user_id", user.Id).
		Page(page, 20).Desc("id"))

	return web.JsonPageData(logs, paging)
}

// 积分排行
func (c *UserController) GetScoreRank() *web.JsonResult {
	users := cache.UserCache.GetScoreRank()
	var results []*model.UserInfo
	for _, user := range users {
		results = append(results, render.BuildUserInfo(&user))
	}
	return web.JsonData(results)
}

// 禁言
func (c *UserController) PostForbidden() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}
	var (
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = params.FormValueIntDefault(c.Ctx, "days", 0)
		reason = params.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return web.JsonErrorMsg("请传入：userId")
	}
	if days == -1 && !user.HasRole(constants.RoleOwner) {
		return web.JsonErrorMsg("无永久禁言权限")
	}
	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return web.JsonError(err)
		}
	}
	return web.JsonSuccess()
}

// PostEmailVerify 请求邮箱验证邮件
func (c *UserController) PostSend_verify_email() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if err := services.UserService.SendEmailVerifyEmail(user.Id); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// PostVerify_email 获取邮箱验证码
func (c *UserController) PostVerify_email() *web.JsonResult {
	token := params.FormValue(c.Ctx, "token")
	if strs.IsBlank(token) {
		return web.JsonErrorMsg("Illegal request")
	}
	var (
		email string
		err   error
	)
	if email, err = services.UserService.VerifyEmail(token); err != nil {
		return web.JsonError(err)
	}
	return web.NewEmptyRspBuilder().Put("email", email).JsonResult()
}
