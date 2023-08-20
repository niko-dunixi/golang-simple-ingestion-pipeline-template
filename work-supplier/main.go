package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib"
	"github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib/envutil"
	_ "github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib/zerologutil"
	"github.com/rs/zerolog"
	"gocloud.dev/pubsub"
)

func main() {
	initCtx, initCtxCancel := context.WithTimeout(context.Background(), time.Second*15)
	initCtx = zerolog.Ctx(initCtx).With().Str("scope", "initialization").Logger().WithContext(initCtx)
	defer initCtxCancel()
	initLog := zerolog.Ctx(initCtx)
	queueURL := envutil.Must(initCtx, "QUEUE_URL")
	topic, err := InitializeQueueSink(initCtx, queueURL)
	if err != nil {
		initLog.Fatal().Err(err).Msg("could not initialize topic")
	}

	r := chi.NewRouter()

	zerologMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := zerolog.Ctx(r.Context()).WithContext(r.Context())
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
	r.Use(zerologMiddleware, middleware.RealIP)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON{
			Data: map[string]any{
				"your_ip": r.RemoteAddr,
			},
		}.Render(w)
	})
	r.Post("/task/{name}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		taskName := chi.URLParam(r, "name")
		payload := lib.PayloadItem{
			ID:       uuid.New(),
			Time:     time.Now(),
			TaskName: taskName,
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			// This should never happen. If it does, something has gone wrong.
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
			log.Panic().Err(err).Msg("could not serialize json")
			return
		}

		err = topic.Send(ctx, &pubsub.Message{
			Body: jsonBytes,
		})
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			render.JSON{
				Data: map[string]any{
					"error": err.Error(),
				},
			}.Render(w)
			return
		}

		response := map[string]any{
			"task_name": taskName,
		}

		log.Info().Any("task_name", taskName).Msg("responding to client for task")
		render.JSON{
			Data: response,
		}.Render(w)
	})
	initLog.Info().Msg("Starting")
	if err := http.ListenAndServe(":3000", r); errors.Is(err, http.ErrServerClosed) {
		initLog.Info().Msg("Gracefullly shutdown")
	} else if err != nil {
		initLog.Fatal().Err(err).Msg("an error occurred and we abruptly shutdown")
	}
}
