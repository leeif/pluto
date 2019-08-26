package avatar

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/leeif/pluto/config"
)

type Avatar struct {
	Bucket          string
	EndPoint        string
	AccessKeyID     string
	AccessKeySecret string
	CDN             string
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type AvatarReader struct {
	Reader io.ReadCloser
	Ext    string
}

func (avatar *Avatar) GetRandomAvatar() (*AvatarReader, error) {
	resp, err := http.Get("https://www.gravatar.com/avatar/" + randToken(8) + "?f=y&d=identicon")
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-type")
	if contentType == "image/png" {
		return &AvatarReader{resp.Body, "png"}, nil
	} else if contentType == "image/jpg" {
		return &AvatarReader{resp.Body, "jpg"}, nil
	}
	return nil, errors.New("Not image content type")
}

func (avatar *Avatar) SaveAvatarImageInS3(reader *AvatarReader) (string, error) {

	client, err := oss.New(avatar.EndPoint, avatar.AccessKeyID, avatar.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(avatar.Bucket)
	if err != nil {
		return "", err
	}
	osskey := fmt.Sprintf("avatar/%s.%s", randToken(8), reader.Ext)
	err = bucket.PutObject(osskey, reader.Reader)
	if err != nil {
		return "", err
	}
	url := ""
	if avatar.CDN == "" {
		url = fmt.Sprintf("https://%s.%s/%s", avatar.Bucket, avatar.EndPoint, osskey)
	} else {
		url = fmt.Sprintf("%s/%s", avatar.CDN, osskey)
	}
	return url, nil
}

func (avatar *Avatar) SaveAvatarImageInLocal(file io.ReadCloser) (string, error) {
	return "", nil
}

func NewAvatar() *Avatar {
	avatar := &Avatar{
		Bucket:          *config.GetConfig().Avatar.Bucket,
		EndPoint:        *config.GetConfig().Avatar.EndPoint,
		AccessKeyID:     *config.GetConfig().Avatar.AccessKeyID,
		AccessKeySecret: *config.GetConfig().Avatar.AccessKeySecret,
		CDN:             *config.GetConfig().Avatar.CDN,
	}
	return avatar
}
