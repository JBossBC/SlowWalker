package infrastructure

type Mobile interface {
	Send(phone string, code string) error
}

func GetMobileProvider() Mobile {
	return getMobileProvider()
}
