package sqlstore

import (
	"errors"
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
	logrus.Info("c.id: ", c.ID)
	if c.ID > 0 {
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
	_, err := r.store.db.Exec(`update public.components set status = 0 where id = $1`, id)
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

	rows, err := r.store.db.Query(`select c.id, c."name", c.photo  from checkpoints c where status = '1'`)
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

func (r *Repo) DeleteCheckpoint(id int) error {

	_, err := r.store.db.Exec(`
	update public.checkpoints  set status = 0 where id = $1 
	`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) AddCheckpoint(name, photo string) error {
	var id int
	err := r.store.db.QueryRow(`
	insert into public.checkpoints ("name", photo) values ($1, $2) returning id
	`, name, photo).Scan(&id)
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec(fmt.Sprintf(`
		CREATE TABLE checkpoints."%d" (
		id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
		component_id int4 NOT NULL,
		quantity int4 NOT NULL,
		CONSTRAINT "%d_pk" PRIMARY KEY (id),
		CONSTRAINT "%d_un" UNIQUE (component_id),
		CONSTRAINT "%d_FK" FOREIGN KEY (component_id) REFERENCES public.components(id))
	`, id, id, id, id))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateCheckpoint(name, photo string, id int) error {
	_, err := r.store.db.Exec(`
	UPDATE public.checkpoints SET "name" = $1, photo = $2 where id = $3
	`, name, photo, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Income(component_id int, quantity float64) error {

	_, err := r.store.db.Exec(`
	insert into income (component_id, quantity, updated ) values ($1, $2, now())
	`, component_id, quantity)
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec(`
	update components set available = available + $1 where id = $2
	`, quantity, component_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) IncomeReport(date1, date2 string) (interface{}, error) {

	type Report struct {
		Code     string  `json:"code"`
		Name     string  `json:"name"`
		Quantity float64 `json:"quantity"`
		Time     string  `json:"time"`
	}
	logrus.Info(fmt.Sprintf(`
	select c.code, c."name", i.quantity, to_char(i."create", 'DD-MM-YYYY hh-mm') time from income i, components c 
	where 
		i."create"::date>=to_date(%s, 'YYYY-MM-DD') 
		and i."create"::date<=to_date(%s, 'YYYY-MM-DD') 
		and c.id = i.component_id 
	`, date1, date2))
	rows, err := r.store.db.Query(`
	select c.code, c."name", i.quantity, to_char(i."create", 'DD-MM-YYYY hh-mm') time from income i, components c 
	where 
		i."create"::date>=to_date($1, 'YYYY-MM-DD') 
		and i."create"::date<=to_date($2, 'YYYY-MM-DD') 
		and c.id = i.component_id 
	`, date1, date2)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	report := []Report{}

	for rows.Next() {
		var comp Report
		if err := rows.Scan(&comp.Code, &comp.Name, &comp.Quantity, &comp.Time); err != nil {
			return nil, err
		}
		report = append(report, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *Repo) Types() (interface{}, error) {
	type Type struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	rows, err := r.store.db.Query(`select * from types`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := []Type{}

	for rows.Next() {
		comp := Type{}
		if err := rows.Scan(&comp.ID, &comp.Name); err != nil {
			return data, err
		}
		data = append(data, comp)
	}
	if err = rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repo) Models() (interface{}, error) {
	type Type struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Code    string `json:"code"`
		Comment string `json:"comment"`
	}

	rows, err := r.store.db.Query(`select m.id, m."name", m.code, m."comment" from models m`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := []Type{}

	for rows.Next() {
		comp := Type{}
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.Code, &comp.Comment); err != nil {
			return data, err
		}
		data = append(data, comp)
	}
	if err = rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repo) Model(id int) (interface{}, error) {
	comp := models.Model{}
	if err := r.store.db.QueryRow(`select m.id, m."name", m.code, m."comment"  from public.models m where id = $1`, id).Scan(&comp.ID, &comp.Name, &comp.Code, &comp.Comment); err != nil {
		return nil, err
	}

	return comp, nil
}

func (r *Repo) InsertUpdateModel(name, code, comment string, id int) error {
	if id > 0 {
		_, err := r.store.db.Exec(`
		update models set "name" = $1, code = $2, comment = $3 where id = $4
		`, name, code, comment, id)
		if err != nil {
			return err
		}
		return nil
	}
	new_id := 0

	if err := r.store.db.QueryRow(`
	insert into public.models ("name", code, "comment")
      values ($1, $2, $3) returning id
	`, name, code, comment).Scan(&new_id); err != nil {
		return err
	}

	_, err := r.store.db.Exec(fmt.Sprintf(`
	CREATE TABLE models."%d" (
    id int NOT NULL GENERATED ALWAYS AS IDENTITY,
    component_id int NOT NULL,
    quantity numeric NOT NULL,
    "comment" varchar NULL,
    "time" timestamp NOT null DEFAULT now()
      )
		`, new_id))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) OutcomeModelCheck(id int, quantity float64) (interface{}, error) {

	type CheckInfo struct {
		Component_id int     `json:"component_id"`
		Quantity     float64 `json:"quantity"`
		Code         string  `json:"code"`
		Name         string  `json:"name"`
		Checkpoint   int     `json:"checkpoint_id"`
		Unit         string  `json:"unit"`
		Available    float64 `json:"available"`
	}

	rows, err := r.store.db.Query(fmt.Sprintf(`
	select m.component_id, m.quantity, c.code, c."name", c."checkpoint", c.unit, c.available from models."%d" m, components c  
	where m.component_id = c.id order by c.code`, id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	check := []CheckInfo{}

	for rows.Next() {
		var comp CheckInfo
		if err := rows.Scan(&comp.Component_id, &comp.Quantity, &comp.Code, &comp.Name, &comp.Checkpoint, &comp.Unit, &comp.Available); err != nil {
			return check, err
		}
		check = append(check, comp)
	}
	if err = rows.Err(); err != nil {
		return check, err
	}

	for i := 0; i < len(check); i++ {
		check[i].Quantity *= quantity
	}

	return check, nil
}

func (r *Repo) OutcomeComponentCheck(id int, quantity float64) (interface{}, error) {

	type CheckComponent struct {
		ID         int     `json:"id"`
		Checkpoint int     `json:"checkpoint_id"`
		Available  float64 `json:"available"`
		Name       string  `json:"name"`
	}

	checkComp := CheckComponent{}

	if err := r.store.db.QueryRow(`
	select c.id, c."checkpoint", c.available, c2."name" from components c, checkpoints c2 where c.id = $1 and c2.id = c."checkpoint"
	`, id).Scan(&checkComp.ID, &checkComp.Checkpoint, &checkComp.Available, &checkComp.Name); err != nil {
		return nil, err
	}

	if checkComp.Available > quantity {
		return checkComp, nil
	}
	return checkComp, errors.New("yetarli emas")
}

func (r *Repo) OutcomeComponentSubmit(component_id, checkpoint_id int, quantity float64) error {

	_, err := r.store.db.Exec(`insert into outcome (component_id, checkpoint_id, quantity) values ($1, $2, $3)`, component_id, checkpoint_id, quantity)
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec(`update public.components set available = available - $1 where id = $2`, quantity, component_id)
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec(fmt.Sprintf(`
	with p_param as (
        select $1::int8 component_id,
              $2::numeric quantity
        ),
        u_checkpoints as (
        update checkpoints."%d" as t
           set quantity = t.quantity + p.quantity
          from p_param p
         where t.component_id = p.component_id
         returning t.*
         ),
        i_checkpoints as (
          insert into checkpoints."%d" (component_id, quantity)
          select p.component_id,
               p.quantity
              from p_param p
            left join checkpoints."%d" t
              on t.component_id = p.component_id
           where t.component_id is null
        returning "%d"
        )
        select i.*, u.*
          from p_param
          left join u_checkpoints u
            on true
          left join i_checkpoints i
            on true
	`, checkpoint_id, checkpoint_id, checkpoint_id, checkpoint_id), component_id, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) OutcomeModelSubmit(model_id int, quantity float64) error {

	type CheckInfo struct {
		Component_id int     `json:"component_id"`
		Quantity     float64 `json:"quantity"`
		Code         string  `json:"code"`
		Name         string  `json:"name"`
		Checkpoint   int     `json:"checkpoint_id"`
		Unit         string  `json:"unit"`
		Available    float64 `json:"available"`
	}

	rows, err := r.store.db.Query(fmt.Sprintf(`
	select m.component_id, m.quantity, c.code, c."name", c."checkpoint", c.unit, c.available from models."%d" m, components c  
	where m.component_id = c.id order by c.code`, model_id))
	if err != nil {
		return err
	}
	defer rows.Close()

	check := []CheckInfo{}

	for rows.Next() {
		var comp CheckInfo
		if err := rows.Scan(&comp.Component_id, &comp.Quantity, &comp.Code, &comp.Name, &comp.Checkpoint, &comp.Unit, &comp.Available); err != nil {
			return err
		}
		check = append(check, comp)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	for i := 0; i < len(check); i++ {
		check[i].Quantity *= quantity
		if check[i].Quantity > check[i].Available {
			return errors.New("komponentlar yetarli emas")
		}
	}
	for _, value := range check {
		_, err := r.store.db.Exec(`insert into outcome (component_id, checkpoint_id, quantity) values ($1, $2, $3)`, value.Component_id, value.Checkpoint, value.Quantity)
		if err != nil {
			return err
		}
		_, err = r.store.db.Exec(`update public.components set available = available - $1 where id = $2`, value.Quantity, value.Component_id)
		if err != nil {
			return err
		}
		_, err = r.store.db.Exec(fmt.Sprintf(`
		with p_param as (
			select $1::int8 component_id,
				  $2::numeric quantity
			),
			u_checkpoints as (
			update checkpoints."%d" as t
			   set quantity = t.quantity + p.quantity
			  from p_param p
			 where t.component_id = p.component_id
			 returning t.*
			 ),
			i_checkpoints as (
			  insert into checkpoints."%d" (component_id, quantity)
			  select p.component_id,
				   p.quantity
				  from p_param p
				left join checkpoints."%d" t
				  on t.component_id = p.component_id
			   where t.component_id is null
			returning "%d"
			)
			select i.*, u.*
			  from p_param
			  left join u_checkpoints u
				on true
			  left join i_checkpoints i
				on true
		`, value.Checkpoint, value.Checkpoint, value.Checkpoint, value.Checkpoint), value.Component_id, value.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) OutcomeReport(date1, date2 string) (interface{}, error) {

	type Report struct {
		Code       string  `json:"code"`
		Name       string  `json:"name"`
		Quantity   float64 `json:"quantity"`
		Checkpoint string  `json:"checkpoint"`
		Time       string  `json:"time"`
	}

	rows, err := r.store.db.Query(`
	select c.code, c."name", i.quantity, c2."name" as checkpoint, to_char(i."create", 'DD-MM-YYYY hh-mm') time from outcome i, components c, checkpoints c2 
	where 
		i."create"::date>=to_date($1, 'YYYY-MM-DD') 
		and i."create"::date<=to_date($2, 'YYYY-MM-DD') 
		and c.id = i.component_id 
		and c2.id = i.checkpoint_id 
	`, date1, date2)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	report := []Report{}

	for rows.Next() {
		var comp Report
		if err := rows.Scan(&comp.Code, &comp.Name, &comp.Quantity, &comp.Checkpoint, &comp.Time); err != nil {
			return nil, err
		}
		report = append(report, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *Repo) BomComponentInfo(id int) (interface{}, error) {

	rows, err := r.store.db.Query(fmt.Sprintf(`
	select t2.id, t2.quantity, t2."comment", t2.component_id, c.code, c."name", c.unit, c.photo, t."name" as type, c.specs, c.weight from models."%d" t2
  	join public.components c
  	on c.id = t2.component_id
  	join public."types" t
  	on t.id = c."type"
	`, id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []models.Component

	for rows.Next() {
		var comp models.Component
		if err := rows.Scan(&comp.ID, &comp.Quantity, &comp.Comment, &comp.Component_id, &comp.Code, &comp.Name, &comp.Unit, &comp.Photo, &comp.Type, &comp.Specs, &comp.Weight); err != nil {
			return components, err
		}
		components = append(components, comp)
	}
	if err = rows.Err(); err != nil {
		return components, err
	}
	return components, nil
}

func (r *Repo) BomComponentAdd(component_id, model_id int, quantity float64, comment string) error {
	_, err := r.store.db.Exec(fmt.Sprintf(`
	insert into models."%d" (component_id, quantity, "comment") values (%d, %f, '%s')
	`, model_id, component_id, quantity, comment))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) BomComponentDelete(component_id, model_id int) error {
	_, err := r.store.db.Exec(fmt.Sprintf(`delete from models."%d" where component_id = $1`, model_id), component_id)
	if err != nil {
		return err
	}
	return nil
}
