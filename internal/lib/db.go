package lib

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	conn *pgxpool.Pool
}

func CreateDB() (*DB, error) {
	connString := os.Getenv("DB_CONN_STRING")
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &DB{conn}, nil
}

func (db *DB) Close() {
	db.conn.Close()
}

func (db *DB) SaveBlockTime(network int, timestamp int64, block int64) error {
	batch := &pgx.Batch{}
	batch.Queue(
		"INSERT INTO time_block (network, time, block)"+
			" VALUES ($1, $2, $3)"+
			" ON CONFLICT (network, time)"+
			" DO NOTHING",
		network, timestamp, block,
	)

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetDailyBlocks(network int) (map[int64]int64, error) {
	rows, err := db.conn.Query(context.Background(),
		"SELECT time, block FROM time_block WHERE network = $1 ORDER BY time ASC",
		network,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timeBlocks := make(map[int64]int64)

	var time int64
	var block int64
	for rows.Next() {
		rows.Scan(&time, &block)
		timeBlocks[time] = block
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return timeBlocks, nil
}

func (db *DB) SaveCovalentTokenHolders(network int, token string, time int64, holders []CovalentHolder) error {
	batch := &pgx.Batch{}

	batch.Queue("DELETE FROM token_holders_x WHERE network = $1 AND token = $2 AND time = $3", network, token, time)

	for _, x := range holders {
		batch.Queue(
			"INSERT INTO token_holders_x(network, token, time, holder, amount) VALUES ($1, $2, $3, $4, $5)",
			network, token, time, x.Address, x.Balance,
		)
	}

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetDailyTokenHolders(network int, token string, fromTime int64) (map[int64]map[string]string, error) {
	m := make(map[int64]map[string]string)

	rows, err := db.conn.Query(context.Background(),
		"SELECT time, holder, amount FROM token_holders_x WHERE network = $1 AND token = $2 AND time >= $3",
		network, token, fromTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timestamp int64
	var holder string
	var amount string
	i := 0
	for rows.Next() {
		i++
		rows.Scan(&timestamp, &holder, &amount)
		if m[timestamp] == nil {
			m[timestamp] = make(map[string]string)
		}
		m[timestamp][holder] = amount
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return m, nil
}

func (db *DB) SaveOpenseaTrades(events []OpenseaTradeEvent) error {
	batch := &pgx.Batch{}

	for _, e := range events {
		t, err := OpenseaEventToTrade(e)
		if err != nil {
			return err
		}

		batch.Queue(
			"INSERT INTO opensea_trades(block, time, tx, maker, taker, price) VALUES ($1, $2, $3, $4, $5, $6)",
			t.block, t.time, t.tx, t.maker, t.taker, t.price,
		)
	}

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}
