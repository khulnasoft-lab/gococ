package sql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Khulnasoft-lab/gococ"
)

// Ping a database.
//
// If the Ping fails, this will return an error.
//
// Params
// - dbname: (required) the name of the database datasource.
//
// Returns:
// - boolean flag set to true if the Ping was successful.
func Ping(cxt gococ.Context, params *gococ.Params) (interface{}, gococ.Interrupt) {
	ok, _ := params.Requires("dbname")
	if !ok {
		e := &gococ.RecoverableError{"Expected a dbname param."}
		return false, e
	}

	dbname := params.Get("dbname", nil).(string)

	db, err := GetDb(cxt, dbname)
	if err != nil {
		return false, err
	}

	err = db.Ping()

	if err != nil {
		return fatalError(err)
	}
	return true, nil
}

// A command that can be used during a shutdown chain.
//
// Params:
// - dbname (required): the name of the db datasource.
func Close(cxt gococ.Context, params *gococ.Params) (interface{}, gococ.Interrupt) {
	ok, _ := params.Requires("dbname")
	if !ok {
		return nil, &gococ.FatalError{"Expected dbname param."}
	}

	dbname := params.Get("dbname", nil).(string)

	db, err := GetDb(cxt, dbname)
	if err != nil {
		return fatalError(err)
	}
	return nil, db.Close()
}

// This is a utility function for executing statements.
//
// While we don't wrap all SQL statements, this particular command is here to
// facilitate creating databases. In other situations, it is assumed that the
// commands will handle SQL internally, and not use high-level commands to run
// each query.
//
// Params:
// - "statement": The statement to execute (as a string)
// - "dbname": The name of the datasource that references the DB.
//
// Returns:
// - database.sql.Result (core Go API)
//
// Example:
//
//	req.Route("install", "Create DB").
//		Does(sql.Execute, "exec").
//		Using("dbname").WithDefault("db").
//		Using("statement").WithDefault("CREATE TABLE IF NOT EXISTS names (id INT, varchar NAME)")
func Execute(cxt gococ.Context, params *gococ.Params) (interface{}, gococ.Interrupt) {
	ok, missing := params.Requires("statement", "dbname")
	if !ok {
		return nil, &gococ.FatalError{fmt.Sprintf("Missing params: %s", strings.Join(missing, ","))}
	}

	dbname := params.Get("dbname", nil).(string)
	statement := params.Get("statement", nil).(string)
	db, err := GetDb(cxt, dbname)
	if err != nil {
		return nil, err
	}

	res, err := db.Exec(statement)
	if err != nil {
		return fatalError(err)
	}

	return res, nil
}

// Utility function to get the database from a datasource.
func GetDb(cxt gococ.Context, dbname string) (*sql.DB, error) {
	dbO, ok := cxt.HasDatasource(dbname)
	if !ok {
		return nil, &gococ.FatalError{fmt.Sprintf("No DB datasource named '%s' found.", dbname)}
	}
	return dbO.(*sql.DB), nil
}

func fatalError(err error) (interface{}, *gococ.FatalError) {
	return nil, &gococ.FatalError{fmt.Sprintf("%v", err)}
}
