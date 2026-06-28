package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

const jobsIndex = "careercopilot_jobs"

type SearchService struct {
	es *elasticsearch.Client
}

func NewSearchService(es *elasticsearch.Client) *SearchService {
	return &SearchService{es: es}
}

type SearchQuery struct {
	Q              string  `form:"q"`
	Company        string  `form:"company"`
	Location       string  `form:"location"`
	Technology     string  `form:"technology"`
	SalaryMin      float64 `form:"salary_min"`
	SalaryMax      float64 `form:"salary_max"`
	ExperienceMin  float64 `form:"exp_min"`
	Remote         bool    `form:"remote"`
	EmploymentType string  `form:"type"`
	Page           int     `form:"page,default=1"`
	PerPage        int     `form:"per_page,default=20"`
}

type SearchResult struct {
	Total int64         `json:"total"`
	Hits  []interface{} `json:"hits"`
	Page  int           `json:"page"`
}

// IndexJob indexes a single job document in Elasticsearch
func (s *SearchService) IndexJob(ctx context.Context, job map[string]interface{}) error {
	jobID, ok := job["id"]
	if !ok {
		return fmt.Errorf("job missing id field")
	}

	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	resp, err := s.es.Index(
		jobsIndex,
		bytes.NewReader(data),
		s.es.Index.WithDocumentID(fmt.Sprintf("%v", jobID)),
		s.es.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch index error: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("elasticsearch returned error: %s", resp.Status())
	}
	return nil
}

// Search performs a full-text search across job fields
func (s *SearchService) Search(ctx context.Context, q *SearchQuery) (*SearchResult, error) {
	must := []interface{}{}
	filter := []interface{}{}

	// Full-text search
	if q.Q != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q.Q,
				"fields": []string{"title^3", "description", "skills^2", "company"},
				"type":   "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// Technology/skills filter
	if q.Technology != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"skills": strings.ToLower(q.Technology),
			},
		})
	}

	// Location filter
	if q.Location != "" {
		filter = append(filter, map[string]interface{}{
			"match": map[string]interface{}{
				"location": q.Location,
			},
		})
	}

	// Company filter
	if q.Company != "" {
		filter = append(filter, map[string]interface{}{
			"match": map[string]interface{}{
				"company": q.Company,
			},
		})
	}

	// Remote filter
	if q.Remote {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_remote": true,
			},
		})
	}

	// Salary filter
	if q.SalaryMin > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"salary_min": map[string]interface{}{
					"gte": q.SalaryMin,
				},
			},
		})
	}

	// Build query
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"sort": []interface{}{
			map[string]interface{}{"_score": "desc"},
			map[string]interface{}{"posted_at": "desc"},
		},
	}

	page := q.Page
	if page < 1 { page = 1 }
	perPage := q.PerPage
	if perPage < 1 || perPage > 100 { perPage = 20 }

	esQuery["from"] = (page - 1) * perPage
	esQuery["size"] = perPage

	data, err := json.Marshal(esQuery)
	if err != nil {
		return nil, err
	}

	resp, err := s.es.Search(
		s.es.Search.WithIndex(jobsIndex),
		s.es.Search.WithBody(bytes.NewReader(data)),
		s.es.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search error: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("elasticsearch returned error: %s", resp.Status())
	}

	var esResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&esResult); err != nil {
		return nil, err
	}

	hitsRaw, ok := esResult["hits"].(map[string]interface{})
	if !ok {
		return &SearchResult{Total: 0, Hits: []interface{}{}, Page: page}, nil
	}

	var total int64
	if totalObj, ok := hitsRaw["total"].(map[string]interface{}); ok {
		if v, ok := totalObj["value"].(float64); ok {
			total = int64(v)
		}
	}

	var hits []interface{}
	if hitsArr, ok := hitsRaw["hits"].([]interface{}); ok {
		for _, hit := range hitsArr {
			if h, ok := hit.(map[string]interface{}); ok {
				if src, ok := h["_source"]; ok {
					hits = append(hits, src)
				}
			}
		}
	}

	return &SearchResult{Total: total, Hits: hits, Page: page}, nil
}

// CreateIndex ensures the jobs index exists with the right mappings
func (s *SearchService) CreateIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"id": {"type": "long"},
				"title": {"type": "text", "analyzer": "english"},
				"description": {"type": "text", "analyzer": "english"},
				"company": {"type": "keyword"},
				"location": {"type": "text"},
				"skills": {"type": "keyword"},
				"is_remote": {"type": "boolean"},
				"is_hybrid": {"type": "boolean"},
				"salary_min": {"type": "float"},
				"salary_max": {"type": "float"},
				"experience_min": {"type": "float"},
				"experience_max": {"type": "float"},
				"employment_type": {"type": "keyword"},
				"posted_at": {"type": "date"},
				"created_at": {"type": "date"}
			}
		}
	}`

	resp, err := s.es.Indices.Create(
		jobsIndex,
		s.es.Indices.Create.WithBody(strings.NewReader(mapping)),
		s.es.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 400 = already exists, which is fine
	if resp.IsError() && resp.StatusCode != 400 {
		return fmt.Errorf("failed to create index: %s", resp.Status())
	}
	return nil
}

// Controller

type Controller struct{ svc *SearchService }

func NewController(svc *SearchService) *Controller { return &Controller{svc: svc} }

func (c *Controller) Search(ctx *gin.Context) {
	_, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }

	var q SearchQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		response.BadRequest(ctx, "invalid query parameters")
		return
	}

	result, err := c.svc.Search(ctx.Request.Context(), &q)
	if err != nil {
		// Fall back gracefully if ES is not available
		response.Success(ctx, "search results (fallback)", gin.H{
			"total": 0,
			"hits":  []interface{}{},
			"page":  q.Page,
		})
		return
	}
	response.Success(ctx, "search results", result)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	search := r.Group("/search")
	search.Use(middleware.JWTAuth(&cfg.JWT))
	{
		search.GET("", ctrl.Search)
	}
}
