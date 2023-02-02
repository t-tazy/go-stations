package handler

import (
	"context"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case "POST":
		var body model.CreateTODORequest
		// http requestを解析
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fmt.Printf("%v", err)
			return
		}
		// 簡易バリデーション
		if body.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("validation error")
			return
		}
		rsp, err := h.Create(ctx, &body)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

	case "PUT":
		var body model.UpdateTODORequest
		// http requestを解析
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fmt.Printf("%v", err)
			return
		}
		// 簡易バリデーション
		if body.ID == 0 || body.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("validation error")
			return
		}
		rsp, err := h.Update(ctx, &body)
		if err != nil {
			if errors.Is(err, &model.ErrNotFound{}) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Println("Not Foundのエラー")
				return
			}
			fmt.Printf("%v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("%v", err)
			return
		}
	case "GET":
		var body model.ReadTODORequest
		//  クエリパラメータ取得
		query := r.URL.Query()
		prevID, _ := strconv.Atoi(query.Get("prev_id"))
		body.PrevID = int64(prevID)
		size, _ := strconv.Atoi(query.Get("size"))
		if size == 0 {
			size = 5
		}
		body.Size = int64(size)

		rsp, err := h.Read(ctx, &body)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("%v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	fmt.Println(todo)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{*todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	fmt.Printf("get後のtodos確認: %+v", todos)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	fmt.Printf("update後のtodo確認: %+v", todo)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{*todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
