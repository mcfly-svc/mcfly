package api

import "fmt"

func (handlers *Handlers) ProjectUpdateWebhook(r *Responder, ctx *RequestContext) {
	sourceProvider := *ctx.SourceProvider
	fmt.Printf("PROJECT UPDATE WEBHOOK FROM %s", sourceProvider.Key())
	r.RespondWithSuccess()
}
