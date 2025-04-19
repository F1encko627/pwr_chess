package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"ust_chess/board"
	"ust_chess/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	game = board.NewGame()
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.New("board.html.templ").Funcs(template.FuncMap{
			"even": func(x, y int) bool {
				return (x+y)%2 == 0
			},
			"string": func(x types.Type) string {
				return string(x)
			},
			"state" : func(m types.State) string {
				switch m {
				case types.WHITE_TURN:
					return "Ходят белые"
				case types.BLACK_TURN:
					return "Ходят черные"
				case types.WHITE_CHECK:
					return "Шах белому королю"
				case types.BLACK_CHECK:
					return "Шах черному королю"
				case types.WHITE_CHECKMATE:
					return "Мат белому королю"
				case types.BLACK_CHECKMATE:
					return "Мат черному королю"
				case types.STALEMATE:
					return "Пат"
				default:
					return "Unknown state"
				}
			},
		}).ParseGlob("web/*.templ")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()
	e.GET("/", Hello)
	e.GET("/move", Move)
	e.GET("/restart", Restart)
	e.Logger.Fatal(e.Start(":1337"))
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "board.html.templ", game)
}

func Move(c echo.Context) error {
	ix, _ := strconv.Atoi(c.QueryParam("ix"))
	iy, _ := strconv.Atoi(c.QueryParam("iy"))
	fx, _ := strconv.Atoi(c.QueryParam("fx"))
	fy, _ := strconv.Atoi(c.QueryParam("fy"))
	fmt.Println(ix, iy, fx, fy)
	err := game.MovePiece(ix, iy, fx, fy)

	if err != nil {
		c.Logger().Error(err)
	}

	return c.Render(http.StatusOK, "board.html.templ", game)
}

func Restart(c echo.Context) error {
	game = board.NewGame()
	c.Logger().Warn("game restated")

	return c.Render(http.StatusOK, "board.html.templ", game)
}