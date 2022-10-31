package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/asstronom/EVO_tech_test/pkg/parse"
	"github.com/gin-gonic/gin"
)

func (srv *Server) getTransactions(c *gin.Context) {
	filters := make(map[string]interface{}, 5)
	terminal_ids, ok := c.GetQuery("terminal_ids")
	if ok {
		split := strings.Split(terminal_ids, ",")
		tids := make([]interface{}, len(split))
		for i := range split {
			id, err := strconv.Atoi(split[i])
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("bad terminal_ids: %s", split[i]))
				return
			}
			tids[i] = interface{}(id)
		}
		filters["terminal_ids"] = tids
	}
	status, ok := c.GetQuery("status")
	if ok {
		filters["status"] = status
	}
	payment_type, ok := c.GetQuery("payment_type")
	if ok {
		filters["payment_type"] = payment_type
	}
	date_from, ok := c.GetQuery("date_post_from")
	if ok {
		date, err := parse.ParseDate(date_from)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("wrong syntaxis of date_post_from: %s", err))
			return
		}
		filters["date_post_from"] = date
	}
	date_to, ok := c.GetQuery("date_post_to")
	if ok {
		date, err := parse.ParseDate(date_to)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("wrong syntaxis of date_post_to: %s", err))
			return
		}
		filters["date_post_to"] = date
	}
	payment_narrative, ok := c.GetQuery("payment_narrative")
	if ok {
		filters["payment_narrative"] = payment_narrative
	}
	c.JSON(http.StatusOK, filters)
}

func (srv *Server) getTransactionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			c.String(http.StatusBadRequest, "bad syntaxis for id")
			return
		}
	}
	c.String(http.StatusOK, "returning transaction id:%d", id)
}
