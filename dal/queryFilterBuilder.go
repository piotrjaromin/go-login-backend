package dal

import (
        "github.com/piotrjaromin/login-template/common"
)

type QueryBuilder struct {
        query Query
}

type Query struct {
        fields      map[string]interface{}
        orFields    map[string]interface{}
        projections map[string]int
        sort        []string
}

type Sort int

const (
        Desc Sort = 1
        Asc Sort = -1
)

//Creates new query builder
func NewQueryBuilder() *QueryBuilder {
        return &QueryBuilder{
                query: Query{
                        fields: map[string]interface{}{},
                        projections: map[string]int{},
                        sort: []string{},
                },
        }
}

func RegexpValue(value string) interface{} {
        return map[string]interface{}{
                "$regexp" : value,
        }
}

//Adds field param
func (qb *QueryBuilder) WithField(field string, value interface{}) *QueryBuilder {
        qb.query.fields[field] = value
        return qb
}

func (qb *QueryBuilder) WithId(value string) *QueryBuilder {
        qb.query.fields["_id"] = value
        return qb
}

func (qb *QueryBuilder) WithFieldInArray(id string, field string, value string) *QueryBuilder {
        qb.query.fields["_id"] = map[string]interface{}{id : map[string]interface{}{"$elemMatch":  map[string]interface{}{field : value }}}
        return qb
}

func (qb *QueryBuilder) WithAnyOfValues(field string, values ...interface{}) *QueryBuilder {
        if len(values) == 0 {
                return qb
        }

        if len(values) == 1 {
                qb.WithField(field, values[0])
                return qb
        } else {
                qb.query.orFields[field] = values
                return qb
        }
}

func (qb *QueryBuilder) SortBy(field string, sort Sort) *QueryBuilder {

        if sort == Asc {
                field = "-" + field
        }

        qb.query.sort = append(qb.query.sort, field)
        return qb
}

func (qb *QueryBuilder) WithProjections(projections... string) *QueryBuilder {

        for _, projection := range projections {
                qb.query.projections[projection] = 1
        }
        return qb
}

func (qb *QueryBuilder) WithSortFields(fields common.SortFields) *QueryBuilder {

        for _, field := range fields {
                if ( field.Order == common.Desc ) {
                        qb.SortBy(field.Name, Desc)
                } else {
                        qb.SortBy(field.Name, Asc)
                }
        }

        return qb
}

//Creates filter field value which can be used with Query
func (qb *QueryBuilder) Build() Query {

        if len(qb.query.orFields) > 0 {
                qb.query.fields["$or"] = qb.query.orFields
        }

        return qb.query
}
