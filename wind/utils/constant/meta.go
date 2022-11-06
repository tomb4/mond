package constant

const (
	//心跳包间隔
	PingPongInterval = 30000 //ms
	//僵死连接间隔
	DeadConnectionInterval = PingPongInterval * 2
	//如果用户掉线超过间隔，则会被判定为下线
	LogoutConnectionInterval = PingPongInterval * 3
	
	ImConnection    = 2 //im连接
	SceneConnection = 1 //场景连接
	
	Login   = 1 // 在线
	Logout  = 2 //退出
	Offline = 3 //掉线
	
	LoginAction   = "login"
	OfflineAction = "offline"
	LeaveAction   = "leave"
	EnterAction   = "enter"
	SitDownAction = "sitDown"
	GetUpAction   = "getUp"
	
	//场景引擎与游戏引擎
	SceneEngine = 1
	//GameEngine  = 2
	//WoodenGame  = 3
	//GoldCoin    = 4
	PassThrough = 5
	
	//DefaultGameLevel            = 5
	//DefaultWoodenGameCapacity   = 100
	//DefaultGoldCoinGameCapacity = 4
	
	GameStagePrepare = 1 //准备阶段
	GameStageRun     = 2 //运行阶段
	GameStageOver    = 3 //结束
	
	//GamePrepareTime = 60 * 1000 //游戏准备时间 ms
	GamePrepareTime = 15 * 1000 //游戏准备时间 ms
)
