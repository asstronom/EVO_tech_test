package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getTransactions(c *gin.Context) {

}

func getTransactionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			c.String(http.StatusBadRequest, "bad syntaxis for id")
		}
	}
	c.String(http.StatusOK, "returning transaction id:%d", id)
}
