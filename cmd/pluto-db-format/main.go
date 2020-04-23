package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/models"

	"github.com/leeif/pluto/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.NewConfig(os.Args, "")

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase(cfg)

	if err != nil {
		log.Fatal(err)
	}

	formatUsersTableData(db)
	formatBindingsTableData(db)

}

func formatUsersTableData(db *sql.DB) {
	users, err := models.Users().All(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.Name == "" {
			user.Name = "user"
		}
		name := fmt.Sprintf("%s-%d", user.Name, user.ID)
		res := fmt.Sprintf("insert into users (id, created_at, updated_at, name, password, avatar, verified) values (%d, `%s`, `%s`, `%s`, `%s`, `%s`, %d)",
			user.ID,
			user.CreatedAt.Time.Local().Format("2006-Jan-02 15:04:05"),
			user.UpdatedAt.Time.Local().Format("2006-Jan-02 15:04:05"),
			name,
			user.Password.String,
			user.Avatar.String,
			1)
		fmt.Println(res)
	}
}

func formatBindingsTableData(db *sql.DB) {
	users, err := models.Users().All(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.Name == "" {
			user.Name = "user"
		}
		res := fmt.Sprintf("insert into bindings (created_at, updated_at, login_type, identify_token, mail, verified, user_id) values (`%s`, `%s`, `%s`, `%s`, `%s`, %d, %d)",
			user.CreatedAt.Time.Local().Format("2006-Jan-02 15:04:05"),
			user.UpdatedAt.Time.Local().Format("2006-Jan-02 15:04:05"),
			user.LoginType,
			user.IdentifyToken,
			user.Mail.String,
			1,
			user.ID)
		fmt.Println(res)
	}
}
