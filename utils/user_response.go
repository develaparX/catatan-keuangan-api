package utils

import (
	"livecode-catatan-keuangan/models/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendSingleResponse sends a single item response
func SendSingleResponse(ctx *gin.Context, data any, code int) {
	ctx.JSON(http.StatusOK, dto.SingleResponse{
		Status: dto.Status{
			ResponseCode: code,
		},
		Data: data,
	})
}

// SendPagingResponse sends a paginated list response
func SendPagingResponse(ctx *gin.Context, data interface{}, paging dto.Paging, code int) {
	dataList, ok := data.([]any)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dto.SingleResponse{
			Status: dto.Status{
				ResponseCode: http.StatusInternalServerError,
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.PagingResponse{
		Status: dto.Status{
			ResponseCode: code,
		},
		Data:   dataList,
		Paging: paging,
	})
}

// SendErrorResponse sends an error response
func SendErrorResponse(ctx *gin.Context, code int) {
	ctx.JSON(code, dto.SingleResponse{
		Status: dto.Status{
			ResponseCode: code,
		},
	})
}
