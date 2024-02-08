package gfly

import (
	"app/core/errors"
	"app/core/log"
	"app/core/utils"
	"app/core/validation"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"io"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// ===========================================================================================================
// 												Context
// ===========================================================================================================

// Ctx HTTP request context
type Ctx struct {
	app    *gFly                // Reference to gFly
	root   *fasthttp.RequestCtx // Reference to root request context of fasthttp.
	router *Router              // Reference to Router.
	data   map[string]any       // Keep data in request context.
}

// Root Get original HTTP request context.
func (c *Ctx) Root() *fasthttp.RequestCtx {
	return c.root
}

// Router Get original Router.
func (c *Ctx) Router() *Router {
	return c.router
}

// ===========================================================================================================
// 										Ctx - Handler
// ===========================================================================================================

// RequestHandler A wrapper of fasthttp.RequestHandler
type RequestHandler func(ctx *Ctx) error

// IHandler Interface a handler request.
type IHandler interface {
	Validate(c *Ctx) error
	Handle(c *Ctx) error
}

// Endpoint Default handler
type Endpoint struct{}

func (e *Endpoint) Validate(c *Ctx) error {
	return nil
}

func (e *Endpoint) Handle(c *Ctx) error {
	return nil
}

// Page Abstract web page
type Page struct {
	Endpoint
}

// Api Abstract api
type Api struct {
	Endpoint
}

// ===========================================================================================================
// 										Ctx - Header Data
// ===========================================================================================================

type IHeader interface {
	// Status sets the response's HTTP code.
	Status(code int) *Ctx
	// ContentType sets the response's HTTP content type.
	ContentType(mime string) *Ctx
	// SetHeader sets the response's HTTP header field to the specified key, value.
	SetHeader(key, val string) *Ctx
	// SetCookie set cookie to the response's HTTP header.
	SetCookie(key, value string) *Ctx
	// GetCookie get cookie from the request's HTTP header.
	GetCookie(key string) string
	// GetReqHeaders returns the HTTP request headers.
	GetReqHeaders() map[string][]string
	// Path returns path URI
	Path() string
}

// Status sets the response's HTTP code.
func (c *Ctx) Status(code int) *Ctx {
	c.root.Response.SetStatusCode(code)

	return c
}

// ContentType sets the response's HTTP content type.
func (c *Ctx) ContentType(mime string) *Ctx {
	c.root.Response.Header.SetContentType(mime)

	return c
}

// SetHeader sets the response's HTTP header field to the specified key, value.
func (c *Ctx) SetHeader(key, val string) *Ctx {
	c.root.Response.Header.Set(key, val)

	return c
}

// SetCookie set cookie to the response's HTTP header.
func (c *Ctx) SetCookie(key, value string) *Ctx {
	cook := fasthttp.Cookie{}
	cook.SetKey(key)
	cook.SetValue(value)
	cook.SetMaxAge(3600000)

	c.root.Response.Header.SetCookie(&cook)

	return c
}

// GetCookie get cookie from the request's HTTP header.
func (c *Ctx) GetCookie(key string) string {
	return string(c.root.Request.Header.Cookie(key))
}

// GetReqHeaders returns the HTTP request headers.
// Returned value is only valid within the handler. Do not store any references.
// Make copies or use the Immutable setting instead.
func (c *Ctx) GetReqHeaders() map[string][]string {
	headers := make(map[string][]string)
	c.root.Request.Header.VisitAll(func(k, v []byte) {
		key := utils.UnsafeString(k)
		headers[key] = append(headers[key], utils.UnsafeString(v))
	})

	return headers
}

// Path returns path URI
func (c *Ctx) Path() string {
	return string(c.root.URI().Path())
}

// ===========================================================================================================
// 										Ctx - Request Data
// ===========================================================================================================

type IResponse interface {
	// Success Response success JSON data.
	Success(data interface{}) error
	// Error Response error JSON data.
	Error(data interface{}) error
	// NoContent Response no content.
	NoContent() error
	// View Load template page.
	View(template string, data ViewData) error
	// JSON sets the HTTP response body for JSON type.
	JSON(data JsonData) error
	// HTML sets the HTTP response body for HTML type.
	HTML(body string) error
	// String sets the HTTP response body for String type.
	String(body string) error
	// Raw sets the HTTP response body without copying it.
	Raw(body []byte) error
	// Stream sets response body stream and optional body size.
	Stream(stream io.Reader, size ...int) error
	// Redirect Send redirect.
	Redirect(path string) error
	// Download transfers the file from path as an attachment.
	Download(file string, filename ...string) error
	// File transfers the file from the given path.
	File(file string, compress ...bool) error
}

// Success Response success JSON data.
func (c *Ctx) Success(data interface{}) error {
	c.root.Response.SetStatusCode(StatusOK)
	return c.JSON(data)
}

// Error Response error JSON data.
func (c *Ctx) Error(data interface{}) error {
	c.root.Response.SetStatusCode(StatusBadRequest)
	return c.JSON(data)
}

// NoContent Response no content.
func (c *Ctx) NoContent() error {
	c.root.Response.SetStatusCode(StatusNoContent)

	return nil
}

type ViewData map[string]any

// View Render from template file.
func (c *Ctx) View(template string, data ViewData) error {
	var tplExample = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s.%s", ViewPath, template, ViewExt)))
	pongoContext := pongo2.Context(data)

	err := tplExample.ExecuteWriter(pongoContext, c.root.Response.BodyWriter())
	if err != nil {
		return err
	}

	c.ContentType(MIMETextHTMLCharsetUTF8)

	return nil
}

type JsonData interface{}

// JSON Response json content.
func (c *Ctx) JSON(data JsonData) error {
	c.root.Response.Header.SetContentType(MIMEApplicationJSONCharsetUTF8)

	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.Raw(marshal)
}

// HTML sets the HTTP response body for HTML types.
// This means no type assertion, recommended for faster performance
func (c *Ctx) HTML(body string) error {
	c.root.Response.Header.SetContentType(MIMETextHTMLCharsetUTF8)

	return c.Raw([]byte(body))
}

// String sets the HTTP response body for string types.
// This means no type assertion, recommended for faster performance
func (c *Ctx) String(body string) error {
	c.root.Response.Header.SetContentType(MIMETextPlain)

	return c.Raw([]byte(body))
}

// Raw sets the HTTP response body without copying it.
// From this point onward, the body argument must not be changed.
func (c *Ctx) Raw(body []byte) error {
	c.root.Response.SetBodyRaw(body)

	return nil
}

// Stream sets response body stream and optional body size.
func (c *Ctx) Stream(stream io.Reader, size ...int) error {
	if len(size) > 0 && size[0] >= 0 {
		c.root.Response.SetBodyStream(stream, size[0])
	} else {
		c.root.Response.SetBodyStream(stream, -1)
	}

	return nil
}

// Redirect Send redirect.
func (c *Ctx) Redirect(path string) error {
	c.root.Redirect(path, StatusMovedPermanently)

	return errors.NA
}

// Download transfers the file from path as an attachment.
// Typically, browsers will prompt the user for download.
// By default, the Content-Disposition header filename= parameter is the filepath (this typically appears in the browser dialog).
// Override this default with the filename parameter.
func (c *Ctx) Download(file string, filename ...string) error {
	var fName string

	if len(filename) > 0 {
		fName = filename[0]
	} else {
		fName = filepath.Base(file)
	}
	c.root.Response.Header.Set(HeaderContentDisposition, `attachment; filename="`+utils.QuoteString(fName)+`"`)

	return c.File(file)
}

var (
	sendFileOnce    sync.Once
	sendFileFS      *fasthttp.FS
	sendFileHandler fasthttp.RequestHandler
)

// File transfers the file from the given path.
// The file is not compressed by default, enable this by passing a 'true' argument
// Sets the Content-Type response HTTP header field based on the filename extension.
func (c *Ctx) File(file string, compress ...bool) error {
	// Save the filename, we will need it in the error message if the file isn't found
	filename := file

	// https://github.com/valyala/fasthttp/blob/c7576cc10cabfc9c993317a2d3f8355497bea156/fs.go#L129-L134
	sendFileOnce.Do(func() {
		const cacheDuration = 10 * time.Second
		sendFileFS = &fasthttp.FS{
			Root:                 "",
			AllowEmptyRoot:       true,
			GenerateIndexPages:   false,
			AcceptByteRange:      true,
			Compress:             true,
			CompressedFileSuffix: c.app.config.CompressedFileSuffix,
			CacheDuration:        cacheDuration,
			IndexNames:           []string{"index.html"},
			PathNotFound: func(ctx *fasthttp.RequestCtx) {
				ctx.Response.SetStatusCode(StatusNotFound)
			},
		}
		sendFileHandler = sendFileFS.NewRequestHandler()
	})

	// Disable compression
	if len(compress) == 0 || !compress[0] {
		// https://github.com/valyala/fasthttp/blob/7cc6f4c513f9e0d3686142e0a1a5aa2f76b3194a/fs.go#L55
		c.root.Request.Header.Del(HeaderAcceptEncoding)
	}
	// copy of https://github.com/valyala/fasthttp/blob/7cc6f4c513f9e0d3686142e0a1a5aa2f76b3194a/fs.go#L103-L121 with small adjustments
	if file == "" || !filepath.IsAbs(file) {
		// extend relative path to absolute path
		hasTrailingSlash := len(file) > 0 && (file[len(file)-1] == '/' || file[len(file)-1] == '\\')

		var err error
		file = filepath.FromSlash(file)
		if file, err = filepath.Abs(file); err != nil {
			return fmt.Errorf("failed to determine abs file path: %w", err)
		}
		if hasTrailingSlash {
			file += "/"
		}
	}
	// convert the path to forward slashes regardless the OS in order to set the URI properly
	// the handler will convert back to OS path separator before opening the file
	file = filepath.ToSlash(file)

	// Restore the original requested URL
	originalURL := utils.CopyString(c.OriginalURL())
	defer c.root.Request.SetRequestURI(originalURL)
	// Set new URI for fileHandler
	c.root.Request.SetRequestURI(file)
	// Save status code
	status := c.root.Response.StatusCode()
	// Serve file
	sendFileHandler(c.root)
	// Get the status code which is set by fasthttp
	fsStatus := c.root.Response.StatusCode()
	// Set the status code set by the user if it is different from the fasthttp status code and 200
	if status != fsStatus && status != StatusOK {
		c.Status(status)
	}
	// Check for error
	if status != StatusNotFound && fsStatus == StatusNotFound {
		return errors.FileNotFound{
			FileName: filename,
			Path:     "",
		}
	}

	return nil
}

// Compress Response compressed content with Gzip|Brotli|Deflate (TODO Need more checking)
func (c *Ctx) Compress(body []byte) error {
	ctx := c.root

	switch {
	case ctx.Request.Header.HasAcceptEncodingBytes([]byte(StrGzip)):
		_, err := fasthttp.WriteGzip(ctx.Response.BodyWriter(), body)
		if err != nil {
			return err
		}
	case ctx.Request.Header.HasAcceptEncodingBytes([]byte(StrBrotli)):
		_, err := fasthttp.WriteBrotli(ctx.Response.BodyWriter(), body)
		if err != nil {
			return err
		}
	case ctx.Request.Header.HasAcceptEncodingBytes([]byte(StrDeflate)):
		_, err := fasthttp.WriteDeflate(ctx.Response.BodyWriter(), body)
		if err != nil {
			return err
		}
	default:
		ctx.Response.SetBodyRaw(body)
	}

	return nil
}

// ===========================================================================================================
// 					Ctx - Request Data (Form|MultipartForm|FormFile, Query, Path, RAW)
// ===========================================================================================================

// UploadedFile uploaded file info.
type UploadedFile struct {
	Field string // Field File name in form
	Name  string // Name Uploaded file name
	Path  string // Path Uploaded file path
	Size  int64  // Size Uploaded file size
}

type IRequestData interface {
	// ParseBody Parse body to struct data type.
	ParseBody(data any) error
	// ParseQuery Parse query string to struct data type.
	ParseQuery(data any) error
	// FormVal Get data from POST|PUT request.
	FormVal(key string) []byte
	// FormInt Get int from POST|PUT request.
	FormInt(key string) (int, error)
	// FormBool Get bool from POST|PUT request.
	FormBool(key string) (bool, error)
	// FormFloat Get float from POST|PUT request.
	FormFloat(key string) (float64, error)
	// FormUpload Process and get uploaded files from POST|PUT request.
	FormUpload(files ...string) ([]UploadedFile, error)
	// Queries Get data from Query string.
	Queries() map[string]string
	// QueryStr Get data from Query string.
	QueryStr(key string) string
	// QueryInt Get int from Query string.
	QueryInt(key string) (int, error)
	// QueryBool Get bool from Query string.
	QueryBool(key string) (bool, error)
	// QueryFloat Get float from Query string.
	QueryFloat(key string) (float64, error)
	// PathVal Get data from Path request.
	PathVal(key string) string
	// OriginalURL Original URL.
	OriginalURL() string
}

// ParseBody Parse body JSON data to struct data type.
func (c *Ctx) ParseBody(data any) error {
	jsonData := c.root.PostBody()

	err := json.Unmarshal(jsonData, data)
	if err != nil {
		return err
	}

	return nil
}

// ParseQuery Parse query string to struct data type.
func (c *Ctx) ParseQuery(out interface{}) error {
	log.Error("===> Not yet implemented <===")

	return nil
}

// FormVal Get data from POST|PUT request.
func (c *Ctx) FormVal(key string) []byte {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return data
}

// FormInt Get int from POST|PUT request.
func (c *Ctx) FormInt(key string) (int, error) {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.Atoi(string(data))
}

// FormBool Get bool from POST|PUT request.
func (c *Ctx) FormBool(key string) (bool, error) {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseBool(string(data))
}

// FormFloat Get float from POST|PUT request.
func (c *Ctx) FormFloat(key string) (float64, error) {
	data := c.root.PostArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseFloat(string(data), 64)
}

// FormUpload Process and get uploaded files from POST|PUT request.
func (c *Ctx) FormUpload(files ...string) ([]UploadedFile, error) {
	var uploadedFiles []UploadedFile

	if len(files) > 0 {
		for _, file := range files {
			// Read file header
			header, err := c.root.FormFile(file)
			if err != nil {
				return uploadedFiles, err
			}

			// Create temporary file.
			tempName := fmt.Sprintf("%s.%s", uuid.New(), fileExt(header.Filename))
			filePath := fmt.Sprintf("%s/%s", TemporaryDir, tempName)

			// Save file
			err = fasthttp.SaveMultipartFile(header, filePath)
			if err != nil {
				return uploadedFiles, err
			}

			uploadedFiles = append(uploadedFiles, UploadedFile{
				Name:  header.Filename,
				Path:  filePath,
				Size:  header.Size,
				Field: file,
			})
		}
	} else {
		form, err := c.root.MultipartForm()
		if err != nil {
			return nil, err
		}

		for _, v := range form.File {
			for _, header := range v {
				// Create temporary file.
				tempName := fmt.Sprintf("%s.%s", uuid.New(), fileExt(header.Filename))
				filePath := fmt.Sprintf("%s/%s", TemporaryDir, tempName)

				err = fasthttp.SaveMultipartFile(header, filePath)
				if err != nil {
					return uploadedFiles, err
				}

				uploadedFiles = append(uploadedFiles, UploadedFile{
					Name:  header.Filename,
					Path:  filePath,
					Size:  header.Size,
					Field: "N/A",
				})
			}
		}
	}

	return uploadedFiles, nil
}

// Queries Get data from Query string.
func (c *Ctx) Queries() map[string]string {
	m := make(map[string]string, c.root.QueryArgs().Len())
	c.root.QueryArgs().VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})
	return m
}

