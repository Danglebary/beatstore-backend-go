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

func randomBeat() db.Beat {
	return db.Beat{
		ID:        int32(util.RandomInt(1, 1000)),
		CreatorID: int32(util.RandomInt(1, 1000)),
		Title:     util.RandomTitle(),
		Genre:     util.RandomGenre(),
		Key:       util.RandomKey(),
		Bpm:       util.RandomBpm(),
		Tags:      util.RandomTags(),
		S3Key:     "not implemented",
	}
}

func randomBeats(n int) []db.Beat {
	var beats []db.Beat

	for i := 0; i < n; i++ {
		beats = append(beats, randomBeat())
	}
	return beats
}

func requireBodyMatchBeat(t *testing.T, body *bytes.Buffer, beat db.Beat) {
	// Read bytes buffer
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// Unmarshal byte data to db.Beat struct
	var gotBeat db.Beat
	err = json.Unmarshal(data, &gotBeat)
	require.NoError(t, err)

	require.Equal(t, beat, gotBeat)
}

func requireBodyMatchBeats(t *testing.T, body *bytes.Buffer, beats []db.Beat) {
	// Read bytes buffer
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// Unmarshal byte data to []db.Beat
	var gotBeats []db.Beat
	err = json.Unmarshal(data, &gotBeats)
	require.NoError(t, err)

	for i := range beats {
		require.Equal(t, beats[i], gotBeats[i])
	}
}

