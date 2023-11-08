package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
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
	defer initLog.Info().Msg("exited")

	shutdownWG := &sync.WaitGroup{}

	queueURL := envutil.Must(initCtx, "QUEUE_URL")
	topic, err := InitializeQueueSink(initCtx, queueURL)
	if err != nil {
		initLog.Fatal().Err(err).Msg("could not initialize topic")
	}
	shutdownWG.Add(1)
	defer func() {
		defer shutdownWG.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		topic.Shutdown(ctx)
	}()

	collectionURL := envutil.Must(initCtx, "COLLECTION_URL")
	collection, err := InitializeCollection(initCtx, collectionURL)
	if err != nil {
		initLog.Fatal().Err(err).Msg("could not initialize document store")
	}
	shutdownWG.Add(1)
	defer func() {
		defer shutdownWG.Done()
		collection.Close()
	}()

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
	r.Get("/task/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		idParamValue := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParamValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ID must be a UUID"))
			log.Warn().Err(err).Str("id", idParamValue).Msg("client provided invalid id")
			return
		}

		// payload := lib.PayloadItem{
		// 	ID: id,
		// }
		// if err := collection.Get(ctx, &payload); err != nil {
		// 	if code := gcerrors.Code(err); code == gcerrors.NotFound {
		// 		w.WriteHeader(http.StatusNotFound)
		// 		w.Write([]byte("document was not found"))
		// 		log.Warn().Err(err).Msg("document was not found")
		// 		return
		// 	}
		// }

		payload := lib.PayloadItem{}
		iter := collection.Query().Where("ID", "=", id.String()).Limit(1).Get(ctx)
		defer iter.Stop()
		if err := iter.Next(ctx, &payload); err != nil {
			if err == io.EOF {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("document was not found"))
				log.Warn().Err(err).Msg("document was not found")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not get document from data store"))
			log.Panic().Err(err).Msg("could not get document from data store")
		}

		render.JSON{
			Data: payload,
		}.Render(w)
	})
	r.Post("/task", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		payload := lib.PayloadItem{
			ID:    uuid.New(),
			Time:  time.Now(),
			State: lib.Pending,
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			// This should never happen. If it does, something has gone wrong.
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
			log.Panic().Err(err).Msg("could not serialize json")
		}

		if err := collection.Create(ctx, &payload); err != nil {
			log.Error().Err(err).Msg("could not save payload to document storage")
			w.WriteHeader(http.StatusFailedDependency)
			render.JSON{
				Data: map[string]any{
					"error": err.Error(),
				},
			}.Render(w)
			return
		}
		if err := topic.Send(ctx, &pubsub.Message{
			Body: jsonBytes,
		}); err != nil {
			log.Error().Err(err).Msg("could not push payload to queue")
			w.WriteHeader(http.StatusFailedDependency)
			render.JSON{
				Data: map[string]any{
					"error": err.Error(),
				},
			}.Render(w)
			return
		}

		taskID := payload.ID
		response := map[string]any{
			"task_id": taskID,
		}

		log.Info().Any("task_id", taskID).Msg("responding to client for task")
		render.JSON{
			Data: response,
		}.Render(w)
	})
	initLog.Info().Msg("Starting")
	if err := http.ListenAndServe(":8080", r); errors.Is(err, http.ErrServerClosed) {
		initLog.Info().Msg("Gracefullly shutdown")
	} else if err != nil {
		initLog.Fatal().Err(err).Msg("an error occurred and we abruptly shutdown")
	}
	shutdownWG.Wait()
}
