package stt

import (
	"fmt"
	"strings"

	h "github.com/ahopo/stt/helper"
)

type STT struct {
	table      interface{}
	table_name string
	offset     int
	limit      int
	fields     []h.Fields
}
type Condition struct {
	offSet int
	limit  int
	query  string
	where  string
}
type WhereVar struct {
	query []string
}
type LogicalOperator struct {
	EQ   string
	GT   string
	GTEQ string
	LT   string
	LTEQ string
	UEQ  string
}

const (
	_select = "SELECT"
	_insert = "INSERT"
	_delete = "DELETE"
	_update = "UPDATE"
	_create = "CREATE"
	_where  = "WHERE"
	_string = "string"
	_int    = "int"
	_float  = "float64"
)

func New(_struct interface{}) *STT {
	st := new(STT)
	st.table = _struct
	st.table_name = fmt.Sprint(h.StructName(_struct), "s")
	st.fields = h.GetFields(_struct)
	return st
}
func (st *STT) Insert(i interface{}) string {
	var flds []string

	for _, field := range st.fields {
		if strings.ToLower(field.Value) == "id" {
			continue
		}

		flds = append(flds, fmt.Sprint(field.Value))

	}
	return fmt.Sprintf("%s INTO %s (%s) VALUES(%s)", _insert, st.table_name, strings.Join(flds, ","), strings.Join(h.GetValues(i), ","))
}
func (st *STT) Create() string {
	var flds []string
	for _, field := range st.fields {
		flds = append(flds, fmt.Sprint(field.Value, " ", field.DataType, " ", field.Identity))
	}
	return fmt.Sprintf("%s %s (\n%s\n);", _create, st.table_name, strings.Join(flds, ",\n"))
}
func (st *STT) Delete() *Condition {
	s := new(Condition)
	s.query = fmt.Sprintf(`%s FROM %s`, _delete, st.table_name)
	return s
}

func (st *STT) Update(i interface{}) *Condition {

	s := new(Condition)
	var str []string
	data := h.GetValues(i)
	for index, field := range st.fields {
		if strings.ToLower(field.Value) == "id" {
			continue
		}
		if index <= len(data) {
			str = append(str, fmt.Sprintf(`%s = %s`, field.Key, h.GetValues(i)[len(data)-index]))
		}

	}
	s.query = fmt.Sprintf("%s %s SET %s", _update, st.table_name, strings.Join(str, ","))
	return s
}
func (st *STT) Select() *Condition {
	s := new(Condition)
	var str []string
	for _, field := range st.fields {
		str = append(str, field.Value)
	}
	s.query = fmt.Sprintf("%s %s FROM %s", _select, strings.Join(str, ","), st.table_name)
	return s
}
func (s *Condition) OffSet(i int) *Condition {
	s.offSet = i
	return s
}
func (s *Condition) Limit(i int) *Condition {
	s.limit = i
	return s
}
func (s *Condition) Where(i string) *Condition {
	s.where = i
	return s
}

func (st *STT) NewWhereVar() WhereVar {
	wv := new(WhereVar)
	return *wv
}
func (w *WhereVar) Normal(condition string, key string, value interface{}) *WhereVar {
	w.query = append(w.query, fmt.Sprintf(" %s %s %v", key, condition, h.GetSqlString(value)))
	return w
}
func (w *WhereVar) OR(condition string, key string, value interface{}) *WhereVar {
	w.query = append(w.query, fmt.Sprintf("OR %s %s %v", key, condition, h.GetSqlString(value)))
	return w
}
func (w *WhereVar) AND(condition string, key string, value interface{}) *WhereVar {
	w.query = append(w.query, fmt.Sprintf("AND %s %s %v", key, condition, h.GetSqlString(value)))
	return w
}
func (w *WhereVar) GetString() string { return strings.Join(w.query, " ") }
func (s *Condition) Build() string {
	if len(s.where) > 0 {
		s.query = fmt.Sprint(s.query, " WHERE ", s.where)
	}
	if s.limit > 0 {
		s.query = fmt.Sprint(s.query, h.GetLimit(s.limit))
	}
	if s.offSet > 0 {
		s.query = fmt.Sprint(s.query, h.GetOffSet(s.offSet))
	}
	return strings.TrimSpace(s.query)
}
func (stt *STT) LO() LogicalOperator {
	lo := new(LogicalOperator)
	lo.EQ = "="
	lo.GT = ">"
	lo.GTEQ = ">="
	lo.LT = "<"
	lo.LTEQ = "<="
	lo.UEQ = "!="
	return *lo
}
