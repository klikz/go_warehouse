package sqlstore

import (
	"time"
	"warehouse/internal/app/models"
)

func (r *Repo) GetAllComponents() (interface{}, error) {
	rows, err := r.store.db.Query("select c.available, c.id, c.code, c.\"name\", c2.\"name\" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c.\"time\", 'DD-MM-YYYY HH24:MI') \"time\", t.\"name\" as type, t.id as type_id, c.weight from components c join checkpoints c2 on c2.id = c.\"checkpoint\" join \"types\" t on t.id  = c.\"type\" order by c.code")
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
	select c.available, c.id, c.code, c."name", c2."name" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c."time", 'DD-MM-YYYY HH24:MI') "time", 
	t."name" as type, t.id as type_id, c.weight from components c join checkpoints c2 on c2.id = c."checkpoint" join "types" t on t.id  = c."type" where c.id = $1 order by c.code`, id).Scan(&comp.Available, &comp.ID, &comp.Code,
		&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit, &comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id, &comp.Weight); err != nil {
		return comp, err
	}
	return comp, nil
}

func (r *Repo) GetLast(line int) ([]models.Last, error) {

	last := []models.Last{}

	rows, err := r.store.db.Query("select p.serial, p.model_id, m.\"name\" as model, p.checkpoint_id, c.\"name\" as line, p.product_id,  to_char(p.\"time\" , 'DD-MM-YYYY HH24:MI') \"time\" from production p, checkpoints c, models m where m.id = p.model_id and c.id = p.checkpoint_id and p.checkpoint_id = $1 ORDER BY p.\"time\" DESC LIMIT 2", line)
	if err != nil {
		return last, err
	}

	defer rows.Close()

	for rows.Next() {
		comp := models.Last{}
		if err := rows.Scan(&comp.Serial,
			&comp.Model_id,
			&comp.Model,
			&comp.Checkpoint_id,
			&comp.Line,
			&comp.Product_id,
			&comp.Time); err != nil {

			return last, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		return last, err
	}
	return last, nil

}

func (r *Repo) GetStatus(line int) (interface{}, error) {
	type Status struct {
		Status byte `json:"status"`
	}
	var last Status
	err := r.store.db.QueryRow("select c.status from checkpoints c where c.id = $1", line).Scan(&last.Status)
	if err != nil {
		return nil, err
	}

	return last, nil
}

func (r *Repo) GetToday(line int) (interface{}, error) {

	type Count struct {
		Count int `json:"count"`
	}
	count := Count{}
	currentTime := time.Now()

	err := r.store.db.QueryRow("select count(*) from production where checkpoint_id = $1 and \"time\"::date=to_date($2, 'YYYY-MM-DD')", line, currentTime).Scan(&count.Count)
	if err != nil {
		return count, err
	}

	return count, nil
}
