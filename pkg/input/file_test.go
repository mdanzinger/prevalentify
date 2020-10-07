package input

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestFromFileHappyPath(t *testing.T) {
	// create fake input file
	tmpFile, err := ioutil.TempFile("", "FromFileInput")
	if err != nil {
		t.Fatalf("unable to create test file: %s", err)
	}

	// populate with fake data
	_, err = tmpFile.WriteString(fakeData)
	if err != nil {
		t.Fatalf("unable to populate test file: %s", err)
	}

	// create input generator from the test file
	inputGenerator, err := FromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("got unexpected error creating input generator")
	}

	// assert it's correctly returning input
	inputCh, _ := inputGenerator(context.Background())
	var result []string
	for input := range inputCh {
		result = append(result, string(input))
	}

	expectedResult := []string{
		"http://i.imgur.com/FApqk3D.jpg",
		"http://i.imgur.com/TKLs9lo.jpg",
		"https://i.redd.it/d8021b5i2moy.jpg",
	}

	if !reflect.DeepEqual(expectedResult, result) {
		t.Fatalf("input gengerator expected to return: \n %v, \n but got: \n %v ", expectedResult, result)
	}
}

var fakeData = `http://i.imgur.com/FApqk3D.jpg
http://i.imgur.com/TKLs9lo.jpg
https://i.redd.it/d8021b5i2moy.jpg`
