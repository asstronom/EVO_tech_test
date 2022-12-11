package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/asstronom/EVO_tech_test/pkg/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionDB struct {
	pool *pgxpool.Pool
}

// open database
func Open(ctx context.Context, url string) (*TransactionDB, error) {
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &TransactionDB{pool: pool}, nil
}

// close database
func (db *TransactionDB) Close() {
	db.pool.Close()
}

// insert single transaction
func (db *TransactionDB) InsertTransaction(ctx context.Context, trx domain.Transaction) error {
	_, err := db.pool.Exec(ctx, `INSERT INTO transactions (id, requestid, terminalid, partnerobjectid,
		 amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
		 dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid,
		 servicetypeid, payeeid, payeenameid, payeebankmfo, payeebankaccount, paymentnarrativeid)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, (SELECT id FROM servicetypes WHERE title=$16),
		 $17, (SELECT id FROM payeenames WHERE title=$18), $19, $20, (SELECT id FROM paymentnarratives WHERE title=$21))`,
		trx.ID, trx.RequestID, trx.TerminalID, trx.PartnerObjectID, trx.AmountTotal, trx.AmountOriginal, trx.CommisionPs, trx.CommisionClient, trx.CommisionProvider,
		trx.DateInput.Time, trx.DatePost.Time, trx.Status, trx.PaymentType, trx.PaymentNumber, trx.ServiceID, trx.Service, trx.PayeeID, trx.PayeeName, trx.PayeeBankMfo,
		trx.PayeeBankAccount, trx.PaymentNarrative,
	)
	if err != nil {
		return err
	}
	return nil
}

// insert many transactions
func (db *TransactionDB) InsertTransactions(ctx context.Context, trxs []domain.Transaction) error {
	if trxs == nil {
		return fmt.Errorf("trxs is nil")
	}
	//paginate trxs because pgx batch can overflow, so we are sending 100 transactions in single batch
	pages := int(math.Ceil(float64(len(trxs)) / 100))
	log.Println("len trxs:", len(trxs))
	log.Println("pages:", pages)
	b := pgx.Batch{}
	for i := range trxs {
		trx := trxs[i]
		b.Queue(`INSERT INTO transactions (id, requestid, terminalid, partnerobjectid,
		 amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
		 dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid,
		 servicetypeid, payeeid, payeenameid, payeebankmfo, payeebankaccount, paymentnarrativeid)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, (SELECT id FROM servicetypes WHERE title=$16),
		 $17, (SELECT id FROM payeenames WHERE title=$18), $19, $20, (SELECT id FROM paymentnarratives WHERE title=$21))`,
			trx.ID, trx.RequestID, trx.TerminalID, trx.PartnerObjectID, trx.AmountTotal, trx.AmountOriginal, trx.CommisionPs, trx.CommisionClient, trx.CommisionProvider,
			trx.DateInput.Time, trx.DatePost.Time, trx.Status, trx.PaymentType, trx.PaymentNumber, trx.ServiceID, trx.Service, trx.PayeeID, trx.PayeeName, trx.PayeeBankMfo,
			trx.PayeeBankAccount, trx.PaymentNarrative)
	}
	bres := db.pool.SendBatch(ctx, &b)
	err := bres.Close()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("duplicate data")
		}
		return fmt.Errorf("error sending batch: %w", err)

	}
	return nil
}

// get single transaction by id
func (db *TransactionDB) GetTransactionByID(ctx context.Context, id int) (*domain.Transaction, error) {
	//use batch to avoid unnessesary network calls
	b := pgx.Batch{}
	//get fields from table
	b.Queue(`SELECT id, requestid, terminalid, partnerobjectid,
	amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
	dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid, payeeid, payeebankmfo, payeebankaccount
	FROM transactions
	WHERE id = $1`, id)
	//get fields that referenced with foreign keys
	b.Queue(`WITH payeenameid AS (SELECT payeenameid FROM transactions WHERE id = $1)
	SELECT title FROM payeenames, payeenameid WHERE payeenames.id = payeenameid`, id)
	b.Queue(`WITH narrativeid AS (SELECT paymentnarrativeid AS narrativeid FROM transactions WHERE id = $1)
	SELECT title FROM paymentnarratives, narrativeid WHERE paymentnarratives.id = narrativeid`, id)
	b.Queue(`WITH servicetypeid AS (SELECT servicetypeid FROM transactions WHERE id = $1)
	SELECT title FROM servicetypes, servicetypeid WHERE servicetypes.id = servicetypeid`, id)
	bres := db.pool.SendBatch(ctx, &b)
	var result domain.Transaction
	row := bres.QueryRow()
	err := row.Scan(&result.ID, &result.RequestID, &result.TerminalID, &result.PartnerObjectID, &result.AmountTotal,
		&result.AmountOriginal, &result.CommisionPs, &result.CommisionClient, &result.CommisionProvider, &result.DateInput.Time,
		&result.DatePost.Time, &result.Status, &result.PaymentType, &result.PaymentNumber, &result.ServiceID, &result.PayeeID,
		&result.PayeeBankMfo, &result.PayeeBankAccount,
	)
	if err != nil {
		return nil, err
	}
	row = bres.QueryRow()
	err = row.Scan(&result.PayeeName)
	if err != nil {
		return nil, err
	}

	row = bres.QueryRow()
	err = row.Scan(&result.PaymentNarrative)
	if err != nil {
		return nil, err
	}

	row = bres.QueryRow()
	err = row.Scan(&result.Service)
	if err != nil {
		return nil, err
	}
	bres.Close()
	return &result, nil
}

