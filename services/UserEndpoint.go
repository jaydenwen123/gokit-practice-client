package services


type UserRequest struct {
	Uid    int64  `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}






