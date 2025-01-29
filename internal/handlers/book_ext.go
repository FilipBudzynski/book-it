package handlers

type QueryType int

const (
	QueryTypeTitle QueryType = iota
	QueryTypeAuthor
	QueryTypeSubject
	QueryTypeISBN
)

func stringToQueryType(queryType string) QueryType {
	switch queryType {
	case "title":
		return QueryTypeTitle
	case "author":
		return QueryTypeAuthor
	case "subject":
		return QueryTypeSubject
	case "isbn":
		return QueryTypeISBN
	default:
		return QueryTypeTitle
	}
}
