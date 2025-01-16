package handlers

type QueryType string

const (
	QueryTypeTitle   = "intitle:"
	QueryTypeAuthor  = "inauthor:"
	QueryTypeSubject = "subject:"
	QueryTypeISBN    = ""
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
