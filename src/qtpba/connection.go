package qtpba

import (
	"encoding/json"
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServerMessage map[string]interface{}

const RECEIVE_TIMEOUT = 5 * time.Minute

func LoadCredentials() (client *twittergo.Client, err error) {
	var credentials []byte
	if credentials, err = ioutil.ReadFile(GetFullPath("/conf/oauth_credentials.conf")); err != nil {
		return nil, err
	}

	lines := strings.Split(string(credentials), "\n")
	if len(lines) < 4 {
		return nil, fmt.Errorf("Invalid format for the oauth_credentials file")
	}

	consumerKey, consumerSecret, accessToken, accessTokenSecret := lines[0], lines[1], lines[2], lines[3]

	config := &oauth1a.ClientConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
	user := oauth1a.NewAuthorizedConfig(accessToken, accessTokenSecret)
	client = twittergo.NewClient(config, user)
	return
}

func StartListeningTweets() {
	//words := []string{"buenos aires", "bsas", "baires", "bs.as.", "bs.as", "bs. as.", "b. aires"}
	words := []string{}
	bounds := []float64{-58.533353, -34.70588, -58.32942, -34.5336}

	URL := generateURL(words, bounds)
	logger.Println("Querying URL", URL)
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		logger.Fatalln(fmt.Sprintf("Error creating HTTP request: %s", err))
	}

	rawLogFile := GetFullPath("logs/raw_stream.log")
	rawLog, err := os.OpenFile(rawLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Println("Unable to open", rawLogFile, "for appending")
	}

	messagesChan := make(chan ServerMessage)

	// TODO: Refactor this, make two separate functions instead of putting all the
	// code inside this function.
	var resp *http.Response
	connected := false
	go func() {
		if rawLog != nil {
			defer rawLog.Close()
		}

		client, err := LoadCredentials()
		if err != nil {
			logger.Fatalln(fmt.Sprintf("Error loading credentials: %s", err))
		}
		client.OAuth.Sign(req, client.User)

		for {
			connected = false
			logger.Println("Opening connection to server...")
			if resp, err = client.HttpClient.Do(req); err != nil {
				logger.Println("Error opening connection:", err, ". Waiting 10 seconds.")
				resp = nil
				time.Sleep(10 * time.Second)
				continue
			}

			logger.Println("Finished connecting")
			connected = true

			message := make(ServerMessage)

			// Redirect output to both the log file and the JSON decoder
			var reader io.Reader
			if rawLog != nil {
				reader = io.TeeReader(resp.Body, rawLog)
			} else {
				reader = resp.Body
			}
			jsonDecoder := json.NewDecoder(reader)
			for {
				err = jsonDecoder.Decode(&message)
				if err != nil {
					if connected {
						logger.Println("Error decoding JSON stream:", err, ". Retrying connection.")
					} else {
						logger.Println("Got disconnected. Retrying connection.")
					}
					time.Sleep(10 * time.Second)
					break
				}
				messagesChan <- message
			}
		}

	}()

	go func() {
		for {
			select {
			case message := <-messagesChan:
				processMessage(&message)
			case <-time.After(RECEIVE_TIMEOUT):
				if connected {
					logger.Println("Timed out while waiting for stream. Disconnecting in 10 seconds.")
					connected = false
					if resp.Body != nil {
						resp.Body.Close()
					}
				}
			}
		}
	}()
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
