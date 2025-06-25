package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"net/http"
	"strings"
)

type Pixel struct {
	X, Y  int
	Color string
}

func main() {
	http.HandleFunc("/generate-avatar", avatarHandler)
	log.Println("🚀 Pixel Art Avatar Microservice started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name param is required", http.StatusBadRequest)
		return
	}

	svg := generateAvatarSVG(name)
	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprint(w, svg)
}

func generateAvatarSVG(name string) string {
	cell := 16
	grid := 32
	width := grid * cell
	height := grid * cell

	bgColor := pickColor(name + "bg")
	skinColor := pickColor(name + "skin")
	hairColor := pickColor(name + "hair")
	beardColor := pickColor(name + "beard")
	eyeColor := "#000000"
	mouthColor := "#c33"

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s"/>`, bgColor))

	cx, cy := grid/2, grid/2
	radius := grid / 3

	// Face oval
	for y := 0; y < grid; y++ {
		for x := 0; x < grid; x++ {
			dx := x - cx
			dy := y - cy
			dist := math.Sqrt(float64(dx*dx + dy*dy))
			if dist < float64(radius) {
				sb.WriteString(rect(x, y, cell, skinColor))
			}
		}
	}

	// Hair
	for y := cy - radius; y < cy-radius/2; y++ {
		for x := cx - radius/2; x <= cx+radius/2; x++ {
			sb.WriteString(rect(x, y, cell, hairColor))
		}
	}
	for x := cx - radius/2; x <= cx+radius/2; x += 1 {
		sb.WriteString(rect(x, cy-radius-1, cell, hairColor))
	}

	// Eyes
	sb.WriteString(rect(cx-4, cy-1, cell, "#fff"))
	sb.WriteString(rect(cx-4, cy-1, cell/2, eyeColor))
	sb.WriteString(rect(cx+2, cy, cell, "#fff"))
	sb.WriteString(rect(cx+2, cy, cell/2, eyeColor))

	// Mouth (crooked)
	for i := -1; i <= 1; i++ {
		sb.WriteString(rect(cx+i, cy+4, cell, mouthColor))
	}

	// Beard
	for x := cx - 3; x <= cx+3; x++ {
		sb.WriteString(rect(x, cy+5, cell, beardColor))
	}

	// Mustache curled
	for x := cx - 2; x <= cx+2; x++ {
		sb.WriteString(rect(x, cy+2, cell, beardColor))
	}
	sb.WriteString(rect(cx-3, cy+1, cell, beardColor))
	sb.WriteString(rect(cx+3, cy+1, cell, beardColor))

	sb.WriteString(`</svg>`)
	return sb.String()
}

func rect(x, y, size int, color string) string {
	return fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x*size, y*size, size, size, color)
}

func pickColor(seed string) string {
	h := fnv.New32a()
	h.Write([]byte(seed))
	hash := h.Sum32()
	return fmt.Sprintf("#%06x", hash&0xFFFFFF)
}
