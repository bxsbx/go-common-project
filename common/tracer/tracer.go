package tracer

import (
	"context"
	beegoConfig "github.com/astaxie/beego/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
	"log"
	"net/http"
	"time"
)

// 链路跟踪
type config struct {
	IsOpenTracing       bool
	ServiceName         string
	HostPort            string
	SamplerType         string
	SamplerParam        float64
	LogSpans            bool
	BufferFlushInterval time.Duration
}

const (
	TRACERSPANKEY = "tracerSpan"
)

var tracerCfg config

// 链路跟踪配置
func Config(cfg beegoConfig.Configer) {
	tracerCfg = config{
		ServiceName:         "test_tracer",
		HostPort:            "127.0.0.0:6831",
		SamplerType:         "const",
		SamplerParam:        1,
		LogSpans:            true,
		BufferFlushInterval: 1 * time.Second,
	}
}

// 创建全局跟踪器
func NewTracer() io.Closer {
	if !tracerCfg.IsOpenTracing {
		return nil
	}
	cfg := jaegerConfig.Configuration{
		ServiceName: tracerCfg.ServiceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  tracerCfg.SamplerType,
			Param: tracerCfg.SamplerParam,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            tracerCfg.LogSpans,
			LocalAgentHostPort:  tracerCfg.HostPort,
			BufferFlushInterval: tracerCfg.BufferFlushInterval,
		},
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("err:%v", err)
	} else {
		// 设置全局的一个跟踪器
		opentracing.SetGlobalTracer(tracer)
	}

	return closer
}

// new a span 并保存再上下文内容上，以便获取
func StarTracerSpan(appCtx context.Context, req *http.Request) {
	if !tracerCfg.IsOpenTracing {
		return
	}
	tracer := opentracing.GlobalTracer()
	opName := req.Host + req.URL.Path
	spanContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	var span opentracing.Span
	if err != nil {
		// 如果没找到上一个SpanCtx，表示是起点，则创建一个root tracer
		if err == opentracing.ErrSpanContextNotFound {
			span = tracer.StartSpan(opName)
		} else {
			log.Printf("err:%v", err)
		}
	} else {
		span = tracer.StartSpan(opName, opentracing.ChildOf(spanContext))
	}
	appCtx = context.WithValue(appCtx, TRACERSPANKEY, span)
}

// 向span写入请求头信息
func InjectTracerSpan(appCtx context.Context, reqHeader http.Header) {
	value := appCtx.Value(TRACERSPANKEY)
	if value == nil {
		return
	}
	span := value.(opentracing.Span)
	err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(reqHeader))
	if err != nil {
		log.Printf("err:%v", err)
	}
}

// 释放span
func FinishSpan(appCtx context.Context) {
	value := appCtx.Value(TRACERSPANKEY)
	if value == nil {
		return
	}
	span := value.(opentracing.Span)
	span.Finish()
}

// 标记span
func SetTagSpan(appCtx context.Context, key string, val interface{}) {
	value := appCtx.Value(TRACERSPANKEY)
	if value == nil {
		return
	}
	span := value.(opentracing.Span)
	span.SetTag(key, val)
}

// 链路跟踪上追加日志
func LogKV(appCtx context.Context, keyValues ...interface{}) {
	value := appCtx.Value(TRACERSPANKEY)
	if value == nil {
		return
	}
	span := value.(opentracing.Span)
	span.LogKV(keyValues...)
}
