package models

type Response struct {
	Error string            `json:"error"`
	Msg   string            `json:"msg"`
	Data  map[string]string `json:"data"`
}
