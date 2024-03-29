package spam

import (
	"calligraphy-forum/model"
	"calligraphy-forum/pkg/errs"
	"calligraphy-forum/services"

	"github.com/dchest/captcha"
)

type CaptchaStrategy struct{}

func (CaptchaStrategy) Name() string {
	return "CaptchaStrategy"
}

func (CaptchaStrategy) CheckTopic(user *model.User, form model.CreateTopicForm) error {
	if services.SysConfigService.GetConfig().TopicCaptcha && !captcha.VerifyString(form.CaptchaId, form.CaptchaCode) {
		return errs.CaptchaError
	}
	return nil
}

func (CaptchaStrategy) CheckArticle(user *model.User, form model.CreateArticleForm) error {
	return nil
}

func (CaptchaStrategy) CheckComment(user *model.User, form model.CreateCommentForm) error {
	return nil
}
