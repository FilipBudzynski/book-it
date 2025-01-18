package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/joho/godotenv"
)

const (
	GoogleBooksAPI          = "https://www.googleapis.com/books/v1/volumes"
	GoogleBooksAPIMaxResult = 40
)

var GoogleAPIKEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	GoogleAPIKEY = os.Getenv("GOOGLE_API_KET")
}

func NewGoogleProvider() handlers.BookProvider {
	return &googleProvider{
		apiUrl:     GoogleBooksAPI,
		maxResults: GoogleBooksAPIMaxResult,
	}
}

type googleProvider struct {
	apiUrl     string
	maxResults int
}

func (p *googleProvider) WithLimit(limit int) handlers.BookProvider {
	p.maxResults = limit
	return p
}

func (p *googleProvider) GetLimit() int {
	return p.maxResults
}

// Google Respnse structs
type BookResponse struct {
	ID         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type TotalItemsResponse struct {
	TotalItems int `json:"totalItems"`
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
	Pages         int      `json:"pageCount"`

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
		resp.Body.Close()
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func (p *googleProvider) GetBook(bookID string) (*models.Book, error) {
	url := fmt.Sprintf(p.apiUrl+"/%s", bookID)

	response, err := p.getResponse(url)
	if err != nil {
		return &models.Book{}, err
	}
	defer response.Body.Close()

	var bookResponse BookResponse
	if err := json.NewDecoder(response.Body).Decode(&bookResponse); err != nil {
		return &models.Book{}, err
	}

	return p.convert(bookResponse), nil
}

func (p *googleProvider) GetTotalForQuery(query string) int {
	url := fmt.Sprintf(p.apiUrl+"?q=%s&maxResults=%d", query, p.maxResults)
	response, err := p.getResponse(url)
	if err != nil {
		return -1
	}
	defer response.Body.Close()

	var totalItems TotalItemsResponse
	if err := json.NewDecoder(response.Body).Decode(&totalItems); err != nil {
		return -1
	}

	return totalItems.TotalItems
}

func (p *googleProvider) GetBooksByQuery(query string, queryType handlers.QueryType, limit, page int) ([]*models.Book, error) {
	startIndex := (page - 1) * limit
	params := url.Values{}
	urlRequest := fmt.Sprintf("%s\"%s\"", queryType, query)

    fmt.Println(urlRequest)

	params.Add("q", urlRequest)
	params.Add("maxResults", fmt.Sprintf("%d", limit))
	params.Add("startIndex", fmt.Sprintf("%d", startIndex))
	encodedUrl := p.apiUrl + "?" + params.Encode()

	response, err := p.getResponse(encodedUrl)
	if err != nil {
		return nil, err
	}
	if response != nil {
		defer response.Body.Close()
	}

	var bookItemsResponse BookItemsResponse
	if err := json.NewDecoder(response.Body).Decode(&bookItemsResponse); err != nil {
		return nil, err
	}

	var books []*models.Book
	for _, bookResp := range bookItemsResponse.Items {
		books = append(books, p.convert(bookResp))
	}

	return books, nil
}

func (p *googleProvider) GetBooksByGenre(genre string, maxResults int) ([]*models.Book, error) {
	query := url.QueryEscape(genre)
	url := fmt.Sprintf(p.apiUrl+"?q=subject:%s&maxResults=%d&key=%s", query, maxResults, GoogleAPIKEY)
	fmt.Printf("google genres url: %s\n", url)

	response, err := p.getResponse(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var bookItemsResponse BookItemsResponse
	if err := json.NewDecoder(response.Body).Decode(&bookItemsResponse); err != nil {
		return nil, err
	}

	var books []*models.Book
	for _, bookResp := range bookItemsResponse.Items {
		books = append(books, p.convert(bookResp))
	}

	return books, nil
}

func (p *googleProvider) convert(bookResponse BookResponse) *models.Book {
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
	book := &models.Book{
		ID:            bookResponse.ID,
		Title:         title,
		Authors:       strings.Join(volumeInfo.Authors, ", "),
		Link:          volumeInfo.ImageLinks.SmallThumbnail,
		Description:   description,
		ImageLink:     volumeInfo.ImageLinks.SmallThumbnail,
		PublishedDate: volumeInfo.PublishedDate,
		Pages:         volumeInfo.Pages,
	}
	if isbn, err := strconv.ParseUint(isbnString, 10, 0); err != nil {
		book.ISBN = uint(isbn)
	}

	return book
}
