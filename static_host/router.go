package static_host

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"path"
)

type responseRouter struct {
	deploymentId string
	config       *config
	server       *server
	rw           http.ResponseWriter
	r            *http.Request
}

func (rr *responseRouter) routeAndRespond() {
	rr.config.matchRule(rr.r.URL.Path)

	s3Req := &s3.GetObjectInput{
		Bucket: &rr.server.s3Bucket,
		Key:    aws.String(path.Join("site-deployments", path.Clean(path.Join(rr.deploymentId, rr.r.URL.Path)))),
	}

	s3Res, err := rr.server.s3svc.GetObjectWithContext(rr.r.Context(), s3Req)
	if err != nil {
		rr.rw.Write([]byte(fmt.Sprint(err)))
		return
	}
	bufR := bufio.NewReader(s3Res.Body)
	_, err = bufR.WriteTo(rr.rw)
	if err != nil {
		rr.rw.Write([]byte("it goofed"))
		return
	}
	return
}
