package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)
type Customer struct {
	Name string
	Age int
	Children []string
	Company *Company
}
type Company struct{
	Name string
	Location string
}
func TestStateEntityJson(t *testing.T){
	customer:=&Customer{
		Name: "Devin",
		Age:  36,
		Children: []string{"a","bb"},
		Company: &Company{Name: "CICC",Location: "Guomao"},
	}
	input,_:=json.Marshal(customer)
	t.Log(string(input))
	m:=make(map[string]interface{})
	json.Unmarshal(input,&m)
	t.Logf("%#v",m)
	entity1:=&ObjectEvidence{
		Owner:     "studyzy",
		Category:  "Customer",
		Timestamp: time.Now().Unix(),
		Object:    m,
	}
	json1,err:=json.Marshal(entity1)
	assert.Nil(t,err)
	t.Log(string(json1))
	entity2:= ObjectEvidence{}
	json.Unmarshal(json1,&entity2)
	t.Logf("%#v",entity2)
	out,_:=json.Marshal(entity2.Object)
	t.Log(string(out))
}
