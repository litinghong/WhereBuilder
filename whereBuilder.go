package whereBuilder

import (
	"fmt"
	"strconv"
	"strings"
)

const condAnd = "and"
const condOr = "or"

type WhereBuilder map[string]interface{}

type Condition struct {
	field      string
	operation  string
	value      interface{}
	conditions []*Condition
	logic      string
}

type prepareBuilder struct {
	k      int
	values []interface{}
}

func (wb WhereBuilder) preProcess() {
	for field, condition := range wb {
		if len(field) >= 2 && field[0:2] == "__" {
			continue
		}
		switch condition.(type) {
		case string:
			wb[field] = Cond(field, "=", condition)
		case int:
			wb[field] = Cond(field, "=", condition)
		case []string:
			wb[field] = Cond(field, "in", condition)
		case []*string:
			wb[field] = Cond(field, "in", condition)
		case []int:
			wb[field] = Cond(field, "in", condition)
		case []*int:
			wb[field] = Cond(field, "in", condition)
		case []float32:
			wb[field] = Cond(field, "in", condition)
		case []*float32:
			wb[field] = Cond(field, "in", condition)
		case []float64:
			wb[field] = Cond(field, "in", condition)
		case []*float64:
			wb[field] = Cond(field, "in", condition)
		}
	}
}

func (wb WhereBuilder) ToPrepare(builder ...*prepareBuilder) (prepare string, values []interface{}) {
	var p *prepareBuilder
	if builder == nil {
		p = &prepareBuilder{
			k:      0,
			values: make([]interface{}, 0),
		}
	} else {
		p = builder[0]
	}

	var sb strings.Builder
	sb.WriteString("(")
	var i = 0
	var max = len(wb)
	wb.preProcess()
	for field, condition := range wb {
		switch condition.(type) {
		case *Condition:
			cond := condition.(*Condition)
			if cond.conditions != nil {
				sb.WriteString("(")
				sb.WriteString(p.groupToPrepare(cond))
				sb.WriteString(")")
			} else {
				sb.WriteString(p.conditionToPrepare(cond))
			}
		case WhereBuilder:
			cond := condition.(WhereBuilder)
			str, _ := cond.ToPrepare(p)
			sb.WriteString(str)
		case *WhereBuilder:
			cond := condition.(*WhereBuilder)
			str, _ := cond.ToPrepare(p)
			sb.WriteString(str)
		case string:
			if len(field) < 2 || field[0:2] != "__" {
				sb.WriteString(" ")
				sb.WriteString(field)
			}
			sb.WriteString(condition.(string))
		}

		if i < max-1 {
			sb.WriteString(" ")
			sb.WriteString(condAnd)
			sb.WriteString(" ")
		}
		i++
	}
	sb.WriteString(")")

	return sb.String(), p.values
}

func (wb WhereBuilder) ToSql() string {
	var sb strings.Builder
	sb.WriteString("(")
	var i = 0
	var max = len(wb)
	wb.preProcess()
	for field, condition := range wb {
		switch condition.(type) {
		case *Condition:
			cond := condition.(*Condition)
			if cond.conditions != nil {
				sb.WriteString("(")
				sb.WriteString(groupToSql(cond))
				sb.WriteString(")")
			} else {
				sb.WriteString(conditionToSql(cond))
			}
		case WhereBuilder:
			cond := condition.(WhereBuilder)
			sb.WriteString(cond.ToSql())
		case *WhereBuilder:
			cond := condition.(*WhereBuilder)
			sb.WriteString(cond.ToSql())
		case string:
			if len(field) < 2 || field[0:2] != "__" {
				sb.WriteString(" ")
				sb.WriteString(field)
			}

			sb.WriteString(condition.(string))
		}

		if i < max-1 {
			sb.WriteString(" ")
			sb.WriteString(condAnd)
			sb.WriteString(" ")
		}
		i++
	}
	sb.WriteString(")")

	return sb.String()
}

