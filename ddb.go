package ddb

import (
	"github.com/jmoiron/sqlx"
)

const (
	TblHitBTCTicker = "HitBTCTicker"
	TblHitBTCOrders = "HitBTCOrders"
)

var (
	db *sqlx.DB
)
