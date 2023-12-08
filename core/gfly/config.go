package gfly

import "time"

type Config struct {
	// App name
	//
	// Default app name is used if left blank.
	AppName string

	// Server name for sending in response headers.
	//
	// Default server name is used if left blank.
	Name string

	// The maximum number of concurrent connections the server may serve.
	//
	// DefaultConcurrency is used if not set.
	//
	// Concurrency only works if you either call Serve once, or only ServeConn multiple times.
	// It works with ListenAndServe as well.
	Concurrency int

	// ReadTimeout is the amount of time allowed to read
	// the full request including body. The connection's read
	// deadline is reset when the connection opens, or for
	// keep-alive connections after the first byte has been read.
	//
	// By default request read timeout is unlimited.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset after the request handler
	// has returned.
	//
	// By default response write timeout is unlimited.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alive is enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used.
	IdleTimeout time.Duration

	// Per-connection buffer size for requests' reading.
	// This also limits the maximum header size.
	//
	// Increase this buffer if your clients send multi-KB RequestURIs
	// and/or multi-KB headers (for example, BIG cookies).
	//
	// Default buffer size is used if not set.
	ReadBufferSize int

	// Per-connection buffer size for responses' writing.
	//
	// Default buffer size is used if not set.
	WriteBufferSize int

	// NoDefaultDate, when set to true, causes the default Date
	// header to be excluded from the Response.
	//
	// The default Date header value is the current date value. When
	// set to true, the Date will not be present.
	NoDefaultDate bool

	// NoDefaultContentType, when set to true, causes the default Content-Type
	// header to be excluded from the Response.
	//
	// The default Content-Type header value is the internal default value. When
	// set to true, the Content-Type will not be present.
	NoDefaultContentType bool

	// Header names are passed as-is without normalization
	// if this option is set.
	//
	// Disabled header names' normalization may be useful only for proxying
	// incoming requests to other servers expecting case-sensitive
	// header names. See https://github.com/valyala/fasthttp/issues/57
	// for details.
	//
	// By default request and response header names are normalized, i.e.
	// The first letter and the first letters following dashes
	// are uppercased, while all the other letters are lowercased.
	// Examples:
	//
	//     * HOST -> Host
	//     * content-type -> Content-Type
	//     * cONTENT-lenGTH -> Content-Length
	DisableHeaderNamesNormalizing bool

	// Whether to disable keep-alive connections.
	//
	// The server will close all the incoming connections after sending
	// the first response to client if this option is set to true.
	//
	// By default keep-alive connections are enabled.
	DisableKeepalive bool

	// Maximum request body size.
	// a zero value means that default values will be honored
	MaxRequestBodySize int

	// NoDefaultServerHeader, when set to true, causes the default Server header
	// to be excluded from the Response.
	//
	// The default Server header value is the value of the Name field or an
	// internal default value in its absence. With this option set to true,
	// the only time a Server header will be sent is if a non-zero length
	// value is explicitly provided during a request.
	NoDefaultServerHeader bool

	// Rejects all non-GET requests if set to true.
	//
	// This option is useful as anti-DoS protection for servers
	// accepting only GET requests and HEAD requests. The request size is limited
	// by ReadBufferSize if GetOnly is set.
	//
	// Server accepts all the requests by default.
	GetOnly bool

	// Aggressively reduces memory usage at the cost of higher CPU usage
	// if set to true.
	//
	// Try enabling this option only if the server consumes too much memory
	// serving mostly idle keep-alive connections. This may reduce memory
	// usage by more than 50%.
	//
	// Aggressive memory usage reduction is disabled by default.
	ReduceMemoryUsage bool

	// StreamRequestBody enables request body streaming,
	// and calls the handler sooner when given body is
	// larger than the current limit.
	StreamRequestBody bool

	// Will not pre parse Multipart Form data if set to true.
	//
	// This option is useful for servers that desire to treat
	// multipart form data as a binary blob, or choose when to parse the data.
	//
	// Server pre parses multipart form data by default.
	DisablePreParseMultipartForm bool

	// When set to true, it will not print out the «gFly» ASCII art and listening address.
	//
	// Default: false
	DisableStartupMessage bool

	// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only)
	// WARNING: When prefork is set to true, only "tcp4" and "tcp6" can be chose.
	//
	// Default: NetworkTCP4
	Network string

	// When set to true, this will spawn multiple Go processes listening on the same port.
	//
	// Default: false
	Prefork bool

	// CompressedFileSuffix adds suffix to the original file name and
	// tries saving the resulting compressed file under the new file name.
	//
	// Default: ".gz"
	CompressedFileSuffix string
}

var DefaultConfig = Config{
	AppName:                       "Laravel inspired web framework written in Go",
	Name:                          "gFly",
	Concurrency:                   256 * 1024,
	ReadTimeout:                   60 * 60 * time.Second, // unlimited
	WriteTimeout:                  60 * 60 * time.Second, // unlimited
	IdleTimeout:                   60 * 60 * time.Second, // unlimited
	ReadBufferSize:                4096 * 10,
	WriteBufferSize:               4096 * 10,
	NoDefaultDate:                 false,
	NoDefaultContentType:          false,
	DisableHeaderNamesNormalizing: false,
	DisableKeepalive:              false,
	MaxRequestBodySize:            4 * 1024 * 1024,
	NoDefaultServerHeader:         false, // True when `Name` Empty
	GetOnly:                       false,
	ReduceMemoryUsage:             false,
	StreamRequestBody:             true,
	DisablePreParseMultipartForm:  true,
	DisableStartupMessage:         false,
	Network:                       NetworkTCP4,
	Prefork:                       false,
	CompressedFileSuffix:          ".gz",
}
