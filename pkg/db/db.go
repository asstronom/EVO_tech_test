package db

import (
	"context"

	"github.com/asstronom/EVO_tech_test/pkg/domain"
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