func TestCreateBeat(t *testing.T) {
	beat := randomBeat()

	testCases := []struct {
		name          string
		body          createBeatRequestParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createBeatRequestParams{
				CreatorID: beat.CreatorID,
				Title:     beat.Title,
				Genre:     beat.Genre,
				Key:       beat.Key,
				Bpm:       beat.Bpm,
				Tags:      beat.Tags,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateBeatParams{
					CreatorID: beat.CreatorID,
					Title:     beat.Title,
					Genre:     beat.Genre,
					Key:       beat.Key,
					Bpm:       beat.Bpm,
					Tags:      beat.Tags,
					S3Key:     beat.S3Key,
				}

				store.EXPECT().
					CreateBeat(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(beat, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeat(t, recorder.Body, beat)
			},
		},
		{
			name: "BadRequest-body",
			body: createBeatRequestParams{
				CreatorID: beat.CreatorID,
				Title:     beat.Title,
				Genre:     beat.Genre,
				Key:       beat.Key,
				Bpm:       beat.Bpm,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBeat(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: createBeatRequestParams{
				CreatorID: beat.CreatorID,
				Title:     beat.Title,
				Genre:     beat.Genre,
				Key:       beat.Key,
				Bpm:       beat.Bpm,
				Tags:      beat.Tags,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateBeatParams{
					CreatorID: beat.CreatorID,
					Title:     beat.Title,
					Genre:     beat.Genre,
					Key:       beat.Key,
					Bpm:       beat.Bpm,
					Tags:      beat.Tags,
					S3Key:     beat.S3Key,
				}

				store.EXPECT().
					CreateBeat(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Beat{}, sql.ErrConnDone)
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
			url := "/beats"
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

func TestUpdateBeat(t *testing.T) {
	beat := randomBeat()

	testCases := []struct {
		name          string
		beatID        int32
		body          updateBeatRequestParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			beatID: beat.ID,
			body: updateBeatRequestParams{
				Title: beat.Title,
				Genre: beat.Genre,
				Key:   beat.Key,
				Bpm:   beat.Bpm,
				Tags:  beat.Tags,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateBeatParams{
					ID:    beat.ID,
					Title: beat.Title,
					Genre: beat.Genre,
					Key:   beat.Key,
					Bpm:   beat.Bpm,
					Tags:  beat.Tags,
					S3Key: beat.S3Key,
				}

				store.EXPECT().
					UpdateBeat(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(beat, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeat(t, recorder.Body, beat)
			},
		},
		{
			name:   "BadRequest-Uri",
			beatID: 0,
			body: updateBeatRequestParams{
				Title: beat.Title,
				Genre: beat.Genre,
				Key:   beat.Key,
				Bpm:   beat.Bpm,
				Tags:  beat.Tags,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequest-Body",
			beatID: beat.ID,
			body: updateBeatRequestParams{
				Title: beat.Title,
				Genre: beat.Genre,
				Key:   beat.Key,
				Bpm:   beat.Bpm,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBeat(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			beatID: beat.ID,
			body: updateBeatRequestParams{
				Title: beat.Title,
				Genre: beat.Genre,
				Key:   beat.Key,
				Bpm:   beat.Bpm,
				Tags:  beat.Tags,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBeat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Beat{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/beats/%d", tc.beatID)
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

func TestGetBeatByIDApi(t *testing.T) {
	beat := randomBeat()

	testCases := []struct {
		name          string
		beatId        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			beatId: beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBeatById(gomock.Any(), gomock.Eq(beat.ID)).
					Times(1).
					Return(beat, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeat(t, recorder.Body, beat)
			},
		},
		{
			name:   "NotFound",
			beatId: beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBeatById(gomock.Any(), gomock.Eq(beat.ID)).
					Times(1).
					Return(db.Beat{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			beatId: beat.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBeatById(gomock.Any(), gomock.Eq(beat.ID)).
					Times(1).
					Return(db.Beat{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "IntvalidID",
			beatId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBeatById(gomock.Any(), gomock.Any()).
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
			url := fmt.Sprintf("/beats/%d", tc.beatId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByIDorderIDApi(t *testing.T) {
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:     "InternalError",
			pageID:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsById(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "BadRequest",
			pageID:   0,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "BadPageID",
			pageID:   0,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "BadPageSize",
			pageID:   1,
			pageSize: 100,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsById(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/beats?page_id=%d&page_size=%d&order=ID", tc.pageID, tc.pageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByIDorderBPMApi(t *testing.T) {
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		bpmMin        int16
		bpmMax        int16
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   1,
			pageSize: 5,
			bpmMin:   60,
			bpmMax:   240,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByBpmRange(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:     "BadBpmMin",
			pageID:   1,
			pageSize: 5,
			bpmMin:   5,
			bpmMax:   100,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByBpmRange(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "BadBpmMax",
			pageID:   1,
			pageSize: 5,
			bpmMin:   60,
			bpmMax:   1000,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByBpmRange(gomock.Any(), gomock.Any()).
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
			bpmMin:   60,
			bpmMax:   100,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByBpmRange(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/beats?page_id=%d&page_size=%d&order=BPM&min=%d&max=%d", tc.pageID, tc.pageSize, tc.bpmMin, tc.bpmMax)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByIDorderKEYApi(t *testing.T) {
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		key           string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   1,
			pageSize: 5,
			key:      "C_SHARP_MINOR",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:     "BadKey",
			pageID:   1,
			pageSize: 5,
			key:      "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByKey(gomock.Any(), gomock.Any()).
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
			key:      "C_SHARP_MINOR",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/beats?page_id=%d&page_size=%d&order=KEY&key=%s", tc.pageID, tc.pageSize, tc.key)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByIDorderGENREApi(t *testing.T) {
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		genre         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   1,
			pageSize: 5,
			genre:    "TRAP",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByGenre(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:     "BadGenre",
			pageID:   1,
			pageSize: 5,
			genre:    "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByGenre(gomock.Any(), gomock.Any()).
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
			genre:    "TRAP",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByGenre(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/beats?page_id=%d&page_size=%d&order=GENRE&genre=%s", tc.pageID, tc.pageSize, tc.genre)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByCreatorIDorderIDApi(t *testing.T) {
	user := randomUser()
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		creatorID     int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorId(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:      "InternalError",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorId(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			pageID:    0,
			pageSize:  5,
			creatorID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorId(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "BadPageID",
			pageID:    0,
			pageSize:  5,
			creatorID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorId(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "BadPageSize",
			pageID:    1,
			pageSize:  100,
			creatorID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorId(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/users/%d/beats?page_id=%d&page_size=%d&order=ID", tc.creatorID, tc.pageID, tc.pageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByCreatorIDorderBPMApi(t *testing.T) {
	user := randomUser()
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		creatorID     int32
		bpmMin        int16
		bpmMax        int16
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			bpmMin:    60,
			bpmMax:    240,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndBpmRange(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:      "BadBpmMin",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			bpmMin:    5,
			bpmMax:    100,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndBpmRange(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "BadBpmMax",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			bpmMin:    60,
			bpmMax:    1000,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndBpmRange(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			bpmMin:    60,
			bpmMax:    100,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndBpmRange(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/users/%d/beats?page_id=%d&page_size=%d&order=BPM&min=%d&max=%d", tc.creatorID, tc.pageID, tc.pageSize, tc.bpmMin, tc.bpmMax)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByCreatorIDorderKEYApi(t *testing.T) {
	user := randomUser()
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		creatorID     int32
		key           string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			key:       "C_SHARP_MINOR",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:      "BadKey",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			key:       "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndKey(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			key:       "C_SHARP_MINOR",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/users/%d/beats?page_id=%d&page_size=%d&order=KEY&key=%s", tc.creatorID, tc.pageID, tc.pageSize, tc.key)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListBeatsByCreatorIDorderGENREApi(t *testing.T) {
	user := randomUser()
	beats := randomBeats(5)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		creatorID     int32
		genre         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			genre:     "TRAP",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndGenre(gomock.Any(), gomock.Any()).
					Times(1).
					Return(beats, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBeats(t, recorder.Body, beats)
			},
		},
		{
			name:      "BadGenre",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			genre:     "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndGenre(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			pageID:    1,
			pageSize:  5,
			creatorID: user.ID,
			genre:     "TRAP",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBeatsByCreatorIdAndGenre(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Beat{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/users/%d/beats?page_id=%d&page_size=%d&order=GENRE&genre=%s", tc.creatorID, tc.pageID, tc.pageSize, tc.genre)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// Server http response
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}
