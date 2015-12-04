package main

import (
	"flag"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"

	"code.google.com/p/graphics-go/graphics"
)

func incColor(c color.RGBA) color.RGBA {
	if c.R == 0 && c.G > 0 {
		c.G--
		c.B++
	} else if c.G == 0 && c.B > 0 {
		c.B--
		c.R++
	} else if c.B == 0 && c.R > 0 {
		c.R--
		c.G++
	} else {
		c = color.RGBA{0xff, 0, 0, 0xff}
	}
	return c
}

func incColorN(c color.RGBA, n int) color.RGBA {
	for i := 0; i < n; i++ {
		c = incColor(c)
	}
	return c
}

type ctx struct {
	g           *gif.GIF
	img         image.Image
	fps         int
	totalFrames int
	currFrame   int
	w, h        int
}

func newCtx(img image.Image, totalFrames, fps int, w, h int) *ctx {
	img = imaging.Fit(img, w, h, imaging.Linear)
	return &ctx{
		g:           &gif.GIF{LoopCount: -1},
		img:         img,
		fps:         fps,
		totalFrames: totalFrames,
		w:           w,
		h:           h,
	}
}

func maxRotatedSize(n int) int {
	nf := float64(n) / 2
	nf2 := nf / math.Cos(math.Pi/4)
	return int(nf2) * 2
}

var addr = flag.String("addr", "", "[host]:port to listen on, if set overrides all other behavior. Other flag values will be used as defaults for calls")
var totalFrames = flag.Int("totalFrames", 20, "total number of frames output gif should have")
var fps = flag.Int("fps", 20, "frames per second the gif should run at")
var width = flag.Int("maxWidth", 128, "max width of the gif to output")
var height = flag.Int("maxHeight", 128, "max height of the gif to output")

var doLog = true

func main() {
	flag.Parse()

	if *addr != "" {
		doLog = false

		http.HandleFunc("/", indexHTTP)
		http.HandleFunc("/partyfy", partyHTTP)

		log.Printf("listening on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}

	if err := partyfy(os.Stdin, *totalFrames, *fps, *width, *height, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func reqInt(r *http.Request, key string, def int) int {
	val := r.FormValue(key)
	if val == "" {
		return def
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return valInt
}

func partyHTTP(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image/gif")
	partyfy(
		f,
		reqInt(r, "totalFrames", *totalFrames),
		reqInt(r, "fps", *fps),
		reqInt(r, "maxWidth", *width),
		reqInt(r, "maxHeight", *height),
		w,
	)
}

func partyfy(r io.Reader, totalFrames, fps, width, height int, w io.Writer) error {
	srcImg, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	c := newCtx(srcImg, totalFrames, fps, width, height)
	for {
		if ok, err := c.nextFrame(); err != nil {
			return err
		} else if !ok {
			break
		}
	}

	return c.writeTo(w)
}

func (c *ctx) modify(img image.Image, col color.Color, angle float64) image.Image {
	out := image.NewRGBA(img.Bounds())
	un := image.NewUniform(col)
	unAlph := image.NewUniform(color.RGBA{0, 0, 0, 100})
	draw.DrawMask(out, out.Bounds(), un, image.ZP, img, image.ZP, draw.Over)
	draw.DrawMask(out, out.Bounds(), img, image.ZP, unAlph, image.ZP, draw.Over)

	out2 := image.NewRGBA(img.Bounds())
	graphics.Rotate(out2, out, &graphics.RotateOptions{angle})
	return out2
}

func (c *ctx) nextFrame() (bool, error) {
	if c.currFrame >= c.totalFrames {
		return false, nil
	}
	if doLog {
		log.Printf("making frame %d/%d", c.currFrame+1, c.totalFrames)
	}

	fract := (1 / float64(c.totalFrames)) * float64(c.currFrame)

	colFrames := int(fract * (0xff * 3))
	col := incColorN(color.RGBA{0xff, 0, 0, 0xff}, colFrames)

	angle := fract * 2 * math.Pi

	c.currFrame++
	return true, c.addFrame(c.modify(c.img, col, angle))
}

var pal = func() color.Palette {
	pal := palette.Plan9
	pal[0] = color.RGBA{0, 0, 0, 0}
	return pal
}()

func (c *ctx) addFrame(img image.Image) error {
	//opts := gif.Options{
	//	NumColors: 256,
	//	Drawer:    draw.FloydSteinberg,
	//}
	b := img.Bounds()

	// More or less taken from the image/gif package
	pimg := image.NewPaletted(b, pal)
	draw.Draw(pimg, b, img, image.ZP, draw.Src)
	//opts.Drawer.Draw(pimg, b, img, image.ZP)

	spf := 1 / float64(c.fps)

	//c.g.Image = append(c.g.Image, g2.Image[0])
	c.g.Image = append(c.g.Image, pimg)
	c.g.Delay = append(c.g.Delay, int(spf*100))
	c.g.Disposal = append(c.g.Disposal, 2)
	return nil
}

func (c *ctx) writeTo(w io.Writer) error {
	return gif.EncodeAll(w, c.g)
}
