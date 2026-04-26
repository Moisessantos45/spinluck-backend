package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type TicketData struct {
	TicketID string `json:"ticket_id"`
	Amount   string `json:"amount"`
	DateTime string `json:"date_time"`
	FullName string `json:"full_name"`
}

type LinearGradient struct {
	Start, End color.RGBA
	Rect       image.Rectangle
}

func (g *LinearGradient) ColorModel() color.Model { return color.RGBAModel }
func (g *LinearGradient) Bounds() image.Rectangle { return g.Rect }
func (g *LinearGradient) At(x, y int) color.Color {
	// Calculamos el ratio basado en la altura relativa del rectángulo del ticket
	ratio := float64(y-g.Rect.Min.Y) / float64(g.Rect.Dy())
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	return color.RGBA{
		uint8(float64(g.Start.R)*(1-ratio) + float64(g.End.R)*ratio),
		uint8(float64(g.Start.G)*(1-ratio) + float64(g.End.G)*ratio),
		uint8(float64(g.Start.B)*(1-ratio) + float64(g.End.B)*ratio),
		255,
	}
}

func addLabel(img *image.RGBA, x, y int, label string, col color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(label)
}

func drawDashedLine(img *image.RGBA, x1, x2, y int, col color.Color) {
	for i := x1; i < x2; i += 10 {
		for j := range 5 {
			if i+j < x2 {
				img.Set(i+j, y, col)
			}
		}
	}
}

func GenerateTicketImage(data TicketData) ([]byte, error) {
	width := 450
	height := 600
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	externalBg := color.RGBA{240, 242, 245, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{externalBg}, image.Point{}, draw.Src)

	ticketRect := image.Rect(40, 40, 410, 560)

	ticketGrad := &LinearGradient{
		Start: color.RGBA{255, 255, 255, 255},
		End:   color.RGBA{245, 248, 255, 255},
		Rect:  ticketRect,
	}

	draw.Draw(img, ticketRect, ticketGrad, ticketRect.Min, draw.Src)

	borderColor := color.RGBA{220, 220, 225, 255}
	for x := ticketRect.Min.X; x < ticketRect.Max.X; x++ {
		img.Set(x, ticketRect.Min.Y, borderColor)
		img.Set(x, ticketRect.Max.Y-1, borderColor)
	}
	for y := ticketRect.Min.Y; y < ticketRect.Max.Y; y++ {
		img.Set(ticketRect.Min.X, y, borderColor)
		img.Set(ticketRect.Max.X-1, y, borderColor)
	}

	textColor := color.RGBA{30, 35, 45, 255}
	labelColor := color.RGBA{130, 140, 150, 255}

	addLabel(img, 180, 100, "COMPROBANTE", textColor)
	drawDashedLine(img, 65, 385, 130, color.RGBA{210, 215, 220, 255})

	addLabel(img, 70, 180, "TICKET ID", labelColor)
	addLabel(img, 70, 205, data.TicketID, textColor)

	addLabel(img, 70, 270, "TOTAL", labelColor)
	addLabel(img, 70, 295, "$"+data.Amount, textColor)

	addLabel(img, 70, 360, "FECHA", labelColor)
	addLabel(img, 70, 385, data.DateTime, textColor)

	addLabel(img, 70, 450, "CLIENTE", labelColor)
	addLabel(img, 70, 475, data.FullName, textColor)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestingGenerate() {
	data := TicketData{
		TicketID: "0120077398910288",
		Amount:   "$3,800.00",
		DateTime: "22 AUG, 2025 | 13:29",
		FullName: "Juan Perez",
	}

	imgBytes, err := GenerateTicketImage(data)
	if err != nil {
		log.Fatal("error generating ticket image: ", err)
	}

	log.Println("bytes generados:", len(imgBytes))
}
