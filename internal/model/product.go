package model

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/util"
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
		sqlClause = append(
			sqlClause,
			"id = $%d",
		)
	}

	if spq.Name != "" {
		params = append(
			params,
			fmt.Sprintf(
				"%%%s%%",
				spq.Name,
			),
		)
		sqlClause = append(
			sqlClause,
			"name ilike $%d",
		)
	}

	if spq.Category != "" &&
		ProductCategory(
			spq.Category,
		).IsValid() {
		params = append(
			params,
			ProductCategory(
				spq.Category,
			).ToDBEnumType(),
		)
		sqlClause = append(
			sqlClause,
			"category = $%d",
		)
	}

	if spq.SKU != "" {
		params = append(params, spq.SKU)
		sqlClause = append(
			sqlClause,
			"sku = $%d",
		)
	}

	if inStock := BooleanString(spq.InStock); inStock.IsValid() {
		params = append(params, 0)
		if inStock == True {
			sqlClause = append(
				sqlClause,
				"stock > $%d",
			)
		} else {
			sqlClause = append(sqlClause, "stock = $%d")
		}
	}

	if isAvailable := BooleanString(spq.IsAvailable); isAvailable.IsValid() {
		params = append(
			params,
			isAvailable.ToBool(),
		)
		sqlClause = append(
			sqlClause,
			"is_available = $%d",
		)
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
	params = append(
		params,
		limit,
		offset,
	)

	return "limit $%d offset $%d", params
}

func (spq SearchProductQuery) BuildOrderByClause() []string {
	var sqlClause []string

	if spq.Price != "" ||
		OrderBy(spq.Price).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf("price %s", spq.Price),
		)
	}

	if spq.CreatedAt != "" ||
		OrderBy(
			spq.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf("created_at %s", spq.CreatedAt),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"created_at desc",
		)
	}

	return sqlClause
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

func (pc ProductCategory) FromDBEnumType(
	enumType string,
) ProductCategory {
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
	case Clothing,
		Accessories,
		Footwear,
		Beverages:
		return true
	default:
		return false
	}
}

type SearchProductResponse struct {
	Category    ProductCategory `json:"category"`
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Notes       string          `json:"notes"`
	ImageURL    string          `json:"imageUrl"`
	Location    string          `json:"location"`
	CreatedAt   string          `json:"createdAt"`
	IsAvailable bool            `json:"isAvailable"`
	Stock       int             `json:"stock"`
	Price       float64         `json:"price"`
}

func (spr *SearchProductResponse) FromProduct(product Product) {
	spr.Category = product.Category
	spr.CreatedAt = util.ToISO8601(product.CreatedAt)
	spr.Name = product.Name
	spr.SKU = product.SKU
	spr.Notes = product.Notes
	spr.ImageURL = product.ImageURL
	spr.Location = product.Location
	spr.IsAvailable = product.IsAvailable
	spr.Stock = product.Stock
	spr.Price = product.Price
	spr.ID = product.ID.String()
}

type Product struct {
	Category    ProductCategory `json:"category"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   time.Time       `json:"deletedAt"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Notes       string          `json:"notes"`
	ImageURL    string          `json:"imageUrl"`
	Location    string          `json:"location"`
	IsAvailable bool            `json:"isAvailable"`
	Stock       int             `json:"stock"`
	Price       float64         `json:"price"`
	CreatedBy   uuid.UUID       `json:"createdBy"`
	UpdatedBy   uuid.UUID       `json:"updatedBy"`
	DeletedBy   uuid.UUID       `json:"deletedBy"`
	ID          uuid.UUID       `json:"id"`
}

func (p *Product) IsValid() bool {
	if nameLen := len(p.Name); nameLen < 1 ||
		nameLen > 30 {
		return false
	}

	if skuLen := len(p.SKU); skuLen < 1 ||
		skuLen > 30 {
		return false
	}

	if _, err := url.ParseRequestURI(p.ImageURL); err != nil {
		log.Println(err)
		return false
	}

	if !p.Category.IsValid() {
		return false
	}

	if notesLen := len(p.Notes); notesLen < 1 ||
		notesLen > 200 {
		return false
	}

	if price := p.Price; price < 1 {
		return false
	}

	if stock := p.Stock; stock < 0 ||
		stock > 100000 {
		return false
	}

	if len(p.Location) < 1 ||
		len(p.Location) > 200 {
		return false
	}

	return true
}

func (p *Product) CompareAndUpdate(product Product) error {
	hasUpdatedData := false
	if product.Name != p.Name {
		hasUpdatedData = true
		p.Name = product.Name
	}
	if product.SKU != p.SKU {
		hasUpdatedData = true
		p.SKU = product.SKU
	}
	if product.Category != p.Category {
		hasUpdatedData = true
		p.Category = product.Category
	}
	if product.Notes != p.Notes {
		hasUpdatedData = true
		p.Notes = product.Notes
	}
	if product.ImageURL != p.ImageURL {
		hasUpdatedData = true
		p.ImageURL = product.ImageURL
	}
	if product.Location != p.Location {
		hasUpdatedData = true
		p.Location = product.Location
	}
	if product.IsAvailable != p.IsAvailable {
		hasUpdatedData = true
		p.IsAvailable = product.IsAvailable
	}
	if product.Stock != p.Stock {
		hasUpdatedData = true
		p.Stock = product.Stock
	}
	if product.Price != p.Price {
		hasUpdatedData = true
		p.Price = product.Price
	}

	if !hasUpdatedData {
		return fmt.Errorf("no data updated")
	}

	return nil
}
