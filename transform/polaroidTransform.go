package transform

import (
	"github.com/Andykaban/pupok-polaroid-bot/utils"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"io/ioutil"
	"unicode/utf8"
)

const (
	BACKGROUND_WIDTH  = 300
	BACKGROUND_HEIGHT = 320

	RESIZED_WITH      = 235
	RESIZED_HEIGHT    = 245

	BLANK_WIDTH = 330
	BLANCH_HEIGHT = 350
)

type PolaroidTransform struct {
	background *image.NRGBA
	ctx        *freetype.Context
}

func New(backgroundPath, fontPath string) (*PolaroidTransform, error) {
	backgroundRaw, err := imaging.Open(backgroundPath)
	if err != nil {
		return nil, err
	}
	background := imaging.Fit(backgroundRaw, BACKGROUND_WIDTH, BACKGROUND_HEIGHT, imaging.Gaussian)

	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	ctx := freetype.NewContext()
	ctx.SetFont(font)
	ctx.SetDPI(80)
	ctx.SetFontSize(16)
	ctx.SetSrc(image.NewUniform(color.RGBA{0, 0, 255, 255}))
	return &PolaroidTransform{background: background, ctx: ctx}, nil
}

func (p *PolaroidTransform) addTextLabel(img *image.NRGBA, x, y int, label string) error {
	p.ctx.SetDst(img)
	p.ctx.SetClip(img.Bounds())
	textLength := utf8.RuneCountInString(label)
	textSize := 2 * textLength
	deltaX := x - textSize
	if deltaX < 0 {
		deltaX = 10
	}
	pt := freetype.Pt(deltaX, y+int(p.ctx.PointToFixed(12)>>6))
	if _, err := p.ctx.DrawString(label, pt); err != nil {
		return err
	}
	return nil
}

func (p *PolaroidTransform) CreatePolaroidImage(srcImagePath, dstImagePath, textLabel string) error {
	sourceImage, err := imaging.Open(srcImagePath)
	if err != nil {
		return err
	}
	resizeImage := imaging.Resize(sourceImage, RESIZED_WITH, RESIZED_HEIGHT, imaging.Gaussian)
	contrastImage := imaging.AdjustContrast(resizeImage, 12.0)
	sharpenImage := imaging.Sharpen(contrastImage, 17.0)

	combinedImage := imaging.Paste(p.background, sharpenImage, image.Point{X: 14, Y: 20})

	p.addTextLabel(combinedImage, 120, 285, textLabel)

	blank := imaging.New(BLANK_WIDTH, BLANCH_HEIGHT, color.White)
	blankBound := blank.Bounds()
	combinedBounds := combinedImage.Bounds()
	toPos := image.Point{X: blankBound.Max.X/2 - combinedBounds.Max.X/2,
		Y: blankBound.Max.Y/2 - combinedBounds.Max.Y/2}
	imageWithBlank := imaging.Paste(blank, combinedImage, toPos)

	angle := float64(utils.GetRandom(-10, 10))
	rotatedImage := imaging.Rotate(imageWithBlank, angle, color.White)

	imaging.Save(rotatedImage, dstImagePath)

	return nil
}
