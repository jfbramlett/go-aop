package messaging

import (
    "fmt"
    "github.com/jfbramlett/go-aop/pkg/jsonutils"
    "testing"
)

type TestPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSerializeDeserialize(t *testing.T) {
	envelope := Envelope{Content: &TestPayload{Name: "John", Age: 50}}

	str, err := jsonutils.ToJSON(envelope)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)

	newEnvelope := Envelope{Content: ContentType()}
	err = jsonutils.FromJSON(str, &newEnvelope)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v", newEnvelope)
}

func ContentType() interface{} {
	return &TestPayload{}
}
