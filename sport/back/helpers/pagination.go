package helpers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationRequest struct {
	// 0-based page
	Page int
	// max number of rows
	Size int
}

// parse pagination info using query string provided by context
func (pi *PaginationRequest) FromQueryString(ctx *gin.Context, maxSize int) error {
	page := ctx.Query("page")
	size := ctx.Query("size")
	var err error
	if page == "" {
		pi.Page = 0 // assuming first page
	} else {
		if pi.Page, err = strconv.Atoi(page); err != nil {
			return err
		}
	}
	if size == "" {
		pi.Size = maxSize
	} else {
		if pi.Size, err = strconv.Atoi(size); err != nil {
			return err
		}
		if maxSize > 0 {
			if pi.Size > maxSize {
				return fmt.Errorf(
					"parse PaginationInfo: requested size %d exceeded max allowed size %d",
					pi.Size,
					maxSize)
			}
		}
	}
	return nil
}

// format pagination info to the "limit ... offset" clause
func (pi *PaginationRequest) ToSqlQueryParams() []interface{} {
	offset := 0
	if pi.Page > 0 && pi.Size > 0 {
		offset = pi.Page * pi.Size
	}
	if pi.Size < 0 {
		return []interface{}{offset}
	}
	return []interface{}{pi.Size, offset}
}

func (pi *PaginationRequest) ToSqlQuery(nextParamIndex *int) string {
	if pi.Size < 0 {
		*nextParamIndex++
		return fmt.Sprintf("offset $%d", *nextParamIndex-1)
	}
	*nextParamIndex += 2
	return fmt.Sprintf("limit $%d offset $%d", *nextParamIndex-2, *nextParamIndex-1)
}

type PaginationResponse struct {
	NPages int
	PaginationRequest
}
