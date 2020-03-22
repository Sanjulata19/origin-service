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
	"github.com/nullserve/static-host/config"
	"go.uber.org/zap"
	"log"
	"net/http"
	"path"
	"strings"
)

type s3Service interface {
	GetObjectWithContext(ctx context.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
	ListObjectsV2PagesWithContext(ctx aws.Context, input *s3.ListObjectsV2Input, fn func(*s3.ListObjectsV2Output, bool) bool, opts ...request.Option) error
}

type server struct {
	config *config.StaticHost
	s3svc  s3Service
	logger *zap.Logger
}

type controlServer struct{}

var (
	errNoSuchSuffix = errors.New("host does not have suffix")
)

func (cs *controlServer) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(200)
	_, _ = rw.Write([]byte("ok."))
}

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
	deploymentId, err := hostToDeploymentId(r.Host, s.config.HostSuffix)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusBadGateway)
		// FIXME: write error response object and headers
		return
	}

	deploymentConfig, err := s.getDeploymentConfig(r.Context(), *deploymentId)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rr := responseRouter{
		deploymentId: *deploymentId,
		config:       deploymentConfig,
		server:       s,
		r:            r,
		rw:           rw,
	}
	rr.routeAndRespond()
	return
}

func (s *server) getDeploymentConfig(context context.Context, deploymentId string) (*siteConfig, error) {
	var err error
	s3Req := &s3.GetObjectInput{
		Bucket: &s.config.S3Source.BucketId,
		Key:    aws.String(path.Join(s.config.S3Source.SiteFolderPrefix, deploymentId, ".well-known", "nullserve.json")),
	}
	s3Res, err := s.s3svc.GetObjectWithContext(context, s3Req)
	// Default config
	cfg := siteConfig{Routes: []route{{UseFilesystem: aws.Bool(true)}}}
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok && aErr.Code() != s3.ErrCodeNoSuchKey {
			s.logger.Error("s3 service error",
				zap.String("error", aErr.Error()))
			return nil, errors.New("s3 service error")
		}
	} else {
		err = json.NewDecoder(s3Res.Body).Decode(&cfg)
		if err != nil {
			s.logger.Error("invalid config, failing",
				zap.String("error", err.Error()))
			return nil, errors.New("invalid config, failing")
		}
	}
	return &cfg, nil
}

func Main(cfg *config.StaticHost) {
	logger := zap.NewExample()
	logger.Info("starting Server")
	var err error
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	s3svc := s3.New(sess)
	srv := http.Server{
		Addr: ":80",
		Handler: &server{
			s3svc:  s3svc,
			logger: logger,
			config: cfg,
		},
	}
	cSrv := http.Server{
		Addr:    ":8080",
		Handler: &controlServer{},
	}
	go func() {
		err = srv.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	err = cSrv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
