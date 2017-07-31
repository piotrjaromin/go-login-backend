package web

import (
	"errors"
	"strconv"
	"github.com/labstack/echo"
)

func parseToInt(value string) (int, error) {
	parsedValue, parsingError := strconv.Atoi(value)
	if parsingError != nil {
		return 0, errors.New("Parameter is incorrect")
	}
	return parsedValue, nil
}

func GetPagination(c echo.Context) Pagination {
	pagination := DefaultPagination()

	page, pageErr := strconv.Atoi( c.QueryParam("page") )
	if pageErr == nil && page > 0  {
		pagination.PageNumber = page
	}

	pageSize, pageSizeErr := strconv.Atoi( c.QueryParam("pageSize") )
	if pageSizeErr == nil && pageSize > 0 {
		pagination.PageSize = pageSize
	}


	return pagination
}