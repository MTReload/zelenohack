package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Game struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ShortName   string `json:"short_name"`
}

func NewGame(ctx context.Context, db *sqlx.DB, game Game) (*Game, error) {
	var err error
	q := `insert into game (name, description, short_name)
values ($1, $2, $3)

returning json_build_object(
        'name', name,
        'description', description,
        'short_name', short_name,
    'id', game_id);`

	var ret Game
	var b []byte

	err = db.QueryRowxContext(ctx, q, game.Name, game.Description, game.ShortName).Scan(&b)
	if err != nil {
		fmt.Printf("%s: can't add new game\n", err.Error())
		return nil, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		fmt.Printf("%s: can't add new game\n", err.Error())
		return nil, err
	}

	return &ret, nil
}

func GameByShortName(ctx context.Context, db *sqlx.DB, shortName string) (*Game, error) {
	var err error
	q := `select json_build_object(
               'id', game_id,
               'short_name', short_name,
               'description', description,
               'name', name
           )
from game
where short_name = $1`

	var ret Game
	var b []byte

	err = db.QueryRowxContext(ctx, q, shortName).Scan(&b)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Printf("%s: can't get game %s\n", err.Error(), shortName)
		return nil, err
	}

	err = json.Unmarshal(b, &ret)

	if err != nil {
		log.WithError(err).Error("can't unmarshal game")
		return nil, err
	}

	return &ret, nil
}

func GameInfo(ctx context.Context, db *sqlx.DB, shortName string) (interface{}, error) {
	var err error
	q := `select json_build_object(
               'game', json_build_object('id', game.game_id, 'name', game.name, 'short_name', game.short_name, 'description', game.description),
               'tasks', array(select json_build_object(
                                             'id', task_id,
                                             'title', task.title,
                                             'description', task.description,
                                             'next_task', task.next_task,
                                             'coords', json_build_object('x', task.coord_x, 'y', task.coord_y)
                                         )
                              from task
                                       join game g on g.game_id = task.game_id
                              where g.short_name = $1),
               'players', array(select json_build_object(
                                               'id', p.player_id,
                                               'name', p.name,
                                               'task_id', pg.task_id
                                           )
                                from player p
                                         join player_game pg on p.player_id = pg.player_id
                                where pg.game_id = (select game_id from game where game.short_name = $1)
                   )
           )
from game
where game.short_name = $1`

	var ret interface{}
	var b []byte

	err = db.GetContext(ctx, &b, q, shortName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.WithError(err).Error("can't get full game info")
		return nil, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		log.WithError(err).Error("can't unmarshal game info")
		return nil, err
	}
	return ret, nil
}

func GetGameInfoForPlayer(ctx context.Context, db *sqlx.DB, gameSN, playerName string) (interface{}, error) {
	var err error

	q := `select json_build_object(
               'game', json_build_object('id', game.game_id, 'name', game.name, 'short_name', game.short_name),
               'tasks', array(select json_build_object(
                                             'id', task_id,
                                             'title', t.title,
                                             'description', t.description,
                                             'next_task', t.next_task,
                                             'coords', json_build_object('x', t.coord_x, 'y', t.coord_y)
                                         )
                ),
               'now_on', pg.task_id
           )
from game
         join player_game pg on game.game_id = pg.game_id
         join player p on p.player_id = pg.player_id
		join task t on game.game_id = t.game_id
where game.short_name = $1
  and p.name = $2`

	var ret interface{}
	var b []byte

	err = db.GetContext(ctx, &b, q, gameSN, playerName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.WithError(err).Error("can't get personal game info")
		return nil, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		log.WithError(err).Error("can't unmarshal personal game info")
		return nil, err
	}
	return ret, nil
}
