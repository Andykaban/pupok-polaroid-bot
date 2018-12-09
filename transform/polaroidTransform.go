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
	backgroundWidth = 300
	backgroundHeight = 320

	resizedWidth = 235
	resizedHeight = 245

	blankWidth = 330
	blankHeight = 350
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
	background := imaging.Fit(backgroundRaw, backgroundWidth, backgroundHeight, imaging.Gaussian)

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
	ctx.SetDPI(90)
	ctx.SetFontSize(16)
	ctx.SetSrc(image.NewUniform(color.RGBA{0, 0, 255, 255}))
	return &PolaroidTransform{background: background, ctx: ctx}, nil
}

func (p *PolaroidTransform) addTextLabel(img *image.NRGBA, x, y int, label string) error {
	var textShift int
	p.ctx.SetDst(img)
	p.ctx.SetClip(img.Bounds())
	textLength := utf8.RuneCountInString(label)
	if textLength > 0 && textLength <= 10 {
		textShift = 3
	} else if textLength > 10 && textLength <=15 {
		textShift = 5
	} else {
		textShift = 9
	}
	deltaX := x - (textShift * textLength)
	if deltaX < 0 {
		deltaX = 15
	}
	pt := freetype.Pt(deltaX, y+int(p.ctx.PointToFixed(12)>>6))
	if _, err := p.ctx.DrawString(label, pt); err != nil {
		return err
	}
	return nil
}

func (p *PolaroidTransform) CreatePolaroidImage(srcImagePath, dstImagePath, textLabel string) error {
	var cropImage *image.NRGBA
	sourceImage, err := imaging.Open(srcImagePath)
	if err != nil {
		return err
	}
	sourceBounds := sourceImage.Bounds()
	if sourceBounds.Max.X > sourceBounds.Max.Y {
		newWidth := 15*sourceBounds.Max.Y/16
		widthDelta := int((sourceBounds.Max.X - newWidth)/2)
		cropImage = imaging.Crop(sourceImage, image.Rect(widthDelta, 0,
			sourceBounds.Max.X-widthDelta, sourceBounds.Max.Y))

	} else {
		cropImage = imaging.Clone(sourceImage)
	}
	resizeImage := imaging.Resize(cropImage, resizedWidth, resizedHeight, imaging.Gaussian)
	contrastImage := imaging.AdjustContrast(resizeImage, 12.0)
	sharpenImage := imaging.Sharpen(contrastImage, 17.0)

	combinedImage := imaging.Paste(p.background, sharpenImage, image.Point{X: 14, Y: 20})

	err = p.addTextLabel(combinedImage, 120, 280, textLabel)
	if err != nil {
		return err
	}
	blank := imaging.New(blankWidth, blankHeight, color.White)
	blankBounds := blank.Bounds()
	combinedBounds := combinedImage.Bounds()
	toPos := image.Point{X: blankBounds.Max.X/2 - combinedBounds.Max.X/2,
		Y: blankBounds.Max.Y/2 - combinedBounds.Max.Y/2}
	imageWithBlank := imaging.Paste(blank, combinedImage, toPos)

	angle := float64(utils.GetRandom(-10, 10))
	rotatedImage := imaging.Rotate(imageWithBlank, angle, color.White)

	err = imaging.Save(rotatedImage, dstImagePath)
	if err != nil {
		return err
	}

	return nil
}
