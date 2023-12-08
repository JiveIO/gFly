/*
Package gothic wraps common behaviour when using Goth. This makes it quick, and easy, to get up
and running with Goth. Of course, if you want complete control over how things flow, in regard
to the authentication process, feel free and use Goth directly.

See https://github.com/markbates/goth/blob/master/examples/main.go to see this in action.
*/
package gothic

import (
	"app/core/log"
	"app/core/oauth"
	"app/core/session"
	"app/core/session/providers/redis"
	"app/core/utils"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SessionName is the key used to access the session store.
var SessionName = "gfly_oauth"

// Store can/should be set by applications using gothic. The default is a cookie store.
var defaultSession *session.Session

type key int

// ProviderParamKey can be used as a key in context when passing in a provider
const ProviderParamKey key = iota

func init() {
	// Create session
	provider, _ := redis.New(redis.Config{
		KeyPrefix:       utils.Getenv("SESSION_OAUTH_KEY", SessionName),
		Addr:            utils.Getenv("SESSION_REDIS_URL", "127.0.0.1:6379"),
		PoolSize:        8,
		ConnMaxIdleTime: 30 * time.Second,
	})

	cfg := session.NewDefaultConfig()
	cfg.EncodeFunc = session.MSGPEncode
	cfg.DecodeFunc = session.MSGPDecode
	defaultSession = session.New(cfg)

	if err := defaultSession.SetProvider(provider); err != nil {
		log.Fatal(err)
	}
}

/*
BeginAuthHandler is a convenience handler for starting the authentication process.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func BeginAuthHandler(ctx *fasthttp.RequestCtx) {
	authURL, err := GetAuthURL(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
	}

	ctx.Redirect(authURL, http.StatusTemporaryRedirect)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx *fasthttp.RequestCtx) string {
	state := GetState(ctx)
	if len(state) > 0 {
		return state
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("gothic: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx *fasthttp.RequestCtx) string {
	return utils.UnsafeString(ctx.QueryArgs().Peek("state"))
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx *fasthttp.RequestCtx) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := oauth.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	authURL, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = StoreInSession(providerName, sess.Marshal(), ctx)
	if err != nil {
		return "", err
	}

	return authURL, err
}

// CompleteUserAuthOptions that affect how CompleteUserAuth works.
type CompleteUserAuthOptions struct {
	// True if CompleteUserAuth should automatically end the user's session.
	//
	// Defaults to True.
	ShouldLogout bool
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all the basic information about the user from the provider.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func CompleteUserAuth(ctx *fasthttp.RequestCtx, options ...CompleteUserAuthOptions) (oauth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return oauth.User{}, err
	}

	provider, err := oauth.GetProvider(providerName)
	if err != nil {
		return oauth.User{}, err
	}

	value, err := GetFromSession(providerName, ctx)
	if err != nil {
		return oauth.User{}, err
	}

	shouldLogout := true
	if len(options) > 0 && !options[0].ShouldLogout {
		shouldLogout = false
	}

	if shouldLogout {
		defer func(ctx *fasthttp.RequestCtx) {
			_ = Logout(ctx)
		}(ctx)
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return oauth.User{}, err
	}

	err = validateState(ctx, sess)
	if err != nil {
		return oauth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, nil
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, &Params{ctx: ctx})
	if err != nil {
		return oauth.User{}, err
	}

	err = StoreInSession(providerName, sess.Marshal(), ctx)
	if err != nil {
		return oauth.User{}, err
	}

	return provider.FetchUser(sess)
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(req *fasthttp.RequestCtx, sess oauth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	reqState := GetState(req)

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
}

// Logout invalidates a user session.
func Logout(ctx *fasthttp.RequestCtx) error {
	return defaultSession.Destroy(ctx)
}

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.

func GetProviderName(ctx *fasthttp.RequestCtx) (string, error) {
	// try to get it from the path param "provider"
	if p := ctx.UserValue("provider"); p != "" {
		return p.(string), nil
	}

	// try to get it from the url param "provider"
	if p := utils.UnsafeString(ctx.QueryArgs().Peek("provider")); p != "" {
		return p, nil
	}

	// try to get it from the context's value of "provider" key
	if p := utils.UnsafeString(ctx.Request.Header.Peek("provider")); p != "" {
		return p, nil
	}

	// try to get it from the Fasthttp context's value of providerContextKey key
	if p := utils.UnsafeString(ctx.Request.Header.Peek(fmt.Sprint(ProviderParamKey))); p != "" {
		return p, nil
	}

	// As a fallback, loop over the used providers, if we already have a valid session for any provider (ie. user has already begun authentication with a provider), then return that provider name
	providers := oauth.GetProviders()
	// Get session store
	store, _ := defaultSession.Get(ctx)

	for _, provider := range providers {
		p := provider.Name()
		value := store.Get(p)
		if _, ok := value.(string); ok {
			return p, nil
		}
	}

	// if not found then return an empty string with the corresponding error
	return "", errors.New("you must select a provider")
}

// StoreInSession stores a specified key/value pair in the session.
func StoreInSession(key, value string, ctx *fasthttp.RequestCtx) error {
	if err := updateSessionValue(key, value, ctx); err != nil {
		return err
	}

	return nil
}

// GetFromSession retrieves a previously-stored value from the session.
// If no value has previously been stored at the specified key, it will return an error.
func GetFromSession(key string, ctx *fasthttp.RequestCtx) (string, error) {
	value, err := getSessionValue(key, ctx)
	if err != nil {
		return "", errors.New("could not find a matching session for this request")
	}

	return value, nil
}

func getSessionValue(key string, ctx *fasthttp.RequestCtx) (string, error) {
	// Get session store
	store, _ := defaultSession.Get(ctx)

	// Get data session store.
	value := store.Get(key)
	if value == nil {
		return "", fmt.Errorf("could not find a matching session for this request")
	}

	// Decompress value
	rdata := strings.NewReader(value.(string))
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func updateSessionValue(key, value string, ctx *fasthttp.RequestCtx) error {
	// Compress value
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	// Get session store
	store, _ := defaultSession.Get(ctx)

	// Auto save session store
	defer func() {
		if err := defaultSession.Save(ctx, store); err != nil {
			log.Fatal(err)
		}
	}()

	// Set data into session store.
	store.Set(key, b.String())

	return nil
}
