package westhttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"unsafe"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/wasmCloud/west/bindings/wasi/http/types"
	"github.com/wasmCloud/west/bindings/wasi/io/poll"
	"github.com/wasmCloud/west/bindings/wasiext/http/ext"
)

func NewFields(h http.Header) types.Fields {
	headers := types.NewFields()
	for name, values := range h {
		for _, v := range values {
			headers.Append(
				types.FieldKey(name),
				types.FieldValue(cm.NewList(
					unsafe.SliceData([]byte(v)),
					uint(len(v)),
				)),
			)
		}
	}
	return headers
}

func NewOutgoingRequest(req *http.Request) (types.OutgoingRequest, func(func(poll.Pollable)) error, error) {
	if req.TLS != nil {
		return 0, nil, errors.New("`http.Request.TLS` is not currently supported")
	}
	res := types.NewOutgoingRequest(NewFields(req.Header))
	if s := req.URL.RequestURI(); s != "" {
		if res.SetPathWithQuery(cm.Some(s)) {
			return 0, nil, fmt.Errorf("failed to set path with query to `%s`", s)
		}
	}
	if s := req.URL.Hostname(); s != "" {
		if res.SetAuthority(cm.Some(s)) {
			return 0, nil, fmt.Errorf("failed to set authority to `%s`", s)
		}
	}
	if s := req.URL.Scheme; s != "" {
		switch s {
		case "http":
			if res.SetScheme(cm.Some(types.SchemeHTTP())) {
				return 0, nil, errors.New("failed to set scheme to HTTP")
			}
		case "https":
			if res.SetScheme(cm.Some(types.SchemeHTTPS())) {
				return 0, nil, errors.New("failed to set scheme to HTTPS")
			}
		default:
			if res.SetScheme(cm.Some(types.SchemeOther(s))) {
				return 0, nil, fmt.Errorf("failed to set scheme to `%s`", s)
			}
		}
	}

	switch req.Method {
	case "", http.MethodGet:
		if res.SetMethod(types.MethodGet()) {
			return 0, nil, errors.New("failed to set method to GET")
		}
	case http.MethodHead:
		if res.SetMethod(types.MethodHead()) {
			return 0, nil, errors.New("failed to set method to HEAD")
		}
	case http.MethodPost:
		if res.SetMethod(types.MethodPost()) {
			return 0, nil, errors.New("failed to set method to POST")
		}
	case http.MethodPut:
		if res.SetMethod(types.MethodPut()) {
			return 0, nil, errors.New("failed to set method to PUT")
		}
	case http.MethodPatch:
		if res.SetMethod(types.MethodPatch()) {
			return 0, nil, errors.New("failed to set method to PATCH")
		}
	case http.MethodDelete:
		if res.SetMethod(types.MethodDelete()) {
			return 0, nil, errors.New("failed to set method to DELETE")
		}
	case http.MethodConnect:
		if res.SetMethod(types.MethodConnect()) {
			return 0, nil, errors.New("failed to set method to CONNECT")
		}
	case http.MethodOptions:
		if res.SetMethod(types.MethodOptions()) {
			return 0, nil, errors.New("failed to set method to OPTIONS")
		}
	case http.MethodTrace:
		if res.SetMethod(types.MethodTrace()) {
			return 0, nil, errors.New("failed to set method to TRACE")
		}
	default:
		if res.SetMethod(types.MethodOther(req.Method)) {
			return 0, nil, fmt.Errorf("failed to set method to `%s`", req.Method)
		}
	}
	if req.Body == nil {
		return res, nil, nil
	}

	resBodyRes := res.Body()
	if err := resBodyRes.Err(); err != nil {
		return 0, nil, errors.New("failed to take outgoing request body")
	}
	resBody := resBodyRes.OK()
	resStreamRes := resBody.Write()
	if err := resStreamRes.Err(); err != nil {
		return 0, nil, errors.New("failed to take outgoing request body stream")
	}
	resStream := resStreamRes.OK()
	return res, func(poll func(poll.Pollable)) error {
		for {
			checkWriteRes := resStream.CheckWrite()
			if err := checkWriteRes.Err(); err != nil {
				if err.Closed() {
					slog.Debug("write stream closed")
					return io.EOF
				}
				return fmt.Errorf("failed to check write buffer capacity: %s", err.LastOperationFailed().ToDebugString())
			}
			wn := *checkWriteRes.OK()
			if wn == 0 {
				p := resStream.Subscribe()
				poll(p)
				p.ResourceDrop()
				continue
			}
			if wn > 4096 {
				wn = 4096
			}
			buf := make([]byte, wn)
			slog.Debug("reading buffer from body stream", "n", wn)
			rn, err := req.Body.Read(buf[:])
			if rn > 0 {
				slog.Debug("writing body stream chunk", "buf", buf[:rn])
				writeRes := resStream.Write(cm.NewList(unsafe.SliceData(buf), uint(rn)))
				if err := writeRes.Err(); err != nil {
					if err.Closed() {
						slog.Debug("write stream closed")
						return io.EOF
					}
					return fmt.Errorf("failed to write buffer: %s", err.LastOperationFailed().ToDebugString())
				}
			}
			if err == nil {
				continue
			}
			if err != io.EOF {
				return fmt.Errorf("failed read buffer from body stream: %w", err)
			}
			if err := req.Body.Close(); err != nil {
				return fmt.Errorf("failed to close request body: %w", err)
			}
			flushRes := resStream.Flush()
			if err := flushRes.Err(); err != nil {
				if err.Closed() {
					slog.Debug("write stream closed")
					return io.EOF
				}
				return fmt.Errorf("failed to flush body stream: %s", err.LastOperationFailed().ToDebugString())
			}
			resStream.ResourceDrop()

			trailers := cm.None[types.Fields]()
			if len(req.Trailer) > 0 {
				trailers = cm.Some(NewFields(req.Trailer))
			}
			slog.Debug("finishing outgoing body")
			resBodyFinishRes := types.OutgoingBodyFinish(*resBody, trailers)
			if err := resBodyFinishRes.Err(); err != nil {
				return fmt.Errorf("failed to finish outgoing body: %v", err)
			}
			return nil
		}
	}, nil
}

