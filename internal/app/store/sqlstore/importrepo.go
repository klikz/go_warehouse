package sqlstore

import (
	"errors"
	"fmt"
	"strings"
	"warehouse/internal/app/models"

	"github.com/google/uuid"
)

func (r *Repo) InsertLot(name, comment string) error {

	rows, err := r.store.db.Exec(`insert into lots ("name", "comment") values ($1, $2)`, name, comment)

	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) DeleteLot(lot_id int) error {

	rows, err := r.store.db.Exec(`update lots set status = false where id = $1`, lot_id)

	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) UpdateLot(name, comment string, lot_id int) error {

	rows, err := r.store.db.Exec(`update lots set "name" = $1, "comment" = $2 where id = $3`, name, comment, lot_id)
	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	return nil
}

func (r *Repo) BlockLot(lot_id int) error {

	rows, err := r.store.db.Exec(`update lots set "blocked" = true where id = $1`, lot_id)
	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	return nil
}

func (r *Repo) UnBlockLot(lot_id int) error {

	rows, err := r.store.db.Exec(`update lots set "blocked" = false where id = $1`, lot_id)
	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	return nil
}

func (r *Repo) ActivateLot(lot_id int) error {

	rows, err := r.store.db.Exec(`update lots set active = true where id = $1`, lot_id)
	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	return nil
}

func (r *Repo) DeActivateLot(lot_id int) error {

	rows, err := r.store.db.Exec(`update lots set active = false where id = $1`, lot_id)
	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	return nil
}

