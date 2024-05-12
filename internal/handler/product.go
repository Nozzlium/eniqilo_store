package handler

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Search(ctx *fiber.Ctx) error {
	var query model.SearchProductQuery
	err := ctx.QueryParser(&query)
	if err != nil {
		log.Println(err)
		return err
	}

	products, err := h.productService.Search(ctx.Context(), query)
	if err != nil {
		log.Println(err)
		return err
	}

	response := make([]model.SearchProductResponse, 0, len(products))
	for _, product := range products {
		var r model.SearchProductResponse
		r.FromProduct(product)
		response = append(response, r)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func (h *ProductHandler) SearchForCustomer(ctx *fiber.Ctx) error {
	var query model.SearchProductQuery
	err := ctx.QueryParser(&query)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("unable to parse query: %v", err.Error()),
		})
	}

	// should only show product that have isAvailable == true
	query.IsAvailable = "true"
	products, err := h.productService.Search(ctx.Context(), query)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("failed to search products: %v", err.Error()),
		})
	}

	response := make([]model.SearchProductResponse, 0, len(products))
	for _, product := range products {
		var r model.SearchProductResponse
		r.FromProduct(product)
		response = append(response, r)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func (h *ProductHandler) Create(ctx *fiber.Ctx) error {
	var product model.Product
	err := ctx.BodyParser(&product)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("unable to process body: %v", err.Error()),
		})
	}

	if !product.IsValid() {
		return HandleError(ctx, ErrorResponse{
			message: "invalid body",
			error:   fmt.Errorf("invalid body"),
		})
	}

	id, createdAt, err := h.productService.Save(ctx.Context(), product)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("unable to save product: %v", err.Error()),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"data": fiber.Map{
			"id":        id,
			"createdAt": createdAt,
		},
	})
}

func (h *ProductHandler) Update(ctx *fiber.Ctx) error {
	var product model.Product
	id := ctx.Params("id")
	err := ctx.BodyParser(&product)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("unable to process body: %v", err.Error()),
		})
	}

	if !product.IsValid() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid body",
		})
	}

	err = h.productService.Update(ctx.Context(), id, product)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("unable to save product: %v", err.Error()),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

func (h *ProductHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := h.productService.Delete(ctx.Context(), id)
	if err != nil {
		return HandleError(ctx,
			ErrorResponse{
				message: "unable to delete product",
				error:   err,
				detail:  fmt.Sprintf("unable to delete product: %v", err.Error()),
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}
