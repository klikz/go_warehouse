package sqlstore

import (
	"fmt"
	"warehouse/internal/app/models"

	"github.com/sirupsen/logrus"
)

type Repo struct {
	store *Store
}

func (r *Repo) Create(u *models.User) error {
	if err := r.store.db.QueryRow(`insert into users (email, encrypted_password, "role") values ($1, $2, $3) returning id`, u.Email, u.EncryptedPassword, u.Role).Scan(&u.ID); err != nil {
		return err
	}
	if err := r.store.db.QueryRow(`insert into routes ("user") values ($1) returning id`, u.Email).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

func (r *Repo) FindByEmail(u *models.User) error {
	if err := r.store.db.QueryRow(`select u.encrypted_password, u."role" from users u where u.email = $1`, u.Email).Scan(&u.EncryptedPassword, &u.Role); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *Repo) CheckRole(route, email string) (bool, error) {

	result := false
	err := r.store.db.QueryRow(fmt.Sprintf(`
	select r."%s" as result from routes r where "user" = '%s'
	 `, route, email)).Scan(&result)
	if err != nil {
		logrus.Error("check role: ", err)
		return false, err
	}
	return result, nil
}
