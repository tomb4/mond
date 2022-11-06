package env

type State string

const (
	Init     State = "init"     //框架初始化中
	Starting State = "starting" //业务代码启动中
	Running  State = "running"  //运行中
	Stopping State = "stopping" //停止中

	defaultCluster string = "master" //默认集群
)

var (
	StateEnv State  //服务状态
	cluster  string //分区环境
	appId    string
	env      string
	port     int32
	podIp    string
	nsId     string
	serverIp string
	hostName string
)

func SetEnv(v string) {
	env = v
}
func GetEnv() string {
	return env
}

func InLocal() bool {
	return env == ""
}
func InDevelop() bool {
	return env == "develop"
}
func InTest() bool {
	return env == "test"
}
func InProd() bool {
	return env == "prod"
}

func SetAppId(id string) {
	appId = id
}
func GetAppId() string {
	return appId
}
func SetCluster(v string) {
	cluster = v
}
func GetCluster() string {
	if cluster == "" {
		if InLocal() || InDevelop() {
			return "shanghai"
		}
		return defaultCluster
	}
	return cluster
}

func GetAppState() State {
	return StateEnv
}

func SetPort(v int32) {
	port = v
}
func GetPort() int32 {
	return port
}

func SetPodIp(v string) {
	podIp = v
}
func GetPodIp() string {
	return podIp
}

func SetNamespaceId(v string) {
	nsId = v
}

func GetNamespaceId() string {
	return nsId
}

func SetServerIp(v string) {

	serverIp = v
}
func GetServerIp() string {
	return serverIp
}

func SetHostName(v string) {
	hostName = v
}

func GetHostName() string {
	return hostName
}
