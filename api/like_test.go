package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/danglebary/beatstore-backend-go/db/mock"
	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/danglebary/beatstore-backend-go/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomLike(userID int32, beatID int32) db.Like {
	return db.Like{
		UserID: userID,
		BeatID: beatID,
	}
}

func randomIDs(n int) []int32 {
	var IDs []int32
	for i := 0; i < n; i++ {
		IDs = append(IDs, int32(util.RandomInt(1, 1000)))
	}
	return IDs
}

func randomLikesForUser(userID int32, beatIDs []int32, n int) []db.Like {
	var likes []db.Like
	for i := 0; i < n; i++ {
		likes = append(likes, randomLike(userID, beatIDs[i]))
	}
	return likes
}

func randomLikesForBeat(beatID int32, userIDs []int32, n int) []db.Like {
	var likes []db.Like
	for i := 0; i < n; i++ {
		likes = append(likes, randomLike(userIDs[i], beatID))
	}
	return likes
}

func requireBodyMatchLike(t *testing.T, body *bytes.Buffer, like db.Like) {
	// Read bytes buffer
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// Unmarshal byte data to db.like struct
	var gotLike db.Like
	err = json.Unmarshal(data, &gotLike)
	require.NoError(t, err)

	require.Equal(t, like, gotLike)
}

func requireBodyMatchLikes(t *testing.T, body *bytes.Buffer, likes []db.Like) {
	// Read bytes buffer
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// Unmarshal byte data to []db.like
	var gotLikes []db.Like
	err = json.Unmarshal(data, &gotLikes)
	require.NoError(t, err)

	for i := range likes {
		require.Equal(t, likes[i], gotLikes[i])
	}
}

func TestCreateLike(t *testing.T) {
	user := randomUser()
	beat := randomBeat()
	like := randomLike(user.ID, beat.ID)

	testCases := []struct {
		name          string
		body          createLikeRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createLikeRequest{
				UserID: user.ID,
				BeatID: beat.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateLikeParams{
					UserID: user.ID,
					BeatID: beat.ID,
				}

				store.EXPECT().
					CreateLike(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(like, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLike(t, recorder.Body, like)
			},
		},
		{
			name: "BadRequest",
			body: createLikeRequest{
				UserID: 0,
				BeatID: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateLike(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: createLikeRequest{
				UserID: user.ID,
				BeatID: beat.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateLike(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Like{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Init controller and store
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			// Build stub
			tc.buildStubs(store)
			// Start test server, build request, and send
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := "/likes"
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetLike(t *testing.T) {
	user := randomUser()
	beat := randomBeat()
	like := randomLike(user.ID, beat.ID)

	testCases := []struct {
		name          string
		userID        int32
		beatID        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			beatID: beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetLikeByUserAndBeatParams{
					UserID: user.ID,
					BeatID: beat.ID,
				}

				store.EXPECT().
					GetLikeByUserAndBeat(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(like, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLike(t, recorder.Body, like)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			beatID: beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLikeByUserAndBeat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Like{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			beatID: beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLikeByUserAndBeat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Like{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			userID: 0,
			beatID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLikeByUserAndBeat(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Init controller and store
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			// Build stub
			tc.buildStubs(store)
			// Start test server, build request, and send
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/likes/%d/%d", tc.userID, tc.beatID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListLikesByBeatID(t *testing.T) {
	n := 5
	beat := randomBeat()
	userIDs := randomIDs(n)
	likes := randomLikesForBeat(beat.ID, userIDs, n)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		beatID        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   1,
			pageSize: 5,
			beatID:   beat.ID,
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListLikesByBeatParams{
					BeatID: beat.ID,
					Limit:  5,
					Offset: 0,
				}

				store.EXPECT().
					ListLikesByBeat(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(likes, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLikes(t, recorder.Body, likes)
			},
		},
		{
			name:     "BadRequest-Uri",
			pageID:   1,
			pageSize: 5,
			beatID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLikesByBeat(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "BadRequest-QueryString",
			pageID:   0,
			pageSize: 1000,
			beatID:   beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLikesByBeat(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			pageID:   1,
			pageSize: 5,
			beatID:   beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListLikesByBeatParams{
					BeatID: beat.ID,
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListLikesByBeat(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Like{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Init controller and store
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			// Build stub
			tc.buildStubs(store)
			// Start test server, build request, and send
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/beats/%d/likes?page_id=%d&page_size=%d", tc.beatID, tc.pageID, tc.pageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListLikesByCreatorID(t *testing.T) {
	n := 5
	user := randomUser()
	beatIDs := randomIDs(n)
	likes := randomLikesForUser(user.ID, beatIDs, n)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		userID        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   1,
			pageSize: 5,
			userID:   user.ID,
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListLikesByUserParams{
					UserID: user.ID,
					Limit:  5,
					Offset: 0,
				}

				store.EXPECT().
					ListLikesByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(likes, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLikes(t, recorder.Body, likes)
			},
		},
		{
			name:     "BadRequest-Uri",
			pageID:   1,
			pageSize: 5,
			userID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLikesByUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "BadRequest-QueryString",
			pageID:   0,
			pageSize: 1000,
			userID:   user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLikesByUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			pageID:   1,
			pageSize: 5,
			userID:   user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListLikesByUserParams{
					UserID: user.ID,
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListLikesByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Like{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Init controller and store
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			// Build stub
			tc.buildStubs(store)
			// Start test server, build request, and send
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d/likes?page_id=%d&page_size=%d", tc.userID, tc.pageID, tc.pageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}
