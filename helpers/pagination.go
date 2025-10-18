package helpers

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

type PaginationLink struct {
	URL    string `json:"url"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil || i < 1 {
		return 1
	}
	return i
}

func TotalPage(totalData int64, perPage int) int {
	if perPage == 0 {
		return 1
	}

	pages := float64(totalData) / float64(perPage)
	result := int(math.Ceil(pages))
	return result
}

func QueryString(search string) string {
	if search == "" {
		return ""
	}
	return fmt.Sprintf("&search=%s", search)
}

func PageURL(baseURL string, page, lastPage int, search string) string {
	if page < 1 || page > lastPage {
		return ""
	}
	return fmt.Sprintf("%s?page=%s%s", baseURL, strconv.Itoa(page), QueryString(search))
}

func GetPaginationParams(c *gin.Context) (search string, page, limit, offset int) {
	search = c.Query("search")
	page = StringToInt(c.DefaultQuery("page", "1"))
	limit = StringToInt(c.DefaultQuery("limit", "5"))
	offset = (page - 1) * limit
	return
}

func BuildBaseURL(c *gin.Context) string {
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.URL.Path)
}

func BuildPaginationLinks(currentPage, lastPage int, baseURL, search string) []PaginationLink {
	links := []PaginationLink{}

	// previous page
	links = append(links, PaginationLink{
		URL:    PageURL(baseURL, currentPage-1, lastPage, search),
		Label:  "&laquo; Previous",
		Active: false,
	})

	for i := 1; i <= lastPage; i++ {
		links = append(links, PaginationLink{
			URL:    fmt.Sprintf("%s?page=%s%s", baseURL, strconv.Itoa(i), QueryString(search)),
			Label:  strconv.Itoa(i),
			Active: i == currentPage,
		})
	}

	links = append(links, PaginationLink{
		URL:    PageURL(baseURL, currentPage+1, lastPage, search),
		Label:  "Next &raquo;",
		Active: false,
	})

	return links
}

func PaginateResponse(c *gin.Context, data any, total int64, page, limit int, baseURL, search, message string) {
	lastPage := TotalPage(total, limit)
	from := (page-1)*limit + 1
	to := from + reflect.ValueOf(data).Len() - 1

	links := BuildPaginationLinks(page, lastPage, baseURL, search)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data: gin.H{
			"current_page":   page,
			"data":           data,
			"first_page_url": fmt.Sprintf("%s?page=1%s", baseURL, QueryString(search)),
			"from":           from,
			"last_page":      lastPage,
			"last_page_url":  fmt.Sprintf("%s?page=%s%s", baseURL, strconv.Itoa(lastPage), QueryString(search)),
			"links":          links,
			"next_page_url":  PageURL(baseURL, page+1, lastPage, search),
			"path":           baseURL,
			"per_page":       limit,
			"prev_page_url":  PageURL(baseURL, page-1, lastPage, search),
			"to":             to,
			"total":          total,
		},
	})
}
