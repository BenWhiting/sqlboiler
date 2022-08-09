package boil

import (
	"testing"

	"github.com/pashagolub/pgxmock"
)

func TestGetSetDB(t *testing.T) {
	t.Parallel()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	SetDB(&BoilerPgxWrap{mock})

	if GetDB() == nil {
		t.Errorf("Expected GetDB to return a database handle, got nil")
	}
}
