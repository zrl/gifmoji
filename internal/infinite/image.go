package infinite

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

const frameCount = 10

type Transformer struct {
}

func (t *Transformer) Transform(args []string) error {
	filename := args[0]
	img, err := loadImage(filename)
	if err != nil {
		return err
	}

	g, err := toGIF(img)
	if err != nil {
		return err
	}

	out := getOutputFilename(filename)
	if err := saveGIF(out, g); err != nil {
		return err
	}
	fmt.Printf("Created %s. Enjoy!\n", out)
	return nil
}

func loadImage(filename string) (image.Image, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(bs))
	return img, err
}

func toGIF(img image.Image) (*gif.GIF, error) {
	x := img.Bounds().Dx()
	deltas := distribute(x, frameCount)
	g := &gif.GIF{
		Image: make([]*image.Paletted, frameCount),
		Delay: make([]int, frameCount),
	}
	for i := 0; i < frameCount; i++ {
		shifted := shiftImage(img, deltas[i], deltas[:i])
		buf := bytes.Buffer{}
		if err := gif.Encode(&buf, shifted, nil); err != nil {
			return nil, err
		}
		tmpimg, err := gif.Decode(&buf)
		if err != nil {
			return nil, err
		}
		g.Image[i] = tmpimg.(*image.Paletted)
	}

	return g, nil
}

func shiftImage(img image.Image, width int, deltas []int) image.Image {
	sum := 0
	for _, d := range deltas {
		sum += d
	}
	x, y := img.Bounds().Dx(), img.Bounds().Dy()
	dx := width + sum
	dst := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.Copy(dst, image.Point{dx, 0}, img, image.Rect(0, 0, x-dx, y), draw.Over, nil)
	draw.Copy(dst, image.Point{0, 0}, img, image.Rect(x-dx, 0, x, y), draw.Over, nil)
	return dst
}

func getOutputFilename(inFilename string) string {
	return fmt.Sprintf("infinite-%s.gif", strings.ReplaceAll(inFilename, ".png", ""))
}

func saveGIF(filename string, g *gif.GIF) error {
	buf := &bytes.Buffer{}
	err := gif.EncodeAll(buf, g)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func distribute(count, buckets int) []int {
	q := count / buckets // quotient, how many whole values per bucket
	r := count % buckets // remainder, how much leftover
	out := make([]int, buckets)
	for i := 0; i < buckets; i++ {
		out[i] = q
		if r != 0 {
			out[i] += 1
			r -= 1
		}
	}
	return out
}
