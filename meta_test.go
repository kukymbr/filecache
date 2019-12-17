package filecache_test

import (
	"gitlab.com/kukymbrgo/filecache"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMetaFromFile(t *testing.T) {
	jsonOk := []byte(
		`{"k":"https://zh.wikipedia.org/wiki/%E5%BA%9E%E5%BE%B7%E4%BC%AF%E9%87%8C",` +
			`"n":"pages","t":1,"c":1576120592,"o":"TestName","f":{"f1":"v1","f2":"v2"}}`)
	jsonBad := []byte(
		`{"bad":"json"}`)

	cachePath, err := ioutil.TempDir("", "kukymbrgo-filecache-test")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = os.RemoveAll(cachePath); err != nil {
			t.Error("failed to clean up after test")
		}
	}()

	err = ioutil.WriteFile(cachePath+"/good.json", jsonOk, 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(cachePath+"/bad.json", jsonBad, 0644)
	if err != nil {
		panic(err)
	}

	metaGood, err := filecache.MetaFromFile(cachePath + "/good.json")
	if err != nil {
		t.Error("failed to create meta from file", err)
		return
	}

	expected := &filecache.Meta{
		Key:          "https://zh.wikipedia.org/wiki/%E5%BA%9E%E5%BE%B7%E4%BC%AF%E9%87%8C",
		Namespace:    "pages",
		OriginalName: "TestName",
		TTL:          1,
		Created:      1576120592,
		Fields: map[string]interface{}{
			"f1": "v1",
			"f2": "v2",
		},
	}

	if !reflect.DeepEqual(metaGood, expected) {
		t.Error("meta is not equal to expected")
		return
	}

	metaBad, err := filecache.MetaFromFile(cachePath + "/bad.json")
	if err == nil {
		t.Error("missing expected error while reading bad JSON data")
		return
	}

	if metaBad != nil {
		t.Error("expected nil, got meta")
		return
	}
}
