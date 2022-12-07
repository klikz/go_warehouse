package sqlstore

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

	// return nil, nil
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

func (r *Repo) GPComponentAddToLine(line_id, component_id int) error {
	_, err := r.store.repo.store.db.Exec(fmt.Sprintf(`
		with p_param as (
		select %v::int8 component_id), i_products as (
		INSERT INTO checkpoints."%v" (component_id, quantity)
		select t.component_id, 1
		from p_param t
		where not exists (select 1 from checkpoints."%v" p where p.component_id  = t.component_id)
		returning checkpoints."%v".*),u_products as (
		update checkpoints."%v" t
		set quantity = quantity + 1
		from p_param p
		where p.component_id = t.component_id
		returning t.*)
		select case when s1.component_id is null then null else 'add' end add_p,
		case when s2.component_id is null then null else 'updated' end edit_p
		from p_param p
		left join i_products s1
		on true
		left join u_products s2
		on true   
	`, component_id, line_id, line_id, line_id, line_id))
	if err != nil {
		fmt.Println("err gp add to line: ", err)
		return err
	}
	return nil
}

func (r *Repo) DecreaseFromLine(line_id, component_id int) error {
	_, err := r.store.db.Exec(fmt.Sprintf(`update checkpoints."%v" set quantity = quantity - 1 where component_id = %v`, line_id, component_id))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) ProductionIncomeSerialsInput(lineIncome int, serial string) error {
	_, err := r.store.db.Exec(`insert into prod_income_serials (serial, checkpoint_id) values ($1, $2)`, serial, lineIncome)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) IncomeInProduction(lineIncome, lineOutcome int, serial string) error {

	serialSlice := serial[0:6]
	componentID := 0

	//check for vacuum hips
	checkVacuum := serialSlice[0:2]
	if checkVacuum == "HI" {
		if lineOutcome != 1 {
			return errors.New("liniya xato")
		}
		//check model
		if err := r.store.db.QueryRow("select v.component_id from \"vacuum\" v where v.serial = $1", serialSlice).Scan(&componentID); err != nil {
			return errors.New("serial xato 1")
		}
	} else {
		//check model
		modelID := 0
		fmt.Println("serialSlice: ", serialSlice)
		if err := r.store.db.QueryRow("select m.id from models m where m.code = $1", serialSlice).Scan(&modelID); err != nil {
			return errors.New("serial xato 2")
		}

		fmt.Println("lineOutcome ", lineOutcome, " modelID", modelID)
		// select component
		if err := r.store.db.QueryRow("select pg.component_id from production_gp pg where pg.checkpoint_id = $1 and pg.model_id = $2", lineOutcome, modelID).Scan(&componentID); err != nil {
			return errors.New("serial xato 3")
		}
	}

	if err := r.GPComponentAddToLine(lineIncome, componentID); err != nil {
		return err
	}

	if err := r.DecreaseFromLine(lineOutcome, componentID); err != nil {
		return err
	}
	return nil
}

func CheckLaboratory(serial string) (string, error) {
	// logrus.Info("Check laboratory")
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
	// return "", nil
}

func PrintLocal(jsonStr []byte, channel chan string, wg *sync.WaitGroup) {
	// logrus.Info("Print local")
	defer wg.Done()
	reprint := true
	count := 0

	for reprint {
		if count > 3 {
			channel <- "qaytadan urinib ko'ring"
			close(channel)
			return
		}
		// logrus.Info("Printing started")
		url := "http://192.168.5.118/BarTender/api/v1/print" //for test
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			channel <- err.Error()
			close(channel)
			return
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			channel <- err.Error()
			close(channel)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(string(body)), &jsonMap)
		// logrus.Info("body: ", string(body))

		if strings.Contains(string(body), "BarTender успешно отправил задание") {
			reprint = false
			channel <- "ok"
			close(channel)
			// logrus.Info("Printing end")
			return
		}
		count++
	}

	// channel <- "ok"
}

