package db

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/asstronom/EVO_tech_test/pkg/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionDB struct {
	pool *pgxpool.Pool
}

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

func (db *TransactionDB) Close() {
	db.pool.Close()
}

func (db *TransactionDB) InsertTransaction(ctx context.Context, trx domain.Transaction) error {
	_, err := db.pool.Exec(ctx, `INSERT INTO transactions (id, requestid, terminalid, partnerobjectid,
		 amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
		 dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid,
		 servicetypeid, payeeid, payeenameid, payeebankmfo, payeebankaccount, paymentnarrativeid)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, (SELECT id FROM servicetypes WHERE title=$16),
		 $17, (SELECT id FROM payeenames WHERE title=$18), $19, $20, (SELECT id FROM paymentnarratives WHERE title=$21))`,
		trx.ID, trx.RequestID, trx.TerminalID, trx.PartnerObjectID, trx.AmountTotal, trx.AmountOriginal, trx.CommisionPs, trx.CommisionClient, trx.CommisionProvider,
		trx.DateInput, trx.DatePost, trx.Status, trx.PaymentType, trx.PaymentNumber, trx.ServiceID, trx.Service, trx.PayeeID, trx.PayeeName, trx.PayeeBankMfo,
		trx.PayeeBankAccount, trx.PaymentNarrative,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *TransactionDB) InsertTransactions(ctx context.Context, trxs []domain.Transaction) error {
	if trxs == nil {
		return fmt.Errorf("trxs is nil")
	}
	pages := len(trxs) / 100
	if pages == 0 {
		pages = 1
	}
	for i := 0; i < pages; i++ {
		b := pgx.Batch{}
		for j := 0; j < 100 && j+i*100 < len(trxs); j++ {
			trx := trxs[i*100+j]
			b.Queue(`INSERT INTO transactions (id, requestid, terminalid, partnerobjectid,
		 amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
		 dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid,
		 servicetypeid, payeeid, payeenameid, payeebankmfo, payeebankaccount, paymentnarrativeid)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, (SELECT id FROM servicetypes WHERE title=$16),
		 $17, (SELECT id FROM payeenames WHERE title=$18), $19, $20, (SELECT id FROM paymentnarratives WHERE title=$21))`,
				trx.ID, trx.RequestID, trx.TerminalID, trx.PartnerObjectID, trx.AmountTotal, trx.AmountOriginal, trx.CommisionPs, trx.CommisionClient, trx.CommisionProvider,
				trx.DateInput, trx.DatePost, trx.Status, trx.PaymentType, trx.PaymentNumber, trx.ServiceID, trx.Service, trx.PayeeID, trx.PayeeName, trx.PayeeBankMfo,
				trx.PayeeBankAccount, trx.PaymentNarrative)
		}
		bres := db.pool.SendBatch(ctx, &b)
		bres.Close()
	}
	return nil
}

func (db *TransactionDB) GetTransactionByID(ctx context.Context, id int) (*domain.Transaction, error) {
	b := pgx.Batch{}
	b.Queue(`SELECT id, requestid, terminalid, partnerobjectid,
	amounttotal, amountoriginal, commisionps, commisionclient, commisionprovider,
	dateinput, datepost, statusid, paymenttype, paymentnumber, serviceid, payeeid, payeebankmfo, payeebankaccount
	FROM transactions
	WHERE id = $1`, id)
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
		&result.AmountOriginal, &result.CommisionPs, &result.CommisionClient, &result.CommisionProvider, &result.DateInput,
		&result.DatePost, &result.Status, &result.PaymentType, &result.PaymentNumber, &result.ServiceID, &result.PayeeID,
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

func (db *TransactionDB) GetTransactions(ctx context.Context, filters map[string]interface{}) ([]domain.Transaction, error) {
	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	trxs := build.Select(`transactions.id, transactions.requestid, transactions.terminalid,
	transactions.partnerobjectid, transactions.amounttotal, transactions.amountoriginal,
	transactions.commisionps, transactions.commisionclient, transactions.commisionprovider,
	transactions.dateinput, transactions.datepost, transactions.statusid,
	transactions.paymenttype, transactions.paymentnumber, transactions.serviceid,
	transactions.servicetypeid, transactions.payeeid, transactions.payeenameid, 
	transactions.payeebankmfo, transactions.payeebankaccount, transactions.paymentnarrativeid`).
		From("transactions")

	if terminal_ids, ok := filters["terminal_ids"]; ok {
		ids, ok := terminal_ids.([]interface{})
		if !ok {
			return nil, fmt.Errorf("wrong type for terminal_ids filter")
		}
		if len(ids) > 0 {
			bb := strings.Builder{}
			bb.WriteString("(" + "?")
			for i := 1; i < len(ids); i++ {
				bb.WriteString(", " + "?")
			}
			bb.WriteString(")")
			trxs = trxs.Where("terminalid IN "+bb.String(), ids...)
		}
	}
	if status, ok := filters["status"]; ok {
		trxs = trxs.Where("statusid = ?", status)
	}
	if payment_type, ok := filters["payment_type"]; ok {
		trxs = trxs.Where("paymenttype = ?", payment_type)
	}
	if date_from, ok := filters["date_post_from"]; ok {
		trxs = trxs.Where("datepost >= ?", date_from)
	}
	if date_to, ok := filters["date_post_to"]; ok {
		trxs = trxs.Where("datepost <= ?", date_to)
	}
	if payment_narrative, ok := filters["payment_narrative"]; ok {
		trxs = trxs.InnerJoin("paymentnarratives ON transactions.paymentnarrativeid = paymentnarratives.id").
			Where("paymentnarratives.title LIKE ?", payment_narrative)
	}

	sql, args, err := trxs.ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(sql)
	fmt.Println(args)

	// withClause := `WITH trxs AS (` + sql + `)` + "\n"
	// selectClause := `SELECT 
	// trxs.id, trxs.requestid, trxs.terminalid,
	// trxs.partnerobjectid, trxs.amounttotal, trxs.amountoriginal,
	// trxs.commisionps, trxs.commisionclient, trxs.commisionprovider,
	// trxs.dateinput, trxs.datepost, trxs.statusid,
	// trxs.paymenttype, trxs.paymentnumber, trxs.serviceid,
	// servicetypes.title, trxs.payeeid, payeenames.title, 
	// trxs.payeebankmfo, trxs.payeebankaccount, paymentnarratives.title
	// FROM payeenames
	// INNER JOIN trxs
	// ON trxs.payeenameid = payeenames.id
	// INNER JOIN paymentnarratives
	// ON paymentnarratives.id = trxs.paymentnarrativeid
	// INNER JOIN servicetypes
	// ON trxs.servicetypeid = servicetypes.id`

	// fmt.Println(withClause + selectClause)

	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	result := []domain.Transaction{}
	for rows.Next() {
		cur := domain.Transaction{}
		rows.Scan(&cur.ID)
		result = append(result, cur)
	}
	return result, nil
}
