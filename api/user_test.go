package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/kelvinator07/golang-bank-microservices/db/mock"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.ComparePasswords(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		// {
		// 	name: "OK",
		// 	body: gin.H{
		// 		"account_name": user.AccountName,
		// 		"password":     password,
		// 		"address":      user.Address,
		// 		"gender":       user.Gender,
		// 		"phone_number": user.PhoneNumber,
		// 		"email":        user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.CreateUserParams{
		// 			AccountName: user.AccountName,
		// 			Address:     user.Address,
		// 			Gender:      user.Gender,
		// 			PhoneNumber: user.PhoneNumber,
		// 			Email:       user.Email,
		// 		}
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
		// 			Times(1).
		// 			Return(user, nil)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		assert.Equal(t, http.StatusOK, recorder.Code)
		// 		requiredBodyMatchUser(t, recorder.Body, user)
		// 	},
		// },
		{
			name: "Internal Error",
			body: gin.H{
				"account_name": user.AccountName,
				"password":     password,
				"address":      user.Address,
				"gender":       user.Gender,
				"phone_number": user.PhoneNumber,
				"email":        user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"account_name": user.AccountName,
				"password":     password,
				"address":      user.Address,
				"gender":       user.Gender,
				"phone_number": user.PhoneNumber,
				"email":        "invalid@email@gmail.com",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Too Short Password",
			body: gin.H{
				"account_name": user.AccountName,
				"password":     "pass",
				"address":      user.Address,
				"gender":       user.Gender,
				"phone_number": user.PhoneNumber,
				"email":        user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			url := "/api/v1/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			assert.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)

	return db.User{
		AccountName:    util.RandomString(10),
		HashedPassword: hashedPassword,
		Address:        util.RandomString(20),
		Gender:         util.RandomGender(),
		PhoneNumber:    util.RandomPhoneNumber(),
		Email:          util.RandomEmail(),
	}, password
}

func requiredBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var expectedUser db.User
	err = json.Unmarshal(data, &expectedUser)

	fmt.Println("expectedUser ", expectedUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, expectedUser.Email)
	// assert.Equal(t, user, expectedUser)
}
