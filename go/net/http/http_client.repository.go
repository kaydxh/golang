/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package http

import (
	"context"
	"fmt"

	context_ "github.com/kaydxh/golang/go/context"
	"github.com/kaydxh/golang/go/encoding/protojson"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/protobuf/proto"
)

func (r *Repository[REQ, RESP]) PostPbJson(ctx context.Context, req *REQ) (resp *RESP, err error) {
	return r.PostPbJsonWithUrl(ctx, r.Url, req)
}

func (r *Repository[REQ, RESP]) PostPbJsonWithUrl(ctx context.Context, url string, req *REQ) (resp *RESP, err error) {

	logger := logs_.GetLogger(ctx)
	tc := time_.New(true)
	summary := func() {
		tc.Tick("PostPbJson")
		respProto, ok := any(resp).(proto.Message)
		if ok {
			logger.WithField("response", reflect_.TruncateBytes(proto.Clone(respProto))).
				WithField("cost", tc.String()).
				Info("recv")
		}
	}
	defer summary()

	reqProto, ok := any(req).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("req is not proto message type")
	}
	logger.WithField("request", reflect_.TruncateBytes(proto.Clone(reqProto))).Info("send")

	reqData, err := protojson.Marshal(reqProto)
	if err != nil {
		logger.WithError(err).WithField("req", req).Errorf("failed to marshal request")
		return resp, err
	}

	var respData []byte
	err = time_.RetryWithContext(ctx, func(ctx context.Context) error {
		ctx, cancel := context_.WithTimeout(ctx, r.Timeout)
		defer cancel()

		respData, err = r.Client.PostJson(ctx, url, nil, reqData)
		if err != nil {
			logger.WithError(err).Errorf("failed to post json")
			return err
		}
		return nil
	}, r.RetryInterval, r.RetryTimes)
	if err != nil {
		return resp, err
	}

	var zeroResp RESP
	resp = &zeroResp
	err = protojson.Unmarshal(respData, any(resp).(proto.Message))
	if err != nil {
		logger.WithError(err).Errorf("failed to unmarshal post response data")
		return nil, err
	}

	return resp, nil
}

func (r *Repository[REQ, RESP]) PostPb(ctx context.Context, req *REQ) (resp *RESP, err error) {

	logger := logs_.GetLogger(ctx)
	tc := time_.New(true)
	summary := func() {
		tc.Tick("PostPbJson")
		respProto, ok := any(resp).(proto.Message)
		if ok {
			logger.WithField("response", reflect_.TruncateBytes(proto.Clone(respProto))).
				WithField("cost", tc.String()).
				Info("recv")
		}
	}
	defer summary()

	reqProto, ok := any(req).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("req is not proto message type")
	}
	logger.WithField("request", reflect_.TruncateBytes(proto.Clone(reqProto))).Info("recv")

	reqData, err := proto.Marshal(reqProto)
	if err != nil {
		logger.WithError(err).WithField("req", req).Errorf("failed to marshal request")
		return resp, err
	}

	var respData []byte
	err = time_.RetryWithContext(ctx, func(ctx context.Context) error {
		ctx, cancel := context_.WithTimeout(ctx, r.Timeout)
		defer cancel()

		respData, err = r.Client.PostPb(ctx, r.Url, nil, reqData)
		if err != nil {
			logger.WithError(err).Errorf("failed to post json, url: %v", r.Url)
			return err
		}
		return nil
	}, r.RetryInterval, r.RetryTimes)
	if err != nil {
		return resp, err
	}

	var zeroResp RESP
	resp = &zeroResp
	err = proto.Unmarshal(respData, any(resp).(proto.Message))
	if err != nil {
		logger.WithError(err).Errorf("failed to unmarshal post response data")
		return nil, err
	}

	return resp, nil
}

/*
func (r *Repository[REQ, RESP]) Do(ctx context.Context, method string, contentType string, req *REQ) (resp *RESP, err error) {

	logger := logs_.GetLogger(ctx)
	tc := time_.New(true)
	summary := func() {
		tc.Tick("PostPbJson")
		respProto, ok := any(resp).(proto.Message)
		if ok {
			logger.WithField("response", reflect_.TruncateBytes(proto.Clone(respProto))).
				WithField("cost", tc.String()).
				Info("recv")
		}
	}
	defer summary()

	reqProto, ok := any(req).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("req is not proto message type")
	}
	logger.WithField("request", reflect_.TruncateBytes(proto.Clone(reqProto))).Info("recv")

	var (
		reqData  []byte
		respData []byte
	)
	switch contentType {
	case binding.MIMEJSON:
		reqData, err = protojson.Marshal(reqProto)
	case binding.MIMEPROTOBUF:
		reqData, err = proto.Marshal(reqProto)

	default:
		reqData, err = protojson.Marshal(reqProto)
	}
	if err != nil {
		logger.WithError(err).WithField("req", req).Errorf("failed to marshal request")
		return resp, err
	}

	err = time_.RetryWithContext(ctx, func(ctx context.Context) error {
		ctx, cancel := context_.WithTimeout(ctx, r.Timeout)
		defer cancel()

		respData, err = r.Client.PostPb(r.Url, nil, reqData)
		if err != nil {
			logger.WithError(err).Errorf("failed to post json")
			return err
		}
		return nil
	}, r.RetryInterval, r.RetryTimes)
	if err != nil {
		return resp, err
	}

	var zeroResp RESP
	resp = &zeroResp
	err = proto.Unmarshal(respData, any(resp).(proto.Message))
	if err != nil {
		logger.WithError(err).Errorf("failed to unmarshal post response data")
		return nil, err
	}

	return resp, nil
}
*/
