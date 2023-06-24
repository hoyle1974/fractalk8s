package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/alitto/pond"
	"github.com/hajimehoshi/ebiten/v2"
)

func calc(client *http.Client, x, y []*big.Float, iter int) []int {

	// fmt.Printf("Calc %d values\n", len(x))

	req := ""
	for i := 0; i < len(x); i++ {
		if i == 0 {
			req = fmt.Sprintf("%s,%s,%d", x[i].String(), y[i].String(), iter)
		} else {
			req = req + fmt.Sprintf("\n%s,%s,%d", x[i].String(), y[i].String(), iter)
		}
	}

	// url := fmt.Sprintf("http://fractalk8s.decepticons.local/iter")
	url := "http://localhost:8080/iter"
	r, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(req)))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ii := strings.Split(string(body), ",")
	if len(ii) != len(x) {
		panic(fmt.Sprintf("We requested %d values, but received %d", len(x), len(ii)))
	}

	ret := make([]int, len(ii))
	for idx, i := range ii {
		iter, _ := strconv.Atoi(i)
		ret[idx] = iter
	}

	return ret
}

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 128
	chunk        = 16
	maxPoolSize  = (screenWidth/chunk)*(screenHeight/chunk) + 1
	maxWorkers   = 32
)

var (
	palette [maxIt]byte
)

func init() {
	for i := range palette {
		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	}
}

func color(it int) (r, g, b byte) {
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
}

func NewGame(cam Camera) *Game {
	g := &Game{
		cam:          cam,
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		offscreenPix: make([]byte, screenWidth*screenHeight*4),
	}
	// Now it is not feasible to call updateOffscreen every frame due to performance.
	g.updateOffscreen()
	return g
}

func (gm *Game) updateOffscreen() {
	pool := pond.New(maxWorkers, maxPoolSize)

	t, _ := gm.cam.Scale.MarshalText()
	fmt.Println(string(t))

	for y := 0; y < screenHeight; y += chunk {
		for x := 0; x < screenWidth; x += chunk {
			pool.Submit(NewWorkerRequest(x, y, gm.cam.X, gm.cam.Y, gm.cam.Scale, gm.offscreenPix))
		}
	}

	go func() {
		pool.StopAndWait()
		fmt.Println("tick")
		gm.cam.Scale.Mul(gm.cam.Scale, big.NewFloat(0.5))
		gm.updateOffscreen()
	}()

}

func (g *Game) Update() error {
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
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mandelbrot (Ebitengine Demo)")

	cam := NewCamera(-1.9, 0.0, 4.0)

	if err := ebiten.RunGame(NewGame(cam)); err != nil {
		log.Fatal(err)
	}
}
