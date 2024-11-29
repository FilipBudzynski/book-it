package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FilipBudzynski/book_it/pkg/schemas"
	"github.com/FilipBudzynski/book_it/pkg/services"
)

const (
	GoogleBooksAPI    = "https://www.googleapis.com/books/v1/volumes"
	DefaultMaxResults = 5
)

// NewGoogleBookService creates new BookService for external api communication
//
// apiUrl: expectes the url schema that will be used for getting the results
// example:  "https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=%d"
//
// maxResults: a number to specify max returned results for a query
func NewGoogleProvider() services.BookProvider {
	return &googleProvider{
		apiUrl:     GoogleBooksAPI,
		maxResults: DefaultMaxResults,
	}
}

type googleProvider struct {
	apiUrl     string
	maxResults int
}

func (p *googleProvider) WithLimit(limit int) services.BookProvider {
	p.maxResults = limit
	return p
}

// Google Respnse structs
type BookResponse struct {
	ID         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type BookItemsResponse struct {
	Items []BookResponse `json:"items"`
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

func (p *googleProvider) getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return resp, nil
}

func (p *googleProvider) GetBook(bookID string) (schemas.Book, error) {
	url := fmt.Sprintf(p.apiUrl+"/%s", bookID)

	response, err := p.getResponse(url)
	if err != nil {
		return schemas.Book{}, err
	}
	defer response.Body.Close()

	var bookResponse BookResponse
	if err := json.NewDecoder(response.Body).Decode(&bookResponse); err != nil {
		return schemas.Book{}, err
	}

	parsedBook, err := p.convert(bookResponse)
	if err != nil {
		return schemas.Book{}, err
	}

	return parsedBook, nil
}

func (p *googleProvider) GetBooksByQuery(query string, limit int) ([]schemas.Book, error) {
	url := fmt.Sprintf(p.apiUrl+"?q=%s&maxResults=%d", query, limit)

	response, err := p.getResponse(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var bookItemsResponse BookItemsResponse
	if err := json.NewDecoder(response.Body).Decode(&bookItemsResponse); err != nil {
		return nil, err
	}

	var books []schemas.Book
	for _, bookResp := range bookItemsResponse.Items {
		parsedBook, err := p.convert(bookResp)
		if err != nil {
			continue
		}
		books = append(books, parsedBook)
	}

	return books, nil
}

func (p *googleProvider) convert(bookResponse BookResponse) (schemas.Book, error) {
	var isbnString string
	volumeInfo := bookResponse.VolumeInfo
	for _, id := range volumeInfo.IndustryIdentifiers {
		if id.Type == "ISBN_13" {
			isbnString = id.Identifier
			break
		}
	}

	// Compose the title and subtitle if subtitle is present
	title := volumeInfo.Title
	if volumeInfo.Subtitle != "" {
		title = fmt.Sprintf("%s: %s", title, volumeInfo.Subtitle)
	}

	// Get description either from or SearchInfo
	description := volumeInfo.Description

	// Create and return a Book instance
	isbn, err := strconv.ParseUint(isbnString, 10, 0)
	if err != nil {
		return schemas.Book{}, err
	}
	return schemas.Book{
		ID:            bookResponse.ID,
		ISBN:          uint(isbn),
		Title:         title,
		Authors:       volumeInfo.Authors,
		Link:          volumeInfo.ImageLinks.SmallThumbnail,
		Description:   description,
		ImageLink:     volumeInfo.ImageLinks.SmallThumbnail,
		PublishedDate: volumeInfo.PublishedDate,
	}, nil
}
