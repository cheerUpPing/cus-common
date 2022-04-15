package c_common

type Msg struct {
	Code string
	Msg  string
	Data interface{}
}

func Success() Msg {
	msg := Msg{
		Code: "00000000",
		Msg:  "success",
		Data: nil,
	}
	return msg
}

func SuccessData(data interface{}) Msg {
	msg := Msg{
		Code: "00000000",
		Msg:  "success",
		Data: data,
	}
	return msg
}

func Error() Msg {
	msg := Msg{
		Code: "00000001",
		Msg:  "system error, please try later",
		Data: nil,
	}
	return msg
}