func PrintMetall(jsonStr []byte, channel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	reprint := true
	count := 0

	for reprint {
		if count > 3 {
			channel <- "qaytadan urinib ko'ring"
			close(channel)
			return
		}
		// logrus.Info("Printing started")
		url := "http://192.168.5.126/BarTender/api/v1/print"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			channel <- err.Error()
			close(channel)
			return
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			channel <- err.Error()
			close(channel)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(string(body)), &jsonMap)
		// logrus.Info("body: ", string(body))

		if strings.Contains(string(body), "BarTender успешно отправил задание") {
			reprint = false
			channel <- "ok"
			close(channel)
			// logrus.Info("Printing end")
			return
		}
		count++
	}

	// channel <- "ok"
	// logrus.Info("Printing: ", string(jsonStr))
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

	rows, err := r.store.db.Query(fmt.Sprintf(`select t.component_id, c.code,  t.quantity, c."name" from checkpoints."%d" t, components c where t.component_id = c.id ORDER BY c."name"`, line))
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

func (r *Repo) GetSectorBalanceGP(line int) (interface{}, error) {

	type Balance struct {
		Component_id int     `json:"component_id"`
		Code         string  `json:"code"`
		Quantity     float32 `json:"quantity"`
		Name         string  `json:"name"`
	}

	rows, err := r.store.db.Query(fmt.Sprintf(`
	select t.component_id, c.code,  t.quantity, c."name" 
	from checkpoints."%v" t, components c 
	where t.component_id = c.id 
	and c."type" = 3
	ORDER BY c."name"
	`, line))
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

func (r *Repo) Last3Defects() (interface{}, error) {
	type Last struct {
		Serial     string `json:"serial"`
		Time       string `json:"time"`
		Checkpoint string `json:"checkpoint"`
		DefectName string `json:"defect_name"`
	}

	rows, err := r.store.db.Query(`select r.serial, to_char(r."input", 'DD-MM-YYYY HH24:MI') as time, c."name" as checkpoint, d.defect_name 
	from remont r, checkpoints c, defects d  
	where r.status = '1' and c.id = r.checkpoint_id and d.id = r.defect_id 
	order by "input"  DESC limit 3`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	last := []Last{}
	for rows.Next() {
		comp := Last{}

		if err := rows.Scan(&comp.Serial, &comp.Time, &comp.Checkpoint, &comp.DefectName); err != nil {
			return nil, err
		}
		last = append(last, comp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return last, nil
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
	select r.id, r.serial, to_char(r."input", 'DD-MM-YYYY HH24-MI') vaqt, r.person_id, c."name" as checkpoint, m."name" as model, d.defect_name as defect, r.photo 
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
	select r.id, r.serial, to_char(r."input", 'DD-MM-YYYY HH24:MI') vaqt, c."name" as checkpoint, m."name" as model, d.defect_name as defect, r.photo 
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

func (r *Repo) GetCountRepaired() (int, error) {
	count := 0

	if err := r.store.db.QueryRow(`select count(*) from remont r where r.status = '0' and r."output" >= current_date`).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
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
		Status     string `json:"status"`
	}

	rows, err := r.store.db.Query(`
	select r.id, r.serial, to_char(r."input", 'DD-MM-YYYY HH24:MI') vaqt, c."name" as checkpoint, m."name" as model, d.defect_name as defect, r.photo, r.status
	 from remont r, checkpoints c, models m, defects d 
	where d.id = r.defect_id 
	and c.id = r.checkpoint_id 
	and m.id = r.model_id 
	and r."input"::date>=to_date($1, 'YYYY-MM-DD') 
	and r."input"::date<=to_date($2, 'YYYY-MM-DD')  
	order by  r.status, r."input"
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
			&comp.Photo,
			&comp.Status); err != nil {
			return nil, err
		}
		list = append(list, comp)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}

	for i := 0; i < len(list); i++ {
		if list[i].Status == "1" {
			list[i].Status = "defected"
		} else {
			list[i].Status = "repaired"
		}
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
	var GPID []int

	rows, err := r.store.db.Query(`
	select pg.component_id  from production_gp pg where pg.checkpoint_id = $1 and pg.model_id = $2
	`, line, modelInfo.id)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var debit int
		rows.Scan(&debit)
		GPID = append(GPID, debit)
	}

	logrus.Info("GPID: ", GPID, "line: ", line)

	switch line {
	//check sborka for ppu
	case 10:
		check := &CheckStation{}
		if err := r.store.db.QueryRow("select product_id from production p where serial = $1 and checkpoint_id = $2", serial, 2).Scan(&check.product_id); err != nil {
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
	case 9:
		// check production to serial
		if err := r.store.db.QueryRow("select product_id from production p where serial = $1 and  checkpoint_id = $2", serial, line).Scan(&prod_id.id); err == nil {
			if _, err := r.store.db.Exec("update production set updated = now() where product_id = $1", prod_id.id); err != nil {
				return err
			}
			modelName := ""
			if err := r.store.db.QueryRow(`select m."name"  from public.models m where code = $1`, serialSlice).Scan(&modelName); err != nil {
				return err
			}

			var wg sync.WaitGroup
			wg.Add(1)
			channel := make(chan string, 1)

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
		}`, modelName, serial))
			PrintMetall(data, channel, &wg)

			wg.Wait()
			errorText1 := <-channel

			if errorText1 != "ok" {
				logrus.Error("error in printing: " + errorText1)
				return errors.New("qaytadan urinib ko'ring")
			}
			return errors.New("serial kiritilgan")
		} else {
			rows, err := r.store.db.Query("insert into production (model_id, serial, checkpoint_id) values ($1, $2, $3)", modelInfo.id, serial, line)

			if err != nil {
				logrus.Error("Case9 Serial Input: ", err)
				return err
			}
			defer rows.Close()
			err = r.debitFromLine(modelInfo.id, line)
			if err != nil {
				logrus.Error("Case9 debitFromLine: ", err)
				return err
			}

			for i := 0; i < len(GPID); i++ {
				err := r.GPComponentAddToLine(line, GPID[i])
				if err != nil {
					logrus.Error("GPComponentAddToLine: ", err)
				}
			}
			return nil
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

		for i := 0; i < len(GPID); i++ {
			err := r.GPComponentAddToLine(line, GPID[i])
			if err != nil {
				logrus.Error("GPComponentAddToLine: ", err)
			}
		}
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
	// logrus.Info("PackingSerialInput started")
	res, err := CheckLaboratory(serial)
	if err != nil {
		return errors.New("check laboratory err")
	}

	s := string(res)
	data := Laboratory{}
	json.Unmarshal([]byte(s), &data)
	logrus.Info("data from lab: ", data)
	if data.Result == "No data" {
		return errors.New("laboratoriyada muammo")
	}
	type ModelId struct {
		id   int
		name string
	}
	var modelId ModelId
	var serialSlice = serial[0:6]
	// logrus.Info("Get id and name")
	if err := r.store.db.QueryRow("select m.id, m.name from models m where m.code = $1", serialSlice).Scan(&modelId.id, &modelId.name); err != nil {
		return errors.New("serial xato")
	}

	var wg sync.WaitGroup

	if retry {
		var check interface{}
		if err := r.store.db.QueryRow(`select g."data" from gs g where model = $1`, modelId.id).Scan(&check); err != nil {
			return err
		}

		code := ""
		if err := r.store.db.QueryRow(`select g."data" from gs g where product = $1`, serial).Scan(&code); err != nil {
			return err
		}
		code = strings.ReplaceAll(code, `"`, ``)
		code = strings.ReplaceAll(code, "", "")

		channel1 := make(chan string, 1)
		channel2 := make(chan string, 1)

		wg.Add(2)

		var data1 = []byte(fmt.Sprintf(`
			{
				"LibraryID": "2de725d4-1952-418e-81cc-450baa035a34",
				"AbsolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premier/%s_1.btw",
				"PrintRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
				"Printer": "Gainscha GS-3405T",
				"DataEntryControls": {
						"GSCodeInput": "%v",
						"SeriaInput": "%s"
				}

				}`, serialSlice, code, serial))
		var data2 = []byte(fmt.Sprintf(`
		{
			"LibraryID": "2de725d4-1952-418e-81cc-450baa035a34",
			"AbsolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premier/%s_2.btw",
			"PrintRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
			"Printer": "Xprinter XP-H500B",
			"DataEntryControls": {
				"SeriaInput": "%s"
			}
		}`, serialSlice, serial))

		go PrintLocal(data1, channel1, &wg)
		go PrintLocal(data2, channel2, &wg)

		wg.Wait()
		errorText1 := <-channel1
		errorText2 := <-channel2

		// if errorText1 == "ok" {
		// 	return nil
		// } else {
		// 	logrus.Error("error in printing: " + errorText1)
		// 	return err
		// }

		if errorText1 == "ok" && errorText2 == "ok" {
			return errors.New("dublicate printed")
		} else {
			logrus.Error("error in printing: " + errorText1 + " " + errorText2)
			return errors.New("qaytadan urinib ko'ring")
		}

	}
	logrus.Info("check gs code")
	var check interface{}
	if err := r.store.db.QueryRow(`select g."data" from gs g where model = $1`, modelId.id).Scan(&check); err != nil {
		logrus.Error("error check: ", err)
		if err == sql.ErrNoRows {
			return errors.New("GS kod tugagan yuklash kerak")
		}
		return err
	}
	logrus.Info("insert serial to db packing")
	rows, err := r.store.db.Query("insert into packing (serial, model_id) values ($1, $2)", serial, modelId.id)
	if err != nil {
		return err
	}
	defer rows.Close()

	type GSCode struct {
		ID   int
		Data string
	}
	codeData := GSCode{}
	logrus.Info("Get GS code")
	if err := r.store.db.QueryRow("select g.id, g.data from gs g where g.model = $1 and g.status = true", modelId.id).Scan(&codeData.ID, &codeData.Data); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("keys not found")
		}
		return err
	}
	logrus.Info("Update gs code")
	_, err = r.store.db.Exec(`update gs set product = $1, status = false where id = $2`, serial, codeData.ID)
	if err != nil {
		return err
	}

	channel1 := make(chan string, 1)
	channel2 := make(chan string, 1)
	codeData.Data = strings.ReplaceAll(codeData.Data, `"`, ``)
	codeData.Data = strings.ReplaceAll(codeData.Data, "", "")
	logrus.Info("Print data")
	var data1 = []byte(fmt.Sprintf(`
			{
				"LibraryID": "2de725d4-1952-418e-81cc-450baa035a34",
				"AbsolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premier/%s_1.btw",
				"PrintRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
				"Printer": "Gainscha GS-3405T",
				"DataEntryControls": {
						"GSCodeInput": "%s",
						"SeriaInput": "%s"
				}
			}`, serialSlice, codeData.Data, serial))
	var data2 = []byte(fmt.Sprintf(`
		{
			"LibraryID": "2de725d4-1952-418e-81cc-450baa035a34",
			"AbsolutePath": "C:/inetpub/wwwroot/BarTender/wwwroot/Templates/premier/%s_2.btw",
			"PrintRequestID": "fe80480e-1f94-4A2f-8947-e492800623aa",
			"Printer": "Xprinter XP-H500B",
			"DataEntryControls": {
				"SeriaInput": "%s"
			}
		}`, serialSlice, serial))

	wg.Add(2)

	go PrintLocal(data1, channel1, &wg)
	go PrintLocal(data2, channel2, &wg)

	wg.Wait()
	errorText1 := <-channel1
	errorText2 := <-channel2

	// if errorText1 == "ok" {
	// 	return nil
	// } else {
	// 	logrus.Error("error in printing: " + errorText1)
	// 	return errors.New("qaytadan urinib ko'ring")
	// }

	if errorText1 == "ok" && errorText2 == "ok" {
		return nil
	} else {
		logrus.Error("error in printing: " + errorText1 + errorText2)
		return errors.New("qaytadan urinib ko'ring")
	}
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

	type Remont struct {
		Kiritgan     string `json:"kiritgan"`
		RemontPerson string `json:"remont_person"`
		Input        string `json:"input"`
		Output       string `json:"output"`
		Checkpoint   string `json:"checkpoint"`
		DefectName   string `json:"defect"`
		Status       int    `json:"status"`
	}

	type Info struct {
		PackingInfo    []Packing
		ProductionInfo []Production
		GalileoInfo    models.Galileo
		Remont         []Remont
	}

	var packing []Packing
	var galileo models.Galileo

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
			if err == sql.ErrNoRows {
				logrus.Error("NO ROWS")
				return nil, errors.New("no data")
			}
			return nil, err
		}
		packing = append(packing, comp)
	}
	if err = rows1.Err(); err != nil {
		return nil, errors.New("no data")
	}

	r.store.db.QueryRow(`
	select g.serial, g.opcode, g.type_freon,  
	round(g.program_quantity, 2) as program_quantity,
	round(g.real_quantity, 2) as real_quantity,
	round(g.contur_pressure, 2) as contur_pressure, 
	round(g.pre_vacuum, 2) as pre_vacuum,
	round(g."vacuum", 2) as "vacuum", 
	round(g.poisk_utechek, 2) as poisk_utechek, 
	round(g.ref_pressure, 2) as ref_pressure, 
	round(g.ref_temp, 2) as ref_pressure,
	to_char(g."time", 'YYYY-MM-DD HH24:MI') as vaqt  
	from galileo g where g.serial = $1
	`, serial).Scan(&galileo.Serial, &galileo.OpCode, &galileo.TypeFreon, &galileo.ProgramQuantity, &galileo.RealQuantity, &galileo.ConturPressure, &galileo.PreVacuum, &galileo.Vacuum, &galileo.PoiskUtechek, &galileo.RefPressure, &galileo.RefTemp, &galileo.Time)

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
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp Production
		if err := rows.Scan(&comp.Checkpoint, &comp.Time); err != nil {
			return nil, err
		}
		production = append(production, comp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var remont []Remont
	rows2, _ := r.store.db.Query(`
	select r.person_id as kiritgan, r.repair_person as remont_person, 
	COALESCE(to_char(r."input" , 'DD-MM-YYYY HH24:MI'), ' ') as input,
	COALESCE(to_char(r."output" , 'DD-MM-YYYY HH24:MI'), ' ') as output, 
	c."name" as checkpoint, d.defect_name, r.status  
	from remont r, checkpoints c, defects d 
	where r.serial = $1
	and c.id = r.checkpoint_id 
	and d.id  = r.defect_id`, serial)

	defer rows2.Close()
	for rows2.Next() {
		var comp Remont
		if err := rows2.Scan(&comp.Kiritgan, &comp.RemontPerson, &comp.Input, &comp.Output, &comp.Checkpoint, &comp.DefectName, &comp.Status); err != nil {
			return nil, err
		}
		remont = append(remont, comp)
	}

	var productInfo Info
	productInfo.PackingInfo = packing
	productInfo.ProductionInfo = production
	productInfo.GalileoInfo = galileo
	productInfo.Remont = remont

	fmt.Println("remont: ", productInfo.Remont)

	// if productInfo.GalileoInfo.Serial == "" {
	// 	return nil, errors.New("no data")
	// }

	return productInfo, nil
}

