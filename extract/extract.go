package extract

import (
	"encoding/base64"
	"strconv"
	"strings"
)

type Thumbnail struct {
	Width        int
	Height       int
	Data         []byte
	size         int
	StartingLine int
	EndingLine   int
}

func ExtractThumbnails(gcode string) ([]Thumbnail, error) {
	lines := strings.Split(gcode, "\n")
	thumbnails := make([]Thumbnail, 0)
	parsingThumbnail := false
	b64buf := strings.Builder{}

	for i, line := range lines {
		if strings.HasPrefix(line, "; thumbnail begin") {
			parsingThumbnail = true
			thumbnail := Thumbnail{}
			thumbnail.StartingLine = i - 1

			dimensions := strings.Split(line, " ")[3]
			size, err := strconv.Atoi(strings.Split(line, " ")[4])
			if err != nil {
				return nil, err
			}
			thumbnail.size = size
			thumbnail.Data = make([]byte, base64.StdEncoding.DecodedLen(thumbnail.size))

			width, err := strconv.Atoi(strings.Split(dimensions, "x")[0])
			if err != nil {
				return nil, err
			}
			thumbnail.Width = width

			height, err := strconv.Atoi(strings.Split(dimensions, "x")[1])
			if err != nil {
				return nil, err
			}
			thumbnail.Height = height

			b64buf = strings.Builder{}

			thumbnails = append(thumbnails, thumbnail)

			continue
		}

		if strings.HasPrefix(line, "; thumbnail end") {
			parsingThumbnail = false
			base64.StdEncoding.Decode(thumbnails[len(thumbnails)-1].Data, []byte(b64buf.String()))
			thumbnails[len(thumbnails)-1].EndingLine = i + 1
			continue
		}

		if parsingThumbnail {
			b64buf.WriteString(strings.Split(line, " ")[1])
		}
	}

	return thumbnails, nil
}
