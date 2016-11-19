package napster

import (
    "fmt"

    "net/http"

    "golang.org/x/net/context"
    "github.com/dsoprea/go-logging"
)

var (
    authLog = log.NewLogger("napster.authorization")
)


// accessTokenResponse Describes the authentication response.
type accessTokenResponse struct {
    AccessToken string      `json:"access_token"`
    RefreshToken string     `json:"refresh_token"`
    ExpiresIn int           `json:"expires_in"`
}


// Authenticator Knows how to get an access token.
type Authenticator struct {
    ctx context.Context
    hc *http.Client

    accessKey string
    secretKey string

    username string
    password string

    atr *accessTokenResponse
}

func NewAuthenticator(ctx context.Context, hc *http.Client, accessKey, secretKey string) *Authenticator {
    return &Authenticator{
        ctx: ctx,
        hc: hc,

        accessKey: accessKey,
        secretKey: secretKey,
    }
}

func (a *Authenticator) SetUserCredentials(username, password string) {
    if username == "" {
        log.Panic(fmt.Errorf("username empty"))
    } else if password == "" {
        log.Panic(fmt.Errorf("password empty"))
    }

    a.username = username
    a.password = password
}

func (a *Authenticator) Authenticate() (accessToken string, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

// TODO(dustin): !! We need to manage refreshing. This will be done where we 
// actually make the calls.
    if a.atr == nil {
        atr, err := a.getToken()
        log.PanicIf(err)

        a.atr = atr
    }

    authLog.Debugf(a.ctx, "Authorization token: [%s]", a.atr.AccessToken)

    return a.atr.AccessToken, nil
}

func (a *Authenticator) getToken() (atr *accessTokenResponse, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    authLog.Debugf(a.ctx, "Requested authentication token: [%s] [%s] [%s] [%s]", a.accessKey, a.secretKey, a.username, a.password)

// TODO(dustin): Add support for other type(s) of authentication later.
    if a.username == "" || a.password == "" {
        log.Panic(fmt.Errorf("only user-credentials authentication currently supported"))
    }

    urlRaw := fmt.Sprintf("%s/oauth/token", apiUrlPrefix)
    verb := "POST"

    data := make(map[string]string)
    data["username"] = a.username
    data["password"] = a.password
    data["grant_type"] = "password"

    hba := &httpBasicAuthentication{
        Username: a.accessKey,
        Password: a.secretKey,
    }

    atr = new(accessTokenResponse)

    if err := doJsonRequest(a.ctx, a.hc, urlRaw, verb, nil, data, nil, hba, atr); err != nil {
        log.Panic(err)
    }

    return atr, nil
}
