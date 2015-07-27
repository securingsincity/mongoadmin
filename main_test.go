package mongoadmin

import (
	// "encoding/json"
	// "net/http"
	// "net/http/httptest"
	// "strings"
	// . "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mainTestHelper() *DatabaseConfig {
	db := &DatabaseConfig{}
	db.ConnectionString = "localhost"
	db.Label = "servertest"
	db.Database = "servertest"
	appConfig = &Config{}
	appConfig.DB = append(appConfig.DB, db)
	return db
}

func TestGetDatabaseAndGetCollection(t *testing.T) {
	db := mainTestHelper()
	sess, dbObj, err := getDatabaseByName(db.Database)
	if err != nil {
		t.Errorf(err.Error())
	}
	if dbObj.Session != sess {
		t.Errorf("db session should match returned session")
	}
	col := getCollectionByName(dbObj, "foo")

	if col.Name != "foo" {
		t.Errorf("Collection shouldn't be null")
	}
}
