package model

import (
	"calligraphy-forum/model/constants"
	"calligraphy-forum/pkg/common"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type CreateTopicForm struct {
	Type        constants.TopicType
	CaptchaId   string
	CaptchaCode string
	NodeId      int64
	Title       string
	Content     string
	HideContent string
	Tags        []string
	ImageList   []ImageDTO
	UserAgent   string
	Ip          string
}

type CreateArticleForm struct {
	Title       string
	Summary     string
	Content     string
	ContentType string
	Cover       *ImageDTO
	Tags        []string
	SourceUrl   string
}

// CreateCommentForm 发表评论
type CreateCommentForm struct {
	EntityType string     `form:"entityType"`
	EntityId   int64      `form:"entityId"`
	Content    string     `form:"content"`
	ImageList  []ImageDTO `form:"imageList"`
	QuoteId    int64      `form:"quoteId"`
	UserAgent  string     `form:"userAgent"`
	Ip         string     `form:"ip"`
}

type ImageDTO struct {
	Url string `json:"url"`
}

func GetCreateTopicForm(ctx iris.Context) CreateTopicForm {
	var (
		topicType = params.FormValueIntDefault(ctx, "type", int(constants.TopicTypeTopic))
	)
	return CreateTopicForm{
		Type:        constants.TopicType(topicType),
		CaptchaId:   params.FormValue(ctx, "captchaId"),
		CaptchaCode: params.FormValue(ctx, "captchaCode"),
		NodeId:      params.FormValueInt64Default(ctx, "nodeId", 0),
		Title:       strings.TrimSpace(params.FormValue(ctx, "title")),
		Content:     strings.TrimSpace(params.FormValue(ctx, "content")),
		HideContent: strings.TrimSpace(params.FormValue(ctx, "hideContent")),
		Tags:        params.FormValueStringArray(ctx, "tags"),
		ImageList:   GetImageList(ctx, "imageList"),
		UserAgent:   common.GetUserAgent(ctx.Request()),
		Ip:          common.GetRequestIP(ctx.Request()),
	}
}

func GetCreateCommentForm(ctx iris.Context) CreateCommentForm {
	form := CreateCommentForm{
		EntityType: params.FormValue(ctx, "entityType"),
		EntityId:   params.FormValueInt64Default(ctx, "entityId", 0),
		Content:    strings.TrimSpace(params.FormValue(ctx, "content")),
		ImageList:  GetImageList(ctx, "imageList"),
		QuoteId:    params.FormValueInt64Default(ctx, "quoteId", 0),
		UserAgent:  common.GetUserAgent(ctx.Request()),
		Ip:         common.GetRequestIP(ctx.Request()),
	}
	return form
}

func GetCreateArticleForm(ctx iris.Context) CreateArticleForm {
	var (
		title   = ctx.PostValue("title")
		summary = ctx.PostValue("summary")
		content = ctx.PostValue("content")
		tags    = params.FormValueStringArray(ctx, "tags")
		cover   = GetImageDTO(ctx, "cover")
	)
	return CreateArticleForm{
		Title:       title,
		Summary:     summary,
		Content:     content,
		ContentType: constants.ContentTypeMarkdown,
		Cover:       cover,
		Tags:        tags,
	}
}

func GetImageList(ctx iris.Context, paramName string) []ImageDTO {
	imageListStr := params.FormValue(ctx, paramName)
	var imageList []ImageDTO
	if strs.IsNotBlank(imageListStr) {
		ret := gjson.Parse(imageListStr)
		if ret.IsArray() {
			for _, item := range ret.Array() {
				url := item.Get("url").String()
				imageList = append(imageList, ImageDTO{
					Url: url,
				})
			}
		}
	}
	return imageList
}

func GetImageDTO(ctx iris.Context, paramName string) (img *ImageDTO) {
	str := params.FormValue(ctx, paramName)
	if strs.IsBlank(str) {
		return
	}
	if err := jsons.Parse(str, &img); err != nil {
		logrus.Error(err)
	}
	return
}
