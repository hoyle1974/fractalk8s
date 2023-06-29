package main

import (
	"log"
	"math"
	"math/big"
	"os"

	"github.com/alitto/pond"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = byte(255)
	chunk        = 16
	maxPoolSize  = (screenWidth/chunk)*(screenHeight/chunk) + 1
	maxWorkers   = 64
)

var (
	palette [maxIt]byte
)

// func calc2(client *http.Client, x, y []*big.Float, iter int) []int {
// 	ret := make([]int, len(x))
// 	for idx, _ := range x {
// 		ret[idx] = mandelbrot.MandelbrotFloat(x[idx], y[idx], iter)

// 		// cx, _ := x[idx].Float64()
// 		// cy, _ := y[idx].Float64()
// 		// ret[idx] = mandelbrot.MandelbrotFast(cx, cy, iter)
// 	}

// 	return ret
// }

func init() {
	for i := range palette {
		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	}
}

func color(it byte) (r, g, b byte) {
	it = maxIt - it

	if it == maxIt {
		return 0xff, 0xff, 0xff
	}
	c := palette[it]
	return (c * 3) % 255, ((c + 25) * 5) % 255, ((c + 128) * 7) % 255
}

type Game struct {
	cam          Camera
	offscreen    *ebiten.Image
	offscreenPix []byte
	keys         []ebiten.Key
	timer        *Metrics
}

func NewGame(cam Camera) *Game {
	g := &Game{
		cam:          cam,
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		offscreenPix: make([]byte, screenWidth*screenHeight*4),
		timer:        NewMetrics(),
	}
	// Now it is not feasible to call updateOffscreen every frame due to performance.
	g.updateOffscreen()
	return g
}

func (gm *Game) updateOffscreen() {
	pool := pond.New(maxWorkers, maxPoolSize)

	for y := 0; y < screenHeight; y += chunk {
		for x := 0; x < screenWidth; x += chunk {
			pool.Submit(NewWorkerRequest(x, y, gm.cam.X, gm.cam.Y, gm.cam.Scale, gm.offscreenPix, gm.timer))
		}
	}

	go func() {
		pool.StopAndWait()

		gm.timer.Reset()

		gm.updateOffscreen()
	}()

}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, key := range g.keys {
		if key == ebiten.KeyArrowLeft {
			// fmt.Println("LEFT")
			g.cam.Left()
		}
		if key == ebiten.KeyArrowRight {
			// fmt.Println("RIGHT")
			g.cam.Right()
		}
		if key == ebiten.KeyArrowUp {
			// fmt.Println("UP")
			g.cam.Up()
		}
		if key == ebiten.KeyArrowDown {
			// fmt.Println("DOWN")
			g.cam.Down()
		}
		if key == ebiten.KeySpace {
			// fmt.Println("ZOOM")
			g.cam.In()
		}
		if key == ebiten.KeyShift {
			// fmt.Println("ZOOM")
			g.cam.Out()
		}
		if key == ebiten.KeyEscape {
			os.Exit(0)
		}
	}

	return nil
}

var t = 0

func (g *Game) Draw(screen *ebiten.Image) {
	t++
	if t > 10 {
		g.offscreen.WritePixels(g.offscreenPix)
		t = 0
	}

	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// temp := new(big.Float).SetInt64(4).SetPrec(1)
	// for {
	// 	temp.Mul(temp, new(big.Float).SetInt64(2))
	// 	fmt.Printf("%d : %d : %s\n", temp.Prec(), temp.MinPrec(), temp.String())
	// }

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mandelbrot (Ebitengine Demo)")

	x := big.NewFloat(-1.94015736280000)
	y := big.NewFloat(-8.66e-7)

	cam := NewCamera(x, y, new(big.Float).SetInt64(4).SetPrec(16))

	if err := ebiten.RunGame(NewGame(cam)); err != nil {
		log.Fatal(err)
	}
}