func conditionToSql(c *Condition) string {
	sb := strings.Builder{}
	sb.WriteString(" ")
	sb.WriteString(c.field)
	sb.WriteString(" ")
	sb.WriteString(c.operation)
	sb.WriteString(" ")

	switch c.value.(type) {
	case string:
		sb.WriteString(`"`)
		str := c.value.(string)
		str = strings.Replace(str, `"`, `\"`, -1)
		sb.WriteString(str)
		sb.WriteString(`"`)
	case int:
		sb.WriteString(strconv.Itoa(c.value.(int)))
	case float32:
		sb.WriteString(fmt.Sprintf("%f", c.value.(float32)))
	case float64:
		sb.WriteString(fmt.Sprintf("%f", c.value.(float32)))
	case bool:
		if c.value.(bool) {
			sb.WriteString("1")
		} else {
			sb.WriteString("0")
		}
	case []int:
		sb.WriteString("(")
		sb.WriteString(JoinInt(c.value.([]int), ","))
		sb.WriteString(")")
	case []float32:
		sb.WriteString("(")
		sb.WriteString(JoinFloat32(c.value.([]float32), ","))
		sb.WriteString(")")
	case []float64:
		sb.WriteString("(")
		sb.WriteString(JoinFloat64(c.value.([]float64), ","))
		sb.WriteString(")")
	case []string:
		sb.WriteString(`("`)
		strList := c.value.([]string)
		for k, v := range strList {
			strList[k] = strings.Replace(v, `"`, `\"`, -1)
		}
		str := strings.Join(strList, `","`)
		sb.WriteString(str)
		sb.WriteString(`")`)
	}

	sb.WriteString(" ")
	return sb.String()
}

func (p *prepareBuilder) getPlace(value interface{}) string {
	p.k++

	switch value.(type) {
	case bool:
		if value.(bool) {
			p.values = append(p.values, 1)
		} else {
			p.values = append(p.values, 0)
		}
	default:
		p.values = append(p.values, value)
	}

	return "?"
}

func groupToSql(group *Condition) string {
	sb := strings.Builder{}
	max := len(group.conditions)
	for k, v := range group.conditions {
		if v.conditions != nil {
			sb.WriteString(" (")
			sb.WriteString(groupToSql(v))
			sb.WriteString(")")
		} else {
			sb.WriteString(conditionToSql(v))
		}

		if k < max-1 {
			sb.WriteString(group.logic)
		}
	}

	return sb.String()
}

func (p *prepareBuilder) conditionToPrepare(c *Condition) string {
	sb := strings.Builder{}
	sb.WriteString(" ")
	sb.WriteString(c.field)
	sb.WriteString(" ")
	sb.WriteString(c.operation)
	sb.WriteString(" ")

	switch c.operation {
	case "in":
		sb.WriteString("(")
		sb.WriteString(p.getPlace(c.value))
		sb.WriteString(")")
	default:
		sb.WriteString(p.getPlace(c.value))
	}

	sb.WriteString(" ")
	return sb.String()
}

func (p *prepareBuilder) groupToPrepare(c *Condition) string {
	sb := strings.Builder{}
	max := len(c.conditions)
	for k, v := range c.conditions {
		if v.conditions != nil {
			sb.WriteString(" (")
			sb.WriteString(p.groupToPrepare(v))
			sb.WriteString(")")
		} else {
			sb.WriteString(p.conditionToPrepare(v))
		}

		if k < max-1 {
			sb.WriteString(c.logic)
		}
	}

	return sb.String()
}

func Cond(field, operation string, value interface{}) *Condition {
	return &Condition{
		field:     field,
		operation: operation,
		value:     value,
	}
}

func And(conditions ...*Condition) *Condition {
	return &Condition{
		conditions: conditions,
		logic:      condAnd,
	}
}

func Or(conditions ...*Condition) *Condition {
	return &Condition{
		conditions: conditions,
		logic:      condOr,
	}
}
