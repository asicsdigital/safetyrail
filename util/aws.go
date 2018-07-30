package util

import (
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

func AwsInit() *session.Session {
	return session.Must(session.NewSession())
}

func GetS3Redirects(bucket string, sess *session.Session) []string {
	redirects := []string{}
	svc := s3.New(sess)

	prefix := viper.GetString("prefix")

	// list bucket contents
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1000),
		Prefix:  aws.String(prefix),
	}

	listResult, err := svc.ListObjects(input)

	if err != nil {
		jww.ERROR.Panic(err)
	}

	for i := range listResult.Contents {
		obj := listResult.Contents[i]
		key := *obj.Key
		// head the object
		input := &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}

		headResult, err := svc.HeadObject(input)

		if err != nil {
			jww.ERROR.Panic(err)
		}

		// get redirect metadata
		location := *headResult.WebsiteRedirectLocation

		// construct modrewrite rule
		if len(location) > 0 {
			u, err := url.Parse(location)

			if err != nil {
				jww.ERROR.Panic(err)
			}

			sourcePath := key
			protocol := u.Scheme
			hostname := u.Hostname()
			destPath := strings.TrimLeft(u.EscapedPath(), "/")

			redirects = append(redirects, FmtModRewrite(sourcePath, protocol, hostname, destPath))
		}
	}

	return redirects
}
