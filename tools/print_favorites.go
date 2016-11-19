package main

import (
    "net/http"

    "fmt"
    "os"
    "sort"

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

    offset := 0
    count := 50

    collection := make(map[string]map[string]map[int]string)

    for {
        l.Infof(ctx, "Reading (%d):(%d)", offset, count)

        favorites, err := amc.GetFavoriteTracks(offset, count)
        log.PanicIf(err)

        len_ := len(favorites)
        if len_ == 0 {
            break
        }

        ids := make([]string, len_)

        for i, info := range favorites {
            ids[i] = info.Id
        }

        l.Infof(ctx, "Retrieving track details.")

        mc := napster.NewMetadataClient(ctx, hc, apiKey)
        tracks, err := mc.GetTrackDetail(ids...)
        log.PanicIf(err)

        for _, track := range tracks {
            artistCollection, foundArtist := collection[track.ArtistName]
            if foundArtist == false {
                artistCollection = make(map[string]map[int]string)
                collection[track.ArtistName] = artistCollection
            }

            albumCollection, foundAlbum := artistCollection[track.AlbumName]
            if foundAlbum == false {
                albumCollection = make(map[int]string)
                artistCollection[track.AlbumName] = albumCollection
            }

            albumCollection[track.Index] = track.Name
        }

        offset += len_
    }

    artists := make([]string, len(collection))
    i := 0
    for artist, _ := range collection {
        artists[i] = artist
        i++
    }

    sortableArtists := sort.StringSlice(artists)
    sortableArtists.Sort()

    for _, artist := range sortableArtists {
        fmt.Printf("%s\n", artist)

        albumCollection := collection[artist]
        albums := make([]string, len(albumCollection))
        i := 0
        for album, _ := range albumCollection {
            albums[i] = album
            i++
        }

        sortableAlbums := sort.StringSlice(albums)
        sortableAlbums.Sort()

        for _, album := range sortableAlbums {
            fmt.Printf("  %s\n", album)

            tracks := albumCollection[album]

            indices := make([]int, len(tracks))
            i := 0
            for index, _ := range tracks {
                indices[i] = index
                i++
            }

            sortableIndices := sort.IntSlice(indices)
            sortableIndices.Sort()

            for _, index := range sortableIndices {
                t := tracks[index]
                fmt.Printf("    %d: %s\n", index, t)
            }

            fmt.Printf("\n")
        }
    }
}
