package napster

import (
    "net/http"

    "fmt"
    "strings"

    "golang.org/x/net/context"
    "github.com/dsoprea/go-logging"
)

// Other
var (
    metadataLog = log.NewLogger("napster.metadata")
)

// Client Napster API client
type MetadataClient struct {
    ctx context.Context
    hc *http.Client
    ak string
}

func NewMetadataClient(ctx context.Context, hc *http.Client, apiKey string) *MetadataClient {
    return &MetadataClient{
        ctx: ctx,
        hc: hc,
        ak: apiKey,
    }
}

/*
{
  "type": "track",
  "id": "Tra.16125739",
  "index": 4,
  "disc": 1,
  "href": "https:\/\/api.napster.com\/v2.0\/tracks\/Tra.16125739",
  "playbackSeconds": 267,
  "explicit": false,
  "name": "Return",
  "isrc": "USAT20703652",
  "shortcut": "needtobreathe\/the-heat\/return",
  "amg": "11714263",
  "blurbs": [
    
  ],
  "artistName": "Needtobreathe",
  "albumName": "The Heat",
  "formats": [
    {
      "type": "format",
      "bitrate": 320,
      "name": "AAC"
    },
    {
      "type": "format",
      "bitrate": 192,
      "name": "AAC"
    },
    {
      "type": "format",
      "bitrate": 64,
      "name": "AAC PLUS"
    }
  ],
  "albumId": "Alb.16115052",
  "contributors": {
    "primaryArtist": "Art.9105561"
  },
  "links": {
    "artists": {
      "ids": [
        "Art.9105561"
      ],
      "href": "https:\/\/api.napster.com\/v2.0\/artists\/Art.9105561"
    },
    "albums": {
      "ids": [
        "Alb.16115052"
      ],
      "href": "https:\/\/api.napster.com\/v2.0\/albums\/Alb.16115052"
    },
    "genres": {
      "ids": [
        "g.454",
        "g.418",
        "g.458"
      ],
      "href": "https:\/\/api.napster.com\/v2.0\/genres\/g.454,g.418,g.458"
    },
    "tags": {
      "ids": [
        "tag.152196499",
        "tag.152196552",
        "tag.152196501"
      ],
      "href": "https:\/\/api.napster.com\/v2.0\/tags\/tag.152196499,tag.152196552,tag.152196501"
    }
  },
  "previewURL": "http:\/\/listen.vo.llnwd.net\/g1\/1\/2\/6\/2\/6\/151362621.mp3",
  "isStreamable": true
},
*/


type MetadataTrackDetailFormat struct {
    Type string         `json:"type"`
    Bitrate int         `json:"bitrate"`
    Name string         `json:"name"`
}


type MetadataTrackDetail struct {
    Type string                             `json:"type"`
    Id string                               `json:"id"`
    Index int                               `json:"index"`
    Disc int                                `json:"disc"`
    Href string                             `json:"href"`
    PlaybackSeconds int                     `json:"playbackSeconds"`
    Explicit bool                           `json:"explicit"`
    Name string                             `json:"name"`
    Isrc string                             `json:"isrc"`
    Shortcut string                         `json:"shortcut"`
    Amg string                              `json:"amg"`
//  blurbs
    ArtistName string                       `json:"artistName"`
    AlbumName string                        `json:"albumName"`
    Formats []MetadataTrackDetailFormat     `json:"formats"`

    AlbumId string                          `json:"albumId"`
    Contributors map[string]string          `json:"contributors"`
//  links
    PreviewURL string                       `json:"previewURL"`
    IsStreamable bool                       `json:"isStreamable"`
}

func (mtd *MetadataTrackDetail) String() string {
    return fmt.Sprintf("MetadataTrackDetail(I=[%s] AR=[%s] AL=[%s] N=[%s])", mtd.Id, mtd.ArtistName, mtd.AlbumName, mtd.Name)
}


type MetadataTrackDetailResult struct {
    Tracks []MetadataTrackDetail        `json:"tracks"`
}


func (mc *MetadataClient) GetTrackDetail(id ...string) (tracks []MetadataTrackDetail, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    if len(id) == 0 {
        log.Panic(fmt.Errorf("no tracks provided"))
    }

    inline := strings.Join(id, ",")
    urlRaw := fmt.Sprintf("%s/tracks/%s", api2UrlPrefix, inline)
    verb := "GET"

    parameters := make(map[string]string)
    parameters["apikey"] = mc.ak

    mtdr := new(MetadataTrackDetailResult)
    if err := doJsonRequest(mc.ctx, mc.hc, urlRaw, verb, parameters, nil, nil, nil, mtdr); err != nil {
        log.Panic(err)
    }

    return mtdr.Tracks, nil
}
