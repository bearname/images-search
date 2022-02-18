package main
//
//import (
//	"fmt"
//	log "github.com/sirupsen/logrus"
//	"time"
//)
//
//type Coordinator struct {
//	maxAttemptsBeforeNotify int
//	pictureRepo             PictureRepo
//	notifier                Notifier
//}
//
//type Notifier interface {
//	Notify(picture Picture) domainerror
//}
//
//type PictureStatus int
//
//const (
//	Success PictureStatus = iota
//	Processing
//	Failed
//)
//
//type PictureRepo interface {
//	UpdateImageHandle(picture Picture) domainerror
//}
//
//type Picture struct {
//	Attempts         int
//	ExecuteAfter     time.Time
//	Status           PictureStatus
//	IsOriginalSaved  bool
//	IsPreviewScaled  bool
//	IsTextRecognized bool
//	IsMobileScaled   bool
//}
//
//func (c *Coordinator) PerformAddImage(picture Picture) domainerror {
//	if !picture.IsOriginalSaved {
//		fmt.Println("IsOriginalSaved")
//	}
//	if !picture.IsPreviewScaled {
//		fmt.Println("IsPreviewSaved")
//	}
//	if !picture.IsTextRecognized {
//		fmt.Println("IsTextRecognized")
//	}
//	if !picture.IsMobileScaled {
//		fmt.Println("IsMobileSaved")
//	}
//
//	picture.Status = Processing
//
//	return c.pictureRepo.UpdateImageHandle(picture)
//}
//
//func (c *Coordinator) handleError(picture Picture) {
//	now := time.Now()
//	picture.Attempts++
//	if picture.Attempts > c.maxAttemptsBeforeNotify {
//		err := c.notifier.Notify(picture)
//		if err != nil {
//			log.Println("failed notify developer")
//			log.Println(err, picture, "err")
//		}
//	}
//
//	picture.Status = Failed
//	picture.ExecuteAfter = now.Add(time.Duration(picture.Attempts*1) * time.Minute)
//	err := c.pictureRepo.UpdateImageHandle(picture)
//	if err != nil {
//		log.Println(err, picture, "save")
//	}
//}
