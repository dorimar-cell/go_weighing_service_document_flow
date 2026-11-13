package repository

import "testing"

type sampleDocument struct {
	ID        int64   `json:"id"`
	SessionID int64   `json:"sessionId"`
	Saved     bool    `json:"saved"`
	Name      string  `json:"name"`
	NDS       float64 `json:"nds"`
	Provider  string  `json:"providerOrgName"`
}

func TestDocumentChangesSkipsServiceFields(t *testing.T) {
	before := sampleDocument{ID: 1, SessionID: 7, Saved: true, Name: "old", NDS: 20, Provider: "ООО"}
	after := sampleDocument{ID: 2, SessionID: 7, Saved: true, Name: "new", NDS: 20, Provider: "ООО"}

	changes := documentChanges(&before, &after)
	if len(changes) != 1 {
		t.Fatalf("expected one change, got %d: %v", len(changes), changes)
	}
	if changes["name"]["old"] != "old" || changes["name"]["new"] != "new" {
		t.Fatalf("unexpected name diff: %v", changes["name"])
	}
	for _, field := range []string{"id", "sessionId", "saved", "nds", "providerOrgName"} {
		if _, ok := changes[field]; ok {
			t.Fatalf("field %s must not be in diff", field)
		}
	}
}

func TestDocumentChangesCreationKeepsFilledFields(t *testing.T) {
	after := sampleDocument{ID: 1, SessionID: 7, Name: "Первый", NDS: 20}
	changes := documentChanges(nil, &after)
	if _, ok := changes["name"]; !ok {
		t.Fatal("filled name must appear in creation diff")
	}
	if _, ok := changes["nds"]; !ok {
		t.Fatal("non-zero nds must appear in creation diff")
	}
	if _, ok := changes["providerOrgName"]; ok {
		t.Fatal("empty provider must be skipped in creation diff")
	}
}
