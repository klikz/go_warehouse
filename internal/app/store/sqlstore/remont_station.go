package sqlstore

import "errors"

func (r *Repo) CheckCountRemont(serial string) (int, error) {

	count := 0

	err := r.store.db.QueryRow(`select count(*) from remont r where r.serial = $1`, serial).Scan(&count)

	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *Repo) BlockProduct(serial string) error {

	result, err := r.store.db.Exec(`insert into blocked_product (serial) values ($1)`, serial)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if affected == 0 || err != nil {
		return errors.New("server error")
	}
	return nil
}

func (r *Repo) CheckBlockedProduct(serial string) int {

	id := 0
	r.store.db.QueryRow(`select bp.id from blocked_product bp where serial = $1 and unblock = false`, serial).Scan(&id)

	return id
}

func (r *Repo) CheckRemontStatusProduct(serial string) int {

	count := 0
	r.store.db.QueryRow(`select count(*) from remont r where r.serial = $1 and r.status = 1`, serial).Scan(&count)

	return count
}
func (r *Repo) GetBlockedProducts() (interface{}, error) {

	type BlockedProducts struct {
		ID     int    `json:"id"`
		Serial string `json:"serial"`
	}

	rows, err := r.store.db.Query(`
	select bp.id, bp.serial from blocked_product bp where bp.unblock = false 
	 `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []BlockedProducts

	for rows.Next() {
		var comp BlockedProducts
		if err := rows.Scan(&comp.ID,
			&comp.Serial); err != nil {
			return nil, err
		}
		list = append(list, comp)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (r *Repo) UnBlockProduct(id int) error {

	result, err := r.store.db.Exec(`update blocked_product set unblock = true where id = $1`, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if affected <= 0 || err != nil {
		return errors.New("server error in unblock product")
	}

	return nil
}
