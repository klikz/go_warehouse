package sqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"warehouse/internal/app/models"

	"github.com/sirupsen/logrus"
)

func (r *Repo) AktInput(account models.Akt, filename string, type_id int) error {
	if err := r.store.db.QueryRow(`select u.id from users u where u.email = $1`, account.UserName).Scan(&account.UserID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("sql.ErrNoRows")
		}
		return err
	}

	// fmt.Println("account: ", account)
	result, err := r.store.db.Exec(`insert into akt (component_id, user_id, "comment", quantity, photo, akt_type_id, checkpoint_id) values ($1, $2, $3, $4, $5, $6, $7)`, account.Component_id, account.UserID, account.Comment, account.Quantity, filename, type_id, account.Checkpoint_id)
	if err != nil {
		logrus.Info("err: ", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected > 0 {
		return nil
	}
	return errors.New("server error")
}

func (r *Repo) AktInputWare(account models.Akt, filename string, type_id, lot_id int) (int, error) {
	if err := r.store.db.QueryRow(`select u.id from users u where u.email = $1`, account.UserName).Scan(&account.UserID); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("sql.ErrNoRows")
		}
		return 0, err
	}

	// fmt.Println("account: ", account)
	id := 0
	err := r.store.db.QueryRow(`insert into akt (component_id, user_id, "comment", quantity, photo, akt_type_id, lot_id, checkpoint_id) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`, account.Component_id, account.UserID, account.Comment, account.Quantity, filename, type_id, lot_id, account.Checkpoint_id).Scan(&id)
	if err != nil {
		logrus.Info("err: ", err)
		return 0, err
	}
	// rowsAffected, _ := result.RowsAffected()

	logrus.Info("id: ", id)

	// if rowsAffected > 0 {
	// 	return nil
	// }
	return id, nil
}