// get transactions with filters
func (db *TransactionDB) GetTransactions(ctx context.Context, filters map[string]interface{}) ([]domain.Transaction, error) {
	//build sql query with filters
	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqltrx := build.Select(`transactions.id, requestid, terminalid, partnerobjectid,
	amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
	dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid, servicetypeid, payeeid,
	payeenameid, payeebankmfo, payeebankaccount, paymentnarrativeid`).From("transactions")
	if f, ok := filters["terminal_ids"]; ok {
		terminal_ids, ok := f.([]interface{})
		if !ok {
			return nil, fmt.Errorf("terminal_ids filter is not of type []interface{}")
		}
		if len(terminal_ids) > 0 {
			b := strings.Builder{}
			b.WriteString("(?")
			for i := 1; i < len(terminal_ids); i++ {
				b.WriteString(", ?")
			}
			b.WriteString(")")
			sqltrx = sqltrx.Where("terminalid IN "+b.String(), terminal_ids...)
		}
	}
	if status, ok := filters["status"]; ok {
		sqltrx = sqltrx.Where("statusid = ?", status)
	}
	if payment_type, ok := filters["payment_type"]; ok {
		sqltrx = sqltrx.Where("paymenttype = ?", payment_type)
	}
	if date_from, ok := filters["date_post_from"]; ok {
		sqltrx = sqltrx.Where("datepost >= ?", date_from)
	}
	if date_to, ok := filters["date_post_to"]; ok {
		sqltrx = sqltrx.Where("datepost <= ?", date_to)
	}
	if payment_narrative, ok := filters["payment_narrative"]; ok {
		sqltrx = sqltrx.Join("paymentnarratives ON transactions.paymentnarrativeid = paymentnarratives.id").Where("paymentnarratives.title LIKE ?", "%"+payment_narrative.(string)+"%")
	}
	sql, args, err := sqltrx.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error converting selectbuilder to string: %w", err)
	}
	fmt.Println(sql)
	fmt.Println(args)

	//use filtered transactions AS t
	withClause := "WITH t AS (\n" + sql + "\n)\n"
	//apply joins with referenced tables and select fields
	selectClause := `SELECT t.id, t.requestid, t.terminalid, t.partnerobjectid,
	t.amounttotal, t.amountoriginal, t.commisionps, t.commisionclient, t.commisionprovider,
	t.dateinput, t.datepost, t.statusid, t.paymenttype, t.paymentnumber, t.serviceid, servicetypes.title, t.payeeid,
	payeenames.title, payeebankmfo, payeebankaccount, paymentnarratives.title FROM servicetypes
	INNER JOIN t
	ON t.servicetypeid = servicetypes.id
	INNER JOIN payeenames
	ON t.payeenameid = payeenames.id
	INNER JOIN paymentnarratives
	ON t.paymentnarrativeid = paymentnarratives.id
	ORDER BY t.id`
	//run query
	rows, err := db.pool.Query(ctx, withClause+selectClause, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying db: %w", err)
	}
	result := make([]domain.Transaction, 0, 10)
	for rows.Next() {
		cur := domain.Transaction{}
		err = rows.Scan(
			&cur.ID, &cur.RequestID, &cur.TerminalID, &cur.PartnerObjectID, &cur.AmountTotal,
			&cur.AmountOriginal, &cur.CommisionPs, &cur.CommisionClient, &cur.CommisionProvider, &cur.DateInput.Time,
			&cur.DatePost.Time, &cur.Status, &cur.PaymentType, &cur.PaymentNumber, &cur.ServiceID, &cur.Service, &cur.PayeeID,
			&cur.PayeeName, &cur.PayeeBankMfo, &cur.PayeeBankAccount, &cur.PaymentNarrative,
		)
		result = append(result, cur)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
	}
	rows.Close()
	return result, nil
}
