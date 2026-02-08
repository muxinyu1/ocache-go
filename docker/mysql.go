package docker

import (
	"database/sql"
	"fmt"
	"iter"
	"math/rand/v2"
	"ocache/impl"
	"ocache/ocache"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlGetter struct {
	db *sql.DB
}

func NewMysqlGetter(dsn string) (*MysqlGetter, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &MysqlGetter{db: db}, nil
}

func (m *MysqlGetter) Get(group string, key string) (ocache.Value, error) {
	var v []byte
	row := m.db.QueryRow("SELECT v FROM kv_store WHERE group_name = ? AND k = ?", group, key)
	if err := row.Scan(&v); err != nil {
		return nil, fmt.Errorf("mysql get error: %v", err)
	}
	// 模拟真实数据库延迟
	timeToSleep := rand.IntN(10) + 1
	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	return impl.NewByteView(v), nil
}

func (m *MysqlGetter) All() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		rows, err := m.db.Query("SELECT k FROM kv_store")
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var k string
			if err := rows.Scan(&k); err != nil {
				return
			}
			if !yield([]byte(k)) {
				return
			}
		}
	}
}
