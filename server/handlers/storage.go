package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"connectrpc.com/connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RemoveFile implements [v1connect.RakerServerHandler].
func (server *RakerServer) RemoveFiles(ctx context.Context, request *v1.RemoveFilesRequest) (*v1.ScrapeResponse, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	var result db.History
	for _, file := range request.Paths {
		err := StorageHandler.Delete(user, PostTypePB2DB(request.Type), request.Owner, file)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		result, err = server.DBClient.UpdateHistoryRemoveFile(ctx, db.UpdateHistoryRemoveFileParams{
			File:      file,
			PostType:  PostTypePB2DB(request.Type),
			Post:      request.Post,
			PostOwner: request.Owner,
			Username:  user.Username,
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	if len(result.Files) == 0 {
		err := server.DBClient.HistoryRemove(ctx, db.HistoryRemoveParams{
			PostType:  PostTypePB2DB(request.Type),
			Post:      request.Post,
			PostOwner: request.Owner,
			Username:  user.Username,
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return &v1.ScrapeResponse{
		Categories: result.Categories,
		PostDate:   timestamppb.New(result.PostDate),
		PostType:   request.Type,
		Files:      result.Files,
		Post:       request.Post,
		PostOwner:  request.Owner,
		Incognito:  result.Incognito,
	}, nil
}

type storageHandler struct {
	root        string
	directories bool
	fileServer  http.Handler
	server      *RakerServer
	fileLocks   sync.Map // map[string]*sync.Mutex - per-file locks to avoid concurrent ops on same file
}

var StorageHandler storageHandler

func (server *RakerServer) NewStorageHandler(root string, directories bool) storageHandler {
	StorageHandler = storageHandler{
		root:        root,
		directories: directories,
		fileServer:  http.FileServer(http.Dir(root)),
		server:      server,
	}
	return StorageHandler
}

func (handler storageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("jwt")
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	u, err := handler.server.GetUserFromCookie(cookie)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if !strings.HasPrefix(request.URL.Path, fmt.Sprintf("/%s/", u.Username)) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	mediaPath := path.Join(handler.root, request.URL.Path)
	mediaPath = cleaner.Path(mediaPath)

	info, err := os.Stat(mediaPath)
	if handler.directories || (err == nil && !info.IsDir()) {
		writer.Header().Set("Cache-Control", "no-store, max-age=0, must-revalidate")
		writer.Header().Set("Pragma", "no-cache")
		writer.Header().Set("Expires", "0")
		handler.fileServer.ServeHTTP(writer, request)
	} else {
		if err != nil {
			escapedURL := cleaner.Line(request.URL.Path)
			log.Error(err, "URL", escapedURL)
		}
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (handler *storageHandler) Save(user db.User, media db.PostType, owner, fileName, URL string, cookies []*http.Cookie) error {
	filePath := path.Join(user.Username, string(media), owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

	_, err := os.Stat(mediaPath)
	if err == nil {
		return fmt.Errorf("file %s already exists", filePath)
	}

	directoryName := path.Dir(mediaPath)
	if _, err := os.Stat(directoryName); err != nil {
		const userGroupReadable = 660
		if err := os.MkdirAll(directoryName, userGroupReadable); err != nil {
			return err
		}
	}

	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}

	switch media {
	case db.PostTypeTiktok:
		request.Header.Add("Range", "bytes=0-")
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
		request.Header.Add("referer", "https://www.tiktok.com/")
	case db.PostTypeVsco:
		request.Header.Add("referer", "https://vsco.co/")
	}

	// request.Header.Add("referer", "https://www.instagram.com/")

	client := shared.NewClient(media == db.PostTypeVsco)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusClass := response.StatusCode / 100
	if statusClass == 4 || statusClass == 5 {
		return fmt.Errorf("response of %d instead of media", response.StatusCode)
	}

	if media == db.PostTypeVsco && response.Header.Get("Content-Type") == "video/MP2T" {
		log.Debugf("decoding stream from %s", media)
		return shared.Stream2MP4(response.Body, mediaPath)
	}

	file, err := os.Create(mediaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}

	log.Debugf("saved %s", filePath)
	return err
}

func (handler *storageHandler) SaveBundle(user db.User, media db.PostType, owner string, fileNames, URLs []string, cookies []*http.Cookie) ([]string, error) {
	if len(URLs) != len(fileNames) {
		return []string{}, errors.New("unequal length URLs & file names slices")
	}

	count := len(URLs)
	var wg sync.WaitGroup
	wg.Add(count)
	var mutex sync.Mutex
	var err error = nil

	for i := 0; i < count; i++ {
		URL := URLs[i]
		fileName := fileNames[i]
		go func(fileName, URL string, i int) {
			if err2 := handler.Save(user, media, owner, fileName, URL, cookies); err2 != nil {
				mutex.Lock()
				err = errors.Join(err, err2)
				fileNames[i] = ""
				mutex.Unlock()
			}
			wg.Done()
		}(fileName, URL, i)
	}

	wg.Wait()

	sucessfulFileNames := make([]string, 0, count)
	for _, fileName := range fileNames {
		if fileName != "" {
			sucessfulFileNames = append(sucessfulFileNames, fileName)
		}
	}

	return sucessfulFileNames, err
}

func (handler *storageHandler) Delete(user db.User, media db.PostType, owner, fileName string) error {
	filePath := path.Join(user.Username, string(media), owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

	// Lock the file to avoid deleting while another operation (crop/rotate) is running.
	muIface, _ := handler.fileLocks.LoadOrStore(mediaPath, &sync.Mutex{})
	mu := muIface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	_, err := os.Stat(mediaPath)
	if err != nil {
		return fmt.Errorf("file %s does not exists", filePath)
	}

	if err := os.Remove(mediaPath); err != nil {
		return err
	}
	log.Warnf("deleted %s", filePath)

	directoryName := path.Dir(mediaPath)
	files, err := os.ReadDir(directoryName)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		if err := os.Remove(directoryName); err != nil {
			return err
		}
		log.Warnf("deleted %s", directoryName)
	}

	return nil
}

func (handler *storageHandler) LocationEXIF(user db.User, media db.PostType, owner, fileName string) (float64, float64) {
	filePath := path.Join(user.Username, string(media), owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

	intfc, err := jpegstructure.NewJpegMediaParser().ParseFile(mediaPath)
	if err != nil {
		// log.Error(err)
		return 0, 0
	}

	sl := intfc.(*jpegstructure.SegmentList)
	rootIfd, _, err := sl.Exif()
	if err != nil {
		// log.Error(err)
		return 0, 0
	}

	gpsIfd, err := rootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)
	if err != nil {
		// log.Error(err)
		return 0, 0
	}

	gpsInfo, err := gpsIfd.GpsInfo()
	if err != nil {
		// log.Error(err)
		return 0, 0
	}

	return gpsInfo.Latitude.Decimal(), gpsInfo.Longitude.Decimal()
}

func (handler *storageHandler) transformJPEG(user db.User, media db.PostType, owner, fileName string, transform func(image.Image) (image.Image, error)) error {
	filePath := path.Join(user.Username, string(media), owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

	// Acquire per-file lock to prevent concurrent operations (crop/rotate/etc.) on the same file.
	// Use sync.Map with *sync.Mutex values to avoid global lock contention.
	muIface, _ := handler.fileLocks.LoadOrStore(mediaPath, &sync.Mutex{})
	mu := muIface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	tempFile, err := os.CreateTemp(path.Dir(mediaPath), fileName+".*.jpg")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	intfc, err := jpegstructure.NewJpegMediaParser().ParseFile(mediaPath)
	if err != nil {
		return err
	}

	sl := intfc.(*jpegstructure.SegmentList)
	exifData, _, _ := sl.Exif()

	file, err := os.Open(mediaPath)
	if err != nil {
		return err
	}

	source, err := jpeg.Decode(file)
	if closeErr := file.Close(); closeErr != nil {
		return closeErr
	}
	if err != nil {
		return err
	}

	transformed, err := transform(source)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, transformed, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	newIntfc, err := jpegstructure.NewJpegMediaParser().ParseBytes(buf.Bytes())
	if err != nil {
		return err
	}
	newSL := newIntfc.(*jpegstructure.SegmentList)

	if err := preserveExifMetadata(newSL, exifData); err != nil {
		log.Warnf("skipping EXIF preservation for %s: %v", filePath, err)
	}

	if err := newSL.Write(tempFile); err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempFile.Name(), mediaPath); err != nil {
		// On Windows, renaming over an existing file can fail. Retry after removing destination.
		if removeErr := os.Remove(mediaPath); removeErr != nil {
			return err
		}
		if err2 := os.Rename(tempFile.Name(), mediaPath); err2 != nil {
			return errors.Join(err, err2)
		}
	}

	return nil
}

func preserveExifMetadata(target *jpegstructure.SegmentList, source *exif.Ifd) (err error) {
	if source == nil {
		return nil
	}

	exifBuilder, err := buildSanitizedExifBuilder(source)
	if err != nil {
		return err
	}

	if err := target.SetExif(exifBuilder); err != nil {
		return err
	}

	return nil
}

func buildSanitizedExifBuilder(source *exif.Ifd) (builder *exif.IfdBuilder, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = fmt.Errorf("sanitize EXIF IFD [%s]: %v", source.IfdIdentity().UnindexedString(), state)
		}
	}()

	builder = exif.NewIfdBuilderWithExistingIfd(source)

	if thumbnailData, thumbnailErr := source.Thumbnail(); thumbnailErr == nil {
		if len(thumbnailData) > 0 {
			if err := builder.SetThumbnail(thumbnailData); err != nil {
				log.Warnf("dropping malformed EXIF thumbnail for %s: %v", source.IfdIdentity().UnindexedString(), err)
			}
		}
	} else if !errors.Is(thumbnailErr, exif.ErrNoThumbnail) {
		log.Warnf("dropping EXIF thumbnail for %s: %v", source.IfdIdentity().UnindexedString(), thumbnailErr)
	}

	children := source.Children()
	for index, entry := range source.Entries() {
		if entry.IsThumbnailOffset() || entry.IsThumbnailSize() {
			continue
		}

		if childIfdPath := entry.ChildIfdPath(); childIfdPath != "" {
			childIfd := findChildIfdForEntry(children, index)
			if childIfd == nil {
				log.Warnf("dropping EXIF child IFD with no matching child for %s tag 0x%04x", source.IfdIdentity().UnindexedString(), entry.TagId())
				continue
			}

			if childIfd.IfdIdentity().TagId() != 0xffff && childIfd.IfdIdentity().TagId() != entry.TagId() {
				log.Warnf("dropping EXIF child IFD with mismatched tag ID for %s tag 0x%04x", source.IfdIdentity().UnindexedString(), entry.TagId())
				continue
			}

			childBuilder, childErr := buildSanitizedExifBuilder(childIfd)
			if childErr != nil {
				log.Warnf("dropping malformed EXIF child IFD for %s tag 0x%04x: %v", source.IfdIdentity().UnindexedString(), entry.TagId(), childErr)
				continue
			}

			if err := builder.AddChildIb(childBuilder); err != nil {
				log.Warnf("dropping malformed EXIF child IFD for %s tag 0x%04x: %v", source.IfdIdentity().UnindexedString(), entry.TagId(), err)
				continue
			}
			continue
		}

		rawBytes, rawErr := entry.GetRawBytes()
		if rawErr != nil {
			log.Warnf("dropping malformed EXIF tag for %s tag 0x%04x: %v", source.IfdIdentity().UnindexedString(), entry.TagId(), rawErr)
			continue
		}

		bt := exif.NewBuilderTag(
			source.IfdIdentity().UnindexedString(),
			entry.TagId(),
			entry.TagType(),
			exif.NewIfdBuilderTagValueFromBytes(rawBytes),
			source.ByteOrder(),
		)

		if err := builder.Set(bt); err != nil {
			log.Warnf("dropping malformed EXIF tag for %s tag 0x%04x: %v", source.IfdIdentity().UnindexedString(), entry.TagId(), err)
			continue
		}
	}

	if nextIfd := source.NextIfd(); nextIfd != nil {
		nextBuilder, nextErr := buildSanitizedExifBuilder(nextIfd)
		if nextErr != nil {
			log.Warnf("dropping malformed EXIF sibling IFD for %s: %v", source.IfdIdentity().UnindexedString(), nextErr)
		} else if err := builder.SetNextIb(nextBuilder); err != nil {
			log.Warnf("dropping malformed EXIF sibling IFD for %s: %v", source.IfdIdentity().UnindexedString(), err)
		}
	}

	return builder, nil
}

func findChildIfdForEntry(children []*exif.Ifd, parentTagIndex int) *exif.Ifd {
	for _, childIfd := range children {
		if childIfd.ParentTagIndex() == parentTagIndex {
			return childIfd
		}
	}

	return nil
}

func (handler *storageHandler) CropImage(user db.User, media db.PostType, owner, fileName string, crop image.Rectangle) error {
	return handler.transformJPEG(user, media, owner, fileName, func(source image.Image) (image.Image, error) {
		bounds := source.Bounds()
		if crop.Min.X < bounds.Min.X || crop.Min.Y < bounds.Min.Y || crop.Max.X > bounds.Max.X || crop.Max.Y > bounds.Max.Y {
			return nil, fmt.Errorf("crop rectangle %v is outside image bounds %v", crop, bounds)
		}

		destination := image.NewRGBA(crop)
		draw.Draw(destination, destination.Bounds(), source, crop.Min, draw.Src)
		return destination, nil
	})
}

func (handler *storageHandler) RotateImage(user db.User, media db.PostType, owner, fileName string, amount int) error {
	return handler.transformJPEG(user, media, owner, fileName, func(source image.Image) (image.Image, error) {
		return imaging.Rotate(source, float64(amount), color.White), nil
	})
}

func (handler *storageHandler) RotateVideo(user db.User, media db.PostType, owner, fileName string, amount int, args ...string) error {
	filePath := path.Join(user.Username, string(media), owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)
	tempPath := strings.Replace(mediaPath, fileName, "_"+fileName, 1)

	transpose := "transpose=1"
	if amount == -90 {
		transpose = "transpose=2"
	}

	// ffmpeg -i input.mp4 -vf "transpose=1" -c:a copy output.mp4

	cmd := exec.Command("ffmpeg", "-i", mediaPath, "-vf", transpose, "-c:a", "copy", tempPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Rename(tempPath, mediaPath); err != nil {
		// On Windows, renaming over an existing file can fail. Retry after removing destination.
		if removeErr := os.Remove(mediaPath); removeErr != nil {
			return err
		}
		if err2 := os.Rename(tempPath, mediaPath); err2 != nil {
			return errors.Join(err, err2)
		}
	}
	return nil
}

// CropFile implements [v1connect.RakerServerHandler].
func (server *RakerServer) CropFile(ctx context.Context, request *v1.CropFileRequest) (*emptypb.Empty, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	crop := image.Rect(int(request.Corner1.X), int(request.Corner1.Y), int(request.Corner2.X), int(request.Corner2.Y))

	err := StorageHandler.CropImage(user, PostTypePB2DB(request.FileRequest.PostType), request.FileRequest.PostOwner, request.FileRequest.File, crop)
	if err != nil {
		log.Error(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return nil, nil
}

// RotateFile implements [v1connect.RakerServerHandler].
func (server *RakerServer) RotateFile(ctx context.Context, request *v1.RotateFileRequest) (*emptypb.Empty, error) {
	user, ok := ctx.Value(authenticatedUserKey).(db.User)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authenticated"))
	}

	switch {
	case strings.HasSuffix(request.FileRequest.File, ".jpg"):
	case strings.HasSuffix(request.FileRequest.File, ".jpeg"):
		err := StorageHandler.RotateImage(user, PostTypePB2DB(request.FileRequest.PostType), request.FileRequest.PostOwner, request.FileRequest.File, -int(request.Amount))
		if err != nil {
			log.Error(err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	case strings.HasSuffix(request.FileRequest.File, ".mp4"):
		err := StorageHandler.RotateVideo(user, PostTypePB2DB(request.FileRequest.PostType), request.FileRequest.PostOwner, request.FileRequest.File, int(request.Amount))
		if err != nil {
			log.Error(err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil, nil
}
