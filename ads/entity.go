package ads

type AdState string

const (
	SUSPENDED         AdState = "SUSPENDED"         // 	Показы объявления остановлены владельцем с помощью метода suspend или в веб-интерфейсе.
	OFF_BY_MONITORING AdState = "OFF_BY_MONITORING" // 	Показы объявления автоматически остановлены мониторингом доступности сайта.
	ON                AdState = "ON"                // 	Объявление активно, принадлежит к активной кампании и может быть показано (при наличии средств на кампании, в соответствии с настройками временного таргетинга и т. п.).
	OFF               AdState = "OFF"               // 	Объявление неактивно (черновик, ожидает модерации, отклонено) или принадлежит к неактивной либо остановленной кампании.
	ARCHIVED          AdState = "ARCHIVED"          // 	Объявление помещено в архив (с помощью метода archive или пользователем в веб-интерфейсе) или принадлежит к архивной кампании.
)

type AdStatus string

const (
	DRAFT       AdStatus = "DRAFT"       // Объявление создано и еще не отправлено на модерацию.
	MODERATION  AdStatus = "MODERATION"  // 	Объявление находится на модерации.
	PREACCEPTED AdStatus = "PREACCEPTED" // 	Объявление допущено к показам автоматически, но будет дополнительно проверено модератором.
	ACCEPTED    AdStatus = "ACCEPTED"    // 	Объявление принято модерацией.
	REJECTED    AdStatus = "REJECTED"    // 	Объявление отклонено модерацией.
)

type AdGroupType string

const (
	TEXT_AD_GROUP         AdGroupType = "TEXT_AD_GROUP"
	SMART_AD_GROUP        AdGroupType = "SMART_AD_GROUP"
	MOBILE_APP_AD_GROUP   AdGroupType = "MOBILE_APP_AD_GROUP"
	DYNAMIC_TEXT_AD_GROUP AdGroupType = "DYNAMIC_TEXT_AD_GROUP"
	CPM_BANNER_AD_GROUP   AdGroupType = "CPM_BANNER_AD_GROUP"
	CPM_VIDEO_AD_GROUP    AdGroupType = "CPM_VIDEO_AD_GROUP"
)

type AdType string

const (
	TEXT_AD                            AdType = "TEXT_AD"
	IMAGE_AD                           AdType = "IMAGE_AD"
	TEXT_IMAGE_AD                      AdType = "TEXT_IMAGE_AD"
	TEXT_AD_BUILDER_AD                 AdType = "TEXT_AD_BUILDER_AD"
	CPC_VIDEO_AD                       AdType = "CPC_VIDEO_AD"
	SMART_AD                           AdType = "SMART_AD"
	MOBILE_APP_AD                      AdType = "MOBILE_APP_AD"
	MOBILE_APP_IMAGE_AD                AdType = "MOBILE_APP_IMAGE_AD"
	MOBILE_APP_AD_BUILDER_AD           AdType = "MOBILE_APP_AD_BUILDER_AD"
	MOBILE_APP_CPC_VIDEO_AD_BUILDER_AD AdType = "MOBILE_APP_CPC_VIDEO_AD_BUILDER_AD"
	DYNAMIC_TEXT_AD                    AdType = "DYNAMIC_TEXT_AD"
	CPM_BANNER_AD                      AdType = "CPM_BANNER_AD"
	CPM_VIDEO_AD                       AdType = "CPM_VIDEO_AD"
)
