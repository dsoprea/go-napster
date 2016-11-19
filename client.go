package napster

import (
    "net/http"
    "net/url"
    "encoding/json"
    "net/http/httputil"

    "fmt"
    "strings"
    "io"
    "bytes"

    "golang.org/x/net/context"

    "github.com/dsoprea/go-logging"
)

// Content types
const (
    CtFormUrlEncoded = "application/x-www-form-urlencoded"
    CtApplicationJson = "application/json"
)

var (
    clientLog = log.NewLogger("napster.client")
)


// httpBasicAuthentication Describes credentials used for basic HTTP 
// authentication
type httpBasicAuthentication struct {
    Username string
    Password string
}

func dumpRequest(r *http.Request) {
    dump, err := httputil.DumpRequestOut(r, true)
    if err != nil {
      fmt.Println(err)
    }

    fmt.Println(string(dump))
}

func doRequest(ctx context.Context, hc *http.Client, urlRaw string, verb string, parameters map[string]string, requestContentType, responseContentType string, data map[string]string, headers map[string]string, hba *httpBasicAuthentication) (response *http.Response, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    if parameters != nil {
        values := url.Values{}
        for k, v := range parameters {
            values.Add(k, v)
        }

        urlRaw += "?" + values.Encode()
    }

    var dataSource io.Reader
    if data != nil {
        values := url.Values{}

        for k, v := range data {
            values.Set(k, v)
        }

        encoded := values.Encode()
        dataSource = strings.NewReader(encoded)
    }

    r, err := http.NewRequest(verb, urlRaw, dataSource)
    if err != nil {
        log.Panic(err)
    }

    h := http.Header{}
    
    if requestContentType != "" {
        h.Add("Content-Type", requestContentType)
    }

    if headers != nil {
        for k, v := range headers {
            h.Add(k, v)
        }
    }

    r.Header = h

    if hba != nil {
        r.SetBasicAuth(hba.Username, hba.Password)
    }

    clientLog.Debugf(ctx, "Call: [%s] [%s]", verb, urlRaw)

    response, err = hc.Do(r)
    if err != nil {
        log.Panic(err)
    } else if response.StatusCode / 100 > 3 {
        b := new(bytes.Buffer)
        b.ReadFrom(response.Body)
        message := b.String()

        log.Panic(fmt.Errorf("API request failed: (%d) [%s]", response.StatusCode, message))
    }

    if responseContentType != "" {
        actualContentType := response.Header.Get("Content-Type")
        if actualContentType == "" {
            log.Panic(fmt.Errorf("no content-type returned"))
        }

        if pivot := strings.Index(actualContentType, ";"); pivot != -1 {
            actualContentType = actualContentType[:pivot]
        }

        if strings.ToLower(actualContentType) != responseContentType {
            log.Panic(fmt.Errorf("content-type of response unexpected: [%s] != [%s]", actualContentType, responseContentType))
        }
    }

    return response, nil
}

func doJsonRequest(ctx context.Context, hc *http.Client, urlRaw string, verb string, parameters map[string]string, data map[string]string, headers map[string]string, hba *httpBasicAuthentication, result interface{}) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    response, err := doRequest(
        ctx, 
        hc, 
        urlRaw, 
        verb, 
        parameters, 
        CtFormUrlEncoded, 
        CtApplicationJson, 
        data, 
        headers, 
        hba)

    log.PanicIf(err)

    // For debugging
    if result == nil {
        b := new(bytes.Buffer)
        b.ReadFrom(response.Body)
        message := b.String()

        clientLog.Debugf(ctx, "Not parsing return:\n%s", message)
    } else {
        d := json.NewDecoder(response.Body)
        if err := d.Decode(result); err != nil {
            log.Panic(err)
        }
    }

    return nil
}
