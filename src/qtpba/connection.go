package qtpba

import (
	"encoding/json"
	"fmt"
	"github.com/kurrik/oauth1a"
	//"github.com/kurrik/twittergo"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const RECEIVE_TIMEOUT = 5 * time.Minute

type ServerMessage map[string]interface{}

type TwitterConnection struct {
	consumerKey, consumerSecret, accessToken, accessTokenSecret string
	clientConfig                                                *oauth1a.ClientConfig
	userConfig                                                  *oauth1a.UserConfig
	request                                                     *http.Request
	response                                                    *http.Response
	httpClient                                                  *http.Client
	oauth                                                       *oauth1a.Service
	rawLogFile                                                  *os.File
	connected                                                   bool
	messagesChan                                                chan (ServerMessage)
}

func NewTwitterConnection() (conn *TwitterConnection, err error) {
	conn = new(TwitterConnection)

	if err = conn.loadCredentials(); err != nil {
		return
	}

	// Set up file logging
	rawLogFile := GetFullPath("logs/raw_stream.log")
	conn.rawLogFile, err = os.OpenFile(rawLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Println("Unable to open", rawLogFile, "for appending")
		err = nil // Non fatal error
	}
	return
}

func (conn *TwitterConnection) loadCredentials() (err error) {
	var credentialsLines []byte
	if credentialsLines, err = ioutil.ReadFile(GetFullPath("/conf/oauth_credentials.conf")); err != nil {
		return
	}

	lines := strings.Split(string(credentialsLines), "\n")
	if len(lines) < 4 {
		return fmt.Errorf("Invalid format for the oauth_credentials file")
	}

	conn.consumerKey, conn.consumerSecret, conn.accessToken, conn.accessTokenSecret = lines[0], lines[1], lines[2], lines[3]

	conn.clientConfig = &oauth1a.ClientConfig{
		ConsumerKey:    conn.consumerKey,
		ConsumerSecret: conn.consumerSecret,
	}
	conn.userConfig = oauth1a.NewAuthorizedConfig(conn.accessToken, conn.accessTokenSecret)
	return
}

func (conn *TwitterConnection) setupClient() (err error) {
	words := []string{}
	bounds := []float64{-58.533353, -34.70588, -58.32942, -34.5336}
	URL := generateURL(words, bounds)
	logger.Println("Querying URL", URL)

	if conn.request, err = http.NewRequest("POST", URL, nil); err != nil {
		return
	}

	transport := new(http.Transport)
	if proxy, _ := http.ProxyFromEnvironment(conn.request); proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}
	conn.httpClient = &http.Client{Transport: transport}

	baseURL := "https://api.twitter.com"
	conn.oauth = &oauth1a.Service{
		RequestURL:   baseURL + "/oauth/request_token",
		AuthorizeURL: baseURL + "/oauth/authorize",
		AccessURL:    baseURL + "/oauth/access_token",
		ClientConfig: conn.clientConfig,
		Signer:       new(oauth1a.HmacSha1Signer),
	}
	conn.oauth.Sign(conn.request, conn.userConfig)

	return
}

func (conn *TwitterConnection) StartListeningTweets() {
	go conn.mainConnectionLoop()
	go conn.messageProcessingLoop()
}

func (conn *TwitterConnection) mainConnectionLoop() {
	conn.messagesChan = make(chan ServerMessage)

	if conn.rawLogFile != nil {
		defer conn.rawLogFile.Close()
	}

	var err error
	for {
		conn.setupClient()
		conn.connected = false

		logger.Println("Opening connection to server...")
		if conn.response, err = conn.httpClient.Do(conn.request); err != nil {
			logger.Println("Error opening connection:", err, ". Waiting 10 seconds.")
			conn.response = nil
			time.Sleep(10 * time.Second)
			continue
		}
		logger.Println("Finished connecting")
		conn.connected = true

		message := make(ServerMessage)

		// Redirect output to both the log file and the JSON decoder, if possible
		var jsonDecoder *json.Decoder
		if conn.rawLogFile != nil {
			jsonDecoder = json.NewDecoder(io.TeeReader(conn.response.Body, conn.rawLogFile))
		} else {
			jsonDecoder = json.NewDecoder(conn.response.Body)
		}

		// Decoding loop
		for {
			if err = jsonDecoder.Decode(&message); err != nil {
				logger.Println("Error decoding JSON stream:", err, ". Retrying connection in 10 seconds.")
				time.Sleep(10 * time.Second)
				break
			} else {
				conn.messagesChan <- message
			}
		}
	}
}

func (conn *TwitterConnection) messageProcessingLoop() {
	for {
		select {
		case message := <-conn.messagesChan:
			processMessage(&message)
		case <-time.After(RECEIVE_TIMEOUT):
			if conn.connected {
				logger.Println("Timed out while waiting for stream. Disconnecting in 10 seconds.")
				conn.connected = false
				if conn.response.Body != nil {
					conn.response.Body.Close()
				}
			}
		}
	}
}

func generateURL(words []string, bounds []float64) (URL string) {
	// Generate the URL
	URL = "https://stream.twitter.com/1.1/statuses/filter.json?"

	// Add the locations
	if len(bounds) > 0 {
		URL += "locations="
		for idx, coord := range bounds {
			strCoord := strconv.FormatFloat(coord, 'f', 5, 64)
			if idx == 0 {
				URL += strCoord
			} else {
				URL += "," + strCoord
			}

		}
	}

	// Add the words
	if len(words) > 0 {
		URL += "&track="
		for idx, word := range words {
			if idx == 0 {
				URL += url.QueryEscape(word)
			} else {
				URL += "," + url.QueryEscape(word)
			}
		}
	}
	return
}