func NewIncomingRequest(req *http.Request) (types.IncomingRequest, func(func(poll.Pollable)) error, error) {
	res, write, err := NewOutgoingRequest(req)
	return ext.NewIncomingRequest(res), write, err
}

func NewIncomingResponse(resp types.IncomingResponse) (*http.Response, error) {
	header := make(http.Header, len(resp.Headers().Entries().Slice()))
	for _, h := range resp.Headers().Entries().Slice() {
		k := string(h.F0)
		header[k] = append(header[k], string(h.F1.Slice()))
	}
	bodyRes := resp.Consume()
	if err := bodyRes.Err(); err != nil {
		return nil, errors.New("failed to get response body")
	}
	body := bodyRes.OK()
	bodyStreamRes := body.Stream()
	if err := bodyStreamRes.Err(); err != nil {
		return nil, errors.New("failed to get response body stream")
	}
	bodyStream := bodyStreamRes.OK()
	var buf []byte
	for {
		bufRes := bodyStream.BlockingRead(4096)
		if err := bufRes.Err(); err != nil {
			if err.Closed() {
				break
			}
			return nil, fmt.Errorf("failed to read response body stream: %s", err.LastOperationFailed().ToDebugString())
		}
		buf = append(buf, bufRes.OK().Slice()...)
	}
	bodyStream.ResourceDrop()
	body.ResourceDrop()

	var trailer http.Header
	return &http.Response{
		Status:     http.StatusText(int(resp.Status())),
		StatusCode: int(resp.Status()),
		Body:       io.NopCloser(bytes.NewReader(buf)),
		Header:     header,
		Trailer:    trailer,
	}, nil
}

func HandleIncomingRequest[I, O ~uint32](f func(I, O), req *http.Request) (*http.Response, error) {
	wr, write, err := NewIncomingRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create new outgoing HTTP request: %w", err)
	}
	if write != nil {
		if err := write(poll.Pollable.Block); err != nil {
			return nil, fmt.Errorf("failed to write body: %w", err)
		}
	}

	out := ext.NewResponseOutparam()
	f(
		I(wr),
		O(out.F0),
	)
	out.F1.Subscribe().Block()
	respOptResRes := out.F1.Get()
	respResRes := respOptResRes.Some()
	if respResRes == nil {
		return nil, errors.New("response missing")
	}
	if err := respResRes.Err(); err != nil {
		return nil, errors.New("failed to get response")
	}

	respRes := respResRes.OK()
	if err := respRes.Err(); err != nil {
		return nil, fmt.Errorf("failed to receive response: %v", err)
	}
	resp, err := NewIncomingResponse(*respRes.OK())
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	resp.Request = req
	return resp, nil
}
