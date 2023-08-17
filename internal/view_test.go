package promo

import (
	"bytes"
	"fmt"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ DBConnexion = &Registry{}

func TestHandlers_HealthCheck(t *testing.T) {
	t.Run("Check server health (no errors)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusOK
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		w := httptest.NewRecorder()
		h.HealthCheck(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
}

func TestHandlers_CreateUser(t *testing.T) {
	t.Run("Check creating user (no errors)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		name := "And"
		surname := "Ersen"
		var position Grade = 1
		project := "Test"
		mock.ExpectExec("INSERT INTO usr").WithArgs(name, surname, dGrades[position], project).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		expected := http.StatusOK
		body := bytes.NewReader([]byte(`{"name": "And","surname": "Ersen", "position": 1, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPost, "/create", body)
		w := httptest.NewRecorder()
		h.CreateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check creating user (unmarshal error)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		body := bytes.NewReader([]byte(`{name: And, surname: Ersen, "position": 1, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPost, "/create", body)
		w := httptest.NewRecorder()
		h.CreateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check creating user (name error)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		body := bytes.NewReader([]byte(`{"name": "A1d", "surname": "Er^en", "position": 1, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPost, "/create", body)
		w := httptest.NewRecorder()
		h.CreateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check creating user (wrong position error)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		body := bytes.NewReader([]byte(`{"name": "And","surname": "Ersen", "position": 8, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPost, "/create", body)
		w := httptest.NewRecorder()
		h.CreateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
}

func TestHandlers_DeleteUser(t *testing.T) {
	t.Run("Check deleting user (no errors)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		mock.ExpectExec("DELETE FROM usr").WithArgs(id).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		expected := http.StatusOK
		req := httptest.NewRequest(http.MethodDelete, "/delete/5", nil)
		w := httptest.NewRecorder()
		h.DeleteUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check deleting user (empty index)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
		w := httptest.NewRecorder()
		h.DeleteUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check deleting user (wrong index type)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		req := httptest.NewRequest(http.MethodDelete, "/delete/abc", nil)
		w := httptest.NewRecorder()
		h.DeleteUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check deleting user (absent index in DB)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		mock.ExpectExec("DELETE FROM usr").WithArgs(id).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))
		expected := http.StatusInternalServerError
		req := httptest.NewRequest(http.MethodDelete, "/delete/5", nil)
		w := httptest.NewRecorder()
		h.DeleteUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
}

func TestHandlers_UpdateUser(t *testing.T) {
	t.Run("Check updating user (no errors)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		mock.ExpectExec("UPDATE usr SET").
			WithArgs("name='And9',surname='Ersen9',position='middle',project='Test9'", id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		expected := http.StatusOK
		body := bytes.NewReader([]byte(`{"name": "And9","surname": "Ersen9", "position": 3, "project": "Test9"}`))
		req := httptest.NewRequest(http.MethodPatch, "/update/5", body)
		w := httptest.NewRecorder()
		h.UpdateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check updating user (empty index)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		req := httptest.NewRequest(http.MethodPatch, "/update", nil)
		w := httptest.NewRecorder()
		h.UpdateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check updating user (wrong index type)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		req := httptest.NewRequest(http.MethodPatch, "/update/txt", nil)
		w := httptest.NewRecorder()
		h.UpdateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check updating user (unmarshal error)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		body := bytes.NewReader([]byte(`{name: And, surname: Ersen, "position": 1, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPatch, "/update/5", body)
		w := httptest.NewRecorder()
		h.UpdateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check updating user (name error)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		body := bytes.NewReader([]byte(`{"name": "A1d", "surname": "Er^en", "position": 1, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPatch, "/update/5", body)
		w := httptest.NewRecorder()
		h.UpdateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check updating user (absent index in DB)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		mock.ExpectExec("UPDATE usr SET").
			WithArgs("name='And9',surname='Ersen9',position='middle',project='Test9'", id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))
		expected := http.StatusInternalServerError
		body := bytes.NewReader([]byte(`{"name": "And9","surname": "Ersen9", "position": 3, "project": "Test9"}`))
		req := httptest.NewRequest(http.MethodPatch, "/update/5", body)
		w := httptest.NewRecorder()
		h.UpdateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
}

func TestHandlers_GetUser(t *testing.T) {
	t.Run("Check getting user (no errors)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		rows := pgxmock.NewRows([]string{"name", "surname", "position", "project"}).
			AddRow("And", "Ersen", "middle", "Test")
		mock.ExpectQuery("SELECT name, surname, position, project FROM").WithArgs(id).
			WillReturnRows(rows)
		expected := http.StatusOK
		expBody := bytes.NewReader([]byte(`{"name": "And","surname": "Ersen", "position": 3, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodGet, "/get/5", nil)
		w := httptest.NewRecorder()
		h.GetUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		gotBody := w.Result().Body
		assert.Equal(t, expBody, gotBody)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check getting user (empty index)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		req := httptest.NewRequest(http.MethodGet, "/get", nil)
		w := httptest.NewRecorder()
		h.GetUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check getting user (wrong index type)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expected := http.StatusBadRequest
		req := httptest.NewRequest(http.MethodGet, "/get/txt", nil)
		w := httptest.NewRecorder()
		h.GetUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
	})
	t.Run("Check getting user (absent index in DB)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expBody := "id error"
		mock.ExpectQuery("SELECT name, surname, position, project FROM").WithArgs(id).
			WillReturnError(fmt.Errorf(expBody))
		expected := http.StatusInternalServerError
		req := httptest.NewRequest(http.MethodGet, "/get/5", nil)
		w := httptest.NewRecorder()
		h.GetUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		gotBody := w.Result().Body
		assert.Equal(t, expBody, gotBody)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check getting user (marshalling error)", func(t *testing.T) {
		id := 5
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		rows := pgxmock.NewRows([]string{"name1", "surname1", "position1", "project1"}).
			AddRow("And", "Ersen", "middle", "Test")
		mock.ExpectQuery("SELECT name, surname, position, project FROM").WithArgs(id).
			WillReturnRows(rows)
		expected := http.StatusInternalServerError
		expBody := "marshalling error"
		req := httptest.NewRequest(http.MethodGet, "/get/5", nil)
		w := httptest.NewRecorder()
		h.GetUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		gotBody := w.Result().Body
		assert.Equal(t, expBody, gotBody)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
}

func TestHandlers_GetUserList(t *testing.T) {
	t.Run("Check getting user list (no errors)", func(t *testing.T) {
		type rec []string
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		Entries := []rec{
			{"And1", "Ersen1", "middle", "Test1"},
			{"And2", "Ersen2", "senior", "Test2"},
			{"And3", "Ersen3", "trainee", "Test3"},
		}
		expEntries := []byte(
			`[{"id": 1, "name": "And1", "surname": "Ersen1", "position": 3, "project": "Test1"},
			{"id": 2, "name": "And2", "surname": "Ersen2", "position": 4, "project": "Test2"},
			{"id": 3, "name": "And3", "surname": "Ersen3", "position": 1, "project": "Test3"}]`)
		rows := pgxmock.NewRows([]string{"id", "name", "surname", "position", "project"})
		for i, entry := range Entries {
			rows.AddRow(i+1, entry[0], entry[1], entry[2], entry[3])
		}
		mock.ExpectQuery("SELECT * FROM").WillReturnRows(rows)

		expected := http.StatusOK
		req := httptest.NewRequest(http.MethodGet, "/getall", nil)
		w := httptest.NewRecorder()
		h.GetUserList(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		gotEntries := w.Result().Body
		assert.ElementsMatch(t, expEntries, gotEntries)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check getting user list (internal DB error)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		expBody := "id error"
		mock.ExpectQuery("SELECT * FROM").WillReturnError(fmt.Errorf(expBody))
		expected := http.StatusInternalServerError
		req := httptest.NewRequest(http.MethodGet, "/getall", nil)
		w := httptest.NewRecorder()
		h.GetUserList(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		gotBody := w.Result().Body
		assert.Equal(t, expBody, gotBody)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
	t.Run("Check getting user list (unmarshalling error)", func(t *testing.T) {
		type rec []string
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening mock", err)
		}
		defer mock.Close()
		r := &Registry{mock}
		h := &Handlers{r}

		Entries := []rec{
			{"And1", "Ersen1", "middle", "Test1"},
			{"And2", "Ersen2", "senior", "Test2"},
			{"And3", "Ersen3", "trainee", "Test3"},
		}
		rows := pgxmock.NewRows([]string{"id", "3name", "surname", "position", "project"})
		for i, entry := range Entries {
			rows.AddRow(i+1, entry[0], entry[1], entry[2], entry[3])
		}
		mock.ExpectQuery("SELECT * FROM").WillReturnRows(rows)

		expected := http.StatusInternalServerError
		req := httptest.NewRequest(http.MethodGet, "/getall", nil)
		w := httptest.NewRecorder()
		h.GetUserList(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
		err = mock.ExpectationsWereMet()
		assert.NoErrorf(t, err, "there were unfulfilled expectations")
	})
}

/*
t.Run("Check server health (no errors)", func(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer s.Close()

	expected := "Hello john"
	req := httptest.NewRequest(http.MethodGet, "/greet?name=john", nil)
	w := httptest.NewRecorder()
	RequestHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if string(data) != expected {
		t.Errorf("Expected Hello john but got %v", string(data))
	}
}
*/
