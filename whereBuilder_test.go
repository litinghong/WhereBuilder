package whereBuilder

import (
	"testing"
)

func TestAnd(t *testing.T) {
	cond := WhereBuilder{}

	cond["a"] = Cond("a", "in", []int{1, 2, 3})
	cond["b"] = Cond("b", "=", 2)
	cond["c"] = And(Cond("c", "=", 2), Cond("d", "=", "x"))
	cond["d"] = Or(Cond("c", "like", "%2%"), Cond("d", "=", "x"))
	//cond["e"] = Or(Cond("f",">",4), Cond("g","=",5), And(Cond("c","=",2), Cond("d","=","x")))
	//cond["a"] = Or(And(Cond("a", "=", "1"), Cond("b", "=", "2")), Or(Cond("c", "=", "1"), Cond("d", "=", "2")))
	//cond["__x"] = "xxx = 1"

	//cond["a"] = WhereBuilder{
	//	"a": " = b",
	//	"b": Cond("b", "=", "B"),
	//}
	//cond["c"] = And(Cond("c","=","C"), Cond("d","=","D"))

	//t.Log(cond.ToSql())
	t.Log(cond.ToPrepare())

	where := WhereBuilder{}
	where["series_id"] = []int{142, 142}
	where["years_num"] = []int{2004, 2005}
	where["brand_id"] = []int{62, 13}
	others := []string{"豪华", "舒适"}
	condList := make([]*Condition, 0, len(others))
	for _, other := range others {
		condList = append(condList, Cond("model_name", "like", "%"+other+"%"))
	}
	where["model_name"] = Or(condList...)

	str, val := where.ToPrepare()
	t.Log(str, val)
}

func TestPreProcess(t *testing.T) {
	wb := WhereBuilder{}
	wb["a"] = "A"
	wb["b"] = []string{"\"B1", "B2"}
	wb["c"] = []float32{3.14, 3.15}
	wb["__c"] = " x = y"
	str, val := wb.ToPrepare()
	t.Log(str, val)
	sql := wb.ToSql()
	t.Log(sql)
}
