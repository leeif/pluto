package avatar

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	b64 "encoding/base64"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
)

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type AvatarGen struct {
}

type AvatarReader struct {
	Reader    io.ReadCloser
	Ext       string
	OriginURL string
}

func (ag *AvatarGen) GenFromBase64String(avatar string) (*AvatarReader, *perror.PlutoError) {
	b, err := b64.StdEncoding.DecodeString(avatar)
	if err != nil {
		return nil, perror.ServerError.Wrapper(fmt.Errorf("input is not in base64 format"))
	}
	ar := &AvatarReader{}
	ar.Reader = ioutil.NopCloser(bytes.NewReader(b))
	ar.Ext = ""
	ar.OriginURL = ""
	return ar, nil
}

func (ag *AvatarGen) GenFromGravatar() (*AvatarReader, *perror.PlutoError) {
	ar := &AvatarReader{}
	originURL := fmt.Sprintf("https://www.gravatar.com/avatar/%s?f=y&d=identicon", randToken(8))
	resp, err := http.Get(originURL)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
	ar.Reader = resp.Body
	ar.OriginURL = originURL
	contentType := resp.Header.Get("Content-type")
	if contentType == "image/png" {
		ar.Ext = "png"
	} else if contentType == "image/jpg" {
		ar.Ext = "jpg"
	} else {
		return nil, perror.ServerError.Wrapper(errors.New("Not support type of avatar " + contentType))
	}
	return ar, nil
}

func (as *AvatarSaver) SaveAvatarImageInOSS(reader *AvatarReader) (string, *perror.PlutoError) {

	if as.AccessKeyID == "" ||
		as.AccessKeySecret == "" ||
		as.Bucket == "" ||
		as.EndPoint == "" {
		return "", perror.ServerError.Wrapper(errors.New("aliyun config is not enough"))
	}

	client, err := oss.New(as.EndPoint, as.AccessKeyID, as.AccessKeySecret)
	if err != nil {
		return "", perror.ServerError.Wrapper(err)
	}
	bucket, err := client.Bucket(as.Bucket)
	if err != nil {
		return "", perror.ServerError.Wrapper(err)
	}
	if reader.Ext == "" {
		reader.Ext = "jpg"
	}
	osskey := fmt.Sprintf("avatar/%s.%s", randToken(8), reader.Ext)
	err = bucket.PutObject(osskey, reader.Reader)
	if err != nil {
		return "", perror.ServerError.Wrapper(err)
	}
	url := ""
	if as.CDN == "" {
		url = fmt.Sprintf("https://%s.%s/%s", as.Bucket, as.EndPoint, osskey)
	} else {
		url = fmt.Sprintf("%s/%s", as.CDN, osskey)
	}
	return url, nil
}

func (as *AvatarSaver) SaveAvatarImageInLocal(file io.ReadCloser) (string, error) {
	return "", nil
}

type AvatarSaver struct {
	Bucket          string
	EndPoint        string
	AccessKeyID     string
	AccessKeySecret string
	CDN             string
}

func NewAvatarSaver(config *config.Config) *AvatarSaver {
	as := &AvatarSaver{
		Bucket:          *config.Avatar.Bucket,
		EndPoint:        *config.Avatar.EndPoint,
		AccessKeyID:     *config.Avatar.AccessKeyID,
		AccessKeySecret: *config.Avatar.AccessKeySecret,
		CDN:             *config.Avatar.CDN,
	}
	return as
}
