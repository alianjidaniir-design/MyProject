package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/teacherSchema"
	"MyProject/models/teachers"
	"context"

	"github.com/gofiber/fiber/v2"
)

type TeacherRepository interface {
	// Create method
	Create(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.InformationSchema]) (res teacherSchema.TeacherSchema, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.PaginationSchema]) (res teacherSchema.ListSchema, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.GetTeacherSchema]) (res teacherSchema.DetailTeacherSchema, errStr string, code int, err error)
	SoftDelete(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]) (res teacherSchema.SoftDeleteTeacherSchema, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]) (res teacherSchema.UpdateTeacherSchema, errStr string, code int, err error)
	Login(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.LoginTeacherRequest], c *fiber.Ctx) (res teacherSchema.EntryStudentSchema, errStr string, code int, err error)
	RefreshToken(ctx context.Context, c *fiber.Ctx) (errStr string, code int, err error)
	Logout(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.LogoutSchema], c *fiber.Ctx) (res teacherSchema.EntryStudentSchema, errStr string, code int, err error)
	InfoTeacher(ctx context.Context, c *fiber.Ctx) (res teacherSchema.InfoTeacherSchema, errStr string, code int, err error)
}

var TeacherRepo TeacherRepository = teachers.GetRepo()
