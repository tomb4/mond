package constant

const (
	ConvTypeSingle   = 1
	ConvTypeTeam     = 2
	ConvTypeChatRoom = 3

	ActionJoin   = 1 //加入某个会话
	ActionBlack  = 2
	ActionForbid = 3

	ServerSenderId = 1
)

// 服务端下发消息
const (
	ServerMsgTypeApplyFriend = 1001 // 好友申请
	ServerMsgTypeSendGift    = 1002 // 礼物消息
)
