package sqlstore

import (
	"fmt"
	"warehouse/internal/app/models"

	"github.com/sirupsen/logrus"
)

func (r *Repo) GetAllComponents() (interface{}, error) {
	rows, err := r.store.db.Query(`
	select c.available, c.id, c.code, c."name", c2."name" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c."time", 'DD-MM-YYYY HH24:MI') "time", 
	t."name" as type, t.id as type_id, c.weight from components c 
	join checkpoints c2 on c2.id = c."checkpoint" join "types" t on t.id  = c."type" 
	where c.status = 1
	order by c.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []models.Component

	for rows.Next() {
		var comp models.Component
		if err := rows.Scan(&comp.Available, &comp.ID, &comp.Code,
			&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit, &comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id, &comp.Weight); err != nil {
			return components, err
		}
		components = append(components, comp)
	}
	if err = rows.Err(); err != nil {
		return components, err
	}
	return components, nil
}

func (r *Repo) GetComponent(id int) (models.Component, error) {
	var comp models.Component
	if err := r.store.db.QueryRow(`
	select c.available, c.id, c.code, c.status, c."name", c2."name" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c."time", 'DD-MM-YYYY HH24:MI') "time", 
	t."name" as type, t.id as type_id, c.weight from components c join checkpoints c2 on c2.id = c."checkpoint" join "types" t on t.id  = c."type" where c.id = $1 order by c.code`, id).Scan(
		&comp.Available, &comp.ID, &comp.Code, &comp.Status,
		&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit,
		&comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id,
		&comp.Weight); err != nil {
		return comp, err
	}
	return comp, nil
}

func (r *Repo) UpdateComponent(c *models.Component) error {
	rows, err := r.store.db.Query(`
	update components set 
	code = $1, 
	"name" = $2, 
	"checkpoint" = $3, 
	unit = $4, 
	photo = $5, 
	specs = $6, 
	"type" = $7, 
	weight = $8 
	where id = $9
	`, c.Code, c.Name, c.Checkpoint_id, c.Unit, c.Photo, c.Specs, c.Type_id, c.Weight, c.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *Repo) AddComponent(c *models.Component) error {

	logrus.Info(fmt.Sprintf(`insert into components 
	(code, "name", "checkpoint", unit, specs, photo, "type", weight) 
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, c.Code, c.Name, c.Checkpoint_id, c.Unit, c.Specs, c.Photo, c.Type_id, c.Weight))

	rows, err := r.store.db.Query(`
	insert into components 
	(code, "name", "checkpoint", unit, specs, photo, "type", weight) 
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, c.Code, c.Name, c.Checkpoint_id, c.Unit, c.Specs, c.Photo, c.Type_id, c.Weight)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *Repo) DeleteComponent(id int) error {

	_, err := r.store.db.Exec(`
	update components set status = 0 where id = $1
	`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetAllCheckpoints() (interface{}, error) {
	type Checkpoints struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
	}

	rows, err := r.store.db.Query(`select c.id, c."name", c.photo  from checkpoints c`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkpoints []Checkpoints

	for rows.Next() {
		var comp Checkpoints
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.Photo); err != nil {
			return checkpoints, err
		}
		checkpoints = append(checkpoints, comp)
	}
	if err = rows.Err(); err != nil {
		return checkpoints, err
	}

	return checkpoints, nil
}
