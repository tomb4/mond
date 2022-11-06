package constant

const (
	BalanceTargetHostName = "balance_target_host_name"
	BalanceTargetLabel    = "balance_target_label"
	ContextCacheKey       = "context_cache_key"
	
	TraceIdKey = "_traceId"
	
	PublishMdCtxKey   = "_publish_md_ctx"
	ConsumerMdCtxKey  = "_consumer_md_ctx"
	MqExchangeKey     = "_mq_exchange"
	AsyncMethodCtxKey = "_async_method_ctx"
	
	UberTraceIdKey = "uber-trace-id"
	
	MongoDbKey         = "_mongo_db"
	MongoCollectionKey = "_mongo_collection"
	MongoOperationKey  = "_mongo_operation"
	
	RedisCmdKey = "_mongo_cmd"
	
	HttpClientHost   = "_http_client_host"
	HttpClientPath   = "_http_client_path"
	HttpClientMethod = "_http_client_method"
	
	GrpcClientAddr = "_grpc_client_addr"
	
	SentinelBreaker = "_sentinel_breaker_scope"
)
