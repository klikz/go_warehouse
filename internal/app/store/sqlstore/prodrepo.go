package sqlstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"warehouse/internal/app/models"

	"github.com/sirupsen/logrus"
)

type PrinStruct struct {
	LibraryID        string `json:"libraryID"`
	AbsolutePath     string `json:"absolutePath"`
	PrintRequestID   string `json:"printRequestID"`
	Printer          string `json:"printer"`
	StartingPosition int    `json:"startingPosition"`
	Copies           int    `json:"copies"`
	SerialNumbers    int    `json:"serialNumbers"`
}

type DataEntryControlsStruct struct {
	Gscode string `json:"code"`
	Model  string `json:"modelName"`
	Serial string `json:"SerialNumber"`
}

func setPin(param, addres string) (interface{}, error) {
	response, err := http.PostForm(addres, url.Values{
		"status": {param}})
	if err != nil {
		return nil, err

	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	return string(body), nil
}
func (r *Repo) debitFromLine(modelId, lineId int) error {
	type Debit struct {
		Component_id int
		Quantity     float64
	}
	rows, err := r.store.db.Query(fmt.Sprintf("select t.component_id, t.quantity  from models.\"%d\" t, public.components c where t.component_id = c.id and c.\"checkpoint\" = %d", modelId, lineId))
	if err != nil {
		return err
	}
	defer rows.Close()
	var debits []Debit
	for rows.Next() {
		var debit Debit
		if err := rows.Scan(&debit.Component_id, &debit.Quantity); err != nil {
			return err
		}
		debits = append(debits, debit)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	for _, x := range debits {
		_, err := r.store.db.Exec(fmt.Sprintf("update checkpoints.\"%d\" set quantity = quantity - %f where component_id = %d", lineId, x.Quantity, x.Component_id))
		if err != nil {
			return err
		}
	}
	return nil
}
func CheckLaboratory(serial string) (string, error) {
	response, err := http.PostForm("http://192.168.5.250:3002/labinfo", url.Values{
		"serial": {serial}})
	if err != nil {
		return "", err

	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PrintLocal(jsonStr []byte) error {
	reprint := true
	count := 0
	for reprint {
		if count > 3 {
			return errors.New("qaytadan urinib ko'ring")
		}
		logrus.Info("Printing started")
		logrus.Info("sending data: ", string(jsonStr))
		url := "http://192.168.5.123/BarTender/api/v1/print"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			return err
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(string(body)), &jsonMap)
		logrus.Info("body: ", string(body))

		if strings.Contains(string(body), "BarTender успешно отправил задание") {
			reprint = false
		}
		count++
	}

	logrus.Info("Printing end")
	return nil
}

func PrintMetall(jsonStr []byte) error {
	reprint := true
	count := 0
	for reprint {
		if count > 3 {
			return errors.New("qaytadan urinib ko'ring")
		}
		logrus.Info("Printing started")
		logrus.Info("sending data: ", string(jsonStr))
		url := "http://192.168.5.139/BarTender/api/v1/print"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			return err
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(string(body)), &jsonMap)
		logrus.Info("body: ", string(body))

		if strings.Contains(string(body), "BarTender успешно отправил задание") {
			reprint = false
		}
		count++
	}

	logrus.Info("Printing end")
	return nil
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

func (r *Repo) AddDefects(serial, name, photo string, checkpoint, defect int) error {

	temp := serial[0:6]
	type Model_ID struct {
		ID int
	}
	var id Model_ID
	err := r.store.db.QueryRow("select m.id from models m where m.code = $1", temp).Scan(&id.ID)
	if err != nil {
		return errors.New("serial xato")
	}
	rows, err := r.store.db.Query("insert into remont (serial, person_id, checkpoint_id, model_id, defect_id, photo) values ($1, $2, $3, $4, $5, $6)", serial, name, checkpoint, id.ID, defect, photo)
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

func (r *Repo) GetCountByDate(date1, date2 string, line int) (interface{}, error) {

	type Count struct {
		Count int `json:"count"`
	}
	count := Count{}
	switch line {
	case 13:
		rows, err := r.store.db.Query(`
		select count(*) from packing 
		where "time"::date>=to_date($1, 'YYYY-MM-DD') 
		and "time"::date<=to_date($2, 'YYYY-MM-DD')`, date1, date2)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&count.Count); err != nil {
				return count, err
			}
		}
		if err = rows.Err(); err != nil {
			return count, err
		}
	default:
		rows, err := r.store.db.Query(`
		select count(*) from production 
		where "time"::date>=to_date($1, 'YYYY-MM-DD') and "time"::date<=to_date($2, 'YYYY-MM-DD') 
		and checkpoint_id = $3`, date1, date2, line)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&count.Count); err != nil {
				return count, err
			}
		}
		if err = rows.Err(); err != nil {
			return count, err
		}
	}
	return count, nil
}

