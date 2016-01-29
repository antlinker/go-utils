package mgoid

import "testing"

func TestGenerateID(t *testing.T) {
	incID, err := NewIncIDWithURL("mongodb://192.168.33.70:27017", "ant")
	if err != nil {
		t.Error(err)
		return
	}
	defer incID.Close()
	id, err := incID.GenerateID("Type1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ID:", id)
}
