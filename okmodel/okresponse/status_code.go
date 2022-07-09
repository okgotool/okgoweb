package okresponse

// response code:
const (
	// OKCode :
	OKCode int = 200
	// FailedCode :
	FailedCode int = 210
	// AuthFailedCode :
	AuthFailedCode int = 401

	// Status codes:
	// StatusOK :
	StatusOK StatusCode = 200

	// StatusBadRequest
	StatusBadRequest       StatusCode = 400
	StatusUnauthorized     StatusCode = 401
	StatusForbidden        StatusCode = 403
	StatusNotFound         StatusCode = 404
	StatusMethodNotAllowed StatusCode = 405
	StatusRequestTimeout   StatusCode = 408

	StatusInternalServerError StatusCode = 500
	StatusNotImplemented      StatusCode = 501
	StatusBadGateway          StatusCode = 502
	StatusServiceUnavailable  StatusCode = 503
	StatusGatewayTimeout      StatusCode = 504
)