func (r *Repo) GetGPCompontents() (interface{}, error) {

	type GPComponent struct {
		ID    int    `json:"id"`
		Code  string `json:"code"`
		Name  string `json:"name"`
		Specs string `json:"specs"`
	}
	rows, err := r.store.db.Query(`
	select c.id, c.code, c."name", c.specs  from components c 
	where c."type" = 3 and c.status = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []GPComponent

	for rows.Next() {
		var comp GPComponent
		if err := rows.Scan(&comp.ID, &comp.Code, &comp.Name, &comp.Specs); err != nil {
			return components, err
		}
		components = append(components, comp)
	}
	if err = rows.Err(); err != nil {
		return components, err
	}
	return components, nil
}

func (r *Repo) GPCompontentsAdded() (interface{}, error) {

	type GPComponent struct {
		ID         int    `json:"id"`
		Checkpoint string `json:"checkpoint"`
		Component  string `json:"component"`
		Code       string `json:"code"`
		Model      string `json:"model"`
	}
	rows, err := r.store.db.Query(`
	select pg.id, c."name" as checkpoint, c2."name" as component, c2.code, m."name" as model 
	from production_gp pg, models m, checkpoints c, components c2 
	where m.id = pg.model_id 
	and c.id = pg.checkpoint_id 
	and c2.id = pg.component_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []GPComponent

	for rows.Next() {
		var comp GPComponent
		if err := rows.Scan(&comp.ID, &comp.Checkpoint, &comp.Component, &comp.Code, &comp.Model); err != nil {
			return components, err
		}
		components = append(components, comp)
	}
	if err = rows.Err(); err != nil {
		return components, err
	}
	return components, nil
}

func (r *Repo) GPCompontentsAdd(line, component, model int) error {
	result, err := r.store.db.Exec(`insert into production_gp (checkpoint_id, component_id, model_id) values ($1, $2, $3)`, line, component, model)
	if err != nil {
		logrus.Info("err: ", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected > 0 {
		return nil
	}
	return errors.New("server error")
}

func (r *Repo) GPCompontentsRemove(id int) error {
	result, err := r.store.db.Exec(`
	delete from production_gp 
	where id = $1
	`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected > 0 {
		return nil
	}
	return errors.New("server error")
}

func (r *Repo) GetAllComponentsOutcome() (interface{}, error) {
	rows, err := r.store.db.Query(`
	select c.available, c.id, c.code, c."name", c2."name" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c."time", 'DD-MM-YYYY HH24:MI') "time", 
	t."name" as type, t.id as type_id, c.weight, c.inner_code 
	from components c 
	join checkpoints c2 on c2.id = c."checkpoint" join "types" t on t.id  = c."type" 
	where c.status = 1 and not c."type" = 3
	order by c.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []models.Component

	for rows.Next() {
		var comp models.Component
		if err := rows.Scan(&comp.Available, &comp.ID, &comp.Code,
			&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit, &comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id, &comp.Weight, &comp.InnerCode); err != nil {
			return components, err
		}
		components = append(components, comp)
	}
	if err = rows.Err(); err != nil {
		return components, err
	}
	return components, nil
}

func (r *Repo) GetAllComponents() (interface{}, error) {
	rows, err := r.store.db.Query(`
	select c.available, c.id, c.code, c."name", c2."name" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c."time", 'DD-MM-YYYY HH24:MI') "time", 
	t."name" as type, t.id as type_id, c.weight, c.inner_code 
	from components c 
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
			&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit, &comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id, &comp.Weight, &comp.InnerCode); err != nil {
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
	t."name" as type, t.id as type_id, c.weight, c.inner_code 
	from components c join checkpoints c2 on c2.id = c."checkpoint" join "types" t on t.id  = c."type" where c.id = $1 order by c.code`, id).Scan(
		&comp.Available, &comp.ID, &comp.Code, &comp.Status,
		&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit,
		&comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id,
		&comp.Weight, &comp.InnerCode); err != nil {
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
	weight = $8, 
	inner_code = $9
	where id = $10
	`, c.Code, c.Name, c.Checkpoint_id, c.Unit, c.Photo, c.Specs, c.Type_id, c.Weight, c.InnerCode, c.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *Repo) AddComponent(c *models.Component) (int, error) {
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
	weight = $8,
	inner_code = $9
	where id = $10
	
	`, c.Code, c.Name, c.Checkpoint_id, c.Unit, c.Photo, c.Specs, c.Type_id, c.Weight, c.InnerCode, c.ID)
		if err != nil {
			return 0, err
		}
		defer rows.Close()
		return 0, nil
	}

	id := 0
	err := r.store.db.QueryRow(`
	insert into components 
	(code, "name", "checkpoint", unit, specs, photo, "type", weight, inner_code) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`, c.Code, c.Name, c.Checkpoint_id, c.Unit, c.Specs, c.Photo, c.Type_id, c.Weight, c.InnerCode).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) AddComponentComponentsIncome(id int) error {
	_, err := r.store.db.Exec(`insert into component_income (component_id) values ($1)`, id)
	if err != nil {
		return err
	}
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
		quantity numeric NOT NULL,
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
	select c.code, c."name", i.quantity, to_char(i."create", 'DD-MM-YYYY HH24-MI') time from income i, components c 
	where 
		i."create"::date>=to_date(%s, 'YYYY-MM-DD') 
		and i."create"::date<=to_date(%s, 'YYYY-MM-DD') 
		and c.id = i.component_id 
	`, date1, date2))
	rows, err := r.store.db.Query(`
	select c.code, c."name", i.quantity, to_char(i."create", 'DD-MM-YYYY HH24-MI') time from income i, components c 
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

	rows, err := r.store.db.Query(`select m.id, m."name", m.code, m."comment" from models m order by m.code `)

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

	_, err = r.store.db.Exec(fmt.Sprintf(`
	insert into metall_serial (model_id) values (%d)
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

	if checkComp.Available >= quantity {
		return checkComp, nil
	}
	return checkComp, errors.New("yetarli emas")
}

func (r *Repo) OutcomeInsert(component_id, checkpoint_id int, quantity float64, comment string) error {
	_, err := r.store.db.Exec(`insert into outcome (component_id, checkpoint_id, quantity, comment) values ($1, $2, $3, $4)`, component_id, checkpoint_id, quantity, comment)
	if err != nil {
		return err
	}
	return nil
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

func (r *Repo) RemoveComponentFromWare(component_id int, quantity float64) error {
	_, err := r.store.db.Exec(`update public.components set available = available - $1 where id = $2`, quantity, component_id)
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
		Comment    string  `json:"comment"`
	}

	rows, err := r.store.db.Query(`
	select c.code, c."name", i.quantity, c2."name" as checkpoint, to_char(i."create", 'DD-MM-YYYY HH12:MI') time, i.comment from outcome i, components c, checkpoints c2 
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
		if err := rows.Scan(&comp.Code, &comp.Name, &comp.Quantity, &comp.Checkpoint, &comp.Time, &comp.Comment); err != nil {
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

func (r *Repo) FileInput(f []models.FileInput) ([]models.FileInput, error) {
	// type CheckComponent struct {
	// 	ID         int     `json:"id"`
	// 	Checkpoint int     `json:"checkpoint_id"`
	// 	Available  float64 `json:"available"`
	// 	Name       string  `json:"name"`
	// }
	// fmt.Println(f)
	// logrus.Info(f)

	for i := 0; i < len(f); i++ {
		_ = r.store.db.QueryRow(`
		select c.id, c2.id as checkpoint_id,  c.available from components c, checkpoints c2  where c.code = $1 and c."checkpoint"  = c2.id 
		`, string(f[i].Code)).Scan(&f[i].ID, &f[i].Checkpoint_id, &f[i].Available)
		if f[i].Quantity > f[i].Available {
			return f, errors.New("komponent_yetarli_emas")
		}

		// logrus.Info("id: ", f[i].ID, " name: ", f[i].Name, " available: ", f[i].Available)
	}
	for i := 0; i < len(f); i++ {
		if err := r.OutcomeComponentSubmit(f[i].ID, f[i].Checkpoint_id, f[i].Quantity); err != nil {
			return nil, err
		}
		// logrus.Info("value.ID: ", f[i].ID, " value.Checkpoint_id: ", f[i].Checkpoint_id, " value.Quantity: ", f[i].Quantity)
	}
	return f, nil
}

func (r *Repo) InsertGsCode(gscode []string, model int) ([]string, error) {
	// logrus.Info("key: ", gscode, " model: ", model)
	var badCode []string
	falseCode := false
	for _, code := range gscode {
		_, err := r.store.db.Exec(`insert into gs ("data", model) values ($1, $2)`, code, model)
		if err != nil {
			badCode = append(badCode, code)
			falseCode = true
		}
	}
	if falseCode {
		return badCode, errors.New("code repeated")
	}

	return badCode, nil
}

func (r *Repo) GetKeys() (interface{}, error) {

	type Report struct {
		Name     string `json:"name"`
		Comment  string `json:"comment"`
		Quantity int    `json:"quantity"`
	}

	rows, err := r.store.db.Query(`
	select m."name", m."comment", count(g.id) as quantity 
	from gs g, models m  
	where g.status = true and m.id = g.model 
	group by m."name", m."comment" `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []Report{}

	for rows.Next() {
		var key Report
		if err := rows.Scan(&key.Name, &key.Comment, &key.Quantity); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *Repo) AktReport(date1, date2 string) (interface{}, error) {

	type Report struct {
		Email     string  `json:"email"`
		Component string  `json:"component"`
		Time      string  `json:"time"`
		Quantity  float64 `json:"quantity"`
		Comment   string  `json:"comment"`
		Photo     string  `json:"photo"`
		AktType   string  `json:"akt_type"`
		ID        int     `json:"id"`
	}

	rows, err := r.store.db.Query(`
	select  u.email, c."name" as component, to_char(a."time" , 'DD-MM-YYYY HH12:MI') time, a.quantity, a."comment", a.photo, at2."name" as akt_type, a.id
	from akt a, components c, users u, akt_type at2  
	where a."time"::date>=to_date($1, 'YYYY-MM-DD') 
	and a."time"::date<=to_date($2, 'YYYY-MM-DD') 
	and u.id = a.user_id 
	and c.id = a.component_id
	and at2.id = a.akt_type_id
	order by a."time" 
	`, date1, date2)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	report := []Report{}

	for rows.Next() {
		var comp Report
		if err := rows.Scan(&comp.Email, &comp.Component, &comp.Time, &comp.Quantity, &comp.Comment, &comp.Photo, &comp.AktType, &comp.ID); err != nil {
			return nil, err
		}
		report = append(report, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}
