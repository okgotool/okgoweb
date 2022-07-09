package okresponse

// common response structs:
type (
	// StatusCode : enum
	StatusCode int

	// Success :
	Success struct {
		Code StatusCode `json:"code"`
		Msg  string     `json:"msg"`
	}

	// SuccessUint :
	SuccessUint struct {
		Code StatusCode `json:"code"`
		Msg  string     `json:"msg"`
		Data uint       `json:"data"`
	}

	// SuccessUint64 :
	SuccessUint64 struct {
		Code StatusCode `json:"code"`
		Msg  string     `json:"msg"`
		Data uint64     `json:"data"`
	}
	// SuccessMap :
	SuccessMap struct {
		Code StatusCode        `json:"code"`
		Msg  string            `json:"msg"`
		Data map[string]string `json:"data"`
	}

	// SuccessMapInt :
	SuccessMapInt struct {
		Code StatusCode       `json:"code"`
		Msg  string           `json:"msg"`
		Data map[string]int64 `json:"data"`
	}

	// SuccessMapArray :
	SuccessMapArray struct {
		Code StatusCode          `json:"code"`
		Msg  string              `json:"msg"`
		Data []map[string]string `json:"data"`
	}
	// ResponseIntMap :
	ResponseIntMap struct {
		Code StatusCode  `json:"code"`
		Msg  string      `json:"msg"`
		Data map[int]int `json:"data"`
	}

	// SuccessArray :
	SuccessArray struct {
		Code StatusCode `json:"code"`
		Msg  string     `json:"msg"`
		Data []string   `json:"data"`
	}

	// Failed :
	Failed struct {
		Code  StatusCode `json:"code"`
		Msg   string     `json:"msg"`
		Error string     `json:"error"`
	}

	ResponseInterface struct {
		Code StatusCode  `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
)
