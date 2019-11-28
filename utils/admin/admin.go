package admin

import (
	"database/sql"

	"github.com/leeif/pluto/utils/salt"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/leeif/pluto/models"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/leeif/pluto/config"
)

var randomToken string

func Init(db *sql.DB, config *config.Config) (string, error) {

	token := false

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	defer func() {
		tx.Rollback()
	}()

	if application, err := models.Applications(qm.Where("name = ?", config.Admin.Application)).One(tx); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		application = &models.Application{}
		application.Name = config.Admin.Application
		if err := application.Insert(tx, boil.Infer()); err != nil {
			return "", err
		}
		token = true
	}

	if role, err := models.RbacRoles(qm.Where("name = ?", config.Admin.Role)).One(tx); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		role = &models.RbacRole{}
		role.Name = config.Admin.Role
		if err := role.Insert(tx, boil.Infer()); err != nil {
			return "", err
		}
		token = true
	}

	if scope, err := models.RbacScopes(qm.Where("name = ?", config.Admin.Scope)).One(tx); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		scope = &models.RbacScope{}
		scope.Name = config.Admin.Scope
		if err := scope.Insert(tx, boil.Infer()); err != nil {
			return "", err
		}
		token = true
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if token {
		randomToken = salt.RandomToken(20)
		return randomToken, nil
	}

	return "", nil
}
