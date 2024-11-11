package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/FilipBudzynski/book_it/pkg/models"
)

// NewBookService creates new BookService for external api communication
//
// apiUrl: expectes the url schema that will be used for getting the results
// example:  "https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=%d"
//
// maxResults: a number to specify max returned results for a query
func NewBookService(apiUrl string, maxResults int) *bookService {
	return &bookService{
		apiUrl:     apiUrl,
		maxResults: maxResults,
	}
}

type bookService struct {
	apiUrl     string
	maxResults int
}

type BooksResponse struct {
	Items []struct {
		VolumeInfo VolumeInfo `json:"volumeInfo"`
	} `json:"items"`
}

type VolumeInfo struct {
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle,omitempty"`
	Authors       []string `json:"authors"`
	PublishedDate string   `json:"publishedDate"`
	Description   string   `json:"description,omitempty"`
	ImageLinks    struct {
		SmallThumbnail string `json:"smallThumbnail"`
		Thumbnail      string `json:"thumbnail"`
	} `json:"imageLinks,omitempty"`
	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`
}

func (s *bookService) GetByTitle(query string, maxResults int) ([]*models.Book, error) {
	encodedQuery := url.QueryEscape(query)
	url := fmt.Sprintf(s.apiUrl, encodedQuery, maxResults)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Decode the JSON response
	var booksResponse BooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&booksResponse); err != nil {
		return nil, err
	}

	var books []*models.Book
	for _, item := range booksResponse.Items {
		parsedBook, err := s.parseBook(item.VolumeInfo)
		if err != nil {
			continue
		}
		books = append(books, &parsedBook)
	}

	return books, nil
}

func (s *bookService) parseBook(item VolumeInfo) (models.Book, error) {
	var isbnString string
	for _, id := range item.IndustryIdentifiers {
		if id.Type == "ISBN_13" {
			isbnString = id.Identifier
			break
		}
	}

	// Compose the title and subtitle if subtitle is present
	title := item.Title
	if item.Subtitle != "" {
		title = fmt.Sprintf("%s: %s", title, item.Subtitle)
	}

	// Get description either from or SearchInfo
	description := item.Description

	// Create and return a Book instance
	isbn, err := strconv.ParseUint(isbnString, 10, 0)
	if err != nil {
		return models.Book{}, err
	}
	return models.Book{
		ISBN:          uint(isbn),
		Title:         title,
		Authors:       item.Authors,
		Link:          item.ImageLinks.SmallThumbnail,
		Description:   description,
		ImageLink:     item.ImageLinks.SmallThumbnail,
		PublishedDate: item.PublishedDate,
	}, nil
}

func (s *bookService) GetMaxResults() int {
	return s.maxResults
}
