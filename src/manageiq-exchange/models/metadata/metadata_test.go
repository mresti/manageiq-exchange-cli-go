package metadata

import (
	"reflect"
	"testing"
)

func TestMetadataInit(t *testing.T) {

	want := Metadata{
		CurrentPage: 2,
		TotalPages:  1,
		TotalCount:  1,
	}

	got := Metadata{}

	var data = map[string]interface{}{
		"current_page": 2,
		"total_pages":  1,
		"total_count":  1,
	}

	got.Init(data)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("MetadataInit returned %+v, want %+v", got, want)
	}
}
