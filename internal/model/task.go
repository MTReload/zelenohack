package model

import (
	"context"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	NextTask    int    `json:"next_task"`
}

func CompleteTask(ctx context.Context, db *sqlx.DB, playerName, gameSN string) error {
	var err error

	q := `update player_on_task
set task_id = (select next_task
    from task
             join player_on_task pot on task.task_id = pot.task_id
             join game g on g.game_id = task.game_id
    join player p on pot.player_id = p.player_id
    where g.short_name = $1 and  p.name = $2 limit 1
) where player_id = (select player_id from player where name=$2)`

	_, err = db.ExecContext(ctx, q)
	if err != nil {
		log.WithError(err).Errorf("can't complete task for %s player in game %s", playerName, gameSN)
		return err
	}

	return nil
}
