package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"ust_chess/internal/board"
	"ust_chess/internal/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	game = board.NewGame([]types.Piece{})
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.New("board.html.tmpl").Funcs(template.FuncMap{
			"even": func(x, y int) bool {
				return (x+y)%2 == 0
			},
			"string": func(x fmt.Stringer) string {
				return x.String()
			},
			"name": func(x types.Figure) string {
				return x.Name()
			},
		}).ParseGlob("./web/*.tmpl")),
	}
}

// TODO:
// Вынести логику серера в отдельный пакет
// Мультипреер по комнатам. Луше сразу по комнатам. Не уверен что получится промежуточно сделать просто два игрока.
// Может быть по приколу отказаться от Echo и сделать самописный сервер на базовом http/net. Вроде неплохое обучение.
// Потом может быть систему аккаунтов примитивную. Авторизацию через битрикс лол.
func main() {

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).With().Timestamp().Logger()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()
	e.GET("/", Hello)
	e.GET("/move", Move)
	e.GET("/restart", Restart)
	e.Logger.Fatal(e.Start(":1337"))
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "board.html.tmpl", game)
}

func Move(c echo.Context) error {
	ix, _ := strconv.Atoi(c.QueryParam("ix"))
	iy, _ := strconv.Atoi(c.QueryParam("iy"))
	fx, _ := strconv.Atoi(c.QueryParam("fx"))
	fy, _ := strconv.Atoi(c.QueryParam("fy"))

	move, err := types.GetMove(ix, iy, fx, fy)
	if err != nil {
		c.Logger().Error(err)
	}

	game.MakeMove(move)

	if err != nil {
		c.Logger().Error(err)
	}

	return c.Render(http.StatusOK, "board.html.tmpl", game.GetForRender())
}

// var test_case = []types.Piece{
// 	types.GP(types.QUEEN, false, types.NewPos(3, 1)),
// 	types.GP(types.QUEEN, false, types.NewPos(1, 2)),
// 	types.GP(types.QUEEN, false, types.NewPos(5, 2)),
// 	types.GP(types.QUEEN, false, types.NewPos(5, 4)),
// 	types.GP(types.QUEEN, false, types.NewPos(5, 6)),
// 	types.GP(types.QUEEN, false, types.NewPos(3, 6)),
// 	types.GP(types.QUEEN, false, types.NewPos(0, 4)),
// 	types.GP(types.QUEEN, false, types.NewPos(1, 6)),
// 	//types.GP(types.PAWN, false, types.NewPos(4, 2)),
// 	//types.GP(types.PAWN, false, types.NewPos(5, 3)),

// 	types.GP(types.KING, true, types.NewPos(3, 4)),
// }

func Restart(c echo.Context) error {
	game = board.NewGame([]types.Piece{})
	c.Logger().Warn("game restated")

	return c.Render(http.StatusOK, "board.html.tmpl", game)
}
