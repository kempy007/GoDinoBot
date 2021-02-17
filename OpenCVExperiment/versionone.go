package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
	"github.com/micmonay/keybd_event"
	"gocv.io/x/gocv"
)

func setWindow_OSWin(WP image.Rectangle) {
	// #### Windows specific
	myClassname, _ := syscall.UTF16PtrFromString("Chrome_WidgetWin_1")
	myWindowname, _ := syscall.UTF16PtrFromString("chrome://dino/ - Google Chrome")
	h := win.HWND(unsafe.Pointer(win.FindWindow(myClassname, myWindowname)))
	win.SetFocus(h)
	win.SetActiveWindow(h)
	win.SetForegroundWindow(h)
	myWindowPlacement := new(win.WINDOWPLACEMENT)
	win.GetWindowPlacement(h, myWindowPlacement)
	// smallest dino game will fit

	// my_xy := int32(100)
	// my_w := int32(520)
	// my_h := int32(240)

	// my_xy := int32(WP.Min.X)
	// my_w := int32(WP.Max.X)
	// my_h := int32(WP.Max.Y)

	myWindowPlacement.RcNormalPosition.Left = int32(WP.Min.X)
	myWindowPlacement.RcNormalPosition.Top = int32(WP.Min.Y)
	myWindowPlacement.RcNormalPosition.Right = int32(WP.Max.X + WP.Min.X)
	myWindowPlacement.RcNormalPosition.Bottom = int32(WP.Max.Y + WP.Min.Y)
	// myWindowPlacement.RcNormalPosition.Left = my_xy
	// myWindowPlacement.RcNormalPosition.Top = my_xy
	// myWindowPlacement.RcNormalPosition.Right = my_w + my_xy
	// myWindowPlacement.RcNormalPosition.Bottom = my_h + my_xy
	win.SetWindowPlacement(h, myWindowPlacement)
}

