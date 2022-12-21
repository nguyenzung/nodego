package eventloop

import "time"

var (
	API_NUM_THREAD                      = 32
	TASK_NUM_THREAD                     = 32
	HTTP_IP                             = "0.0.0.0"
	HTTP_PORT                           = 9090
	WEBSOCKET_IP                        = "0.0.0.0"
	WEBSOCKET_PORT                      = 9091
	WS_READ_BUFFER_SIZE                 = 1024
	WS_WRITE_BUFFER_SIZE                = 1024
	WS_TIMEOUT_IN_SECONDS time.Duration = 3000
)
