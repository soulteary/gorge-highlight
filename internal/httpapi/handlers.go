package httpapi

import (
	"net/http"

	"github.com/soulteary/gorge-highlight/internal/highlight"

	"github.com/labstack/echo/v4"
)

type Deps struct {
	Highlighter *highlight.Highlighter
	Token       string
	MaxBytes    int
}

type apiResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *apiError `json:"error,omitempty"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RegisterRoutes(e *echo.Echo, deps *Deps) {
	e.GET("/", healthPing())
	e.GET("/healthz", healthPing())

	g := e.Group("/api/highlight")
	g.Use(tokenAuth(deps))

	g.POST("/render", renderHighlight(deps))
	g.GET("/languages", listLanguages(deps))
}

func tokenAuth(deps *Deps) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if deps.Token == "" {
				return next(c)
			}
			token := c.Request().Header.Get("X-Service-Token")
			if token == "" {
				token = c.QueryParam("token")
			}
			if token == "" || token != deps.Token {
				return c.JSON(http.StatusUnauthorized, &apiResponse{
					Error: &apiError{Code: "ERR_UNAUTHORIZED", Message: "missing or invalid service token"},
				})
			}
			return next(c)
		}
	}
}

func healthPing() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func respondOK(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, &apiResponse{Data: data})
}

func respondErr(c echo.Context, status int, code, msg string) error {
	return c.JSON(status, &apiResponse{
		Error: &apiError{Code: code, Message: msg},
	})
}

type highlightRequest struct {
	Source   string `json:"source"`
	Language string `json:"language"`
}

func renderHighlight(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req highlightRequest
		if err := c.Bind(&req); err != nil {
			return respondErr(c, http.StatusBadRequest, "ERR_BAD_REQUEST", err.Error())
		}

		if req.Source == "" {
			return respondOK(c, &highlight.Result{HTML: "", Language: req.Language})
		}

		if deps.MaxBytes > 0 && len(req.Source) > deps.MaxBytes {
			return respondErr(c, http.StatusRequestEntityTooLarge,
				"ERR_TOO_LARGE", "source exceeds maximum allowed size")
		}

		result, err := deps.Highlighter.Highlight(req.Source, req.Language)
		if err != nil {
			return respondErr(c, http.StatusInternalServerError,
				"ERR_HIGHLIGHT_FAILED", err.Error())
		}

		return respondOK(c, result)
	}
}

func listLanguages(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		return respondOK(c, deps.Highlighter.Languages())
	}
}
