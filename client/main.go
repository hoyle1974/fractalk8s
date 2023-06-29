package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/alitto/pond"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 255
	chunk        = 32
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

func color(it int) (r, g, b byte) {
	it = maxIt - it

	if it == maxIt {
		return 0xff, 0xff, 0xff
	}
	c := palette[it]
	return (c * 3) % 255, ((c + 25) * 5) % 255, ((c + 128) * 7) % 255
}

type Game struct {
	cam          *Camera
	offscreen    *ebiten.Image
	offscreenPix []byte
	keys         []ebiten.Key
	metrics      *Metrics
}

func NewGame(cam *Camera) *Game {
	g := &Game{
		cam:          cam,
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		offscreenPix: make([]byte, screenWidth*screenHeight*4),
		metrics:      NewMetrics(),
	}
	// Now it is not feasible to call updateOffscreen every frame due to performance.
	g.updateOffscreen()
	return g
}

func (gm *Game) updateOffscreen() {
	pool := pond.New(maxWorkers, maxPoolSize)

	for y := 0; y < screenHeight; y += chunk {
		for x := 0; x < screenWidth; x += chunk {
			pool.Submit(NewWorkerRequest(x, y, gm.cam.X, gm.cam.Y, gm.cam.Scale, gm.offscreenPix, gm.metrics))
		}
	}

	go func() {
		pool.StopAndWait()

		gm.metrics.Reset()
		fmt.Printf("X,Y = %s %s\n", gm.cam.X.Text('f', int(gm.cam.X.MinPrec())), gm.cam.Y.Text('f', int(gm.cam.Y.MinPrec())))
		fmt.Printf("S = %s\n", gm.cam.Scale.Text('f', int(gm.cam.Scale.MinPrec())))

		for {
			if gm.cam.IsDirty() {
				break
			}
			time.Sleep(time.Second)
		}
		gm.updateOffscreen()
	}()

}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, key := range g.keys {
		if key == ebiten.KeyArrowLeft {
			g.cam.Left()
		}
		if key == ebiten.KeyArrowRight {
			g.cam.Right()
		}
		if key == ebiten.KeyArrowUp {
			g.cam.Up()
		}
		if key == ebiten.KeyArrowDown {
			g.cam.Down()
		}
		if key == ebiten.KeySpace {
			g.cam.In()
		}
		if key == ebiten.KeyShift {
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
		for x := -10; x < 10; x++ {
			for y := -10; y < 10; y++ {
				if x == 0 || y == 0 {
					p := 4 * ((((screenHeight / 2) + y) * screenWidth) + ((screenWidth / 2) + x))
					g.offscreenPix[p] = 0xff - g.offscreenPix[p]
					g.offscreenPix[p+1] = 0xff - g.offscreenPix[p+1]
					g.offscreenPix[p+2] = 0xff - g.offscreenPix[p+2]
					g.offscreenPix[p+3] = 0xff

					g.offscreen.WritePixels(g.offscreenPix)

				}
			}
		}
		t = 0
	}

	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mandelbrot (Ebitengine Demo)")

	// x := big.NewFloat(-1.94015736280000)
	// y := big.NewFloat(-8.66e-7)

	x, _, _ := new(big.Float).Parse("-1.940158125066515110802223095951715258162359949311875725036274748603942405522637670856056502088904380798339843750", 10)
	y, _, _ := new(big.Float).Parse("-0.000000862709948779723462565049404001370272990386055543158893226625069661217537486663787", 10)

	x.SetPrec(1000)
	y.SetPrec(1000)

	cam := NewCamera(x, y, new(big.Float).SetInt64(4).SetPrec(1000))

	if err := ebiten.RunGame(NewGame(cam)); err != nil {
		log.Fatal(err)
	}
}
