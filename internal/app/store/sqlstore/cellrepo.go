package sqlstore

import (
	"errors"
	"warehouse/internal/app/models"
)

// func (r *Repo) CellIncome(component_id, lot_id int, quantity float64) error {

// 	result, err := r.store.db.Exec(`
// 	insert into cell (component_id, quantity, lot_id, "time") values ($1, $2, $3, now()
// 	`, component_id, quantity, lot_id)
// 	if err != nil {
// 		return err
// 	}
// 	affected, _ := result.RowsAffected()
// 	if affected <= 0 {
// 		return errors.New("IncomeCell server error")
// 	}

// 	return nil
// }

func (r *Repo) CellAddComponent(quantity float64, component_id, lot_id, cell_id int) error {
	print("quantity: ", quantity, " component_id: ", component_id, " lot_id: ", lot_id, " cell_id: ", cell_id)

	result, err := r.store.db.Exec(`
	update cell set quantity = quantity + $1, component_id = $2, lot_id = $3, "time" = now() where id = $4
	`, quantity, component_id, lot_id, cell_id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected <= 0 {
		return errors.New("CellAddComponent server error")
	}

	return nil
}

func (r *Repo) CellRemoveComponent(id int, quantity float64) error {

	result, err := r.store.db.Exec(`
	update cell set quantity = quantity - $1 where id = $2
	`, quantity, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected <= 0 {
		return errors.New("IncomeCell server error")
	}

	return nil
}

func (r *Repo) CellInfo(id int, quantity float64) (models.CellInfoModel, error) {

	getInfo := models.CellInfoModel{}

	err := r.store.db.QueryRow(`
	update cell set quantity = quantity + $1 where id = $2
	`, quantity, id).Scan(&getInfo.Component_id, &getInfo.Quantity, &getInfo.Lot_id)
	if err != nil {
		return getInfo, err
	}

	return getInfo, nil
}

func (r *Repo) CellCheckEmpty(id int) error {

	var quantity float64

	err := r.store.db.QueryRow(`
	select c.quantity from cell c where c.id = $1
	`, id).Scan(&quantity)
	if err != nil {
		return err
	}

	if quantity <= 0 {
		r.store.db.Exec(`
		update cell set component_id = null, quantity = 0, lot_id = null, "time" = null where id= $1
		`, id)
	}

	return nil
}

func (r *Repo) CellGetEmpty(id, component_id int) (interface{}, error) {

	type EmptyCell struct {
		ID       int     `json:"id"`
		Adress   string  `json:"adress"`
		Code     string  `json:"component_code"`
		Quantity float64 `json:"quantity"`
		Lot_id   int     `json:"lot_id"`
	}

	// print("id: ", id, " component_id: ", component_id)

	rows, err := r.store.db.Query(`
	select c.id, c.adres, 
	case when (select c3.code from components c3 where c3.id = c.component_id) is null then 'empty' else (select c3.code from components c3 where c3.id = c.component_id) end, 
	c.quantity, 
	case when c.lot_id is null then 0 else c.lot_id end
	from cell c
	where c.lot_id = $1 and c.component_id = $2 or c.quantity = 0
	order by c.adres 
	`, id, component_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var emptyCells []EmptyCell

	for rows.Next() {
		var comp EmptyCell
		if err := rows.Scan(&comp.ID, &comp.Adress, &comp.Code, &comp.Quantity, &comp.Lot_id); err != nil {
			return nil, err
		}
		emptyCells = append(emptyCells, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emptyCells, nil
}

func (r *Repo) CellGetAll() (interface{}, error) {

	type EmptyCell struct {
		ID       int     `json:"id"`
		Adress   string  `json:"adress"`
		Lot_Name string  `json:"lot_name"`
		Code     string  `json:"component_code"`
		Quantity float64 `json:"quantity"`
		Time     string  `json:"time"`
	}

	rows, err := r.store.db.Query(`
	select c.id, 
	c.adres, 
	case when (select l."name" from lots l where l.id = c.lot_id) is null then 'empty' else (select l."name" from lots l where l.id = c.lot_id) end,
	case when (select c2.code from components c2 where c.component_id = c2.id) is null then 'empty' else (select c2.code from components c2 where c.component_id = c2.id) end,
	c.quantity,
	case when c."time" is null then 'empty' else to_char(c."time", 'YYYY-MM-DD HH24:MI') end
	from cell c
	order by c.adres
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var emptyCells []EmptyCell

	for rows.Next() {
		var comp EmptyCell
		if err := rows.Scan(&comp.ID, &comp.Adress, &comp.Lot_Name, &comp.Code, &comp.Quantity, &comp.Time); err != nil {
			return nil, err
		}
		emptyCells = append(emptyCells, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emptyCells, nil
}

func (r *Repo) CellGetByComponent(component_id int) (interface{}, error) {

	type EmptyCell struct {
		ID       int     `json:"id"`
		Adress   string  `json:"adress"`
		Code     string  `json:"component_code"`
		Name     string  `json:"component_name"`
		Lot_Name string  `json:"lot_name"`
		Quantity float64 `json:"quantity"`
		Time     string  `json:"time"`
	}

	rows, err := r.store.db.Query(`
	select c.id, c.adres, c2.code, c2."name", l."name" as lot, c.quantity, c."time" from cell c, components c2, lots l  
	where c.component_id = $1
	and c.component_id = c2.id 
	and l.id = c.lot_id 
	and l."blocked" = false
	order by c."time" 
	limit 1
	`, component_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var emptyCells []EmptyCell

	for rows.Next() {
		var comp EmptyCell
		if err := rows.Scan(&comp.ID, &comp.Adress, &comp.Code, &comp.Name, &comp.Lot_Name, &comp.Quantity, &comp.Time); err != nil {
			return nil, err
		}
		emptyCells = append(emptyCells, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emptyCells, nil
}

func (r *Repo) CellGetByComponentAll(component_id int) (interface{}, error) {

	type EmptyCell struct {
		ID       int     `json:"cell_id"`
		Adress   string  `json:"adress"`
		Code     string  `json:"component_code"`
		Name     string  `json:"component_name"`
		Lot_Name string  `json:"lot_name"`
		Lot_ID   string  `json:"lot_id"`
		Quantity float64 `json:"quantity"`
		Time     string  `json:"time"`
	}

	rows, err := r.store.db.Query(`
	select c.id, c.adres, c2.code, c2."name", l."name" as lot, l.id as lot_id, c.quantity, c."time" from cell c, components c2, lots l  
	where c.component_id = $1
	and c.component_id = c2.id 
	and l.id = c.lot_id 
	and l."blocked" = false
	order by c."time" 
	`, component_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var emptyCells []EmptyCell

	for rows.Next() {
		var comp EmptyCell
		if err := rows.Scan(&comp.ID, &comp.Adress, &comp.Code, &comp.Name, &comp.Lot_Name, &comp.Lot_ID, &comp.Quantity, &comp.Time); err != nil {
			return nil, err
		}
		emptyCells = append(emptyCells, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emptyCells, nil
}
