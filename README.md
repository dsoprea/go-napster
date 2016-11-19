## Overview

This is a Go client for the Napster/Rhapsody API.

Implemented functionalities:

- Username/password-based authorization (still requires an API key).
- Retrieve favorite tracks
- Retrieve track detail for one or more tracks

**I have only implemented the functionalities that I required. I would be appreciative and more than happy to merge pull-requests from other developers that wish to extend this project.** See the [documentation](https://developer.rhapsody.com/api) for more of what the API supports.


## Example Usage

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
```

Output:

```
$ NAPSTER_API_KEY="APIHERE" NAPSTER_SECRET_KEY="SECRETHERE" NAPSTER_USERNAME="USERNAMEHERE" NAPSTER_PASSWORD="PASSWORDHERE" GOPATH=`pwd`/../../../../.. go run example.go
2016/11/19 00:43:43 main: Getting favorites.
2016/11/19 00:43:44 main: Retrieving track details.
0: MetadataTrackDetail(I=[Tra.6870863] AR=[3 Doors Down] AL=[Seventeen Days] N=[Never Will I Break])
1: MetadataTrackDetail(I=[Tra.6870860] AR=[3 Doors Down] AL=[Seventeen Days] N=[Landing In London])
2: MetadataTrackDetail(I=[Tra.3212111] AR=[3 Doors Down] AL=[Away From The Sun] N=[This Time])
3: MetadataTrackDetail(I=[Tra.3212108] AR=[3 Doors Down] AL=[Away From The Sun] N=[Changes])
4: MetadataTrackDetail(I=[Tra.3212104] AR=[3 Doors Down] AL=[Away From The Sun] N=[Running Out Of Days])
5: MetadataTrackDetail(I=[Tra.3212103] AR=[3 Doors Down] AL=[Away From The Sun] N=[Ticket To Heaven])
6: MetadataTrackDetail(I=[Tra.3212101] AR=[3 Doors Down] AL=[Away From The Sun] N=[Away From The Sun])
7: MetadataTrackDetail(I=[Tra.1868845] AR=[Matchbox Twenty] AL=[Yourself Or Someone Like You] N=[Shame])
8: MetadataTrackDetail(I=[Tra.1868838] AR=[Matchbox Twenty] AL=[Yourself Or Someone Like You] N=[Push])
9: MetadataTrackDetail(I=[Tra.1868837] AR=[Matchbox Twenty] AL=[Yourself Or Someone Like You] N=[3 a.m.])
```


## Tools

There is a [script](tools/print_favorites.go) to read, group, sort, and print a list of favorites. Set the same environment variables as above to pass the necessary arguments.

Output:

```
OMC
  How Bizarre
    2: How Bizarre

Of Mice and Men
  Cold World
    1: Game of War
    3: Real
    4: Like a Ghost
    5: Contagious

Orchestral Manoeuvres in the Dark (OMD)
  Crush
    1: So In Love

Otherwise
  Peace at All Costs
    5: Coming for the Throne

PFR
  Great Lengths
    1: Great Lengths  (Great Lengths Album Version)
    3: Merry Go Round
    5: It's You Jesus
    8: See the Sun Again
```


## Reasoning

The reason I wrote this tool is because of how broken Napster/Rhapsody is:

- The website only shows the first five hundred entries in my favorites.
- The Android app now no longer shows which songs are favorited unless I am actually looking at the favorites list (and I have to browse the whole thing since there is no way to search through it from the device).
- There is no general way to see which artists or albums I have any favorites in (before it stopped working completely; it only shows which songs you have downloaded which usually agrees unless you change your phone).
- The ability to download all of your favorites is available but will cause the app to stop responding. Even though it will still be downloading your songs, it will only download the top N undownloaded songs in the list and will eventually leave everything queued but stop downloading. There is no apparent way to jumpstart it (like canceling the current item and re-adding).

This is only the most-important problem. I stopped discussing the other ones because Napster Support kept confusing them.


Their support department is ineffective:

- Since the support department does not have a contact-number, you are forced to wait on their responses. 
- A different person will respond to active threads every time.
- Even when you have a scheduled time for them to call, they will call you at all hours of the day and not leave a message telling you *why* they are calling. 
- I gave them my permission to reset my password once in order to observe what I was seeing, and they reset it several times after that without my permission. 
- They have suggested moving my favorited songs to a playlist and then moving it back. They tried it with the first five-hundred of my favorites before I had the chance to ask why that would even work. 
- I have been struggling with this for four months and I have to explain the problem in different ways every time in an attempt to get the issue escalated to development. Once, he told me that he would, and he never did. 

I would complain, except that the supervisor was one of the two guys that are working on the ticket. Who do you complain to when the supervisor is dangerously incompetent?

**I have been a customer for most of the past thirteen years** and I have spent the last several years systematically listening through every album of every artist I have ever liked and I am scared that Napster is going to screw something up permanently. Backing-up the list is the first step. Then, I will elect the same via the API of the next service.
