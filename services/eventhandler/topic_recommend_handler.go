package eventhandler

import (
	"calligraphy-forum/model/constants"
	"calligraphy-forum/pkg/event"
	"calligraphy-forum/pkg/msg"
	"calligraphy-forum/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicRecommendEvent{}), handleTopicRecommend)
}

func handleTopicRecommend(i interface{}) {
	e := i.(event.TopicRecommendEvent)

	if e.Recommend {
		sendTopicRecommendMsg(e.TopicId)
	}
}

// sendTopicRecommendMsg 话题被设为推荐
func sendTopicRecommendMsg(topicId int64) {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return
	}
	var (
		from         int64 = 0
		to                 = topic.UserId
		title              = "你的话题被设为推荐"
		quoteContent       = "《" + topic.GetTitle() + "》"
	)
	services.MessageService.SendMsg(from, to, msg.TypeTopicRecommend, title, "", quoteContent,
		&msg.TopicRecommendExtraData{
			TopicId: topicId,
		})
}
