package http_handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/service"
)

type LanguageHandler struct {
	languageService service.LanguageService
}

var (
	ErrMissingRequiredParams = errors.New("missing required params")
	ErrAlreadyExist          = errors.New("already exist")
)

func NewLanguageHandler(languageService service.LanguageService) *LanguageHandler {
	return &LanguageHandler{languageService: languageService}
}

func (h *LanguageHandler) HandleUpdateRequest(w http.ResponseWriter, r *http.Request) {
	post := &model.PostRequest{}
	resp := &model.Response{}

	// get and parse post data
	if err := r.ParseForm(); err != nil {
		WriteErrorResponse(w, resp, err)
		return
	}

	if err := json.Unmarshal([]byte(r.Form.Get("Data")), &post); err != nil {
		WriteErrorResponse(w, resp, err)
		return
	}

	for i := range post.Changes {
		ch, err := h.handleChanges(post.Changes[i])
		if err != nil {
			log.Println("error handle changes", err)
			if errors.Is(err, ErrAlreadyExist) || errors.Is(err, ErrMissingRequiredParams) {
				resp.IO.Result = -1
				resp.IO.Message += err.Error() + "\n"
			}
		}

		resp.Changes = append(resp.Changes, ch)
	}

	WriteResponse(w, resp)
}

func (h *LanguageHandler) HandleGetAllLanguage(w http.ResponseWriter, _ *http.Request) {
	languages := h.languageService.GetAllLanguage()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	jsonB, err := json.Marshal(languages)
	if err != nil {
		log.Println("Err", err)
	}

	ret := []byte("{\"Body\":[" + string(jsonB) + "]}")

	_, _ = w.Write(ret)
}

func (h *LanguageHandler) handleChanges(row *model.Changes) (changedRow *model.ChangedRow, err error) {
	changedRow = new(model.ChangedRow)
	changedRow.Id = row.Id
	changedRow.Color = "rgb(255, 255, 166)"
	defer func() {
		if err != nil {
			// error accured so highlight row with red color
			changedRow.Color = "rgb(255,0,0)"
		}
	}()
	var lang *model.Language
	switch {
	case row.Added == 1:
		lang, err = h.languageService.AddNewLanguage(row)
		if err != nil {
			return
		}

		// on success set id received from store
		changedRow.Added = 1
		changedRow.Id = strconv.FormatInt(int64(lang.Id), 10)
	case row.Deleted == 1:
		var id int64

		id, _ = strconv.ParseInt(row.Id, 10, 0)
		if err = h.languageService.DeleteLanguage(id); err != nil {
			return
		}

		changedRow.Deleted = 1
	case row.Changed == 1:
		if _, err = h.languageService.UpdateLanguage(row); err != nil {
			return
		}

		changedRow.Changed = 1
	}

	return
}

func WriteErrorResponse(w http.ResponseWriter, resp *model.Response, err error) {
	log.Println("Err", err)
	resp.IO.Result = -1
	resp.IO.Message = err.Error()

	// write response with error
	WriteResponse(w, resp)
}

// WriteResponse writes json response
func WriteResponse(w http.ResponseWriter, resp *model.Response) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", " application/json")
	w.WriteHeader(http.StatusOK)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("Err", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Println("Err", err)
	}
}