func main() {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")
	}
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_SPACE)
	// remember the window edges of few px and title bars ~40px etc can eat into positions of things
	WindowPosition := image.Rectangle{Min: image.Point{X: 100, Y: 100}, Max: image.Point{X: 520, Y: 240}}
	WindowFocus := image.Rectangle{Min: image.Point{X: 110, Y: 180}, Max: image.Point{X: 500, Y: 150}}
	//bgFocus := image.Rectangle{Min: image.Point{X: 120, Y: 190}, Max: image.Point{X: 5, Y: 5}}
	// roiStartFocus := image.Rectangle{Min: image.Point{X: 200, Y: 285}, Max: image.Point{X: 50, Y: 35}}
	// TODO: fix as these two not fit correctly
	roiStartFocus := image.Rectangle{Min: image.Point{X: 196, Y: 279}, Max: image.Point{X: 50, Y: 40}}
	roiStartBounding := image.Rectangle{Min: image.Point{X: 85, Y: 103}, Max: image.Point{X: 135, Y: 140}} // 50, 40
	setWindow_OSWin(WindowPosition)

	// webcam, _ := gocv.VideoCaptureDevice(0)
	//screen, _ := screenshot.Capture(0, 0, 200, 50)
	fullWindow := gocv.NewWindow("Let's Play Dino")
	fullWindow.MoveWindow(620, 100)
	fullWindow.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowAutosize) // This does not work
	// img := gocv.NewMat()

	roiWindow := gocv.NewWindow("ROI")
	roiWindow.MoveWindow(620, 350)
	roiWindow.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowAutosize)

	dino := gocv.IMRead("../imageDino/dino.jpg", gocv.IMReadGrayScale)
	w_dino := dino.Rows()
	h_dino := dino.Cols()
	i := 0
	di := 0

	roiPriorScreen, _ := screenshot.Capture(int(roiStartFocus.Min.X), int(roiStartFocus.Min.Y), int(roiStartFocus.Max.X), int(roiStartFocus.Max.Y))
	var buff3 bytes.Buffer
	png.Encode(&buff3, roiPriorScreen)
	roiPriorImg, _ := gocv.IMDecode(buff3.Bytes(), gocv.IMReadColor)

	for {
		// capture screen and encode into buffer
		// screen, _ := screenshot.Capture(int(my_xy), int(my_xy), int(my_w), int(my_h))
		// screen, _ := screenshot.Capture(int(WindowPosition.Min.X+10), int(WindowPosition.Min.Y+80), int(WindowPosition.Max.X-20), int(WindowPosition.Max.Y-90)) // Trimmed down
		sceneScreen, _ := screenshot.Capture(int(WindowFocus.Min.X), int(WindowFocus.Min.Y), int(WindowFocus.Max.X), int(WindowFocus.Max.Y)) // Trimmed down
		var buff bytes.Buffer
		png.Encode(&buff, sceneScreen)
		// decode from buffer into gocv image format
		sceneImg, _ := gocv.IMDecode(buff.Bytes(), gocv.IMReadColor)
		defer sceneImg.Close()
		// webcam.Read(&img)

		// create copy of full image and crop > is pita so going with screenshot again
		roiScreen, _ := screenshot.Capture(int(roiStartFocus.Min.X), int(roiStartFocus.Min.Y), int(roiStartFocus.Max.X), int(roiStartFocus.Max.Y))
		var buff2 bytes.Buffer
		png.Encode(&buff2, roiScreen)
		roiImg, _ := gocv.IMDecode(buff2.Bytes(), gocv.IMReadColor)
		roiWindow.IMShow(roiImg)
		roiWindow.WaitKey(1)
		gocv.Rectangle(&sceneImg, roiStartBounding, color.RGBA{R: 0, G: 255, B: 0, A: 0}, 1)

		// TODO: compare previous and current image for diff.
		grayRoiImg := gocv.NewMat()
		defer grayRoiImg.Close()
		gocv.CvtColor(roiImg, &grayRoiImg, gocv.ColorRGBAToGray)

		grayPriorRoiImg := gocv.NewMat()
		defer grayPriorRoiImg.Close()
		gocv.CvtColor(roiPriorImg, &grayPriorRoiImg, gocv.ColorRGBAToGray)

		absdiffImg := gocv.NewMat()
		defer absdiffImg.Close()
		gocv.AbsDiff(grayRoiImg, grayPriorRoiImg, &absdiffImg)
		if absdiffImg.Mean().Val1 > 0 {
			gocv.PutText(&sceneImg, "jump", image.Point{X: 30, Y: 30}, gocv.FontHersheyPlain, 1, color.RGBA{R: 0, G: 0, B: 255, A: 0}, 2)
			// TODO: send keyinput here
			err = kb.Launching()
			if err != nil {
				panic(err)
			}
		}

		// // Do we need to capture bg colour???
		// subtractImg := gocv.NewMat()
		// defer subtractImg.Close()
		// gocv.Subtract(roiPriorImg, roiImg, &subtractImg)
		// subsMean := subtractImg.Mean()
		// if subsMean.Val1 > 0 {
		// 	fmt.Print(subsMean) // Prints: {0 0 0 0}
		// 	// Then copy current image to prior
		// 	//roiImg.CopyTo(&roiPriorImg) // TODO: ************ Crashes, the changing window could be issue with copy
		// 	roiPriorScreen, _ = screenshot.Capture(int(roiStartFocus.Min.X), int(roiStartFocus.Min.Y), int(roiStartFocus.Max.X), int(roiStartFocus.Max.Y))
		// 	var buff4 bytes.Buffer
		// 	png.Encode(&buff4, roiPriorScreen)
		// 	roiPriorImg, _ = gocv.IMDecode(buff4.Bytes(), gocv.IMReadColor)
		// }

		grayImg := gocv.NewMat()
		defer grayImg.Close()
		resultDino := gocv.NewMat()
		defer resultDino.Close()

		gocv.CvtColor(sceneImg, &grayImg, gocv.ColorRGBAToGray)
		// gocv.CvtColor(img, &mask, gocv.ColorRGBAToGray)
		// window.IMShow(mask)
		// window.WaitKey(1)
		// window.IMShow(dino) // dino crashes this, was not pathed right :FIXED
		// window.WaitKey(1)

		// gocv.MatchTemplate(img, dino, &res_dino, gocv.TmCcoeffNormed, mask)
		// gocv.MatchTemplate(gray_img, dino, &res_dino, gocv.TmCcoeffNormed, mask) // TODO this crashes here???
		mask := gocv.NewMat()
		gocv.MatchTemplate(grayImg, dino, &resultDino, gocv.TmCcoeffNormed, mask)
		mask.Close()
		_, maxConfidence, minLoc, maxLoc := gocv.MinMaxLoc(resultDino)
		if maxConfidence < 0.935 {
			fmt.Println("Max confidence of " + fmt.Sprintf("%f", maxConfidence) + " is too low") //("Max confidence of %f is too low. MatchTemplate could not find template in scene. \\n", maxConfidence)
			// we are dead perhaps press a key to reset
			di++
			if di > 100 {
				roiStartFocus.Max.X = 40
				roiStartBounding.Max.X = 140
				di = 0
				err = kb.Launching()
				if err != nil {
					panic(err)
				}
			}
			fullWindow.IMShow(sceneImg)
			fullWindow.WaitKey(1)
		} else {
			// alive it is above 95% else below 85%
			fmt.Println("Max confidence is " + fmt.Sprintf("%f", maxConfidence))
			// we are alive
			// DinoPosition := image.Rectangle{Min: image.Point{X: minLoc.X, Y: minLoc.Y}, Max: image.Point{X: maxLoc.X, Y: maxLoc.Y}}
			DinoPosition := image.Rectangle{Min: image.Point{X: maxLoc.X, Y: maxLoc.Y}, Max: image.Point{X: maxLoc.X + dino.Cols(), Y: maxLoc.Y + dino.Rows()}}
			// gocv.Rectangle(&resultDino, DinoPosition, color.RGBA{R: 255, G: 255, B: 255, A: 0}, 1)
			gocv.Rectangle(&sceneImg, DinoPosition, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
			fullWindow.IMShow(sceneImg)
			fullWindow.WaitKey(1)
			i++
			// fmt.Print(i)
			// We grow the ROI to account for speed increase, 10 seems to fast and jumps early
			// 10 = ~100
			// 20 = ~600-706
			// 18 = 490-716
			if i > 18 {
				roiStartFocus.Max.X++
				roiStartBounding.Max.X++
				i = 0
			}
		}

		// update the priorscreen here so that hopefully it is right size for next loop and does not break the absDiff or subtract func
		roiPriorScreen, _ = screenshot.Capture(int(roiStartFocus.Min.X), int(roiStartFocus.Min.Y), int(roiStartFocus.Max.X), int(roiStartFocus.Max.Y))
		var buff4 bytes.Buffer
		png.Encode(&buff4, roiPriorScreen)
		roiPriorImg, _ = gocv.IMDecode(buff4.Bytes(), gocv.IMReadColor)

		//threshold := 0.8
		//gocv.BoundingRect() // look at instead of bbox from python or gocv.findNonZero() as per https://stackoverflow.com/questions/58763007/opencv-equivalent-of-np-where
		// dinoLocation := gocv.NewMat()
		// defer dinoLocation.Close()
		// gocv.FindNonZero(resultDino, &dinoLocation)
		// window.IMShow(dinoLocation) // crashes here
		// window.WaitKey(1)
		// _, _, minLoc, maxLoc := gocv.MinMaxLoc(dinoLocation)
		// dinoX = dinoLocation.Cols()
		// dinoH = dinoLocation.Rows()
		//fmt.Print(threshold, w_dino, h_dino)

		if 1 == 0 {
			fmt.Print(w_dino, h_dino)
			fmt.Print(minLoc, maxLoc)
		}
		//fmt.Print(minLoc, maxLoc)

		// window.IMShow(img)
		// window.WaitKey(1)

	}
}
