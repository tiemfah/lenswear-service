//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Users = newUsersTable("public", "users", "")

type usersTable struct {
	postgres.Table

	//Columns
	UserID     postgres.ColumnString
	RoleID     postgres.ColumnString
	Username   postgres.ColumnString
	Password   postgres.ColumnString
	CreateDate postgres.ColumnTimestamp
	UpdateDate postgres.ColumnTimestamp
	DeleteDate postgres.ColumnTimestamp
	CreateBy   postgres.ColumnString
	UpdateBy   postgres.ColumnString
	DeleteBy   postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UsersTable struct {
	usersTable

	EXCLUDED usersTable
}

// AS creates new UsersTable with assigned alias
func (a UsersTable) AS(alias string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UsersTable with assigned schema name
func (a UsersTable) FromSchema(schemaName string) *UsersTable {
	return newUsersTable(schemaName, a.TableName(), a.Alias())
}

func newUsersTable(schemaName, tableName, alias string) *UsersTable {
	return &UsersTable{
		usersTable: newUsersTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newUsersTableImpl("", "excluded", ""),
	}
}

func newUsersTableImpl(schemaName, tableName, alias string) usersTable {
	var (
		UserIDColumn     = postgres.StringColumn("user_id")
		RoleIDColumn     = postgres.StringColumn("role_id")
		UsernameColumn   = postgres.StringColumn("username")
		PasswordColumn   = postgres.StringColumn("password")
		CreateDateColumn = postgres.TimestampColumn("create_date")
		UpdateDateColumn = postgres.TimestampColumn("update_date")
		DeleteDateColumn = postgres.TimestampColumn("delete_date")
		CreateByColumn   = postgres.StringColumn("create_by")
		UpdateByColumn   = postgres.StringColumn("update_by")
		DeleteByColumn   = postgres.StringColumn("delete_by")
		allColumns       = postgres.ColumnList{UserIDColumn, RoleIDColumn, UsernameColumn, PasswordColumn, CreateDateColumn, UpdateDateColumn, DeleteDateColumn, CreateByColumn, UpdateByColumn, DeleteByColumn}
		mutableColumns   = postgres.ColumnList{RoleIDColumn, UsernameColumn, PasswordColumn, CreateDateColumn, UpdateDateColumn, DeleteDateColumn, CreateByColumn, UpdateByColumn, DeleteByColumn}
	)

	return usersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:     UserIDColumn,
		RoleID:     RoleIDColumn,
		Username:   UsernameColumn,
		Password:   PasswordColumn,
		CreateDate: CreateDateColumn,
		UpdateDate: UpdateDateColumn,
		DeleteDate: DeleteDateColumn,
		CreateBy:   CreateByColumn,
		UpdateBy:   UpdateByColumn,
		DeleteBy:   DeleteByColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
