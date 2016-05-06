package api

import "encoding/json"

func (r *Responder) DecodeRequest(v interface{}) error {
	if err := json.NewDecoder(r.Request.Body).Decode(v); err != nil {
		r.RespondWithError(NewInvalidJsonErr())
		return err
	}
	return nil
}
