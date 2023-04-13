package static

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"sport/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// append string to baseurl path
func (c *Config) AppendToBaseUrl(p string) string {
	u, err := url.Parse(c.BaseUrl)
	if err != nil {
		// already validated that this parse should succeed no matter what
		panic(err)
	}
	u.Path = path.Join(u.Path, p)
	return u.String()
}

func ValidateExt(ext string) error {
	var err error
	switch ext {
	case ".jpeg":
		break
	case ".jpg":
		break
	case ".png":
		break
	default:
		return fmt.Errorf("extension %s is not supported", ext)
	}
	return err
}

func ValidateImage(fh *multipart.FileHeader) error {
	buff := make([]byte, 512)
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Read(buff); err != nil {
		return err
	}

	mime := http.DetectContentType(buff)

	if !strings.HasPrefix(mime, "image") {
		return fmt.Errorf("invalid content type: " + mime)
	}

	return nil
}

func (ctx *Ctx) ValidateImgAndGetRelpath(
	h *multipart.FileHeader, prefix string,
) (string, error) {
	if h.Size == 0 {
		return "", fmt.Errorf("no image file was provided")
	}
	if err := ValidateImage(h); err != nil {
		return "", err
	}
	ext := path.Ext(h.Filename)
	if err := ValidateExt(ext); err != nil {
		return "", err
	}
	return path.Join(prefix, uuid.New().String()+ext), nil
}

func (ctx *Ctx) DeleteImg(relpath string) error {
	p := path.Join(ctx.Config.Basepath, relpath)
	if err := os.Remove(p); err != nil {
		return err
	}
	return nil
}

func (ctx *Ctx) UpsertImage(
	h *multipart.FileHeader,
	newImgRelPath, oldImgPath string) error {

	fullpath := path.Join(ctx.Config.Basepath, newImgRelPath)

	src, err := h.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return err
	}

	if oldImgPath != "" {
		path := path.Join(ctx.Config.Basepath, oldImgPath)
		if err := os.Remove(path); err != nil {
			fmt.Fprintf(os.Stderr,
				"UpsertImage: failed to remove %s, err was: %v\n", path, err)
		}
	}

	return nil
}

var allowedDirectories = map[string]bool{}

/*
	creates directory and adds it into
	"allowedDirectories" so it can be accessed by static server.
	this function is not thread safe.
*/
func (ctx *Ctx) RegisterDir(dir string) {
	if err := helpers.MkdirAllIfNotExists(
		path.Join(ctx.Config.Basepath, dir),
	); err != nil {
		panic(err)
	}
	if allowedDirectories[dir] {
		panic(fmt.Errorf("RegisterDir: %s already exists", dir))
	}
	allowedDirectories[dir] = true
}

func findPathTraversal(prms gin.Params) error {

	if _, ok := allowedDirectories[prms[0].Value]; !ok {
		goto badDirectory
	}

	{
		f := prms[1].Value
		parts := strings.Split(f, ".")
		if len(parts) != 2 {
			goto badFile
		}
		if _, err := uuid.Parse(parts[0]); err != nil {
			return err
		}
		if err := ValidateExt("." + parts[1]); err != nil {
			return err
		}
	}

	return nil

badDirectory:
	return errors.New("directory not found: " + prms[0].Value)

badFile:
	return errors.New("invalid file format")
}
