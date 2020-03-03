package static_host

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"path"
	"strings"
)

type s3Service interface {
	GetObjectWithContext(ctx context.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
	ListObjectsV2PagesWithContext(ctx aws.Context, input *s3.ListObjectsV2Input, fn func(*s3.ListObjectsV2Output, bool) bool, opts ...request.Option) error
}

type server struct {
	hostSuffix string
	s3Bucket   string
	s3svc      s3Service
}

var (
	errNoSuchSuffix = errors.New("host does not have suffix")
)

func hostToDeploymentId(host, suffix string) (*string, error) {
	if strings.HasSuffix(host, "."+suffix) {
		trimmed := strings.TrimSuffix(host, "."+suffix)
		return &trimmed, nil
	} else {
		return nil, errNoSuchSuffix
	}
}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err error
	deploymentId, err := hostToDeploymentId(r.Host, s.hostSuffix)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusBadGateway)
		// FIXME: write error response object and headers
		return
	}

	config, err := s.getDeploymentConfig(r.Context(), *deploymentId)
	if err != nil {
		return
		// FIXME
	}

	rr := responseRouter{
		deploymentId: *deploymentId,
		config:       config,
		server:       s,
		r:            r,
		rw:           rw,
	}
	rr.routeAndRespond()
	return
}

func (s *server) getDeploymentConfig(context context.Context, deploymentId string) (*config, error) {
	var err error
	s3Req := &s3.GetObjectInput{
		Bucket: &s.s3Bucket,
		Key:    aws.String(path.Join("site-deployments", deploymentId, ".well-known", "nullserve.json")),
	}
	s3Res, err := s.s3svc.GetObjectWithContext(context, s3Req)
	// Default config
	cfg := config{}
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok && aErr.Code() != s3.ErrCodeNoSuchKey {
			// FIXME: handle
			return nil, nil
		}
	} else {
		err = json.NewDecoder(s3Res.Body).Decode(&cfg)
		if err != nil {
			// FIXME: handle
			return nil, nil
		}
	}
	return &cfg, nil
}

func Main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	s3svc := s3.New(sess)
	srv := http.Server{
		Addr: ":80",
		Handler: &server{
			s3svc:    s3svc,
			s3Bucket: "nullserve-api-site-deployments20191125172523931100000001",
		},
	}
	_ = srv.ListenAndServe()
}
