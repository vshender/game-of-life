package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const CELL_SIZE = 5

type Cell struct {
	X, Y int
}

func (cell Cell) ForEachNeighbor(f func(cell Cell)) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i != 0 || j != 0 {
				f(Cell{cell.X + i, cell.Y + j})
			}
		}
	}
}

type LifeState map[Cell]bool

func (ls *LifeState) String() string {
	var buf strings.Builder
	buf.WriteString("[")
	for cell, alive := range *ls {
		if alive {
			buf.WriteString(fmt.Sprintf("(%d, %d), ", cell.X, cell.Y))
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (ls *LifeState) Size() int {
	return len(*ls)
}

func (ls *LifeState) Set(cell Cell, alive bool) {
	(*ls)[cell] = alive
}

func (ls *LifeState) Alive(cell Cell) bool {
	return (*ls)[cell]
}

func (ls *LifeState) Next() *LifeState {
	nextState := make(map[Cell]bool)

	calcNextCellState := func(cell Cell) {
		if _, found := nextState[cell]; found {
			// already calculated
			return
		}

		alive := 0
		cell.ForEachNeighbor(func(cell Cell) {
			if ls.Alive(cell) {
				alive += 1
			}
		})

		nextState[cell] = alive == 3 || alive == 2 && ls.Alive(cell)
	}

	for cell, alive := range *ls {
		if alive {
			calcNextCellState(cell)
			cell.ForEachNeighbor(calcNextCellState)
		}
	}

	return (*LifeState)(&nextState)
}

func (ls *LifeState) Draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.Black

	for cell, alive := range *ls {
		if alive {
			imd.Push(pixel.V(float64(cell.X*CELL_SIZE), float64(cell.Y*CELL_SIZE)))
			imd.Push(pixel.V(float64((cell.X+1)*CELL_SIZE-1), float64((cell.Y+1)*CELL_SIZE-1)))
			imd.Rectangle(0)
		}
	}
}

func NewLifeState(w, h int) *LifeState {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	state := make(map[Cell]bool)

	for i := 0; i < w*h/8; i++ {
		state[Cell{r.Intn(w), r.Intn(h)}] = true
	}

	return (*LifeState)(&state)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "The Game of Life",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	ls := NewLifeState(1024/CELL_SIZE, 768/CELL_SIZE)

	imd := imdraw.New(nil)
	ls.Draw(imd)

	ticker := time.Tick(time.Second)
	iter := 1
	fmt.Printf("Iter #%d: %d live cells\n", iter, ls.Size())

	for !win.Closed() {
		win.Clear(colornames.Darkgray)
		imd.Draw(win)
		win.Update()

		select {
		case <-ticker:
			ls = ls.Next()
			iter++
			fmt.Printf("Iter #%d: %d live cells\n", iter, ls.Size())

			imd.Clear()
			ls.Draw(imd)

		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