func (r *Repo) GetAllLot() (interface{}, error) {

	type Lots struct {
		LotID   int    `json:"lot_id"`
		Name    string `json:"name"`
		Comment string `json:"comment"`
		Blocked bool   `json:"blocked"`
		Active  bool   `json:"active"`
	}

	rows, err := r.store.db.Query(`select l.id, l."name", l."comment", "blocked", active  from lots l where status = true order by l.id `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var allData []Lots

	for rows.Next() {
		var comp Lots
		if err := rows.Scan(
			&comp.LotID,
			&comp.Name,
			&comp.Comment,
			&comp.Blocked,
			&comp.Active); err != nil {
			return nil, err
		}
		allData = append(allData, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allData, nil
}
func (r *Repo) GetAllLotActive() (interface{}, error) {

	type Lots struct {
		LotID   int    `json:"lot_id"`
		Name    string `json:"name"`
		Comment string `json:"comment"`
		Blocked bool   `json:"blocked"`
		Active  bool   `json:"active"`
	}

	rows, err := r.store.db.Query(`select l.id, l."name", l."comment", "blocked", active  from lots l where status = true and active = true order by l.id `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var allData []Lots

	for rows.Next() {
		var comp Lots
		if err := rows.Scan(
			&comp.LotID,
			&comp.Name,
			&comp.Comment,
			&comp.Blocked,
			&comp.Active); err != nil {
			return nil, err
		}
		allData = append(allData, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allData, nil
}

//Repo for Batches

func (r *Repo) InsertBatch(lot_id int, name, comment string) error {

	rows, err := r.store.db.Exec(`insert into batch (lot_id, "name", "comment") values ($1, $2, $3)`, lot_id, name, comment)

	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) DeleteBatch(id int) error {

	rows, err := r.store.db.Exec(`update batch set status = false where id = $1`, id)

	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) UpdateBatch(name, comment string, batch_id int) error {

	rows, err := r.store.db.Exec(`update batch set "name" = $1, "comment" = $2, "time" = now() where id = $3`, name, comment, batch_id)
	if err != nil {
		return err
	}

	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetBatchByLot(lot_id int) (interface{}, error) {

	type Batchs struct {
		Batch_id   int    `json:"batch_id"`
		Batch_name string `json:"batch_name"`
		Comment    string `json:"comment"`
		Lot_name   string `json:"lot_name"`
	}

	rows, err := r.store.db.Query(`
	select b.id, b."name" as batch, b."comment", l."name" as lot
	from batch b, lots l  
	where lot_id = $1
	and l.id = b.lot_id  
	order by b.id`, lot_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var allData []Batchs

	for rows.Next() {
		var comp Batchs
		if err := rows.Scan(
			&comp.Batch_id,
			&comp.Batch_name,
			&comp.Comment,
			&comp.Lot_name); err != nil {
			return nil, err
		}
		allData = append(allData, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allData, nil
}

//repo for container

func (r *Repo) InsertContainer(name, comment string, lot_id, batch_id int) error {

	rows, err := r.store.db.Exec(`
	insert into container ("name", lot_id, batch_id, "comment") 
	values ($1, $2, $3, $4)`, name, lot_id, batch_id, comment)

	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) DeleteContainer(id int) error {

	rows, err := r.store.db.Exec(`update container set status = false where id = $1`, id)

	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) UpdateContainer(name, comment string, container_id int) error {

	rows, err := r.store.db.Exec(`update container set "name" = $1, "comment" = $2 where id = $3`, name, comment, container_id)

	changed, _ := rows.RowsAffected()
	if err != nil {
		return err
	}

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) ContainerComponentsDelete(component_id int) error {
	print(component_id)

	rows, err := r.store.db.Exec(`delete from import_income where id = $1`, component_id)
	if err != nil {
		return err
	}
	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) ContainerComponentsUpdate(r_quantity float64, comment string, component_id int) error {

	print(" r_quantity: ", r_quantity, " ", "comment: ", comment, "\n")
	rows, err := r.store.db.Exec(`update import_income set r_quantity = $1, "comment" = $2, income_time = now() where id = $3`, r_quantity, comment, component_id)
	if err != nil {
		return err
	}

	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetContainerByBatch(id int) (interface{}, error) {

	type Info struct {
		ContainerID   int    `json:"container_id"`
		ContainerName string `json:"container_name"`
		LotName       string `json:"lot_name"`
		BatchName     string `json:"batch_name"`
		Comment       string `json:"comment"`
	}

	rows, err := r.store.db.Query(`
	select c.id, c."name" as container, l."name" as lot, b."name" as batch, c."comment" 
	from container c, lots l, batch b
	where batch_id =$1
	and c.lot_id = l.id 
	and c.batch_id = b.id 
	and c.status = true `, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var allData []Info

	for rows.Next() {
		var comp Info
		if err := rows.Scan(
			&comp.ContainerID,
			&comp.ContainerName,
			&comp.LotName,
			&comp.BatchName,
			&comp.Comment); err != nil {
			return nil, err
		}
		allData = append(allData, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if allData == nil {
		return nil, errors.New("containers not found")
	}

	return allData, nil
}

func (r *Repo) GetContainerComponents(container_id int) (interface{}, error) {

	type Info struct {
		ID          int     `json:"id"`
		Code        string  `json:"component_code"`
		Name        string  `json:"component_name"`
		Izm         string  `json:"izm"`
		Quantity    float64 `json:"quantity"`
		R_Quantity  float64 `json:"r_quantity"`
		Unit        string  `json:"unit"`
		Income_Time string  `json:"income_time"`
		Comment     string  `json:"comment"`
	}

	rows, err := r.store.db.Query(`
	select i.id, c.code,  c."name" as component_name, c.unit, i.quantity, i.r_quantity, 
	case when i.unit is null then ' ' else i.unit end, 
	case when i.income_time is null then '2000-01-01 00:00:00.000' else to_char(i.income_time, 'YYYY-MM-DD HH24:MI') end, i.comment
	from import_income i, components c 
	where i.container_id = $1
	and c.id = i.component_id
	order by i.id `, container_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var allData []Info

	for rows.Next() {
		var comp Info
		if err := rows.Scan(
			&comp.ID,
			&comp.Code,
			&comp.Name,
			&comp.Izm,
			&comp.Quantity,
			&comp.R_Quantity,
			&comp.Unit,
			&comp.Income_Time,
			&comp.Comment); err != nil {
			return nil, err
		}
		allData = append(allData, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if allData == nil {
		return nil, errors.New("components not found")
	}

	return allData, nil
}

func (r *Repo) ImportIncomeRegister(lot_id, batch_id, container_id, component_id int, quantity float64, comment, unit string) error {

	id := uuid.New()
	if unit == "" {
		unit = id.String()
	}
	rows, err := r.store.db.Exec(`
	insert into import_income (lot_id, batch_id, container_id, component_id, quantity, "time", "comment", unit) values ($1, $2, $3, $4, $5, now(), $6, $7)`,
		lot_id, batch_id, container_id, component_id, quantity, comment, unit)

	if err != nil {
		return err
	}

	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) AddComponentsImport(lot_id, batch_id, container_id, component_id int, quantity float64, comment, unit string) error {

	id := uuid.New()
	if unit == "" {
		unit = id.String()
	}
	rows, err := r.store.db.Exec(`
	insert into import_income (lot_id, batch_id, container_id, component_id, quantity, "time", "comment", unit) values ($1, $2, $3, $4, $5, now(), $6, $7)`,
		lot_id, batch_id, container_id, component_id, quantity, comment, unit)

	if err != nil {
		return err
	}

	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

func (r *Repo) ImportIncomeRegisterFromFile(lot_id, batch_id, container_id int, file []models.FileInput2) error {

	hasError := []string{}

	for i := 0; i < len(file); i++ {
		var component_id int
		if file[i].Code == "" {
			return nil
		}
		err := r.store.db.QueryRow(`select c.id from components c where c.code = $1`, file[i].Code).Scan(&component_id)
		if err != nil {
			// logrus.Error("error check: ", err)
			hasError = append(hasError, file[i].Code)
			// return err
		} else {
			err := r.ImportIncomeRegister(lot_id, batch_id, container_id, component_id, file[i].Quantity, file[i].Comment, file[i].Unit)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "повторяющееся значение") {

				} else {
					return err
				}

			}
		}
	}

	if len(hasError) > 0 {
		errorString := ""
		for i := 0; i < len(hasError); i++ {
			errorString += fmt.Sprintf("%s \n", hasError[i])
		}
		return errors.New(string(errorString))
	}

	return nil
}

func (r *Repo) ImportIncomeAdd(quantity float64, income_id int) error {

	rows, err := r.store.db.Exec(`
	update import_income set r_quantity = $1, income_time = now() where id = $2`,
		quantity, income_id)

	if err != nil {
		return err
	}

	changed, _ := rows.RowsAffected()

	if changed == 0 {
		return errors.New("0 rows changed")
	}
	return nil
}

// func (r *Repo) ComponentsIncomeAdd(quantity float64, component_id int) error {

// 	rows, err := r.store.db.Exec(`
// 	update import_income set r_quantity = $1, income_time = now() where id = $2`,
// 		quantity, income_id)

// 	if err != nil {
// 		return err
// 	}

// 	changed, _ := rows.RowsAffected()

// 	if changed == 0 {
// 		return errors.New("0 rows changed")
// 	}
// 	return nil
// }
