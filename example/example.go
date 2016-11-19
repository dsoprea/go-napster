package main

import (
    "net/http"

    "fmt"
    "os"

    "golang.org/x/net/context"
    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-napster"
)

// Parameters
var (
    apiKey = os.Getenv("NAPSTER_API_KEY")
    secretKey = os.Getenv("NAPSTER_SECRET_KEY")
    username = os.Getenv("NAPSTER_USERNAME")
    password = os.Getenv("NAPSTER_PASSWORD")
)

// Other
var (
    l = log.NewLogger("main")
)

func main() {
    cla := log.NewConsoleLogAdapter()
    log.AddAdapter("console", cla)

    ctx := context.Background()
    hc := new(http.Client)

    l.Infof(ctx, "Getting favorites.")

    a := napster.NewAuthenticator(ctx, hc, apiKey, secretKey)
    a.SetUserCredentials(username, password)

    amc := napster.NewAuthenticatedMemberClient(ctx, hc, a)

    offset := 50
    count := 10

    trackInfo, err := amc.GetFavoriteTracks(offset, count)
    if err != nil {
        log.Panic(err)
    }

    ids := make([]string, len(trackInfo))

    for i, info := range trackInfo {
        ids[i] = info.Id
    }

    l.Infof(ctx, "Retrieving track details.")

    mc := napster.NewMetadataClient(ctx, hc, apiKey)
    tracks, err := mc.GetTrackDetail(ids...)
    log.PanicIf(err)

    for i, track := range tracks {
        fmt.Printf("%d: %s\n", i, track.String())
    }
}
