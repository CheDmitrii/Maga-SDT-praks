package task

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RoutesV1() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)          // GET /tasks
	r.Post("/", h.create)       // POST /tasks
	r.Get("/{id}", h.get)       // GET /tasks/{id}
	r.Put("/{id}", h.update)    // PUT /tasks/{id}
	r.Delete("/{id}", h.delete) // DELETE /tasks/{id}
	return r
}

func (h *Handler) RoutesV2() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.listWithFilter) // GET /tasks
	r.Post("/", h.create)        // POST /tasks
	r.Get("/{id}", h.get)        // GET /tasks/{id}
	r.Put("/{id}", h.update)     // PUT /tasks/{id}
	r.Delete("/{id}", h.delete)  // DELETE /tasks/{id}
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	all := h.repo.List()
	params := r.URL.Query()
	//pageV, isPage := params["q"]
	//limitV, isLimit := params["limit"]
	doneStr := params.Get("done")
	done, err := strconv.ParseBool(doneStr)
	var result []*Task
	if err == nil {
		result = lo.Filter(all, func(item *Task, _ int) bool {
			return item.Done == done
		})
	} else {
		result = all
	}

	pageStr := params.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err == nil {
		limitStr := params.Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 7
		}
		result = paginate(result, page, limit)
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) listWithFilter(w http.ResponseWriter, r *http.Request) {
	all := h.repo.List()
	params := r.URL.Query()
	//pageV, isPage := params["q"]
	//limitV, isLimit := params["limit"]
	doneStr := params.Get("done")
	done, err := strconv.ParseBool(doneStr)
	var result []*Task
	if err == nil {
		result = lo.Filter(all, func(item *Task, _ int) bool {
			return item.Done == done
		})
	} else {
		result = all
	}

	pageStr := params.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err == nil {
		limitStr := params.Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 7
		}
		result = paginate(result, page, limit)
	}
	writeJSON(w, http.StatusOK, result)
}

func paginate(items []*Task, page, pageSize int) []*Task {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // значение по умолчанию
	}

	start := (page - 1) * pageSize
	if start > len(items) {
		return []*Task{} // пустой срез, если страница вне диапазона
	}

	end := start + pageSize
	if end > len(items) {
		end = len(items) // не выйти за пределы среза
	}

	return items[start:end]
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	t, err := h.repo.Get(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

type createReq struct {
	Title string `json:"title"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}
	if len(req.Title) < 3 || len(req.Title) > 100 {
		httpError(w, http.StatusUnprocessableEntity, "invalid title")
		return
	}
	t := h.repo.Create(req.Title)
	writeJSON(w, http.StatusCreated, t)
}

type updateReq struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}
	t, err := h.repo.Update(id, req.Title, req.Done)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// helpers

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, "invalid id")
		return 0, true
	}
	return id, false
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
