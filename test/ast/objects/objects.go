package objects

type OutsideFileResponse2 struct {
	V string `json:"v"`
}

type OutsideFileResponse struct {
	Example string                 `json:"example"`
	Wuwu    OutsideFileResponse2   `json:"wuwu"`
	WuwuArr []OutsideFileResponse2 `json:"wuwu_arr"`
}
