package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
	"go.uber.org/zap"
)

type Currencies_psq struct {
	Conn *sql.DB
	l    *zap.SugaredLogger
}

func NewCurrencyRepo(Conn *sql.DB, logger *zap.SugaredLogger) *Currencies_psq {
	return &Currencies_psq{
		Conn: Conn,
		l:    logger,
	}
}

func (c *Currencies_psq) fetch(ctx context.Context, stmt string, filters []interface{}) (ts []*domain.CurrencyRepository, err error) {
	currencies := []*domain.CurrencyRepository{}
	query := `SELECT * FROM currency ` + stmt
	rows, err := c.Conn.QueryContext(ctx, query, filters...)
	if err != nil {
		c.l.Error(err.Error())
		return nil, errors.New("query_context")
	}
	defer rows.Close()

	for rows.Next() {
		currency := &domain.CurrencyRepository{}
		err := rows.Scan(&currency.ID, &currency.CreatedAt, &currency.UpdatedAt, &currency.Code, &currency.Value)
		if err != nil {
			c.l.Error(err.Error())
			return nil, errors.New(err.Error())
		}
		currencies = append(currencies, currency)
	}
	if err := rows.Err(); err != nil {
		c.l.Error(err.Error())
		return nil, errors.New("row_corrupt")
	}

	return currencies, nil
}

func (c *Currencies_psq) GetByCode(ctx context.Context, code string, f domain.Filter) (t []*domain.CurrencyRepository, err error) {
	c.l.Info("💾 Getting by ID")
	cont := 1
	query := fmt.Sprintf("WHERE code=$%d ", cont)
	filter := []interface{}{code}
	if code == "ALL"{
		query = fmt.Sprintf("WHERE code!=$%d ", cont)
	}
	if f.Finit != nil {
		cont += 1
		query += fmt.Sprintf("AND created_at > $%d ", cont)
		filter = append(filter, f.Finit)
	}
	if f.Fend != nil {
		cont += 1
		query += fmt.Sprintf("AND created_at < $%d ", cont)
		filter = append(filter, f.Fend)
	}
	query = fmt.Sprintf("%s ORDER BY created_at ASC LIMIT $%d OFFSET $%d",
		query, cont+1, cont+2)
	filter = append(filter, f.Limit)
	filter = append(filter, f.Offset)

	currencies, err := c.fetch(ctx, query, filter)
	if err != nil {
		c.l.Info(err.Error())
		return nil, err
	}

	if len(currencies) == 0 {
		c.l.Infof("Record %s Not Found", code)
		return nil, errors.New("not_found")
	}

	return currencies, nil
}

func (c *Currencies_psq) InsertBulk(ctx context.Context, ta domain.Currencies) (err error) {
	c.l.Infoln("Insert Bulk", ta)

	values_string := []string{} // ($1, $2, $3, $4)
	values := []interface{}{}
	created_at := time.Now()
	count := 0
	for _, v := range ta {
		fields := fmt.Sprintf("($%d,$%d,$%d,$%d)",
			count+1, count+2, count+3, count+4)
		count += 4
		values_string = append(values_string, fields)
		values = append(values, v.Code, v.Value, created_at, created_at)
	}

	query := "INSERT INTO currency(code, value, created_at, updated_at) VALUES " + strings.Join(values_string, ",") + ";"
	fmt.Println(query)

	stmt, err := c.Conn.PrepareContext(ctx, query)
	if err != nil {
		c.l.Error(err.Error())
		return errors.New("query_prepare_ctx")
	}

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		c.l.Error(err.Error())
		return errors.New("query_exec")
	}
	affect, err := res.RowsAffected()
	if err != nil {
		c.l.Error(err.Error())
		return errors.New("query_exec")
	}
	c.l.Infof("%d currencies saved", affect)
	return nil
}
