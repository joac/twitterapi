package twitterapi

import (
    "bufio"
    "encoding/json"
    "github.com/garyburd/go-oauth/oauth"
    "io"
    "net/http"
    "net/url"
    "strconv"
)

type AccessCredentialSetting struct {
    Api, Access *oauth.Credentials
}

type User struct {
    Screen_name string
}

type Tweet struct {
    Text string
    User User
}

const (
    StatusStreamUrl = "https://stream.twitter.com/1.1/statuses/filter.json"
)

var (
    oauthClient = oauth.Client{
        TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
        ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
        TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
    }
    accessCredential = oauth.Credentials{}
    Setting          = AccessCredentialSetting{
        Api:    &oauthClient.Credentials,
        Access: &accessCredential,
    }
)

func PostValues(terms string) (form url.Values) {
    form = url.Values{}
    form.Add("track", terms)
    form.Add("delimited", "length")
    return form
}

func Listen(terms string, outChannel chan *Tweet) (err error) {
    v := PostValues(terms)
    resp, err := oauthClient.Post(http.DefaultClient, &accessCredential,
        StatusStreamUrl, v)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    err = parseBody(resp.Body, outChannel)
    if err != nil {
        return err
    }
    return nil
}

func parseBody(body io.Reader, outChannel chan *Tweet) (err error) {
    var (
        dataLength uint64
        tweet      Tweet
    )
    bufReader := bufio.NewReader(body)
    for {
        dataLength, err = ReadBlobLength(bufReader)
        if err != nil {
            return err
        }
        tweetBlob := make([]byte, dataLength)
        _, err = bufReader.Read(tweetBlob)
        if err != nil {
            return err
        }
        json.Unmarshal(tweetBlob, &tweet)
        outChannel <- &tweet
    }
}

func ReadBlobLength(r *bufio.Reader) (length uint64, err error) {
    for {
        blob, _, err := r.ReadLine()
        if err != nil {
            return 0, err
        }
        if len(blob) > 0 {
            length, err = strconv.ParseUint(string(blob), 10, 32)
            if err != nil {
                return 0, err
            }
            return length, nil
        }
    }
}
