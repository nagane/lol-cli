package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"
)

type SummonerNameResponse struct {
	Id            int    `json:id`
	AccountId     int    `json:accountId`
	Name          string `json:name`
	ProfileIconId int    `json:profileIconId`
	RevisionDate  int64  `json:revisionDate`
	SummonerLevel int    `json:summonerLevel`
}

var (
	apiServerUrl string = "https://jp1.api.riotgames.com"
	apiKey       string = viper.GetString("apikey")
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(httpClient *http.Client) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(apiServerUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "faild to parse url: %s", apiServerUrl)
	}
	client := &Client{
		URL:        parsedURL,
		HTTPClient: httpClient,
		ApiKey:     apiKey,
	}
	return client, nil
}

func (client *Client) newRequest(ctx context.Context, method string, subPath string, body io.Reader) (*http.Request, error) {
	endpointURL := *client.URL
	endpointURL.Path = path.Join(client.URL.Path, subPath)

	req, err := http.NewRequest(method, endpointURL.String(), body)
	if err != nil {
		return nil, err
	}

	// create query
	values := url.Values{}
	values.Add("api_key", viper.GetString("apikey"))
	req.URL.RawQuery = values.Encode()

	//println(req.URL.String())
	req = req.WithContext(ctx)
	req.Header.Set("Content-Text", "application/x-www-form-urlencoded")

	var userAgent = fmt.Sprintf("Goclient/ (%s)", runtime.Version())
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func init() {

	viper.SetConfigName("apikey")
	viper.AddConfigPath("$HOME/go/src/lol-cli")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config error: %s", err))
	}

}

func main() {
	httpClient := &http.Client{}
	client, _ := NewClient(httpClient)

	subPath := "/lol/summoner/v3/summoners/by-name/Kotan"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpRequest, err := client.newRequest(ctx, "GET", subPath, nil)
	if err != nil {
		print("newRequest error")
	}

	httpResponse, err := client.HTTPClient.Do(httpRequest)

	var apiResponce SummonerNameResponse
	if decerr := decodeBody(httpResponse, &apiResponce); decerr != nil {
		errors.Wrap(decerr, "decorde error:")
	}

	fmt.Printf("id:%d \nname:%s\niconId:%d\n", apiResponce.Id, apiResponce.Name, apiResponce.ProfileIconId)

}
