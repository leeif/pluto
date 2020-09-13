package admin

import (
	"database/sql"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"

	"github.com/MuShare/pluto/utils/general"
	"github.com/MuShare/pluto/utils/mail"
	"github.com/MuShare/pluto/utils/salt"

	"github.com/MuShare/pluto/datatype/request"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/manage"

	"github.com/MuShare/pluto/config"
	plog "github.com/MuShare/pluto/log"
)

func Init(db *sql.DB, config *config.Config, bundle *i18n.Bundle) *perror.PlutoError {

	if config.Admin.Mail == "" {
		return nil
	}

	logger, err := plog.NewLogger(config)

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	manager, err := manage.NewManager(db, config, nil)

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	apps, perr := manager.ListApplications()

	if perr != nil {
		return perr
	}

	for _, app := range apps {
		// skip if the pluto application already exists
		if app.Name == general.PlutoApplication {
			return nil
		}
	}

	ca := request.CreateApplication{}
	ca.Name = general.PlutoApplication
	application, perr := manager.CreateApplication(ca)
	if perr != nil && perr.PlutoCode == perror.ServerError.PlutoCode {
		return perr
	}

	cr := request.CreateRole{}
	cr.Name = general.PlutoAdminRole
	cr.AppID = application.ID

	adminRole, perr := manager.CreateRole(cr)
	if perr != nil && perr.PlutoCode == perror.ServerError.PlutoCode {
		return perr
	}

	cr = request.CreateRole{}
	cr.Name = general.PlutoUserRole
	cr.AppID = application.ID

	userRole, perr := manager.CreateRole(cr)
	if perr != nil && perr.PlutoCode == perror.ServerError.PlutoCode {
		return perr
	}

	ar := request.ApplicationRole{}
	ar.AppID = application.ID
	ar.RoleID = userRole.ID

	if perr := manager.ApplicationDefaultRole(ar); perr != nil {
		return perr
	}

	cs := request.CreateScope{}
	cs.Name = general.PlutoAdminScope
	cs.AppID = application.ID
	adminScope, perr := manager.CreateScope(cs)
	if perr != nil && perr.PlutoCode == perror.ServerError.PlutoCode {
		return perr
	}

	cs = request.CreateScope{}
	cs.Name = general.PlutoUserScope
	cs.AppID = application.ID
	userScope, perr := manager.CreateScope(cs)
	if perr != nil && perr.PlutoCode == perror.ServerError.PlutoCode {
		return perr
	}

	password := salt.RandomToken(20)
	mr := request.MailRegister{}
	mr.Mail = config.Admin.Mail
	name, perr := manager.RandomUserName("pluto_admin_user")
	if perr != nil {
		return perr
	}
	mr.Name = name
	mr.Password = password
	user, perr := manager.RegisterWithEmail(mr, true)
	if perr != nil && perr.PlutoCode != perror.MailIsAlreadyRegister.PlutoCode {
		return perr
	}

	if err == nil {

		mailBody := fmt.Sprintf("Your Pluto Admin Mail : %s, Password : %s", mr.Mail, mr.Password)

		log.Println(mailBody)

		go func() {
			ml, err := mail.NewMail(config, bundle)
			if err != nil {
				logger.Error("smtp server is not set, can't send the mail")
			}
			if err := ml.SendPlainText(mr.Mail, "[Pluto]Admin Password", mailBody); err != nil {
				logger.Error("send mail failed: " + err.LogError.Error())
			} else {
				logger.Info("Mail with your admin login info has been sent")
			}
		}()
	}

	rsu := request.RoleScopeUpdate{}
	rsu.RoleID = adminRole.ID
	rsu.Scopes = []uint{adminScope.ID}

	if err := manager.RoleScopeUpdate(rsu); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	rsu = request.RoleScopeUpdate{}
	rsu.RoleID = userRole.ID
	rsu.Scopes = []uint{userScope.ID}

	if err := manager.RoleScopeUpdate(rsu); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	rs := request.RoleScope{}
	rs.RoleID = adminRole.ID
	rs.ScopeID = adminScope.ID

	if err := manager.RoleDefaultScope(rs); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	rs = request.RoleScope{}
	rs.RoleID = userRole.ID
	rs.ScopeID = userScope.ID

	if err := manager.RoleDefaultScope(rs); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	ur := request.UserRole{}
	ur.AppID = application.ID
	ur.RoleID = adminRole.ID
	ur.UserID = user.ID

	if err := manager.SetUserRole(ur); err != nil && err.PlutoCode == perror.ServerError.PlutoCode {
		return err
	}

	return nil
}
