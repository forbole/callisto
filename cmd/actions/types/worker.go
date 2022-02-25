package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ActionsWorker represents the worker that is used to handle Hasura actions queries
type ActionsWorker struct {
	mux     *http.ServeMux
	context *Context
}

// NewActionsWorker returns a new ActionsWorker instance
func NewActionsWorker(context *Context) *ActionsWorker {
	return &ActionsWorker{
		mux:     http.NewServeMux(),
		context: context,
	}
}

// RegisterHandler registers the provided handler to be used on each call to the provided path
func (w *ActionsWorker) RegisterHandler(path string, handler ActionHandler) {
	w.mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		// Set the content type
		writer.Header().Set("Content-Type", "application/json")

		// Read the body
		reqBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "invalid payload", http.StatusBadRequest)
			return
		}
		defer request.Body.Close()

		// Get the actions payload
		var payload Payload
		err = json.Unmarshal(reqBody, &payload)
		if err != nil {
			http.Error(writer, "invalid payload: failed to unmarshal json", http.StatusInternalServerError)
			return
		}

		// Handle the request
		res, err := handler(w.context, &payload)
		if err != nil {
			w.handleError(writer, err)
			return
		}

		// Marshal the response
		data, err := json.Marshal(res)
		if err != nil {
			w.handleError(writer, err)
			return
		}

		// Write the response
		writer.Write(data)
	})
}

// handleError allows to handle the given error by writing it to the provided writer
func (w *ActionsWorker) handleError(writer http.ResponseWriter, err error) {
	errorObject := GraphQLError{
		Message: err.Error(),
	}
	errorBody, err := json.Marshal(errorObject)
	if err != nil {
		panic(err)
	}

	writer.WriteHeader(http.StatusBadRequest)
	writer.Write(errorBody)
}

// Start starts the worker
func (w *ActionsWorker) Start(port uint) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), w.mux)
}
