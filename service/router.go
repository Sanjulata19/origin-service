package service

import (
	"bufio"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strconv"
)

type responseRouter struct {
	deploymentId string
	config       *appConfig
	server       *server
	rw           http.ResponseWriter
	r            *http.Request
}

func (rr *responseRouter) routeAndRespond() {
	var err error
	rr.server.logger.Info("routing", zap.String("path", rr.r.URL.Path))
	action, err := rr.config.matchRoute(rr.r.URL.Path)
	if err != nil {
		rr.server.logger.Error("no route match",
			zap.String("error", err.Error()))
		rr.rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rr.rw.WriteHeader(http.StatusBadGateway)
		// TODO: write a message to users
		//_, err = rr.rw.Write([]byte(""))
		return
	}

	s3Req := &s3.GetObjectInput{
		Bucket: &rr.server.config.S3Source.BucketId,
		Key:    aws.String(path.Join("site-deployments", path.Clean(path.Join(rr.deploymentId, action.Destination)))),
	}

	s3Res, err := rr.server.s3svc.GetObjectWithContext(rr.r.Context(), s3Req)
	if err != nil {
		// TODO: 500, 502, 503 here
		//_, err = rr.rw.Write([]byte(fmt.Sprint(err)))
		// TODO: ALERT HERE. can't do much else because we can't notify user
		return
	}

	// Write headers
	rr.rw.Header().Add("Content-Type", *s3Res.ContentType)
	rr.rw.Header().Add("Content-Length", strconv.FormatInt(*s3Res.ContentLength, 10))
	rr.rw.Header().Add("Server", "NullServe")
	rr.rw.Header().Add("Etag", *s3Res.ETag)
	rr.rw.Header().Add("Last-Modified", s3Res.LastModified.UTC().Format(http.TimeFormat))
	for _, header := range action.Headers {
		if header.Overwrite != nil && *header.Overwrite {
			rr.rw.Header().Set(header.Key, header.Value)
		}
		rr.rw.Header().Add(header.Key, header.Value)
	}
	rr.rw.WriteHeader(int(action.StatusCode))

	bufR := bufio.NewReader(s3Res.Body)
	_, err = bufR.WriteTo(rr.rw)
	if err != nil {
		// TODO: ALERT HERE. already too late to send an error to user
		return
	}
}
