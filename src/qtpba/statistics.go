package qtpba

import (
	"bufio"
	"database/sql"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type blacklistType []string
type tagCloudType map[string]int

var blacklist *blacklistType
var cleanupRegexp = regexp.MustCompile(`[[:punct:][:space:]]`)

func loadBlacklist(file string) (*blacklistType, error) {
	handler, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer handler.Close()

	blacklist := make(blacklistType, 0)

	scanner := bufio.NewScanner(handler)
	for scanner.Scan() {
		blacklist = append(blacklist, scanner.Text())
	}

	sort.Strings(blacklist)

	return &blacklist, nil
}

func (blacklist blacklistType) IsBlacklisted(word string) bool {
	word = strings.ToLower(word)
	idx := sort.SearchStrings(([]string)(blacklist), word)
	return (idx < len(blacklist) && blacklist[idx] == word)
}

func initBlacklist() {
	var err error

	blacklistFile := path.Join(GetBaseDir(), "static/blacklist.txt")
	logger.Println("Loading blacklist file...")
	if blacklist, err = loadBlacklist(blacklistFile); err != nil {
		logger.Println("Error loading blacklists file", blacklistFile, ". Proceeding without blacklists!")
	}
	logger.Println("Finished loading blacklist")
}

func cleanupWord(word string) string {
	return cleanupRegexp.ReplaceAllString(word, "")
}

func extractValidWords(message string) []string {
	words := make([]string, 0)
	// Create a scanner to extract the individual words
	scanner := bufio.NewScanner(strings.NewReader(message))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		word = cleanupWord(word)

		if !isValidWord(word) {
			continue
		} else if blacklist.IsBlacklisted(word) {
			continue
		}

		words = append(words, word)
	}

	return words
}

func isValidWord(word string) bool {
	if len(word) < 5 {
		return false
	}

	return true
}

func isValidMessage(message string) bool {
	// TODO
	return true
}

type TagCount struct {
	word  string
	count int
}

type TagCountSlice []TagCount

// Implements the sort.Sort interface
func (tc TagCountSlice) Len() int {
	return len(tc)
}

func (tc TagCountSlice) Less(i, j int) bool {
	return tc[i].count >= tc[j].count // Sorting in DESCENDING order!
}

func (tc TagCountSlice) Swap(i, j int) {
	tc[i], tc[j] = tc[j], tc[i]
}

func (tagCloud tagCloudType) TopN(n int) []TagCount {
	tagCounts := make(TagCountSlice, len(tagCloud))
	var i int

	for word, count := range tagCloud {
		tagCounts[i] = TagCount{word, count}
		i++
	}

	sort.Sort(tagCounts)

	if n >= len(tagCounts) {
		n = len(tagCounts)
	}
	return tagCounts[0:n]
}

func GetTweetStats() tagCloudType {
	tagCloud := make(tagCloudType)

	var err error
	var rows *sql.Rows
	if rows, err = db.Query("SELECT message FROM tweets ORDER BY rowid DESC LIMIT 1000"); err != nil {
		logger.Println(err)
		return nil
	}

	for rows.Next() {
		var message string
		rows.Scan(&message)

		if !isValidMessage(message) {
			continue
		}

		words := extractValidWords(message)
		for _, word := range words {
			tagCloud[word]++
		}
	}
	return tagCloud
}

/*
import (
	"os/exec"
	"strconv"
	"time"
)

func runStatsDaemon() {
	retryCount := 3
	for {
		// Me fijo las últimas horas
		timestamp := time.Now().Add(-2 * time.Hour).UTC().Unix()
		cmd := exec.Command("/home/marcelo/Documents/Programacion/qtpba/process_data.R", strconv.FormatInt(timestamp, 10))
		logger.Println("Running statistics with timestamp", timestamp, "...")

		if err := cmd.Run(); err != nil {
			logger.Println("Error running stats:", err)
			if retryCount--; retryCount > 0 {
				logger.Println("Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			} else {
				logger.Println("Giving up this run.")
			}
		} else {
			logger.Println("Successfully ran stats")
		}

		retryCount = 3
		time.Sleep(30 * time.Minute)
	}
}
*/
