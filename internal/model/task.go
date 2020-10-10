package model

import (
	"context"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	NextTask    int    `json:"next_task"`
	Coords      Point  `json:"coords"`
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

	_, err = db.ExecContext(ctx, q, gameSN, playerName)
	if err != nil {
		log.WithError(err).Errorf("can't complete task for %s player in game %s", playerName, gameSN)
		return err
	}

	return nil
}

func NewTaskBatch(ctx context.Context, db *sqlx.DB, gameSN string, tasks []Task) error {
	var gid int
	var err error

	err = db.GetContext(ctx, &gid, `select game_id from game where short_name = $1`, gameSN)
	if err != nil {
		log.WithError(err).Error("task batch")
		return err
	}

	q := "insert into task (game_id, title, description, coord_x, coord_y) values ($1,$2,$3,$4,$5) returning task_id"

	for i := len(tasks) - 1; i >= 0; i++ {
		var tid int
		err = db.GetContext(ctx, &tid, q, gid, tasks[i].Title, tasks[i].Description, tasks[i].Coords.X, tasks[i].Coords.Y)
		if err != nil {
			log.WithError(err).Error("task batch")
			return err
		}
	}

	return nil
}
