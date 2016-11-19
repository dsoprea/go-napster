## Overview

This is a Go client for the Napster/Rhapsody API.

Implemented functionalities:

- Username/password-based authorization (still requires an API key).
- Retrieve favorite tracks
- Retrieve track detail for one or more tracks

**I have only implemented the functionalities that I required. I would be appreciative and more than happy to merge pull-requests from other developers that wish to extend this project.** See the [documentation](https://developer.rhapsody.com/api) for more of what the API supports.


## Example Usage

This file is available as [example.go](example/example.go):

```go
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

    l.Debugf(ctx, "Getting favorites.")

    a := napster.NewAuthenticator(ctx, hc, apiKey, secretKey)
    a.SetUserCredentials(username, password)

    amc := napster.NewAuthenticatedMemberClient(ctx, hc, a)

    offset := 100
    count := 10

    trackInfo, err := amc.GetFavoriteTracks(offset, count)
    if err != nil {
        log.Panic(err)
    }

    ids := make([]string, len(trackInfo))

    for i, info := range trackInfo {
        ids[i] = info.Id
    }

    l.Debugf(ctx, "Retrieving track details.")

    mc := napster.NewMetadataClient(ctx, hc, apiKey)
    tracks, err := mc.GetTrackDetail(ids...)
    log.PanicIf(err)

    for i, track := range tracks {
        fmt.Printf("%d: %s\n", i, track.String())
    }
}
```
