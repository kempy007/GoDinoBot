package main

import (
	// "image"
	// "image/color"
	// "math"
	// "reflect"
	// "testing"
	"testing"

	"gocv.io/x/gocv"
)

func TestMatchTemplate(t *testing.T) {
	imgScene := gocv.IMRead("testScene.jpg", gocv.IMReadGrayScale)
	if imgScene.Empty() {
		t.Error("Invalid read of testScene.jpg in MatchTemplate test")
	}
	defer imgScene.Close()

	imgTemplate := gocv.IMRead("testDetect.jpg", gocv.IMReadGrayScale)
	if imgTemplate.Empty() {
		t.Error("Invalid read of testDetect.jpg in MatchTemplate test")
	}
	defer imgTemplate.Close()

	result := gocv.NewMat()
	defer result.Close()
	m := gocv.NewMat()
	gocv.MatchTemplate(imgScene, imgTemplate, &result, gocv.TmCcoeffNormed, m)
	m.Close()
	_, maxConfidence, _, _ := gocv.MinMaxLoc(result)
	if maxConfidence < 0.95 {
		t.Errorf("Max confidence of %f is too low. MatchTemplate could not find template in scene.", maxConfidence)
	}
}
