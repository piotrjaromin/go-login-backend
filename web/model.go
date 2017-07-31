package web


type Pagination struct {
        PageNumber     int  `url:"pageNumber,omitempty"`
        PageSize       int  `url:"pageSize,omitempty"`
        WithTotalCount bool `url:"totalCount,omitempty"`
}

type ErrorType string

const (
        MissingField ErrorType = "MISSING_FIELD"
        InvalidField ErrorType = "INVALID_FIELD"
)

type ErrorDetails struct {
        Field   string    `json:"field,omitempty"`
        Type    ErrorType `json:"type"`
        Message string    `json:"message,omitempty"`
}

type Error struct {
        Message      string  `json:"message,omitempty"`
        ErrorDetails []ErrorDetails  `json:"details,omitempty"`
        Status       int `json:"status"`
}

func (err Error) Error() string {
        return err.Message
}

func AppendErrorDetails(errors []ErrorDetails, field string, msg string, t ErrorType) []ErrorDetails {
        detail := ErrorDetails{
                Field: field,
                Message: msg,
                Type: t,
        }

        return append(errors, detail)
}

func DefaultPagination() Pagination {
        return Pagination{
                PageNumber:     1,
                PageSize:       15,
                WithTotalCount: false,
        }
}
