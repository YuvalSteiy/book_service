package dal

type UserDater interface {
	InitClient()
	AddActivity(username string, req string)
	CreateNewClientActivity(username string, req string)
	UpdateClientActivity(username string, req string)
	GetUserActivity(username string) ([]interface{})
}

func NewUserDater() UserDater {
	var r Redis
	r.InitClient()
	return &r
}
