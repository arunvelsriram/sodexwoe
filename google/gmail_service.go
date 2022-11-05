package google

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/arunvelsriram/sodexwoe/constants"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Debug("unable to identify the home directory")
		return nil, err
	}
	tokFile := filepath.Join(homeDir, constants.GOOGLE_TOKEN_FILE)
	log.WithField("tokenFile", tokFile).Debug("trying to use token from file")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.WithField("tokenFile", tokFile).Info("token file not found so getting token from web")
		tok, err := getTokenFromWeb(config)
		if err != nil {
			log.Debug("failed to get token from web")
			return nil, err
		}

		log.WithField("tokenFile", tokFile).Info("saving token in file")
		err = saveToken(tokFile, tok)
		if err != nil {
			log.Debug("failed to save the token")
			return nil, err
		}

		return config.Client(context.Background(), tok), nil
	}

	log.WithField("tokenFile", tokFile).Info("used token from file")

	return config.Client(context.Background(), tok), nil
}

func getAuthCode() (string, error) {
	type authCodeResponseHolder struct {
		authCode string
		err      error
	}
	ch := make(chan authCodeResponseHolder)

	addr := "localhost:7080"
	server := http.Server{Addr: addr}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("got callback")
		var authCodeRes authCodeResponseHolder
		q := r.URL.Query()
		if q.Has("code") {
			log.Debug("callback has code in query params")
			authCodeRes = authCodeResponseHolder{authCode: q.Get("code"), err: nil}
		} else {
			log.Debug("unable to get auth code from query params")
			authCodeRes = authCodeResponseHolder{err: fmt.Errorf("unable to get auth code from query params: %v", q)}
		}
		w.Write([]byte("you may close this tab now!"))
		ch <- authCodeRes
	})

	go func() {
		log.WithField("addr", addr).Info("starting server")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.WithField("addr", addr).Debug("failed to start server")
			ch <- authCodeResponseHolder{err: fmt.Errorf("failed to start server: %v", err)}
		}
	}()

	a := <-ch

	log.WithField("addr", addr).Info("shutting down the server")
	if err := server.Shutdown(context.Background()); err != nil {
		log.WithField("addr", addr).Debug("failed to shutdown server")
		return "", fmt.Errorf("failed to shutdown server: %v", err)
	}

	return a.authCode, a.err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Infof("Opening Auth URL: %v", authURL)
	if err := browser.OpenURL(authURL); err != nil {
		log.WithField("authURL", authURL).Debug("failed to open URL in browser")
		return nil, err
	}

	authCode, err := getAuthCode()
	if err != nil {
		log.Debugf("unable to get authCode")
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Debugf("Unable to retrieve token from web")
		return nil, err
	}

	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		log.WithField("file", file).Debugf("failed to get token from file")
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	log.WithField("path", path).Info("saving credential file")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Debug("unable to cache oauth token")
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}

func NewGmailService(googleAPICredentials string) (*gmail.Service, error) {
	b, err := base64.StdEncoding.DecodeString(googleAPICredentials)
	if err != nil {
		log.Debug("failed to deocde Google API credentials")
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Debug("failed get config from JSON")
		return nil, err
	}

	httpClient, err := getClient(config)
	if err != nil {
		log.Debug("failed crreate gmail client")
		return nil, err
	}

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		log.Debug("unable to create gmail service")
		return nil, err
	}

	return srv, nil
}
