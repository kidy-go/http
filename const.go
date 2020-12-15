package http

const (
	ProtocolTLS    = "h2"
	ProtocolTCP    = "h2c"
	ClientPreface  = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
	defaultBufSize = 4096
	minReadBufSize = len(ClientPreface)

	maxConcurrentStreams   = 1<<31 - 1
	maxInitialWindowSize   = 1<<31 - 1
	maxFrameSizeLowerBound = 1 << 14
	maxFrameSizeUpperBound = 1<<24 - 1

	defaultHeaderTableSize      = 4096
	defaultEnablePush           = 1
	defaultMaxConcurrentStreams = maxConcurrentStreams
	defaultInitialWindowSize    = 65535
	defaultMaxFrameSize         = maxFrameSizeLowerBound
)

type ErrCode uint32

const (
	ErrCodeNo                 ErrCode = 0x0
	ErrCodeProtocol           ErrCode = 0x1
	ErrCodeInternal           ErrCode = 0x2
	ErrCodeFlowControl        ErrCode = 0x3
	ErrCodeSettingsTimeout    ErrCode = 0x4
	ErrCodeStreamClosed       ErrCode = 0x5
	ErrCodeFrameSize          ErrCode = 0x6
	ErrCodeRefusedStream      ErrCode = 0x7
	ErrCodeCancel             ErrCode = 0x8
	ErrCodeCompression        ErrCode = 0x9
	ErrCodeConnect            ErrCode = 0xa
	ErrCodeEnhanceYourCalm    ErrCode = 0xb
	ErrCodeInadequateSecurity ErrCode = 0xc
	ErrCodeHTTP11Required     ErrCode = 0xd
)
