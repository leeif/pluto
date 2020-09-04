package manage

import (
	"database/sql"
	"time"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/datatype/request"
	"github.com/MuShare/pluto/models"
	"github.com/MuShare/pluto/utils/jwt"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (m *Manager) RefreshAccessToken(rat request.RefreshAccessToken) (*GrantResult, *perror.PlutoError) {
	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	application, perr := m.getApplication(tx, rat.AppID)
	if perr != nil {
		return nil, perr
	}

	rt, err := models.RefreshTokens(qm.Where("refresh_token = ?", rat.RefreshToken)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.InvalidRefreshToken
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if time.Now().After(rt.ExpireAt) {
		return nil, perror.RefreshTokenExpired
	}

	da, err := models.DeviceApps(qm.Where("id = ?", rt.DeviceAppID)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.InvalidRefreshToken
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if da.AppID != application.ID {
		return nil, perror.InvalidRefreshToken
	}

	scopes, perr := getValidScopes(tx, rat.Scopes, rt.UserID, application.Name)

	if perr != nil {
		return nil, perr
	}

	if scopes == "" && !rt.Scopes.IsZero() {
		scopes = rt.Scopes.String
	}

	if perr := m.updateRefreshToken(tx, rt, scopes); perr != nil {
		return nil, perr
	}

	grantResult, perr := m.grantToken(rt.UserID, rt, scopes, rat.AppID, m.config.Token.AccessTokenExpire)
	if perr != nil {
		return nil, perr
	}

	tx.Commit()
	return grantResult, nil
}

func (m *Manager) VerifyAccessToken(accessToken string) (*jwt.AccessPayload, *perror.PlutoError) {
	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	jwtToken, perr := jwt.VerifyRS256JWT(accessToken)

	if perr != nil {
		return nil, perror.InvalidJWTToken
	}

	accessPayload := &jwt.AccessPayload{}

	if perr := jwtToken.UnmarshalPayload(accessPayload); perr != nil {
		return nil, perr
	}

	if accessPayload.Type != jwt.ACCESS {
		return nil, perror.InvalidAccessToken
	}

	if time.Now().Unix() > accessPayload.Expire {
		return nil, perror.JWTTokenExpired
	}

	return accessPayload, nil
}
