package conns

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

func traceHandler() func(*request.Request) {
	return func(r *request.Request) {
		tracer := opentracing.GlobalTracer()
		var sp opentracing.Span

		sp, ctx := opentracing.StartSpanFromContext(r.Context(), r.Operation.Name)

		ext.SpanKindRPCClient.Set(sp)

		ext.Component.Set(sp, "go-aws")
		ext.HTTPMethod.Set(sp, r.Operation.HTTPMethod)
		ext.HTTPUrl.Set(sp, r.HTTPRequest.URL.String())
		ext.PeerService.Set(sp, r.ClientInfo.ServiceName)
		if r.Context() == nil || fmt.Sprintf("%v", r.Context()) == "context.Background" {
			sp.SetTag("missing.parentspan", true)
		}

		r.SetContext(ctx)

		_ = inject(tracer, sp, r.HTTPRequest.Header)

		r.Handlers.Complete.PushBack(func(req *request.Request) {
			if req.HTTPResponse != nil {
				ext.HTTPStatusCode.Set(sp, uint16(req.HTTPResponse.StatusCode))
			} else {
				ext.Error.Set(sp, true)
			}
			sp.Finish()
		})

		r.Handlers.Retry.PushBack(func(req *request.Request) {
			sp.LogFields(log.String("event", "retry"))
		})
	}
}

func inject(tracer opentracing.Tracer, span opentracing.Span, header http.Header) error {
	return tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
}
