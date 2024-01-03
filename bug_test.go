package bug

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"entgo.io/bug/ent"
	"entgo.io/bug/ent/enttest"
	"entgo.io/bug/ent/user"

	_ "entgo.io/bug/ent/runtime"
)

func TestBugSQLite(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	test(t, client)
}

func TestBugMySQL(t *testing.T) {
	for version, port := range map[string]int{"56": 3306, "57": 3307, "8": 3308} {
		addr := net.JoinHostPort("localhost", strconv.Itoa(port))
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugPostgres(t *testing.T) {
	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5432, "13": 5433, "14": 5434} {
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugMaria(t *testing.T) {
	for version, port := range map[string]int{"10.5": 4306, "10.2": 4307, "10.3": 4308} {
		t.Run(version, func(t *testing.T) {
			addr := net.JoinHostPort("localhost", strconv.Itoa(port))
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()
	client.User.Delete().ExecX(ctx)

	log.Println("creating user")
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		t.Errorf("unexpected error when creating user: %v", err)
	}

	log.Println("try update one user")
	_, err = u.Update().SetAge(20).Save(ctx)
	if err == nil {
		t.Errorf("expected error updating one user")
	}

	log.Println("try update many user")
	i, err := client.User.Update().Where(user.ID(u.ID)).SetAge(20).Save(ctx)
	if err != nil {
		t.Errorf("unexpected error when updating many users: %v", err)
	}
	if i != 0 {
		t.Errorf("expected update many to apply privacy filter")
	}

	log.Println("try delete many user")
	i, err = client.User.Delete().Where(user.ID(u.ID)).Exec(ctx)
	if err != nil {
		t.Errorf("unexpected error when deleting many users: %v", err)
	}
	if i != 0 {
		t.Errorf("expected delete many to apply privacy filter")
	}

	log.Println("try delete one user")
	err = client.User.DeleteOneID(u.ID).Exec(ctx)
	if err == nil {
		t.Errorf("expected error deleting one user")
	}
}
