package define



type SubscribeEntity interface {
	Unsubscribe()
}


type Subscribe func([]*Instance)
