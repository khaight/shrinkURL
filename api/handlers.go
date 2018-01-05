package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const alphabet = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
const length = int64(len(alphabet))

// ShortURL object
type ShortURL struct {
	ShortURL string `json:"shortURL"`
	LongURL  string `json:"longURL"`
	Visits   int64  `json:"visits"`
	Created  int64  `json:"created"`
}

func (s *ShortURL) marshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ShortURL) unmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// ShortURLRequest object
type shortURLRequest struct {
	*ShortURL
	URL string `json:"url" validate:"required"`
}

// ShortURLResponse object
type shortURLResponse struct {
	*ShortURL
}

// Render method for shortURL Response
func (rd *shortURLResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newShortURLResponse(shortURL *ShortURL) *shortURLResponse {
	resp := &shortURLResponse{ShortURL: shortURL}

	return resp
}

func (a *API) createNewShortURL(w http.ResponseWriter, r *http.Request) {
	data := &shortURLRequest{}

	if err := render.DecodeJSON(r.Body, data); err != nil {
		render.Render(w, r, errBadRequest(errors.New("Bad request. Must contain valid payload")))
		return
	}

	// validate request
	if err := data.validateShortURLRequest(); err != nil {
		render.Render(w, r, errBadRequest(err))
		return
	}

	// save url to db
	shortURL, err := a.dbCreateShortURL(data.URL)

	if err != nil {
		render.Render(w, r, errBadRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, newShortURLResponse(shortURL))

}

func (a *API) loadShortURL(w http.ResponseWriter, r *http.Request) {
	shortURL, err := a.dbLoadShortURL(chi.URLParam(r, "slug"))

	if err != nil {
		render.Render(w, r, errNotFound(errors.New("Unable to find short url")))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, newShortURLResponse(shortURL))

}

func (a *API) redirectShortURL(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	shortURL, err := a.dbLoadShortURL(slug)

	if err != nil {
		render.Render(w, r, errNotFound(errors.New("Unable to find short url")))
		return
	}

	shortURL.Visits = shortURL.Visits + 1
	buf, _ := shortURL.marshalBinary()

	// save short url to redis hash by key
	go a.Client.Set(slug, buf, 0)

	http.Redirect(w, r, shortURL.LongURL, http.StatusMovedPermanently)
}

func (s *shortURLRequest) validateShortURLRequest() error {
	// Validate structure
	err := validate.Struct(s)
	if err != nil {
		return errors.New("invalid request")
	}

	if len(s.URL) == 0 {
		return errors.New("empty url")
	}

	// validate we have a valid prefix
	if !strings.HasPrefix(s.URL, "http://") && !strings.HasPrefix(s.URL, "https://") {
		return errors.New("invalid url")
	}

	url, err := url.Parse(s.URL)

	if err != nil {
		return err
	}

	s.URL = url.String()

	return nil
}

func (a *API) dbCreateShortURL(url string) (*ShortURL, error) {
	// increment counter to shorten url
	ctr, _ := a.Client.IncrCounter()
	key := encode(ctr)

	// create short url
	sURL := new(ShortURL)
	sURL.Created = time.Now().UnixNano()
	sURL.LongURL = url
	sURL.ShortURL = fmt.Sprintf("http://%s/%s", a.hostname, key)
	sURL.Visits = 0

	buf, _ := sURL.marshalBinary()

	// save short url to redis hash by key
	go a.Client.Set(key, buf, 0)

	return sURL, nil
}

func (a *API) dbLoadShortURL(key string) (*ShortURL, error) {
	data, err := a.Client.Get(key).Result()
	if err != nil {
		return nil, errors.New("Unable to find key: " + key)
	}

	sURL := new(ShortURL)
	sURL.unmarshalBinary([]byte(data))
	return sURL, nil
}

func encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	s := ""
	for ; n > 0; n = n / length {
		s = string(alphabet[n%length]) + s
	}

	return s
}