func (r *Repo) GalileoInput(g *models.Galileo) error {

	type InputInfo struct {
		id int
	}
	var modelInfo InputInfo

	var serialSlice = g.Serial[0:6]

	if g.RealQuantity != 0 {

		if err := r.store.db.QueryRow("select m.id from models m where m.code = $1", serialSlice).Scan(&modelInfo.id); err != nil {
			return err
		}

		rows, err := r.store.db.Query("insert into galileo (serial, opcode, type_freon, program_quantity, real_quantity, pre_vacuum, model_id, contur_pressure, \"vacuum\", poisk_utechek, ref_pressure, ref_temp, galileo_time) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", g.Serial, g.OpCode, g.TypeFreon, g.ProgramQuantity, g.RealQuantity, g.PreVacuum, modelInfo.id, g.ConturPressure, g.Vacuum, g.PoiskUtechek, g.RefPressure, g.RefTemp, g.Time)
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

	var wg sync.WaitGroup
	wg.Add(1)
	channel := make(chan string, 1)

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
	PrintMetall(data, channel, &wg)

	wg.Wait()
	errorText1 := <-channel

	if errorText1 != "ok" {
		logrus.Error("error in printing: " + errorText1)
		return errors.New("qaytadan urinib ko'ring")
	}

	if err := r.SerialInput(9, countString); err != nil {
		return err
	}
	return nil
}

func (r *Repo) SectorBalanceUpdate(line, component_id int, quantity float64) error {

	result, err := r.store.db.Exec(fmt.Sprintf("update checkpoints.\"%d\" set quantity = %f where component_id = %d", line, quantity, component_id))
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected > 0 {
		return nil
	}

	return errors.New("affected 0")
}

func (r *Repo) SectorBalanceUpdateByQuantity(line, component_id int, quantity float64) error {

	result, err := r.store.db.Exec(fmt.Sprintf("update checkpoints.\"%d\" set quantity = quantity - %f where component_id = %d", line, quantity, component_id))
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected > 0 {
		return nil
	}

	return errors.New("affected 0")
}
