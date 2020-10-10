package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewPlayer(ctx context.Context, db *sqlx.DB, playerName, gameShortName string) (interface{}, error) {
	var err error
	qNewPlayer := `insert into player (name)
values ($1)
on conflict DO NOTHING
returning json_build_object(
        'player_id', player_id,
        'name', name
    )`

	qInitPlayerOnGame := `insert into player_game (player_id, game_id) (select player_id, game_id
                                              from player,
                                                   game
                                              where player.name = $1
                                                and game.short_name = $2)`

	qSetPlayerTasks := `insert into player_on_task (player_id, task_id)
    (select p.player_id, task_id
     from player p
              join player_game pg on p.player_id = pg.player_id
              join task t on pg.game_id = t.game_id
    join game g on pg.game_id = g.game_id
     where p.name = $1
       and g.short_name = $2
     order by t.task_id
     limit 1)`

	var b []byte

	var ret Player

	err = db.QueryRowxContext(ctx, qNewPlayer, playerName).Scan(&b)
	if err != nil {
		fmt.Printf("%s: can't add new player\n", err.Error())
		return nil, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		log.WithError(err).Error("can't unmarshal player")
		return nil, err
	}

	tx, err := db.BeginTx(ctx, nil)

	_, err = tx.ExecContext(ctx, qInitPlayerOnGame, playerName, gameShortName)
	if err != nil {
		fmt.Printf("%s: can't init player on game\n", err.Error())
		return nil, err
	}

	_, err = tx.ExecContext(ctx, qSetPlayerTasks, playerName, gameShortName)
	if err != nil {
		fmt.Printf("%s: can't set new player tasks\n", err.Error())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.WithError(err).Error()
		err = tx.Rollback()
		if err != nil {
			log.WithError(err).Error()
		}
		return nil, err
	}

	return GetGameInfoForPlayer(ctx, db, gameShortName, playerName)
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
	var b []byte

	err = db.QueryRowxContext(ctx, q, name).Scan(&b)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Printf("%s: can't get player %s\n", err.Error(), name)
		return nil, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		log.WithError(err).Error("can't unmarshal player")
		return nil, err
	}

	return &ret, nil
}
