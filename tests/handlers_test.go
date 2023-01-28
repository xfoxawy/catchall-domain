package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mailgun/catchall"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/fx"

	handlers "github.com/xfoxawy/catchall-domain/app/interface/http"
	"github.com/xfoxawy/catchall-domain/app/repository"
)

func TestHttpHandlers(t *testing.T) {
	router := gin.Default()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("testing events", func(mt *mtest.T) {
		args := handlers.EventsApiArgs{
			fx.In{},
			router,
			handlers.NewHttpHandlers(
				repository.NewEventsCounterRepository(mt.Client.Database("test")),
			),
		}
		handlers.RegisterEventsApi(args)

		t.Run("test put bounced", func(t *testing.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			ev := catchall.Event{
				Type:   catchall.TypeBounced,
				Domain: "example.com",
			}

			var buff bytes.Buffer
			err := json.NewEncoder(&buff).Encode(ev)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/events/example.com/bounced", &buff)
			if err != nil {
				t.Fatal(err)
			}

			router.ServeHTTP(w, req)
			o, err := json.Marshal(ev)
			if err != nil {
				t.Fatal(err)
			}

			require.JSONEq(t, string(o), w.Body.String())
			assert.Equal(t, http.StatusAccepted, w.Code)
		})

		t.Run("test put delivered", func(t *testing.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			ev := catchall.Event{
				Type:   catchall.TypeDelivered,
				Domain: "example.com",
			}

			var buff bytes.Buffer
			err := json.NewEncoder(&buff).Encode(ev)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/events/example.com/delivered", &buff)
			if err != nil {
				t.Fatal(err)
			}
			router.ServeHTTP(w, req)
			o, err := json.Marshal(ev)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusAccepted, w.Code)
			require.JSONEq(t, string(o), w.Body.String())
		})

		t.Run("test get not catchall status", func(t *testing.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", primitive.NewObjectID()},
				{"delivered", 1},
				{"bounced", 1},
				{"domain", "example.com"},
			}))

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/domains/example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			require.JSONEq(t, `{"status": "not catch-all"}`, w.Body.String())
		})

		t.Run("test get unknown status", func(t *testing.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", primitive.NewObjectID()},
				{"delivered", 888},
				{"domain", "example.com"},
			}))

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/domains/example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			require.JSONEq(t, `{"status": "unknown"}`, w.Body.String())
		})

		t.Run("test get catch-all status", func(t *testing.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", primitive.NewObjectID()},
				{"delivered", 1001},
				{"domain", "example.com"},
			}))

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/domains/example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			require.JSONEq(t, `{"status": "catch-all"}`, w.Body.String())
		})
	})
}
