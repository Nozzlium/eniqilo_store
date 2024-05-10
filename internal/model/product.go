package model

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type OrderBy string

const (
	Asc  OrderBy = "asc"
	Desc OrderBy = "desc"
)

func (o OrderBy) IsValid() bool {
	switch o {
	case Asc, Desc:
		return true
	default:
		return false
	}
}

type BooleanString string

const (
	True  BooleanString = "true"
	False BooleanString = "false"
)

func (b BooleanString) IsValid() bool {
	switch b {
	case True, False:
		return true
	default:
		return false
	}
}

func (b BooleanString) ToBool() bool {
	return b == True
}

type SearchProductQuery struct {
	ID          string `query:"id"`
	Name        string `query:"name"`
	Category    string `query:"category"`
	SKU         string `query:"sku"`
	Price       string `query:"price"`
	CreatedAt   string `query:"createdAt"`
	IsAvailable string `query:"isAvailable"`
	InStock     string `query:"inStock"`
	Limit       int    `query:"limit"`
	Offset      int    `query:"offset"`
}

func (spq SearchProductQuery) BuildWhereClauseAndParams() ([]string, []interface{}) {
	var (
		sqlClause []string
		params    []interface{}
	)

	if spq.ID != "" {
		params = append(params, spq.ID)
		sqlClause = append(sqlClause, "id = $%d")
	}

	if spq.Name != "" {
		params = append(params, fmt.Sprintf("%%%s%%", spq.Name))
		sqlClause = append(sqlClause, "name ilike $%d")
	}

	if spq.Category != "" && ProductCategory(spq.Category).IsValid() {
		params = append(params, ProductCategory(spq.Category).ToDBEnumType())
		sqlClause = append(sqlClause, "category = $%d")
	}

	if spq.SKU != "" {
		params = append(params, spq.SKU)
		sqlClause = append(sqlClause, "sku = $%d")
	}

	if inStock := BooleanString(spq.InStock); inStock.IsValid() {
		params = append(params, 0)
		if inStock == True {
			sqlClause = append(sqlClause, "stock > $%d")
		} else {
			sqlClause = append(sqlClause, "stock = $%d")
		}
	}

	if isAvailable := BooleanString(spq.IsAvailable); isAvailable.IsValid() {
		params = append(params, isAvailable.ToBool())
		sqlClause = append(sqlClause, "is_available = $%d")
	}

	return sqlClause, params
}

func (spq SearchProductQuery) BuildPagination() (string, []interface{}) {
	var params []interface{}

	limit := 5
	offset := 0
	if spq.Limit > 0 {
		limit = spq.Limit
	}
	if spq.Offset > 0 {
		offset = spq.Offset
	}
	params = append(params, limit, offset)

	return "limit $%d offset $%d", params
}

func (spq SearchProductQuery) BuildOrderByClause() ([]string, []interface{}) {
	var (
		sqlClause []string
		params    []interface{}
	)

	if spq.CreatedAt != "" || OrderBy(spq.CreatedAt).IsValid() {
		params = append(params, OrderBy(spq.CreatedAt))
		sqlClause = append(sqlClause, "created_at $%d")
	}

	if spq.Price != "" || OrderBy(spq.Price).IsValid() {
		params = append(params, OrderBy(spq.Price))
		sqlClause = append(sqlClause, "price $%d")
	}

	return sqlClause, params
}

// product category
type ProductCategory string

const (
	Clothing    ProductCategory = "Clothing"
	Accessories ProductCategory = "Accessories"
	Footwear    ProductCategory = "Footwear"
	Beverages   ProductCategory = "Beverages"
)

func (pc ProductCategory) ToDBEnumType() string {
	switch pc {
	case Clothing:
		return "clothing"
	case Accessories:
		return "accessories"
	case Footwear:
		return "footwear"
	case Beverages:
		return "beverages"
	default:
		return ""
	}
}

func (pc ProductCategory) FromDBEnumType(enumType string) ProductCategory {
	switch enumType {
	case "clothing":
		return Clothing
	case "accessories":
		return Accessories
	case "footwear":
		return Footwear
	case "beverages":
		return Beverages
	default:
		return ""
	}
}

func (pc ProductCategory) IsValid() bool {
	switch pc {
	case Clothing, Accessories, Footwear, Beverages:
		return true
	default:
		return false
	}
}

type Product struct {
	Category    ProductCategory `json:"category"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   time.Time       `json:"deleted_at"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Notes       string          `json:"notes"`
	ImageURL    string          `json:"image_url"`
	Location    string          `json:"location"`
	IsAvailable bool            `json:"is_available"`
	Stock       int             `json:"stock"`
	Price       float64         `json:"price"`
	ID          uuid.UUID       `json:"id"`
}

func (p *Product) IsValid() bool {
	if nameLen := len(p.Name); nameLen < 1 || nameLen > 30 {
		return false
	}

	if skuLen := len(p.SKU); skuLen < 1 || skuLen > 30 {
		return false
	}

	if _, err := url.ParseRequestURI(p.ImageURL); err != nil {
		log.Println(err)
		return false
	}

	if !p.Category.IsValid() {
		return false
	}

	if notesLen := len(p.Notes); notesLen < 1 || notesLen > 200 {
		return false
	}

	if price := p.Price; price < 1 {
		return false
	}

	if stock := p.Stock; stock < 0 || stock > 100000 {
		return false
	}

	if len(p.Location) < 1 || len(p.Location) > 200 {
		return false
	}

	return true
}
