package decoder

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Black       = 0
	White       = 1
	Transparent = 2
)

type Layer []int

type LayerStats struct {
	Zs, Ones, Twos int
}

type ImgData struct {
	Width, Height int
	layers        []Layer
	lStats        []*LayerStats
	CRC           int
}

type Img struct {
	Width, Height int
	pixels        []int
}

func (img *ImgData) Decode() *Img {
	nLayers := len(img.layers)

	pxls := make([]int, img.Width*img.Height)
	for i := 0; i < img.Width*img.Height; i++ {
		for j := 0; j < nLayers; j++ {
			if pxl := img.layers[j][i]; pxl == Black || pxl == White {
				pxls[i] = pxl
				break
			}
		}
	}

	return &Img{
		Width:  img.Width,
		Height: img.Height,
		pixels: pxls,
	}
}

func (img *Img) String() string {
	sb := &strings.Builder{}
	sb.WriteString("\n")
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			sb.WriteString(fmt.Sprintf("%s ", func(d int) string {
				ret := ""
				switch d {
				case Black:
					ret = " "
				case White:
					ret = "*"
				}
				return ret
			}(img.pixels[(y*img.Width)+x])))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func CreateImage(imgStream string, h, w int) (*ImgData, error) {
	nLayers := int(len(imgStream) / (h * w))
	img := &ImgData{
		Width:  w,
		Height: h,
		layers: make([]Layer, nLayers),
		lStats: make([]*LayerStats, nLayers),
	}

	sr := strings.NewReader(imgStream)

	parseLayer := func() (*LayerStats, Layer, error) {
		layerStats := &LayerStats{Zs: 0, Ones: 0, Twos: 0}
		lyr := make(Layer, h*w)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				ch, _, err := sr.ReadRune()
				if err != nil {
					return nil, nil, err
				}

				p, err := strconv.Atoi(string(ch))
				if err != nil {
					return nil, nil, err
				}

				lyr[(y*w + x)] = p
				switch p {
				case 0:
					layerStats.Zs++
				case 1:
					layerStats.Ones++
				case 2:
					layerStats.Twos++
				}
			}
		}

		return layerStats, lyr, nil

	}

	minZs, p := int((^uint(0))>>1), 0
	for i := 0; i < nLayers; i++ {
		ls, l, err := parseLayer()
		if err != nil {
			return nil, err
		}
		img.layers[i] = l
		img.lStats[i] = ls

		if ls.Zs < minZs {
			minZs = ls.Zs
			p = ls.Ones * ls.Twos
		}
	}

	img.CRC = p
	return img, nil
}

// Utils
func (img *ImgData) String() string {
	lyr := &strings.Builder{}
	for i, layer := range img.layers {
		ls := img.lStats[i]
		w, h := img.Width, img.Height
		lyr.WriteString(fmt.Sprintf("Layer %d\n", i+1))
		lyr.WriteString(fmt.Sprintf("%s\n", layer.String(h, w)))
		lyr.WriteString(fmt.Sprintf("%s\n", ls))
	}
	return lyr.String()
}

func (layer Layer) String(h, w int) string {
	lyr := &strings.Builder{}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			lyr.WriteString(fmt.Sprintf("%d  ",
				layer[y*w+x]))
		}
		lyr.WriteString("\n")
	}
	return lyr.String()
}

func (ls LayerStats) String() string {
	return fmt.Sprintf("LayerStats = [Zeros= %d, Ones= %d, Twos= %d]",
		ls.Zs, ls.Ones, ls.Twos)
}
