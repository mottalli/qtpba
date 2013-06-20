package qtpba

import (
	"encoding/json"
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ServerMessage map[string]interface{}

func LoadCredentials() (client *twittergo.Client, err error) {
	var credentials []byte
	if credentials, err = ioutil.ReadFile(GetBaseDir() + "/conf/oauth_credentials"); err != nil {
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
	bounds := []int{-59, -35, -58, -34}

	URL := generateURL(words, bounds)
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		logger.Fatalln(fmt.Sprintf("Error creating HTTP request: %s", err))
	}

	// TODO: Clean this up!
	client, err := LoadCredentials()
	if err != nil {
		logger.Fatalln(fmt.Sprintf("Error loading credentials: %s", err))
	}
	client.OAuth.Sign(req, client.User)

	go func() {
		for {
			resp, err := client.HttpClient.Do(req)
			message := make(ServerMessage)
			jsonDecoder := json.NewDecoder(resp.Body)
			for {
				if err = jsonDecoder.Decode(&message); err != nil {
					logger.Println("Error decoding JSON stream:", err, ". Retrying connection.")
					break
				}

				processMessage(&message)
			}
		}

	}()

	/*
		// Get the username and the password from the environment
		username, password := os.Getenv("QTPBA_USERNAME"), os.Getenv("QTPBA_PASSWORD")
		if username == "" || password == "" {
			logger.Fatal("QTPBA_USERNAME or QTPBA_PASSWORD environment variables not set")
		}

		go func() {
			var err error
			var req *http.Request

			client := new(http.Client)

			for {
				if req, err = http.NewRequest("POST", URL, nil); err != nil {
					panic(err) // Should not happen
				}
				req.SetBasicAuth(username, password)

				var resp *http.Response
				logger.Println("Connecting to", URL)
				if resp, err = client.Do(req); err != nil {
					logger.Println("Error requesting URL:", err, ". Trying again in 10 seconds...")
					time.Sleep(10 * time.Second)
					continue
				}
				logger.Println("Connected")

				message := make(ServerMessage)
				jsonDecoder := json.NewDecoder(resp.Body)
				for {
					if err = jsonDecoder.Decode(&message); err != nil {
						logger.Println("Error decoding JSON stream:", err, ". Retrying connection.")
						break
					}

					processMessage(&message)
				}
			}
		}()*/
}

func generateURL(words []string, bounds []int) (URL string) {
	// Generate the URL
	URL = "https://stream.twitter.com/1.1/statuses/filter.json?"

	// Add the locations
	if len(bounds) > 0 {
		URL += "locations="
		for idx, coord := range bounds {
			if idx == 0 {
				URL += strconv.Itoa(coord)
			} else {
				URL += "," + strconv.Itoa(coord)
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
