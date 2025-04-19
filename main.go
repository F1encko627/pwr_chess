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
		}).ParseGlob("web/*.templ")),
	}
}

func main() {
	fmt.Println(game)
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
	game.MovePiece(ix, iy, fx, fy)

	return c.Render(http.StatusOK, "board.html.templ", game)
}

func Restart(c echo.Context) error {
	game = board.NewGame()

	return c.Render(http.StatusOK, "board.html.templ", game)
}