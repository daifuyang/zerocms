package types

func (base *Rest) Success(msg string) (res Rest) {
	res.Code = 1
	res.Msg = msg
	return
}

func (base *Rest) Error(msg string) (res Rest) {
	res.Code = 0
	res.Msg = msg
	return
}
