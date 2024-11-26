package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FilipBudzynski/book_it/pkg/schemas"
)

const (
	GoogleBooksAPI    = "https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=%d"
	DefaultMaxResults = 5
)

// NewGoogleBookService creates new BookService for external api communication
//
// apiUrl: expectes the url schema that will be used for getting the results
// example:  "https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=%d"
//
// maxResults: a number to specify max returned results for a query
func NewGoogleBookService() *googleBookService {
	return &googleBookService{
		apiUrl:     GoogleBooksAPI,
		maxResults: DefaultMaxResults,
	}
}

type googleBookService struct {
	apiUrl     string
	maxResults int
}

type BooksResponse struct {
	Items []struct {
		ID         string     `json:"id"`
		VolumeInfo VolumeInfo `json:"volumeInfo"`
	} `json:"items"`
}

type VolumeInfo struct {
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle,omitempty"`
	Authors       []string `json:"authors"`
	PublishedDate string   `json:"publishedDate"`
	Description   string   `json:"description,omitempty"`

	ImageLinks struct {
		SmallThumbnail string `json:"smallThumbnail"`
		Thumbnail      string `json:"thumbnail"`
	} `json:"imageLinks,omitempty"`

	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`
}

func (s *googleBookService) GetByQuery(query string, maxResults int) ([]*schemas.Book, error) {
	url := fmt.Sprintf(s.apiUrl, query, maxResults)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var booksResponse BooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&booksResponse); err != nil {
		return nil, err
	}

	var books []*schemas.Book
	for _, item := range booksResponse.Items {
		parsedBook, err := s.convert(item.VolumeInfo, item.ID)
		if err != nil {
			continue
		}
		books = append(books, &parsedBook)
	}

	return books, nil
}

func (s *googleBookService) GetMaxResults() int {
	return s.maxResults
}

func (s *googleBookService) convert(responseVolume VolumeInfo, googleId string) (schemas.Book, error) {
	var isbnString string
	for _, id := range responseVolume.IndustryIdentifiers {
		if id.Type == "ISBN_13" {
			isbnString = id.Identifier
			break
		}
	}

	// Compose the title and subtitle if subtitle is present
	title := responseVolume.Title
	if responseVolume.Subtitle != "" {
		title = fmt.Sprintf("%s: %s", title, responseVolume.Subtitle)
	}

	// Get description either from or SearchInfo
	description := responseVolume.Description

	// Create and return a Book instance
	isbn, err := strconv.ParseUint(isbnString, 10, 0)
	if err != nil {
		return schemas.Book{}, err
	}
	return schemas.Book{
		ID:            googleId,
		ISBN:          uint(isbn),
		Title:         title,
		Authors:       responseVolume.Authors,
		Link:          responseVolume.ImageLinks.SmallThumbnail,
		Description:   description,
		ImageLink:     responseVolume.ImageLinks.SmallThumbnail,
		PublishedDate: responseVolume.PublishedDate,
	}, nil
}
