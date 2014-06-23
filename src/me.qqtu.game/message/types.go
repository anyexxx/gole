package message

const (
	MSG_TEST     uint16 = 0
	MSG_CG_LOGIN uint16 = 1
	MSG_GC_LOGIN uint16 = 2
	//notifier
	MSG_NOTI_SERVER_LOGIN  uint16 = 10001
	MSG_NOTI_SERVER_LOGOUT uint16 = 10002
	MSG_NOTI_USER_LOGIN    uint16 = 10003
	MSG_NOTI_USER_LOGOUT   uint16 = 10004
	MSG_NOTI_SENDUSERMSG   uint16 = 10005
	MSG_NOTI_BROADCAST     uint16 = 10006
)

type CGLoginData struct {
	Account string
	Pwd     string
}

type GCLoginData struct {
	Ret int
}

type NotiServerLogin struct {
	SId int
}

type NotiServerLogout struct {
	SId int
}

type NotiUserLogin struct {
	UId int64
	SId int
}

type NotiUserLogout struct {
	UId int64
	SId int
}
