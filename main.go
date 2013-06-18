package main

import (
	"fmt"
	"os"
	"log"
	"flag"
	"net/http"
	"code.google.com/p/google-api-go-client/drive/v2"
	"code.google.com/p/goauth2/oauth"
)

var config = &oauth.Config{
  ClientId:     "1009616161174.apps.googleusercontent.com",
  ClientSecret: "TVxYsdNfl2VYcp2bb0rLQ1db",
  Scope:        "https://www.googleapis.com/auth/drive",
  RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
  AuthURL:      "https://accounts.google.com/o/oauth2/auth",
  TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

var (
	srv       *drive.Service
	transport oauth.Transport
)

var (
	doInit    = flag.Bool("init", false, "retrieve a new token")
	tokenFile = flag.String("tokenfile", getTokenFile(), "path to the token file")
)

func AllFiles(d *drive.Service) ([]*drive.File, error) {
  var fs []*drive.File
  pageToken := ""
  for {
    q := d.Files.List()
    // If we have a pageToken set, apply it to the query
    if pageToken != "" {
      q = q.PageToken(pageToken)
    }
    r, err := q.Do()
    if err != nil {
      fmt.Printf("An error occurred: %v\n", err)
      return fs, err
    }
    fs = append(fs, r.Items...)
    pageToken = r.NextPageToken
    if pageToken == "" {
      break
    }
  }
  return fs, nil
}

func getToken() {
	cache := oauth.CacheFile(*tokenFile)
	authUrl := config.AuthCodeURL("state")
	fmt.Printf("Go to the following link in your browser: %v\n", authUrl)
	t := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}

	// Read the code, and exchange it for a token.
	fmt.Printf("Enter verification code: ")
	var code string
	fmt.Scanln(&code)
	tok, err := t.Exchange(code)
	if err != nil {
		fmt.Printf("An error occurred exchanging the code: %v\n", err)
	}
	err = cache.PutToken(tok)
	if err != nil {
		log.Fatalln("Failed to save token:", err)
	}
}

func connect() {
	cache := oauth.CacheFile(*tokenFile)
	tok, err := cache.Token()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read token:", err)
		fmt.Fprintln(os.Stderr, "Did you run with -init?")
		os.Exit(1)
	} else {
		transport.Token = tok
	}
	srv, err = drive.New(transport.Client())
	if err != nil {
		log.Fatalln("Failed to create drive service:", err)
	}
}

func getTokenFile() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home := os.Getenv("HOME")
		if home == "" {
			log.Fatalln("Failed to determine token location (neither HOME nor" +
				" XDG_DATA_HOME are set)")
		}
		return home + "/.local/share/fusion/token"
	}
	return dataHome + "/drivefs/token"
}

func main(){
	transport.Config = config
	if *doInit {
		getToken()
		os.Exit(1)
	}

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("No Params")
		os.Exit(2)
	}
	connect()

	// Create a new authorized Drive client.
	svc, err := drive.New(transport.Client())
	if err != nil {
		fmt.Printf("An error occurred creating Drive client: %v\n", err)
	}

	files, err := AllFiles(svc)
	if err != nil {
		log.Fatalln("Failed to create drive service:", err)
	}
	file := files[1]
	fmt.Printf("%+v\n", file)
	// for _, value := range files {
	// 	fmt.Printf("File %+v\n", value)
	// }
	Mount(flag.Arg(0))
}