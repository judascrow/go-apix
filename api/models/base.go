package models

type SwagGetBase struct {
	Success bool   `json:"success" example:"true"`       // ผลการเรียกใช้งาน
	Status  int    `json:"status" example:"200"`         // HTTP Status Code
	Message string `json:"message" example:"Data Found"` // ข้อความตอบกลับ
}

type SwagCreateBase struct {
	Success bool   `json:"success" example:"true"`                 // ผลการเรียกใช้งาน
	Status  int    `json:"status" example:"201"`                   // HTTP Status Code
	Message string `json:"message" example:"Created Successfully"` // ข้อความตอบกลับ
}

type SwagUpdateBase struct {
	Success bool   `json:"success" example:"true"`                 // ผลการเรียกใช้งาน
	Status  int    `json:"status" example:"200"`                   // HTTP Status Code
	Message string `json:"message" example:"Updated Successfully"` // ข้อความตอบกลับ
}

type SwagDeleteBase struct {
	Success bool        `json:"success" example:"true"`                 // ผลการเรียกใช้งาน
	Status  int         `json:"status" example:"200"`                   // HTTP Status Code
	Message string      `json:"message" example:"Deleted Successfully"` // ข้อความตอบกลับ
	Data    interface{} `json:"data" `                                  // ข้อมูล
}

type SwagPageMeta struct {
	CurrentItemsCount int    `json:"currentItemsCount" example:"1"`
	CurrentPageNumber int    `json:"currentPageNumber" example:"1"`
	HasNextPage       bool   `json:"hasNextPage" example:"false"`
	HasPrevPage       bool   `json:"hasPrevPage" example:"false"`
	NextPageNumber    int    `json:"nextPageNumber" example:"1"`
	NextPageUrl       string `json:"nextPageUrl" example:"/api/v1/users/addresses?page=1&pageSize=10"`
	Offset            int    `json:"offset" example:"0"`
	PrevPageNumber    int    `json:"prevPageNumber" example:"1"`
	PrevPageUrl       string `json:"prevPageUrl" example:"/api/v1/users/addresses?page=1&pageSize=10"`
	RequestedPageSize int    `json:"requestedPageSize" example:"10"`
	TotalPagesCount   int    `json:"totalPagesCount" example:"1"`
	TotalItemsCount   int    `json:"totalItemsCount" example:"1"`
}

type SwagError400 struct {
	Success bool        `json:"success" example:"false"`       // ผลการเรียกใช้งาน
	Status  int         `json:"status" example:"400"`          // HTTP Status Code
	Message string      `json:"message" example:"Bad Request"` // ข้อความตอบกลับ
	Data    interface{} `json:"data" `                         // ข้อมูล
}

type SwagError404 struct {
	Success bool        `json:"success" example:"false"`     // ผลการเรียกใช้งาน
	Status  int         `json:"status" example:"404"`        // HTTP Status Code
	Message string      `json:"message" example:"Not Found"` // ข้อความตอบกลับ
	Data    interface{} `json:"data" `                       // ข้อมูล
}

type SwagError500 struct {
	Success bool        `json:"success" example:"false"`                 // ผลการเรียกใช้งาน
	Status  int         `json:"status" example:"500"`                    // HTTP Status Code
	Message string      `json:"message" example:"Internal Server Error"` // ข้อความตอบกลับ
	Data    interface{} `json:"data" `                                   // ข้อมูล
}

type SwagID struct {
	ID uint `json:"id" example:"1"` // ID
}

type SwagLogin struct {
	Success bool     `json:"success" example:"true"`                        // ผลการเรียกใช้งาน
	Status  int      `json:"status" example:"200"`                          // HTTP Status Code
	Message string   `json:"message" example:"User logged in successfully"` // ข้อความตอบกลับ
	Token   string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsI"`         // token
	Data    SwagUser `json:"data" `
}
