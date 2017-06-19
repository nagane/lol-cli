package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"runtime"
)

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
	client := &Client{
		URL:        apiServerUrl,
		HTTPClient: httpClient,
		ApiKey:     apikey,
	}
	return client, nil
}

func (client *Client) newRequest(ctx context.Context, method string, subPath string, body io.Reader) (*http.Request, error) {
	endpointURL := *client.URL
	endpointURL.Path = path.Join(client.URL, subPath)

	req, err := http.NewRequest(method, endpointURL.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Text", "application/x-www-form-urlencoded")

	var userAgent = fmt.Sprintf("Goclient/%s (%s)", version, runtime.Version())

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

	httpClient = &http.Client{}
	client := NewClient(httpClient)

}
