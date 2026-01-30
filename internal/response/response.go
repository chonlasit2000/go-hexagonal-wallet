package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code    int         `json:"code"`             // Business Code (เช่น 1000=Success, 999=Error)
	Message string      `json:"message"`          // ข้อความอ่านง่าย
	Data    interface{} `json:"data,omitempty"`   // ข้อมูล (ถ้ามี)
	Error   string      `json:"error,omitempty"`  // Error message (ถ้ามี)
	Errors  interface{} `json:"errors,omitempty"` // รายละเอียด Error (ถ้ามี)
}

type PagedResponse struct {
	Data       interface{} `json:"items"`
	Pagination Meta        `json:"pagination"`
}

type Meta struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int   `json:"total_page"`
}

// 1. Helper สำหรับตอบ Success แบบปกติ
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    1000, // รหัสที่เรากำหนดเองว่าแปลว่า "สำเร็จ"
		Message: "Success",
		Data:    data,
	})
}

// 2. Helper สำหรับตอบ Success แบบมี Pagination
func SuccessWithPage(c *fiber.Ctx, data interface{}, page, size, totalPage int, totalItem int64) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    1000,
		Message: "Success",
		Data: PagedResponse{
			Data: data,
			Pagination: Meta{
				Page:      page,
				Size:      size,
				TotalItem: totalItem,
				TotalPage: totalPage,
			},
		},
	})
}

// 3. Helper สำหรับตอบ Error
func Error(c *fiber.Ctx, httpStatus int, msg string, errStr string) error {
	return c.Status(httpStatus).JSON(Response{
		Code:    9999, // รหัส Generic Error
		Message: msg,
		Error:   errStr,
	})
}

// 4. Helper สำหรับตอบ Validation Error
func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Code:    400, // หรือจะใช้ domain.ErrBadRequest ถ้าเป็น int
		Message: "Validation Failed",
		Errors:  errors, // ส่ง Array เข้าไปตรงๆ เลย ไม่ต้อง Marshal เป็น string
	})
}
