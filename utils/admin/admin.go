package admin

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/leeif/pluto/utils/general"
	"github.com/leeif/pluto/utils/mail"
	"github.com/leeif/pluto/utils/salt"

	"github.com/leeif/pluto/datatype/request"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/manage"

	"github.com/leeif/pluto/config"
)

func Init(db *sql.DB, config *config.Config) *perror.PlutoError {

	if config.Admin.Mail == "" || config.Admin.Name == "" {
		return nil
	}

	manager := manage.NewManager(db, config, nil)

	ca := request.CreateApplication{}
	ca.Name = general.PlutoAdminApplication
	application, err := manager.CreateApplication(ca)
	if err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	cr := request.CreateRole{}
	cr.Name = general.PlutoAdminRole
	cr.AppID = application.ID

	role, err := manager.CreateRole(cr)
	if err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	cs := request.CreateScope{}
	cs.Name = general.PlutoAdminScope
	cs.AppID = application.ID
	scope, err := manager.CreateScope(cs)
	if err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	password := salt.RandomToken(20)
	mr := request.MailRegister{}
	mr.Mail = config.Admin.Mail
	mr.Name = config.Admin.Name
	mr.Password = password
	user, err := manager.RegisterWithEmail(mr)
	if err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	if err == nil {

		user.Verified.SetValid(true)
		if _, err := user.Update(db, boil.Infer()); err != nil {
			return perror.ServerError.Wrapper(err)
		}

		mailBody := fmt.Sprintf("Your Pluto Admin Mail : %s, Password : %s", mr.Mail, mr.Password)

		log.Println(mailBody)

		go func() {
			ml, err := mail.NewMail(config)
			if err != nil {
				log.Println("smtp server is not set, mail can't be send")
			}
			if err := ml.SendPlainText(mr.Mail, "[Pluto]Admin Password", mailBody); err != nil {
				log.Println("send mail failed: " + err.LogError.Error())
			}
			log.Println("Mail with your admin login info is send")
		}()
	}

	rsu := request.RoleScopeUpdate{}
	rsu.RoleID = role.ID
	rsu.Scopes = []uint{scope.ID}

	if err := manager.RoleScopeUpdate(rsu); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	rs := request.RoleScope{}
	rs.RoleID = role.ID
	rs.ScopeID = scope.ID

	if err := manager.RoleDefaultScope(rs); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	ur := request.UserRole{}
	ur.AppID = application.ID
	ur.RoleID = role.ID
	ur.UserID = user.ID

	if err := manager.SetUserRole(ur); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	return nil
}
