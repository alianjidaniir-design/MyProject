package route

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) map[string]string {
	return mergeMaps(
		SetupUserRoute(app),
		SetupCourseRoutes(app),
		SetupEnrollmentRoute(app),
		SetupTeacherRoute(app),
		SetupDepartmentRoutes(app),
	)
}

func mergeMaps(maps ...map[string]string) map[string]string {
	mergeMap := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			mergeMap[k] = v
		}
	}
	return mergeMap
}
