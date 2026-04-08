package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RishatShay/url-shortener/internal/http-server/lib/api/responce"
	"github.com/RishatShay/url-shortener/internal/http-server/lib/random"
	"github.com/RishatShay/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Responce struct {
	responce.Responce
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

const (
	aliasLength = 6
)

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, responce.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, responce.Error("invalid request"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.GetAlias(aliasLength)
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, responce.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", slog.String("error", err.Error()))

			render.JSON(w, r, responce.Error("failed to add url"))

			return
		}

		log.Info("url successsfully added", slog.String("alias", alias))

		render.JSON(w, r, Responce{
			Responce: responce.OK(),
			Alias:    alias,
		})

	}
}
