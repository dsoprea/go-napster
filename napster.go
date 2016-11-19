package napster

import (
    "fmt"
    "time"

    "net/http"

    "golang.org/x/net/context"

    "github.com/dsoprea/go-logging"
)

var (
    napsterLog = log.NewLogger("napster.napster")
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

    napsterLog.Debugf(a.ctx, "Authenticate")

// TODO(dustin): !! We need to manage refreshing. This will be done where we 
// actually make the calls.
    if a.atr == nil {
        atr, err := a.getToken()
        log.PanicIf(err)

        a.atr = atr
    }

    return a.atr.AccessToken, nil
}

func (a *Authenticator) getToken() (atr *accessTokenResponse, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    napsterLog.Debugf(a.ctx, "Requested authentication token: [%s] [%s] [%s] [%s]", a.accessKey, a.secretKey, a.username, a.password)

// TODO(dustin): Add support for other type(s) of authentication later.
    if a.username == "" || a.password == "" {
        log.Panic(fmt.Errorf("only user-credentials authentication currently supported"))
    }

    urlRaw := "https://api.napster.com/oauth/token"
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


// Client Napster API client
type Client struct {
    ctx context.Context
    hc *http.Client
    a *Authenticator
}

func NewClient(ctx context.Context, hc *http.Client, a *Authenticator) *Client {
    return &Client{
        ctx: ctx,
        hc: hc,
        a: a,
    }
}


//{"id":"tra.29611615","date":"2016-11-18T14:56:51.883Z","type":"favorite","links":{"tracks":{"ids":["tra.29611615"],"type":"track"}}},

type TrackResult struct {
    Id string               `json:"id"`
    Timestamp time.Time     `json:"date"`
    Type string             `json:"type"`
    Links interface{}       `json:"links"`
}

func (tr *TrackResult) String() string {
    return fmt.Sprintf("TrackResult(I=[%s] T=[%s])", tr.Id, tr.Timestamp.String())
}


type FavoriteTracksResult struct {
    Favorites []TrackResult     `json:"favorites"`
}


func (c *Client) GetFavoriteTracks(offset, limit int) (tracks []TrackResult, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    napsterLog.Debugf(c.ctx, "GetFavoriteTracks")

    accessToken, err := c.a.Authenticate()
    log.PanicIf(err)

    urlRaw := "https://api.napster.com/v2.0/me/favorites"
    verb := "GET"

    parameters := make(map[string]string)
    parameters["rights"] = "0"
    parameters["filter"] = "track"
    parameters["offset"] = fmt.Sprintf("%d", offset)
    parameters["limit"] = fmt.Sprintf("%d", limit)

    headers := make(map[string]string)
    headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)

    ftr := new(FavoriteTracksResult)
    if err := doJsonRequest(c.ctx, c.hc, urlRaw, verb, parameters, nil, headers, nil, ftr); err != nil {
        log.Panic(err)
    }

    return ftr.Favorites, nil
}
