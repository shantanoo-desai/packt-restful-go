package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/shantanoo-desai/packt-restful-go/user"
	"gopkg.in/mgo.v2/bson"
)

func TestBody(t *testing.T) {
	validUser := &user.User {
		ID: bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}
	validUser2 := &user.User {
		ID: validUser.ID,
		Name: "John",
		Role: "developer",
	}
	js, err := json.Marshal(validUser)
	if err != nil {
		t.Errorf("Error marshalling a valid user: %s", validUser)
		t.FailNow()
	}


	testCases := []struct{
		txt string
		r *http.Request
		u *user.User
		err bool
		exp *user.User
	}{
		{ txt: "nil Request", err: true },
		{ txt: "empty user", r: &http.Request { Body: io.NopCloser(bytes.NewBufferString("{}")) }, err: true },
		{ txt: "malformed data", r: &http.Request { Body: io.NopCloser(bytes.NewBufferString(`{"id":12}`)) }, u: &user.User{}, err: true},
		{ txt: "valid request", r: &http.Request { Body: io.NopCloser(bytes.NewBuffer(js)) }, u: &user.User{}, exp: validUser },
		{ txt: "valid partial request", r: &http.Request { Body: io.NopCloser(bytes.NewBufferString(`{"role":"developer", "age":22}`)) }, u: validUser, exp: validUser2 },
	}

	for _, tc := range testCases {
		t.Log(tc.txt)
		err := bodyToUser(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("expected error, got none")
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}

		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different:")
			t.Error(tc.u)
			t.Error(tc.exp)
		}
	}
}
