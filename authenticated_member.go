package napster

import (
    "net/http"

    "time"
    "fmt"

    "golang.org/x/net/context"
    "github.com/dsoprea/go-logging"
)

// Other
var (
    amLog = log.NewLogger("napster.authenticated_member")
)

// Client Napster API client
type AuthenticatedMemberClient struct {
    ctx context.Context
    hc *http.Client
    a *Authenticator
}

func NewAuthenticatedMemberClient(ctx context.Context, hc *http.Client, a *Authenticator) *AuthenticatedMemberClient {
    return &AuthenticatedMemberClient{
        ctx: ctx,
        hc: hc,
        a: a,
    }
}


//{"id":"tra.29611615","date":"2016-11-18T14:56:51.883Z","type":"favorite","links":{"tracks":{"ids":["tra.29611615"],"type":"track"}}},

type AuthenticatedMemberTrackReferenceResult struct {
    Id string               `json:"id"`
    Timestamp time.Time     `json:"date"`
    Type string             `json:"type"`
    Links interface{}       `json:"links"`
}

func (tr *AuthenticatedMemberTrackReferenceResult) String() string {
    return fmt.Sprintf("AuthenticatedMemberTrackReferenceResult(I=[%s] T=[%s])", tr.Id, tr.Timestamp.String())
}


type AuthenticatedMemberFavoriteTracksResult struct {
    Favorites []AuthenticatedMemberTrackReferenceResult     `json:"favorites"`
}


func (amc *AuthenticatedMemberClient) GetFavoriteTracks(offset, limit int) (tracks []AuthenticatedMemberTrackReferenceResult, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    accessToken, err := amc.a.Authenticate()
    log.PanicIf(err)

    urlRaw := fmt.Sprintf("%s/me/favorites", api2UrlPrefix)
    verb := "GET"

    parameters := make(map[string]string)
    parameters["rights"] = "0"
    parameters["filter"] = "track"
    parameters["offset"] = fmt.Sprintf("%d", offset)
    parameters["limit"] = fmt.Sprintf("%d", limit)

    headers := make(map[string]string)
    headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)

    ftr := new(AuthenticatedMemberFavoriteTracksResult)
    if err := doJsonRequest(amc.ctx, amc.hc, urlRaw, verb, parameters, nil, headers, nil, ftr); err != nil {
        log.Panic(err)
    }

    return ftr.Favorites, nil
}
