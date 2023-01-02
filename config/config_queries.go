package config

import (
	"fmt"
	"strings"

	"ariga.io/atlas/sql/schema"
)

type Queries struct {
	schema *Schema             `json:"-"`
	Custom []*Query            `json:"custom"` // user defined - not overwritten by schema
	Tables map[string][]*Query `json:"tables"` // auto generated by schema
}

func (t *Queries) init(schema *Schema) {
	t.schema = schema
	t.Custom = make([]*Query, 0, 10)
	t.Tables = make(map[string][]*Query)
}

func (t *Queries) InitQueryTables(tables []*schema.Table) error {
	for _, table := range tables {
		err := t.initQueryTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Queries) initQueryTable(table *schema.Table) error {
	// insert all
	questionare := strings.Repeat("?, ", len(table.Columns))
	questionare = questionare[:len(questionare)-2]
	t.AddQueryTables(table.Name, &Query{
		Name:    "insert",
		Comment: "default query - insert all",
		Sql:     fmt.Sprintf("INSERT INTO %s VALUES (%s)", table.Name, questionare),
	})

	// select all
	t.AddQueryTables(table.Name, &Query{
		Name:    "select",
		Comment: "default query - select all",
		Sql:     fmt.Sprintf("SELECT * FROM %s", table.Name),
	})

	// TODO: select where by index

	// TODO: update

	// delete
	t.AddQueryTables(table.Name, &Query{
		Name:    "delete",
		Comment: "default query - delete all",
		Sql:     fmt.Sprintf("DELETE FROM %s", table.Name),
	})

	return nil
}

func (t *Queries) AddQueryTables(tableName string, query *Query) {
	query.Schema = t.schema
	if t.Tables == nil {
		t.Tables = make(map[string][]*Query, 10)
	}
	if _, ok := t.Tables[tableName]; ok == false {
		t.Tables[tableName] = make([]*Query, 0, 10)
	}
	t.Tables[tableName] = append(t.Tables[tableName], query)
}

func (t *Queries) AddQueryCustom(query *Query) {
	query.Schema = t.schema
	if t.Custom == nil {
		t.Custom = make([]*Query, 0, 10)
	}
	t.Custom = append(t.Custom, query)
}

//------------------------------------------------------------------------------------------------//
// query

type Query struct {
	Name    string  `json:"name"`
	Comment string  `json:"comment,omitempty"`
	Sql     string  `json:"sql"`
	Schema  *Schema `json:"-"` // 쿼리 제작을 위한 전체 스키마 정보

	// options
	SelectFieldTypes []*SelectFieldType `json:"fields,omitempty"`
	InsertMulti      bool               `json:"insert_multi,omitempty"`
	UpdateNullIgnore bool               `json:"update_null_ignore,omitempty"`
	ErrQuery         string             `json:"-"`
	ErrParser        string             `json:"-"`
}

// select 만 field type이 있는 이유
// select query 는 bp.json 의 schema type 을 통해 타입을 지정할 수 없기 때문에
// 직접 쿼리를 select 를 하고 결과를 추출해 타입에 넣음
// snum, uint 등의 custom type 은 여기서 처리
type SelectFieldType struct {
	Name    string `json:"name"`
	TypeGen string `json:"type"`
}

//------------------------------------------------------------------------------------------------//
// query

func (t *Query) Init(name, sql string) {
	t.Name = name
	t.Sql = sql
	t.SelectFieldTypes = make([]*SelectFieldType, 0, 10)
}

func (t *Query) AddFieldType(name string, typeGen string) {
	if t.SelectFieldTypes == nil {
		t.SelectFieldTypes = make([]*SelectFieldType, 0, 10)
	}
	t.SelectFieldTypes = append(t.SelectFieldTypes, &SelectFieldType{
		Name:    name,
		TypeGen: typeGen,
	})
}

func (t *Query) GetFieldType(name string) (genType string) {
	for _, pt := range t.SelectFieldTypes {
		if pt.Name == name {
			return pt.TypeGen
		}
	}
	return ""
}
