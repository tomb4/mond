package dynamicConfig

var (
	dynamicConfigInstance *DynamicConfig = &DynamicConfig{}
)

func GetDynamicConfig() *DynamicConfig {
	return dynamicConfigInstance
}

func SetDynamicConfig(instance *DynamicConfig) {
	dynamicConfigInstance = instance
}

type DynamicConfig struct {
	Id string `json:"id"`
}
