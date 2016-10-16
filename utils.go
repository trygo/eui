package eui

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
)

func SaveImage(filename string, img image.Image) error {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return errors.New("no file extension name")
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch filename[idx:] {
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
	case ".png":
		return png.Encode(f, img)
	}
	log.Fatalln("not supported file type, ", filename[idx:])
	return nil
}

func LoadImage(filename string) (image.Image, error) {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return nil, errors.New("no file extension name")
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch filename[idx:] {
	case ".jpg", ".jpeg":
		return jpeg.Decode(f)
	case ".png":
		return png.Decode(f)
	}
	return nil, errors.New("not supported file type, " + filename[idx:])
}

func FormatRect(r *image.Rectangle) *image.Rectangle {
	rect := image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
	return &rect
}

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func Minf(v1, v2 REAL) REAL {
	if v1 < v2 {
		return v1
	}
	return v2
}

func Maxf(v1, v2 REAL) REAL {
	if v1 > v2 {
		return v1
	}
	return v2
}

func Min(v1, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func Max(v1, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func Mins(vs ...int) (mv int, idx int) {
	mv = vs[0]
	for i := 1; i < len(vs); i++ {
		if mv > vs[i] {
			mv = vs[i]
			idx = i
		}
	}
	return
}

func Maxs(vs ...int) (mv int, idx int) {
	mv = vs[0]
	for i := 1; i < len(vs); i++ {
		if mv < vs[i] {
			mv = vs[i]
			idx = i
		}
	}
	return
}

func Union(r *image.Rectangle, s *image.Rectangle) *image.Rectangle {
	x, y, x3, y3 := r.Min.X, r.Min.Y, r.Max.X, r.Max.Y
	x1, y1, x2, y2 := s.Min.X, s.Min.Y, s.Max.X, s.Max.Y
	xs, _ := Mins(x, x1, x2, x3)
	ys, _ := Mins(y, y1, y2, y3)
	xe, _ := Maxs(x, x1, x2, x3)
	ye, _ := Maxs(y, y1, y2, y3)
	return &image.Rectangle{Min: image.Point{xs, ys}, Max: image.Point{xe, ye}}
}

func IsIntersect(r *image.Rectangle, s *image.Rectangle) bool {
	rx, ry, rx1, ry1 := r.Min.X, r.Min.Y, r.Max.X, r.Max.Y
	if rx > rx1 {
		rx, rx1 = rx1, rx
	}
	if ry > ry1 {
		ry, ry1 = ry1, ry
	}
	sx, sy, sx1, sy1 := s.Min.X, s.Min.Y, s.Max.X, s.Max.Y
	if sx > sx1 {
		sx, sx1 = sx1, sx
	}
	if sy > sy1 {
		sy, sy1 = sy1, sy
	}

	x, y, x1, y1 := rx, ry, rx1, ry1
	if rx < sx {
		x = sx
	}
	if ry < sy {
		y = sy
	}
	if rx1 > sx1 {
		x1 = sx1
	}
	if ry1 > sy1 {
		y1 = sy1
	}
	if x > x1 || y > y1 || (x == 0 && y == 0 && x1 == 0 && y1 == 0) {
		return false
	}

	return true
}

//type TickerChan chan struct{}
type Ticker struct {
	handle chan struct{}
}

func NewTicker(d time.Duration, f func(t time.Time)) *Ticker {
	ticker := &Ticker{handle: make(chan struct{}, 1)}
	go func() {
		timer := time.NewTicker(d)
		defer timer.Stop()
		for {
			select {
			case timestamp := <-timer.C:
				f(timestamp)
			case <-ticker.handle:
				return
			}
		}
	}()
	return ticker
}

func (this *Ticker) Close() {
	close(this.handle)
}

type timespender struct {
	name  string
	start int64
}

func NewTimespender(name string) *timespender {
	return &timespender{name: name, start: time.Now().UnixNano()}
}

func (this *timespender) Spendtime() int64 {
	return time.Now().UnixNano() - this.start

}

func (this *timespender) Print(brake ...float64) {
	spendtime := float64(time.Now().UnixNano()-this.start) / 1000000
	if len(brake) == 0 || spendtime > brake[0] {
		fmt.Println(this.name, ",spend time:", spendtime)
	}
}

func (this *timespender) Reset(name string) {
	this.name = name
	this.start = time.Now().UnixNano()
}
