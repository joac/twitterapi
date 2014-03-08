package main

import (
    "encoding/json"
    "fmt"
    "github.com/garyburd/go-oauth/oauth"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
)

type AccessCredentialSetting struct {
    Api, Access *oauth.Credentials
}

var (
    oauthClient = oauth.Client{
        TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
        ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
        TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
    }
    accessCredential = oauth.Credentials{}
    setting          = AccessCredentialSetting{
        Api:    &oauthClient.Credentials,
        Access: &accessCredential,
    }
)

func readCredentials() error {
    b, err := ioutil.ReadFile("tokens.json")
    if err != nil {
        return err
    }
    return json.Unmarshal(b, &setting)
}

func main() {
    if err := readCredentials(); err != nil {
        log.Fatal(err)
    }
    fmt.Println(setting.Api.Token)
    fmt.Println(setting.Access.Token)
    fmt.Println(accessCredential.Token)
    v := url.Values{}
    v.Add("track", "go, golang")
    resp, err := oauthClient.Post(http.DefaultClient, &accessCredential,
        "https://stream.twitter.com/1.1/statuses/filter.json", v)

    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Hello World")
}