// QueryStr Get data from Query string.
func (c *Ctx) QueryStr(key string) string {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return string(data)
}

// QueryInt Get int from Query string.
func (c *Ctx) QueryInt(key string) (int, error) {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.Atoi(string(data))
}

// QueryBool Get bool from Query string.
func (c *Ctx) QueryBool(key string) (bool, error) {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseBool(string(data))
}

// QueryFloat Get float from Query string.
func (c *Ctx) QueryFloat(key string) (float64, error) {
	data := c.root.QueryArgs().Peek(key)
	if data == nil {
		data = c.root.FormValue(key)
	}

	return strconv.ParseFloat(string(data), 64)
}

// PathVal Get data from Path request.
func (c *Ctx) PathVal(key string) string {
	val := c.root.UserValue(key)

	if val == nil {
		return ""
	}

	return val.(string)
}

// OriginalURL contains the original request URL.
// Returned value is only valid within the handler. Do not store any references.
// Make copies or use the Immutable setting to use the value outside the Handler.
func (c *Ctx) OriginalURL() string {
	return string(c.root.Request.Header.RequestURI())
}

// ===========================================================================================================
// 											Ctx - Data
// ===========================================================================================================

type IData interface {
	// Validate Validate data struct type.
	Validate(structData interface{}, msgForTag ...validation.MsgForTagFunc) (map[string][]string, error)
	// SetData Keep data in request context Ctx.
	SetData(key string, data interface{})
	// GetData Get data from request context Ctx.
	GetData(key string) interface{}
	// SetSession Keep data in request context Ctx.
	SetSession(key string, data interface{})
	// GetSession Get data from request context Ctx.
	GetSession(key string) interface{}
}

// Validate Validate data struct type.
func (c *Ctx) Validate(structData interface{}, msgForTagFunc ...validation.MsgForTagFunc) (map[string][]string, error) {
	// Default message tag function.
	fn := validation.MsgForTag

	if len(msgForTagFunc) > 0 {
		fn = msgForTagFunc[0]
	}

	return validation.Validate(structData, fn)
}

// SetData Keep data in request context Ctx.
func (c *Ctx) SetData(key string, data interface{}) {
	c.data[key] = data
}

// GetData Get data from request context Ctx.
func (c *Ctx) GetData(key string) interface{} {
	return c.data[key]
}

// SetSession Keep data in session.
func (c *Ctx) SetSession(key string, data interface{}) {
	store, err := serverSession.Get(c.root)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := serverSession.Save(c.root, store); err != nil {
			log.Fatal(err)
		}
	}()

	store.Set(key, data)
}

// GetSession Get data from session.
func (c *Ctx) GetSession(key string) interface{} {
	store, err := serverSession.Get(c.root)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := serverSession.Save(c.root, store); err != nil {
			log.Fatal(err)
		}
	}()

	return store.Get(key)
}
