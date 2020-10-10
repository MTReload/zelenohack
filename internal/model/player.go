package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewPlayer(ctx context.Context, db *sqlx.DB, name string) (*Player, error) {
	var err error
	q := `insert into player (name)
values ($1)
on conflict do nothing 
returning json_build_object(
        'player_id', player_id,
        'name', name,
        );`

	var ret Player
	err = db.QueryRowxContext(ctx, q, name).Scan(&ret)
	if err != nil {
		fmt.Printf("%s: can't add new player\n", err.Error())
		return nil, err
	}

	return &ret, nil
}

func PlayerByName(ctx context.Context, db *sqlx.DB, name string) (*Player, error) {
	var err error
	q := `select json_build_object(
               'id', player_id,
               'name', name
           )
from player
where name = $1`

	var ret Player
	err = db.QueryRowxContext(ctx, q, name).Scan(&ret)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Printf("%s: can't get player %s\n", err.Error(), name)
		return nil, err
	}

	return &ret, nil
}
