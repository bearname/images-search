package transport

import (
    "archive/zip"
    "aws_rekognition_demo/internal/domain"
    "aws_rekognition_demo/internal/domain/picture"
    "fmt"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

type PictureController struct {
    BaseController
    service    picture.Service
    downloader domain.Downloader
}

func NewPictureController(service picture.Service, downloader domain.Downloader) *PictureController {
    c := new(PictureController)
    c.service = service
    c.downloader = downloader
    return c
}

func (c *PictureController) DetectImageFromArchive() func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if (*r).Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        fmt.Println("read file")
        value := r.FormValue("eventId")
        if len(value) == 0 {
            log.Println("Failed get eventId")
            http.Error(w, "Failed get eventId", http.StatusBadRequest)
            return
        }

        eventId, err := strconv.Atoi(value)
        if err != nil {
            log.Println(err)
            http.Error(w, "Failed get eventId", http.StatusBadRequest)
            return
        }

        fileReader, header, err := r.FormFile("file")
        if err != nil {
            log.Println(err)
            http.Error(w, "Failed get file", http.StatusBadRequest)
            return
        }

        fmt.Println(header.Filename)
        fmt.Println(header.Size)
        fmt.Println(header.Header)

        f, err := os.OpenFile("./test.zip", os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            log.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        defer f.Close()
        _, err = io.Copy(f, fileReader)
        if err != nil {
            log.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        c.detectImage(w, "test", eventId)
    }
}

func (c *PictureController) DetectImageFromDropboxUrl() func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if (*r).Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        query := r.URL.Query()
        dropboxUrl := query.Get("path")

        if len(dropboxUrl) == 0 {
            log.Println("err")
            http.Error(w, "Failed get url to dropbox zip file", http.StatusBadRequest)
            return
        }

        eventIdValue := query.Get("eventId")
        if len(eventIdValue) == 0 {
            log.Println("Failed get eventId")
            http.Error(w, "Failed get eventId", http.StatusBadRequest)
            return
        }

        eventId, err := strconv.Atoi(eventIdValue)
        if err != nil {
            log.Println(err)
            http.Error(w, "Invalid event id", http.StatusBadRequest)
            return
        }

        start := time.Now()
        err = c.downloader.Download(dropboxUrl)
        if err != nil {
            log.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        dest := "test"
        src := "./test.zip"
        err = unzip(src, dest)

        if err != nil {
            fmt.Println(os.RemoveAll("./test.zip"))
            fmt.Println(RemoveContents("test"))
            log.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        c.detectImage(w, dest, eventId)
        fmt.Println(os.RemoveAll("./test.zip"))
        fmt.Println(RemoveContents("test"))
        end := time.Now()
        fmt.Println("\n\n\n\n\n\n\n\n\n\n\nSeconds:")
        fmt.Println(end.Sub(start).Seconds())
        fmt.Println("Minutes:")
        fmt.Println(end.Sub(start).Minutes())
    }
}

func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}

func unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()

    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        fpath := filepath.Join(dest, f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, f.Mode())
        } else {
            var fdir string
            if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
                fdir = fpath[:lastIndex]
            }

            err = os.MkdirAll(fdir, f.Mode())
            if err != nil {
                log.Fatal(err)
                return err
            }
            f, err := os.OpenFile(
                fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer f.Close()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func (c *PictureController) detectImage(w http.ResponseWriter, dest string, eventId int) {
    //zipReader, err := zip.NewReader(bytes.NewReader(*body), int64(len(*body)))
    //body = nil
    //if err != nil {
    //    log.Println(err)
    //    http.Error(w, err.Error(), http.StatusBadRequest)
    //    return
    //}
    err := c.service.DetectImageFromArchive(dest, 0, int64(eventId))
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (c *PictureController) SearchPictures() func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if (*r).Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        query := r.URL.Query()
        participantNumberStr := query.Get("number")
        participantNumber := picture.ValueNotSetted
        var err error
        if len(participantNumberStr) != 0 {
            participantNumber, err = strconv.Atoi(participantNumberStr)
            if err != nil {
                log.Println(err)
                http.Error(w, "Invalid 'number' query parameter. 'Number' must be number", 400)
                return
            }
        }

        confidenceParameter := query.Get("confidence")
        confidence := 85
        if len(confidenceParameter) != 0 {
            confidence, err = strconv.Atoi(confidenceParameter)
            if err != nil {
                log.Println(err)
                log.Println("Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]")
                http.Error(w, "Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]", 400)
                return
            }
            if confidence < 0 || confidence > 100 {
                log.Println("Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]")
                http.Error(w, "Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]", 400)
                return
            }
        }
        limitParameter := query.Get("limit")
        limit := 20
        if len(limitParameter) != 0 {
            limit, err = strconv.Atoi(limitParameter)
            if err != nil {
                log.Println(err)
                log.Println("Invalid 'limit' query parameter. 'limit' must be in range [0, 100]")
                http.Error(w, "Invalid 'limit' query parameter. 'limit' must be in range [0, 100]", 400)
                return
            }
            if limit < 0 || limit > 100 {
                log.Println("Invalid 'limit' query parameter. 'limit' must be in range [0, 100]")
                http.Error(w, "Invalid 'limit' query parameter. 'limit' must be in range [0, 100]", 400)
                return
            }
        }
        offsetParameter := query.Get("offset")
        offset := 0
        if len(offsetParameter) != 0 {
            offset, err = strconv.Atoi(offsetParameter)
            if err != nil {
                log.Println(err)
                log.Println("Invalid 'offset' query parameter. 'offset' must be in range [0, 100]")
                http.Error(w, "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]", 400)
                return
            }
            if offset < 0 {
                log.Println("Invalid 'offset' query parameter. 'offset' must be in range [0, 100]")
                http.Error(w, "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]", 400)
                return
            }
        }

        eventIdParameter := query.Get("eventId")
        eventId := picture.ValueNotSetted
        if len(eventIdParameter) != 0 {
            eventId, err = strconv.Atoi(eventIdParameter)
            if err != nil {
                log.Println(err)
                log.Println("Invalid 'eventId' query parameter. 'offset' must be in range [0, 100]")
                http.Error(w, "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]", 400)
                return
            }
            if eventId < 0 {
                log.Println("Invalid 'eventId' query parameter. 'eventId' must be more than 0")
                http.Error(w, "Invalid 'eventId' query parameter. 'eventId' must be more than 0", 400)
                return
            }
        }

        dto := picture.NewSearchPictureDto(participantNumber, confidence, eventId, domain.Page{Limit: limit, Offset: offset})
        searchDto, err := c.service.Search(dto)
        if err != nil {
            fmt.Println(err.Error())
            log.Println(err)
            http.Error(w, "Failed found searchDto", 400)
            return
        }

        fmt.Println(searchDto)
        if searchDto.Pictures == nil {
            searchDto.Pictures = make([]picture.SearchPictureItem, 0)
        }

        c.WriteJsonResponse(w, searchDto)
    }
}

func (c *PictureController) DeletePicture() func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if (*r).Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        vars := mux.Vars(r)
        idString := vars["id"]
        var err error
        if len(idString) == 0 {
            log.Println(err)
            http.Error(w, "Invalid 'number' query parameter. 'Number' must be number", 400)
            return
        }

        err = c.service.Delete(idString)
        if err != nil {
            fmt.Println(err.Error())
            log.Println(err)
            http.Error(w, "Failed found pictures", 400)
            return
        }

        c.WriteJsonResponse(w, responseWithoutData{1, "Success delete "})
    }
}
