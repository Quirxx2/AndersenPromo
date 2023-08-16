package promo

import (
	"bytes"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

		expected := http.StatusOK
		body := bytes.NewReader([]byte(`{"name": "And","surname": "Ersen", "position": 1, "project": "Test"}`))
		req := httptest.NewRequest(http.MethodPost, "/create", body)
		w := httptest.NewRecorder()
		h.CreateUser(w, req)
		got := w.Result().StatusCode
		assert.Equal(t, expected, got)
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

		expected := http.StatusInternalServerError
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

		mock.ExpectExec("UPDATE usr").WithArgs("And", "Ersen", 3, "Test", id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		expected := http.StatusOK
		body := bytes.NewReader([]byte(`{"name": "And","surname": "Ersen", "position": 3, "project": "Test"}`))
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

		mock.ExpectExec("UPDATE usr").WithArgs("And", "Ersen", 3, "Test", id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))
		expected := http.StatusInternalServerError
		body := bytes.NewReader([]byte(`{"name": "And","surname": "Ersen", "position": 3, "project": "Test"}`))
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
