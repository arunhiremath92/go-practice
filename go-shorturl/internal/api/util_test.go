package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeJson(t *testing.T) {
	type SampleJsonObj struct {
		IntVar    int     `json:"intVal"`
		FloatVar  float32 `json:"floatVal"`
		StringVar string  `json:"stringVal"`
	}
	input := `{"intVal":20,"floatVal":32.4,"stringVal": "hello_world"}`

	expectedOutput := SampleJsonObj{
		IntVar:    20,
		FloatVar:  32.4,
		StringVar: "hello_world",
	}

	var generatedOutput SampleJsonObj

	req, err := http.NewRequest("POST", "/test", strings.NewReader(input))
	if err != nil {
		log.Fatal("failed to make a client request")
	}
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	err = decodeJSON(w, req, &generatedOutput)
	if err != nil {
		log.Fatal("failed to decode the json object ", err)
	}

	if !reflect.DeepEqual(expectedOutput, generatedOutput) {
		t.Errorf("expected %+v received %+v", expectedOutput, generatedOutput)
	}
}
