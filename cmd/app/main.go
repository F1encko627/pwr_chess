package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"ust_chess/internal/board"
	"ust_chess/internal/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		templates: template.Must(template.New("board.html.templ").Funcs(template.FuncMap{
			"even": func(x, y int) bool {
				return (x+y)%2 == 0
			},
			"string": func(x types.Type) string {
				return string(x)
			},
			"state": func(m types.State) string {
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
		}).ParseGlob("./web/*.templ")),
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

var test_case = []types.Piece{
	types.GP(types.PAWN, false, types.NewPos(3, 1)),
	types.GP(types.PAWN, false, types.NewPos(1, 2)),
	types.GP(types.PAWN, false, types.NewPos(5, 2)),
	types.GP(types.PAWN, false, types.NewPos(5, 4)),
	types.GP(types.PAWN, false, types.NewPos(5, 6)),
	types.GP(types.PAWN, false, types.NewPos(3, 6)),
	types.GP(types.PAWN, false, types.NewPos(0, 4)),
	types.GP(types.PAWN, false, types.NewPos(1, 6)),
	//types.GP(types.PAWN, false, types.NewPos(4, 2)),
	//types.GP(types.PAWN, false, types.NewPos(5, 3)),

	types.GP(types.QUEEN, true, types.NewPos(3, 4)),
}

func Restart(c echo.Context) error {
	game = board.NewGame(test_case)
	c.Logger().Warn("game restated")

	return c.Render(http.StatusOK, "board.html.templ", game)
}
