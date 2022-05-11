package lib

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/pgtype"
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

func (db *DB) SaveLastTokenHolders(network int, token string, holders []CovalentHolder) error {
	batch := &pgx.Batch{}

	batch.Queue("DELETE FROM token_holders_last WHERE network = $1 AND token = $2", network, token)

	for _, x := range holders {
		batch.Queue(
			"INSERT INTO token_holders_last(network, token, holder, amount) VALUES ($1, $2, $3, $4)",
			network, token, x.Address, x.Balance,
		)
	}

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetLastTokenHolders(network int, token string) ([]NFTHolder, error) {
	holders := []NFTHolder{}

	rows, err := db.conn.Query(context.Background(),
		`SELECT T0.holder, T0.amount, 
			CASE
				WHEN T3.total_balance_usd IS NULL THEN 0.0
				ELSE T3.total_balance_usd
			END as total_balance_usd
		FROM token_holders_last T0
		LEFT JOIN (
			SELECT T1.holder, SUM(T2.amount_usd) as total_balance_usd
			FROM token_holders_last T1
			LEFT JOIN (
				SELECT * FROM balances
			) T2 on T1.holder = T2.address
			WHERE T1.network = $1 AND T1.token = $2 
			GROUP BY T1.holder
		) T3 on T0.holder = T3.holder
		WHERE T0.network = $1 AND T0.token = $2`,
		network, token,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holder string
	var amountStr string
	var balance float64
	for rows.Next() {
		rows.Scan(&holder, &amountStr, &balance)

		var amount, err = strconv.Atoi(amountStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %v", err)
		}
		holders = append(holders, NFTHolder{Address: holder, Amount: int64(amount), TotalBalanceUsd: balance})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return holders, nil
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

func (db *DB) SaveOpenseaTrades(network int64, events []BitqueryEvent) error {
	batch := &pgx.Batch{}

	for _, e := range events {
		t, err := ParseOpenseaTradeBitqueryEvent(e)
		if err != nil {
			return err
		}

		batch.Queue(
			"INSERT INTO opensea_trades(network, block, time, tx, buy_hash, sell_hash, maker, taker, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
			network, t.block, t.time, t.tx, t.buyHash, t.sellHash, t.maker, t.taker, t.price,
		)
	}

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SaveAaveDeposits(network int64, events []BitqueryEvent) error {
	batch := &pgx.Batch{}

	for _, e := range events {
		d, err := ParseAaveDepositBitqueryEvent(e)
		if err != nil {
			return err
		}

		batch.Queue(
			"INSERT INTO aave_deposits(network, block, time, tx, \"user\", on_behalf_of, reserve, amount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			network, d.block, d.time, d.tx, d.user, d.onBehalfOf, d.reserve, d.amount,
		)
	}

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SaveLooksRareTrades(network int64, events []BitqueryEvent) error {
	batch := &pgx.Batch{}

	for _, e := range events {
		t, err := ParseLooksRareTradeBitqueryEvent(e)
		if err != nil {
			return err
		}

		batch.Queue(
			"INSERT INTO looksrare_trades(network, block, time, tx, maker, taker, currency, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			network, t.block, t.time, t.tx, t.maker, t.taker, t.currency, t.price,
		)
	}

	br := db.conn.SendBatch(context.Background(), batch)
	err := br.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAllNFTs() ([]NFTCollection, error) {
	nfts := []NFTCollection{}

	tokens := []string{
		"0x026224a2940bfe258d0dbe947919b62fe321f042",
		"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
		"0x1A92f7381B9F03921564a437210bB9396471050C",
		"0x938e5ed128458139a9c3306ace87c60bcba9c067",
		"0x23581767a106ae21c074b2276D25e5C3e136a68b",
	}
	rows, err := db.conn.Query(context.Background(),
		"SELECT address, name, symbol, logo FROM tokens WHERE address = ANY ($1) ORDER BY id",
		tokens,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var address string
	var name pgtype.Text
	var symbol pgtype.Text
	var logo pgtype.Text
	for rows.Next() {
		rows.Scan(&address, &name, &symbol, &logo)
		nfts = append(nfts, NFTCollection{Address: address, Name: name.String, Symbol: symbol.String, Logo: logo.String})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return nfts, nil
}

func (db *DB) GetToken(network int64, token string) (*Token, error) {
	rows, err := db.conn.Query(context.Background(),
		"SELECT name, symbol, logo FROM tokens WHERE network = $1 AND address = $2",
		network, token,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var name pgtype.Text
		var symbol pgtype.Text
		var logo pgtype.Text
		rows.Scan(&name, &symbol, &logo)
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		token := &Token{
			Network: network, Address: token,
			Name: name.String, Symbol: symbol.String, Logo: logo.String,
		}
		return token, nil
	} else {
		return nil, errors.New("no rows")
	}
}

func (db *DB) GetNFTTokenHoldings(network int, token string, nft bool, limit int) ([]HoldingRow, error) {
	holdings := []HoldingRow{}

	var tableName string
	if nft {
		tableName = "nft_holdings"
	} else {
		tableName = "token_holdings"
	}
	rows, err := db.conn.Query(context.Background(),
		fmt.Sprintf(
			`SELECT T1.holding_token as token_address, T2.name as token_name, T2.symbol as token_symbol, T2.decimals as token_decimals, T2.logo as token_logo, T1.holders, T1.holders_share 
		FROM %s T1
		LEFT JOIN (
			SELECT * FROM tokens
		) T2 on T1.holding_token = T2.address
		WHERE T1.network = $1 AND T1.token = $2 
		ORDER BY holders 
		DESC LIMIT $3`,
			tableName,
		),
		network, token, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokenAddress string
	var tokenName pgtype.Text
	var tokenSymbol pgtype.Text
	var tokenDecimals pgtype.Int8
	var tokenLogo pgtype.Text
	var holders int64
	var share float64
	for rows.Next() {
		rows.Scan(&tokenAddress, &tokenName, &tokenSymbol, &tokenDecimals, &tokenLogo, &holders, &share)
		t := Token{
			Network:  1,
			Address:  tokenAddress,
			Name:     tokenName.String,
			Symbol:   tokenSymbol.String,
			Logo:     tokenLogo.String,
			Decimals: tokenDecimals.Int,
		}
		holdings = append(holdings, HoldingRow{Token: t, Holders: holders, Share: share})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return holdings, nil
}

func (db *DB) GetNFTProtocols(network int, token string, limit int) ([]ProtocolUsage, error) {
	list := []ProtocolUsage{}

	rows, err := db.conn.Query(context.Background(), `
			SELECT p.name as protocolName,
						 p.logo as protocolLogo,
						 p.url as protocolUrl,
						 pu.users_last_month as usersLastMonth,
						 pu.users_in_total as usersInTotal
			FROM protocols_usage pu
							 JOIN protocols p ON p.id = pu.protocol_id
			WHERE pu.nft_id = (
					SELECT id
					FROM tokens
					WHERE network = $1 AND address = $2 
			)
			ORDER BY pu.users_last_month DESC,
         			 pu.users_in_total DESC
		  LIMIT $3
		`, network, token, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var protocolName pgtype.Text
	var protocolLogo pgtype.Text
	var protocolUrl pgtype.Text
	var usersLastMonth int64
	var usersInTotal int64
	for rows.Next() {
		rows.Scan(&protocolName, &protocolLogo, &protocolUrl, &usersLastMonth, &usersInTotal)
		list = append(list, ProtocolUsage{
			Name:           protocolName.String,
			Logo:           protocolLogo.String,
			Url:            protocolLogo.String,
			UsersLastMonth: usersLastMonth,
			UsersInTotal:   usersInTotal,
		})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return list, nil
}

func (db *DB) GetNFTStats(token string) (*NFTStats, error) {
	rows, err := db.conn.Query(context.Background(), `
		SELECT SUM(ha.day::int) as day, SUM(ha.week::int) as week, SUM(ha.month::int) as month, COUNT(th.holder) as total 
 		FROM token_holders_last th
			INNER JOIN holders_activity ha ON th.holder = ha.holder
		WHERE th.token = $1`,
		token,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var stats *NFTStats
		var day pgtype.Int8
		var week pgtype.Int8
		var month pgtype.Int8
		var total pgtype.Int8
		rows.Scan(&day, &week, &month, &total)
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		stats = &NFTStats{Active1d: day.Int, Active7d: week.Int, Active30d: month.Int, Total: total.Int}
		return stats, nil
	} else {
		return nil, errors.New("no rows")
	}
}
