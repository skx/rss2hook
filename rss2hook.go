package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron"
)

// RSSEntry describes a single RSS feed and the corresponding hook
// to POST to.
type RSSEntry struct {
	feed string
	hook string
}

// Loaded contains the loaded feeds + hooks, as read from the specified
// configuration file
var Loaded []RSSEntry

// loadConfig loads the named configuration file and populates our
// `Loaded` list of RSS-feeds & Webhook addresses
func loadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening %s - %s\n", filename, err.Error())
		return
	}
	defer file.Close()

	//
	// Process it line by line.
	//
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tmp := scanner.Text()
		tmp = strings.TrimSpace(tmp)

		//
		// Skip lines that begin with a comment.
		//
		if (tmp != "") && (!strings.HasPrefix(tmp, "#")) {

			//
			// Otherwise find the feed + post-point
			//
			parser := regexp.MustCompile("^(.*)=([^=]+)")
			match := parser.FindStringSubmatch(tmp)
			if len(match) == 3 {
				entry := RSSEntry{feed: strings.TrimSpace(match[1]),
					hook: strings.TrimSpace(match[2])}
				Loaded = append(Loaded, entry)
			}

		}
	}

}

// fetchFeed fetches a feed from the remote URL.
func fetchFeed(url string) (string, error) {

	client := &http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "rss2email (https://github.com/skx/rss2email)")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// isNew returns TRUE if this feed-item hasn't been notified about
// previously.
func isNew(parent string, item *gofeed.Item) bool {

	hasher := sha1.New()
	hasher.Write([]byte(parent))
	hasher.Write([]byte(item.GUID))
	hashBytes := hasher.Sum(nil)

	// Hexadecimal conversion
	hexSha1 := hex.EncodeToString(hashBytes)

	if _, err := os.Stat(os.Getenv("HOME") + "/.rss2hook/seen/" + hexSha1); os.IsNotExist(err) {
		return true
	}
	return false
}

// recordSeen ensures that we won't re-announce a given feed-item.
func recordSeen(parent string, item *gofeed.Item) {

	hasher := sha1.New()
	hasher.Write([]byte(parent))
	hasher.Write([]byte(item.GUID))
	hashBytes := hasher.Sum(nil)

	// Hexadecimal conversion
	hexSha1 := hex.EncodeToString(hashBytes)

	dir := os.Getenv("HOME") + "/.rss2hook/seen"
	os.MkdirAll(dir, os.ModePerm)

	_ = ioutil.WriteFile(dir+"/"+hexSha1, []byte(item.Link), 0644)

}

// checkFeeds is our work-horse.
//
// For each available feed it looks for new entries, and when founds
// triggers `notify` upon the resulting entry
func checkFeeds() {

	for _, monitor := range Loaded {

		content, err := fetchFeed(monitor.feed)

		if err != nil {
			fmt.Printf("Error fetching %s - %s\n",
				monitor.feed, err.Error())
			continue
		}

		// Now we have the content - parse the feed
		fp := gofeed.NewParser()
		feed, err := fp.ParseString(content)
		if err != nil {
			fmt.Printf("Error parsing %s contents: %s\n", monitor.feed, err.Error())
			continue
		}

		// For each entry in the feed ..
		for _, i := range feed.Items {

			// If we've not already notified about this one.
			if isNew(monitor.feed, i) {

				err := notify(monitor.hook, i)
				if err == nil {
					recordSeen(monitor.feed, i)
				}
			}
		}
	}
}

// notify actually submits the specified item to the remote webhook.
//
// The RSS-item is submitted as a JSON-object.
func notify(hook string, item *gofeed.Item) error {
	jsonValue, err := json.Marshal(item)
	if err != nil {
		fmt.Printf("notify: Failed to encode JSON:%s\n", err.Error())
		return err
	}

	//
	// Post to purppura
	//
	res, err := http.Post(hook,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("notify: Failed to POST to %s - %s\n",
			hook, err.Error())
		return err
	}

	//
	// OK now we've submitted the post.
	//
	// We should retrieve the status-code + body, if the status-code
	// is "odd" then we'll show them.
	//
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	status := res.StatusCode

	if status != 200 {
		fmt.Printf("notify: Warning - Status code was not 200: %d\n", status)
	}
	return nil
}

// main is our entry-point
func main() {

	// Parse the command-line flags
	config := flag.String("config", "", "The path to the configuration-file to read")
	flag.Parse()

	if *config == "" {
		fmt.Printf("Please specify a configuration-file to read\n")
		return
	}

	//
	// Load the configuration file
	//
	loadConfig(*config)

	// Show the things we're monitoring
	for _, ent := range Loaded {
		fmt.Printf("Monitoring feed %s\nPosting to %s\n\n",
			ent.feed, ent.hook)
	}

	// Make the initial load
	checkFeeds()

	// Now repeat that every five minutes
	c := cron.New()
	c.AddFunc("@every 5m", func() { checkFeeds() })
	c.Start()

	// Wait to be terminated.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = <-sigs
		done <- true
	}()
	<-done
}
