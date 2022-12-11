package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/asstronom/EVO_tech_test/pkg/parse"
	"github.com/gin-gonic/gin"
)

// handler that returns transactions with applied filters
func (srv *Server) transactions(c *gin.Context) {
	//retrieving filters
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
		if status != "accepted" && status != "declined" {
			c.String(http.StatusBadRequest, "bad status filter")
			return
		}
		filters["status"] = status
	}
	payment_type, ok := c.GetQuery("payment_type")
	if ok {
		if payment_type != "cash" && payment_type != "card" {
			c.String(http.StatusBadRequest, "bad payment_type filter")
			return
		}
		filters["payment_type"] = payment_type
	}

	date_from, ok := c.GetQuery("date_post_from")
	if ok {
		date_from = strings.ReplaceAll(date_from, "T", " ")
		date, err := parse.ParseDate(date_from)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("wrong syntaxis of date_post_from: %s", err))
			return
		}
		filters["date_post_from"] = date
	}

	date_to, ok := c.GetQuery("date_post_to")
	if ok {
		date_to = strings.ReplaceAll(date_to, "T", " ")
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
	//running database query
	trxs, err := srv.service.GetTransactions(context.Background(), filters)
	if err != nil {
		c.String(http.StatusNotFound, "error getting transactions: ", err)
		return
	}
	c.IndentedJSON(http.StatusOK, trxs)
}

// handler that returns transaction by its id
func (srv *Server) transactionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			c.String(http.StatusBadRequest, "bad syntaxis for id")
			return
		}
	}
	trx, err := srv.service.GetTransactionByID(context.Background(), id)
	if err != nil {
		c.String(http.StatusNotFound, "error getting trx: %s", err)
		return
	}
	c.IndentedJSON(http.StatusOK, trx)
}

// handler that recieves csv file, parses it and stores data in db
func (srv *Server) uploadCSV(c *gin.Context) {
	fileh, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, "error uploading file")
		return
	}
	log.Println(fileh.Filename)
	file, err := fileh.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "error opening file: %s", err)
		return
	}
	err = srv.service.InsertTransactions(context.Background(), file)
	if err != nil {
		c.String(http.StatusInternalServerError, "error inserting data: %s", err)
		return
	}
	c.String(http.StatusCreated, "success!")
}