func (r *Repo) GetByDateModels(date1, date2 string, line int) (interface{}, error) {

	type ByModel struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    string `json:"count"`
	}
	var byModel []ByModel

	switch line {
	case 12:
		rows, err := r.store.db.Query(`
		select p.model_id, m."name", COUNT(*) FROM galileo p, models m 
		where p."time"::date>=to_date($1, 'YYYY-MM-DD') and p."time"::date<=to_date($2, 'YYYY-MM-DD') 
		and m.id = p.model_id 
		group by m."name", p.model_id`, date1, date2)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
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
	case 13:
		rows, err := r.store.db.Query(`
		select p.model_id, m."name", COUNT(*) FROM packing p, models m 
		where p."time"::date>=to_date($1, 'YYYY-MM-DD') 
		and p."time"::date<=to_date($2, 'YYYY-MM-DD') 
		and m.id = p.model_id 
		group by m."name", p.model_id`, date1, date2)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

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
	default:
		rows, err := r.store.db.Query(`
		select p.model_id, m."name", COUNT(*) FROM production p, models m 
		where p."time"::date>=to_date($1, 'YYYY-MM-DD') 
		and p."time"::date<=to_date($2, 'YYYY-MM-DD') 
		and checkpoint_id = $3 
		and m.id = p.model_id 
		group by m."name", p.model_id`, date1, date2, line)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
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
	}
	return byModel, nil
}

