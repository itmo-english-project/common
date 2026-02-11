package httperrors

type Error struct {
	Message string      `json:"error_message"`
	Code    string      `json:"error_code"`
	Details interface{} `json:"error_details"`
}
