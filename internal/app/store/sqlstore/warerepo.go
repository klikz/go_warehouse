package sqlstore

import (
	"errors"
	"fmt"
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

func (r *Repo) GetTodayModels(line int) (interface{}, error) {

	type ByModel struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    string `json:"count"`
	}

	currentTime := time.Now()

	rows, err := r.store.db.Query(`
	select p.model_id, m."name", COUNT(*) FROM production p, models m 
	where p.checkpoint_id = $1 and p."time"::date>=to_date($2, 'YYYY-MM-DD') 
	and m.id = p.model_id group by m."name", p.model_id order by m."name"`, line, currentTime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var byModel []ByModel

	for rows.Next() {
		var comp ByModel
		if err := rows.Scan(&comp.Model_id,
			&comp.Name,
			&comp.Count); err != nil {
			return byModel, err
		}
		byModel = append(byModel, comp)
	}
	if err = rows.Err(); err != nil {
		return byModel, err
	}
	return byModel, nil
}

func (r *Repo) GetSectorBalance(line int) (interface{}, error) {

	type Balance struct {
		Component_id int     `json:"component_id"`
		Code         string  `json:"code"`
		Quantity     float32 `json:"quantity"`
		Name         string  `json:"name"`
	}

	rows, err := r.store.db.Query(fmt.Sprintf(`select t.component_id, c.code,  t.quantity, c."name" from checkpoints."%d" t, components c where t.component_id = c.id ORDER BY t.quantity`, line))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var balance []Balance
	for rows.Next() {
		var comp Balance
		if err := rows.Scan(&comp.Component_id,
			&comp.Code,
			&comp.Quantity,
			&comp.Name); err != nil {
			return balance, err
		}
		balance = append(balance, comp)
	}
	if err = rows.Err(); err != nil {
		return balance, err
	}
	return balance, nil
}

func (r *Repo) GetPackingLast() (interface{}, error) {

	type PackingLast struct {
		ID      int    `json:"id"`
		Serial  string `json:"serial"`
		Packing string `json:"packing"`
		Time    string `json:"time"`
	}

	rows, err := r.store.db.Query(`select p.id, p.serial, p.packing, to_char(p."time" , 'DD-MM-YYYY HH24:MI') "time" from packing p ORDER BY p."time" DESC LIMIT 3`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var last []PackingLast
	for rows.Next() {
		var comp PackingLast
		if err := rows.Scan(&comp.ID,
			&comp.Serial,
			&comp.Packing,
			&comp.Time); err != nil {
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		return last, err
	}
	return last, nil
}

func (r *Repo) GetPackingToday() (interface{}, error) {

	type PackingToday struct {
		Count int `json:"count"`
	}
	currentTime := time.Now()
	var last PackingToday
	err := r.store.db.QueryRow(`select count(*) from packing where "time"::date=to_date($1, 'YYYY-MM-DD')`, currentTime).Scan(&last.Count)
	if err != nil {
		return nil, err
	}
	return last, nil
}

func (r *Repo) GetPackingTodaySerial() (interface{}, error) {

	type PackingTodaySerial struct {
		Serial  string `json:"serial"`
		Packing string `json:"packing"`
		Time    string `json:"time"`
	}
	currentTime := time.Now()
	rows, err := r.store.db.Query(`
	select serial, packing, to_char("time" , 'DD-MM-YYYY HH24:MI') "time" from packing 
	where "time"::date=to_date($1, 'YYYY-MM-DD') order by serial`, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var last []PackingTodaySerial
	for rows.Next() {
		var comp PackingTodaySerial
		if err := rows.Scan(&comp.Serial,
			&comp.Packing,
			&comp.Time); err != nil {
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		return last, err
	}

	return last, nil
}

func (r *Repo) GetPackingTodayModels() (interface{}, error) {

	type PackingTodayModels struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    int    `json:"count"`
	}
	currentTime := time.Now()
	rows, err := r.store.db.Query(`
	select p.model_id, m."name", COUNT(*) FROM packing p, models m 
	where p."time"::date>=to_date($1, 'YYYY-MM-DD') and m.id = p.model_id group by m."name", p.model_id order by m."name" `, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var last []PackingTodayModels
	for rows.Next() {
		var comp PackingTodayModels
		if err := rows.Scan(&comp.Model_id, &comp.Name, &comp.Count); err != nil {
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return last, nil
}

func (r *Repo) GetLines() (interface{}, error) {

	type Lines struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	rows, err := r.store.db.Query(`select c.id, c."name"  from checkpoints c where c.status = '1' `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var last []Lines

	for rows.Next() {
		var comp Lines
		if err := rows.Scan(&comp.ID,
			&comp.Name); err != nil {
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		return last, err
	}

	return last, nil
}

func (r *Repo) GetDefectsTypes() (interface{}, error) {

	type defectsTypes struct {
		ID          int    `json:"id"`
		Defect_name string `json:"defect_name"`
		Line_id     string `json:"line_id"`
		Name        string `json:"name"`
	}

	rows, err := r.store.db.Query(`select r.id, r.defect_name, r.line_id, c."name"  from defects r, checkpoints c where c.id = r.line_id and r.status = '1' order by line_id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var last []defectsTypes

	for rows.Next() {
		var comp defectsTypes
		if err := rows.Scan(&comp.ID,
			&comp.Defect_name,
			&comp.Line_id,
			&comp.Name); err != nil {
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		return last, err
	}

	return last, nil
}

func (r *Repo) DeleteDefectsTypes(id int) error {
	rows, err := r.store.db.Query(`update defects set status = '0' where id = $1`, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *Repo) AddDefectsTypes(id int, name string) error {
	rows, err := r.store.db.Query("insert into defects (defect_name, line_id) values ($1, $2)", name, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *Repo) AddDefects(serial, name string, checkpoint, defect int) error {

	temp := serial[0:6]
	type Model_ID struct {
		ID int
	}
	var id Model_ID
	err := r.store.db.QueryRow("select m.id from models m where m.code = $1", temp).Scan(&id.ID)
	if err != nil {
		return errors.New("serial xato")
	}
	rows, err := r.store.db.Query("insert into remont (serial, person_id, checkpoint_id, model_id, defect_id) values ($1, $2, $3, $4, $5)", serial, name, checkpoint, id.ID, defect)
	fmt.Println(fmt.Sprint("insert into remont (serial, person_id, checkpoint_id, model_id, defect_id) values ($1, $2, $3, $4, $5)", serial, name, checkpoint, id.ID, defect))
	if err != nil {
		return err

	}
	defer rows.Close()
	return nil
}

func (r *Repo) GetByDateSerial(date1, date2 string) (interface{}, error) {
	type Serial struct {
		Serial string `json:"serial"`
		Model  string `json:"model"`
		Time   string `json:"time"`
		Sector string `json:"sector"`
	}
	var serial []Serial
	// rows, err := r.store.db.Query("(select p.serial, m.\"name\" as model, p.\"time\", c.\"name\" as sector  from packing p, models m, checkpoints c  where p.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p.model_id and c.id = p.checkpoint_id  order by p.model_id) union ALL (select p2.serial, m.\"name\" as model, p2.\"time\", c.\"name\" as sector  from production p2, models m, checkpoints c where p2.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p2.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p2.model_id and c.id = p2.checkpoint_id order by p2.model_id, p2.checkpoint_id)", date1, date2)
	rows, err := r.store.db.Query(`
	(select p.serial, m."name" as model, to_char(p."time" , 'DD-MM-YYYY HH24:MI') "time" , c."name" as sector  from packing p, models m, checkpoints c
	where p."time"::date>=to_date($1, 'YYYY-MM-DD') and p."time"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p.model_id and c.id = p.checkpoint_id  order by p.model_id)
	union all
	(select p.serial, m."name" as model, to_char(p."time" , 'DD-MM-YYYY HH24:MI') "time" , c."name" as sector  from galileo p, models m, checkpoints c
	where p."time"::date>=to_date($3, 'YYYY-MM-DD') and p."time"::date<=to_date($4, 'YYYY-MM-DD') and m.id = p.model_id and c.id = p.checkpoint_id  order by p.model_id)
	union all
	(select p2.serial, m."name" as model, to_char(p2."time" , 'DD-MM-YYYY HH24:MI') "time", c."name" as sector  from production p2, models m, checkpoints c
	where p2."time"::date>=to_date($5, 'YYYY-MM-DD') and p2."time"::date<=to_date($6, 'YYYY-MM-DD') and m.id = p2.model_id and c.id = p2.checkpoint_id order by p2.model_id, p2.checkpoint_id)`,
		date1, date2, date1, date2, date1, date2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp Serial
		if err := rows.Scan(&comp.Serial, &comp.Model, &comp.Time, &comp.Sector); err != nil {
			return serial, err
		}
		serial = append(serial, comp)
	}
	if err = rows.Err(); err != nil {
		return serial, err
	}
	return serial, nil
}