func (r *Repo) GetRemont() (interface{}, error) {

	type Remont struct {
		ID         int    `json:"id"`
		Serial     string `json:"serial"`
		Vaqt       string `json:"vaqt"`
		Person     string `json:"person"`
		Checkpoint string `json:"checkpoint"`
		Model      string `json:"model"`
		Defect     string `json:"defect"`
		Photo      string `json:"photo"`
	}

	rows, err := r.store.db.Query(`
	select r.id, r.serial, to_char(r."input", 'DD-MM-YYYY') vaqt, r.person_id, c."name" as checkpoint, m."name" as model, d.defect_name as defect, r.photo 
	from remont r, checkpoints c, models m, defects d 
	where r.status = 1 and d.id = r.defect_id and c.id = r.checkpoint_id and m.id = r.model_id order by r."input"
	 `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []Remont

	for rows.Next() {
		var comp Remont
		if err := rows.Scan(&comp.ID,
			&comp.Serial,
			&comp.Vaqt,
			&comp.Person,
			&comp.Checkpoint,
			&comp.Model,
			&comp.Defect,
			&comp.Photo); err != nil {
			return nil, err
		}
		list = append(list, comp)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (r *Repo) GetRemontToday() (interface{}, error) {

	type Remont struct {
		ID         int    `json:"id"`
		Serial     string `json:"serial"`
		Vaqt       string `json:"vaqt"`
		Checkpoint string `json:"checkpoint"`
		Model      string `json:"model"`
		Defect     string `json:"defect"`
		Photo      string `json:"photo"`
	}

	currentTime := time.Now()

	rows, err := r.store.db.Query(`
	select r.id, r.serial, to_char(r."input", 'DD-MM-YYYY') vaqt, c."name" as checkpoint, m."name" as model, d.defect_name as defect, r.photo 
	from remont r, checkpoints c, models m, defects d 
	where r.status = 1 and d.id = r.defect_id and c.id = r.checkpoint_id and m.id = r.model_id and r.input::date=to_date($1, 'YYYY-MM-DD')  order by r."input"
	 `, currentTime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []Remont

	for rows.Next() {
		var comp Remont
		if err := rows.Scan(&comp.ID,
			&comp.Serial,
			&comp.Vaqt,
			&comp.Checkpoint,
			&comp.Model,
			&comp.Defect,
			&comp.Photo); err != nil {
			return nil, err
		}
		list = append(list, comp)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (r *Repo) GetRemontByDate(date1, date2 string) (interface{}, error) {

	type Remont struct {
		ID         int    `json:"id"`
		Serial     string `json:"serial"`
		Vaqt       string `json:"vaqt"`
		Checkpoint string `json:"checkpoint"`
		Model      string `json:"model"`
		Defect     string `json:"defect"`
		Photo      string `json:"photo"`
	}

	rows, err := r.store.db.Query(`
	select r.id, r.serial, to_char(r."input", 'DD-MM-YYYY') vaqt, c."name" as checkpoint, m."name" as model, d.defect_name as defect, r.photo
	 from remont r, checkpoints c, models m, defects d 
	where r.status = 1 
	and d.id = r.defect_id 
	and c.id = r.checkpoint_id 
	and m.id = r.model_id 
	and r."input"::date>=to_date($1, 'YYYY-MM-DD') 
	and r."input"::date<=to_date($2, 'YYYY-MM-DD')  
	order by r."input"
	 `, date1, date2)
	if err != nil {
		fmt.Println("GetRemont err: ", err)
		return nil, err
	}

	defer rows.Close()
	var list []Remont

	for rows.Next() {
		var comp Remont
		if err := rows.Scan(&comp.ID,
			&comp.Serial,
			&comp.Vaqt,
			&comp.Checkpoint,
			&comp.Model,
			&comp.Defect,
			&comp.Photo); err != nil {
			return nil, err
		}
		list = append(list, comp)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (r *Repo) UpdateRemont(name string, id int) error {
	rows, err := r.store.db.Query(`
	update remont set status = 0, repair_person = $1, "output" = now() where id = $2
	 `, name, id)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r *Repo) SerialInput(line int, serial string) error {
	type InputInfo struct {
		id      int
		address string
	}
	var modelInfo InputInfo
	var serialSlice = serial[0:6]
	//check address of station
	if err := r.store.db.QueryRow("select address from checkpoints where id = $1", line).Scan(&modelInfo.address); err != nil {
		return errors.New("sector address topilmadi")
	}
	//check model
	if err := r.store.db.QueryRow("select m.id from models m where m.code = $1", serialSlice).Scan(&modelInfo.id); err != nil {
		req, err := setPin("0", modelInfo.address)
		if err != nil {
			return err
		}
		logrus.Info("from raspberry: ", req)
		return errors.New("serial xato")
	}
	type product_id struct {
		id int
	}
	var prod_id product_id
	//check stations before
	type CheckStation struct {
		product_id int
	}
	switch line {
	//check sborka for ppu
	case 10:
		check := &CheckStation{}
		if err := r.store.db.QueryRow("select product_id from production p where serial = $1 and  checkpoint_id = $2", serial, 2).Scan(&check.product_id); err != nil {
			req, err := setPin("0", modelInfo.address)
			if err != nil {
				return err
			}
			logrus.Info("from raspberry: ", req)
			return errors.New("sborkada reg qilinmagan")
		}
	case 11:
		type Laboratory struct {
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
			Duration  string `json:"duration"`
			Model     string `json:"model"`
			Result    string `json:"result"`
		}

		res, err := CheckLaboratory(serial)
		if err != nil {
			return errors.New("check laboratory err")
		}
		s := string(res)
		data := Laboratory{}
		json.Unmarshal([]byte(s), &data)
		logrus.Info("laboratory: ", data.Result)
		if data.Result != "Good" {
			logrus.Info("Not Good")
			_, err := setPin("0", modelInfo.address)
			if err != nil {
				return err
			}
			return errors.New("laboratoriyada muammo")
		}
		if data.Result == "No data" {
			logrus.Info("No Data")
			_, err := setPin("0", modelInfo.address)
			if err != nil {
				return err
			}
			return errors.New("laboratoriyada muammo")
		}
	}

	// check production to serial
	if err := r.store.db.QueryRow("select product_id from production p where serial = $1 and  checkpoint_id = $2", serial, line).Scan(&prod_id.id); err == nil {
		if _, err := r.store.db.Exec("update production set updated = now() where product_id = $1", prod_id.id); err != nil {
			return err
		}
		req, err := setPin("1", modelInfo.address)
		if err != nil {
			return err
		}
		logrus.Info("from raspberry: ", req)
		return errors.New("serial kiritilgan")
	} else {
		rows, err := r.store.db.Query("insert into production (model_id, serial, checkpoint_id) values ($1, $2, $3)", modelInfo.id, serial, line)

		if err != nil {
			logrus.Info("SerialInput Setpin err: ", err)
		}
		defer rows.Close()
		err2 := r.debitFromLine(modelInfo.id, line)
		if err2 != nil {
			logrus.Info("inputSerial debit err: ", err2)
		}
		req, err := setPin("1", modelInfo.address)
		if err != nil {
			logrus.Info("SerialInput rasp err: ", err)
			return err
		}
		logrus.Info("from raspberry: ", req)
	}
	return nil
}

func (r *Repo) PackingSerialInput(serial string, retry bool) error {

	type Laboratory struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Duration  string `json:"duration"`
		Model     string `json:"model"`
		Result    string `json:"result"`
	}

	res, err := CheckLaboratory(serial)
	if err != nil {
		return errors.New("check laboratory err")
	}

	s := string(res)
	data := Laboratory{}
	json.Unmarshal([]byte(s), &data)
	if data.Result == "No data" {
		return errors.New("laboratoriyada muammo")
	}
	type ModelId struct {
		id   int
		name string
	}
	var modelId ModelId
	var serialSlice = serial[0:6]

	if err := r.store.db.QueryRow("select m.id, m.name from models m where m.code = $1", serialSlice).Scan(&modelId.id, &modelId.name); err != nil {
		return errors.New("serial xato")
	}

	rows, err := r.store.db.Query("insert into packing (serial, model_id) values ($1, $2)", serial, modelId.id)
	if err != nil {
		if retry {
			errString := err.Error()
			if strings.Contains(errString, "duplicate key") || strings.Contains(errString, "packing_un_serial") {
				codeWithError := ""
				if err := r.store.db.QueryRow(`select g."data" from gs g where product = $1`, serial).Scan(&codeWithError); err != nil {
					return err
				}
				code := strings.ReplaceAll(codeWithError, `"`, ``)
				var data1 = []byte(fmt.Sprintf(`
				{
					"libraryID": "2de725d4-1952-418e-81cc-450baa035a34",
					"absolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premer/%s_1.btw",
					"printRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
					"printer": "Gainscha GS-3405T",
					"startingPosition": 0,
					"copies": 0,
					"serialNumbers": 0,
					"dataEntryControls": {
							"GSCodeInput": "%v",
							"SeriaInput": "%s"
					}

					}`, serialSlice, code, serial))
				logrus.Info("serialSlice: ", serialSlice, "codeData.Data: ", code, "serial: ", serial)
				// var data2 = []byte(fmt.Sprintf(`
				// {
				// 	"libraryID": "2de725d4-1952-418e-81cc-450baa035a34",
				// 	"absolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premer/%s_2.btw",
				// 	"printRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
				// 	"printer": "Gainscha GS-3405T 2",
				// 	"startingPosition": 0,
				// 	"copies": 0,
				// 	"serialNumbers": 0,
				// 	"dataEntryControls": {
				// 			"SeriaInput": "%s"
				// 	}

				// 	}`, serialSlice, serial))
				PrintLocal(data1)
				// PrintLocal(data2)

				return errors.New("dublicate printed")
			}
		}
		logrus.Info("ProdRepo: err, retry = false")
		return err
	}
	defer rows.Close()

	type GSCode struct {
		ID   int
		Data string
	}
	codeData := GSCode{}

	if err := r.store.db.QueryRow("select g.id, g.data from gs g where g.model = $1 and g.status = true", modelId.id).Scan(&codeData.ID, &codeData.Data); err != nil {
		return errors.New("keys not found")
	}
	_, err = r.store.db.Exec(`update gs set product = $1, status = false where id = $2`, serial, codeData.ID)
	if err != nil {
		logrus.Info("update error: ", err)
		return err
	}
	var data1 = []byte(fmt.Sprintf(`
			{
				"libraryID": "2de725d4-1952-418e-81cc-450baa035a34",
				"absolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premer/%s_1.btw",
				"printRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
				"printer": "Gainscha GS-3405T",
				"startingPosition": 0,
				"copies": 0,
				"serialNumbers": 0,
				"dataEntryControls": {
						"GSCode": "%s",
						"SeriaInput": "%s"
				}
			}`, serialSlice, codeData.Data, serial))
	// var data2 = []byte(fmt.Sprintf(`
	// 		{
	// 			"libraryID": "2de725d4-1952-418e-81cc-450baa035a34",
	// 			"absolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premer/%s_2.btw",
	// 			"printRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
	// 			"printer": "Gainscha GS-3405T 2",
	// 			"startingPosition": 0,
	// 			"copies": 0,
	// 			"serialNumbers": 0,
	// 			"dataEntryControls": {
	// 					"SeriaInput": "%s"
	// 			}

	// 			}`, serialSlice, serial))
	PrintLocal(data1)
	// PrintLocal(data2)

	return nil
}

func (r *Repo) GetInfoBySerial(serial string) (interface{}, error) {
	type Packing struct {
		Ref_serial     string `json:"ref_serial"`
		Packing_serial string `json:"packing_serial"`
		Packing_time   string `json:"packing_time"`
	}
	type Production struct {
		Checkpoint string `json:"checkpoint"`
		Time       string `json:"time"`
	}
	type Info struct {
		PackingInfo    []Packing
		ProductionInfo []Production
	}

	var packing []Packing

	rows1, err := r.store.db.Query(`
	select p.serial as ref_serial, p.packing as packing_serial, to_char(p."time" , 'DD-MM-YYYY HH24:MI') "time" from packing p
	where p.serial = $1 `, serial)
	if err != nil {
		return nil, errors.New("no data")
	}
	defer rows1.Close()
	for rows1.Next() {
		var comp Packing
		if err := rows1.Scan(&comp.Ref_serial, &comp.Packing_serial, &comp.Packing_time); err != nil {
			return packing, errors.New("no data")
		}
		packing = append(packing, comp)
	}
	if err = rows1.Err(); err != nil {
		return nil, errors.New("no data")
	}

	// err := r.store.db.QueryRow(fmt.Sprintf(`
	// select p.serial as ref_serial, p.packing as packing_serial, to_char(p."time" , 'DD-MM-YYYY HH24:MI') "time" from packing p
	// where p.serial = '%s' `, serial)).Scan(&packing.Ref_serial, &packing.Packing_serial, &packing.Packing_time)
	// if err != nil {
	// 	fmt.Println("GetInfoBySerial get packing info err: ", err)
	// 	return nil, errors.New("no data")
	// }
	var production []Production
	// rows, err := r.store.db.Query("(select p.serial, m.\"name\" as model, p.\"time\", c.\"name\" as sector  from packing p, models m, checkpoints c  where p.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p.model_id and c.id = p.checkpoint_id  order by p.model_id) union ALL (select p2.serial, m.\"name\" as model, p2.\"time\", c.\"name\" as sector  from production p2, models m, checkpoints c where p2.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p2.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p2.model_id and c.id = p2.checkpoint_id order by p2.model_id, p2.checkpoint_id)", date1, date2)
	rows, err := r.store.db.Query(`
	(select c."name" as checkpoint , to_char(p2."time" , 'DD-MM-YYYY HH24:MI') "time"  from production p2, checkpoints c  
	where p2.serial = $1
	and p2.checkpoint_id = c.id)
	union all 
	select c."name" as checkpoint , to_char(g."time" , 'DD-MM-YYYY HH24:MI') "time"  from galileo g, checkpoints c  
	where g.serial = $2
	and g.checkpoint_id = c.id`, serial, serial)
	if err != nil {
		return nil, errors.New("no data")
	}
	defer rows.Close()
	for rows.Next() {
		var comp Production
		if err := rows.Scan(&comp.Checkpoint, &comp.Time); err != nil {
			return production, errors.New("no data")
		}
		production = append(production, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New("no data")
	}

	var productInfo Info
	productInfo.PackingInfo = packing
	productInfo.ProductionInfo = production

	return productInfo, nil
}

func (r *Repo) GalileoInput(g *models.Galileo) error {

	type InputInfo struct {
		id int
	}
	var modelInfo InputInfo

	var serialSlice = g.Barcode[0:6]

	if g.Quantity != "0" {

		if err := r.store.db.QueryRow("select m.id from models m where m.code = $1", serialSlice).Scan(&modelInfo.id); err != nil {
			return err
		}

		rows, err := r.store.db.Query("insert into galileo (serial, opcode, \"type\", progquantity, quantity, cycletotaltime, model_id, \"time\", \"result\") values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", g.Barcode, g.OpCode, g.TypeFreon, g.ProgQuantity, g.Quantity, g.CycleTotalTime, modelInfo.id, g.Time, g.Result)
		if err != nil {
			return err
		}
		defer rows.Close()
	}
	return nil
}

func (r *Repo) GalileoTodayModels() (interface{}, error) {

	type ByModel struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    string `json:"count"`
	}

	currentTime := time.Now()

	rows, err := r.store.db.Query(`
	select p.model_id, m."name", COUNT(*) FROM galileo p, models m 
	where p."time"::date>=to_date($1, 'YYYY-MM-DD') 
	and m.id = p.model_id 
	group by m."name", p.model_id 
	order by m."name" `, currentTime)
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

func (r *Repo) Metall_Serial(id int) error {

	type Data struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	info := Data{}

	count := 0

	if err := r.store.db.QueryRow(`select m2.code, m2."name" from public.models m2 where m2.id = $1`, id).Scan(&info.Code, &info.Name); err != nil {
		return err
	}

	//update count
	if err := r.store.db.QueryRow(`update metall_serial set "last" = "last" + 1 where model_id = $1 returning "last" `, id).Scan(&count); err != nil {
		return err
	}

	countString := ""

	switch {
	case count < 10:
		countString = fmt.Sprintf(`%s000000%d`, info.Code, count)

	case count < 100:
		countString = fmt.Sprintf(`%s00000%d`, info.Code, count)
	case count < 1000:
		countString = fmt.Sprintf(`%s0000%d`, info.Code, count)
	case count < 10000:
		countString = fmt.Sprintf(`%s000%d`, info.Code, count)
	case count < 100000:
		countString = fmt.Sprintf(`%s00%d`, info.Code, count)
	case count < 1000000:
		countString = fmt.Sprintf(`%s0%d`, info.Code, count)
	case count > 1000000:
		countString = fmt.Sprintf(`%s%d`, info.Code, count)
	}

	data := []byte(fmt.Sprintf(`
	{
		"libraryID": "986278f7-755f-4412-940f-a89e893947de",
		"absolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premier/serial.btw",
		"printRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
		"printer": "Gainscha GS-3405T",
		"startingPosition": 0,
		"copies": 0,
		"serialNumbers": 0,
		"dataEntryControls": {
				"Printer": "Gainscha GS-3405T",
				"ModelInput": "%s",
				"SerialInput": "%s"
		}
	}`, info.Name, countString))
	PrintMetall(data)
	// logrus.Info(string(data))
	return nil
}
