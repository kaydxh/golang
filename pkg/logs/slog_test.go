package logs

import (
	"github.com/sirupsen/logrus"
	"log/slog"
	"os"
	"testing"
)

// Address 地址结构体
type Address struct {
	Province string
	City     string
	Street   string
	ZipCode  string
}

// Contact 联系方式结构体
type Contact struct {
	Phone   string
	Email   string
	Address *Address // 嵌套指针
}

// Person 人员结构体
type Person struct {
	Name    string
	Age     int
	Contact *Contact // 嵌套指针
}

// Company 公司结构体
type Company struct {
	Name      string
	CEO       *Person  // 嵌套指针
	Address   *Address // 嵌套指针
	EmployNum int
}

func TestSlog(t *testing.T) {
	// 创建一个 JSON 格式的 logger
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	//	Level: slog.LevelInfo,
	//}))

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 创建嵌套指针对象并赋值
	company := &Company{
		Name: "科技有限公司",
		CEO: &Person{
			Name: "张三",
			Age:  45,
			Contact: &Contact{
				Phone: "13800138000",
				Email: "zhangsan@example.com",
				Address: &Address{
					Province: "广东省",
					City:     "深圳市",
					Street:   "科技园南路1号",
					ZipCode:  "518000",
				},
			},
		},
		Address: &Address{
			Province: "广东省",
			City:     "深圳市",
			Street:   "高新技术产业园区",
			ZipCode:  "518057",
		},
		EmployNum: 500,
	}

	// 使用 slog 打印结构化日志
	logger.Info("公司信息",
		slog.String("company_name", company.Name),
		slog.Int("employee_count", company.EmployNum),
		slog.Group("ceo",
			slog.String("name", company.CEO.Name),
			slog.Int("age", company.CEO.Age),
			slog.Group("contact",
				slog.String("phone", company.CEO.Contact.Phone),
				slog.String("email", company.CEO.Contact.Email),
				slog.Group("address",
					slog.String("province", company.CEO.Contact.Address.Province),
					slog.String("city", company.CEO.Contact.Address.City),
					slog.String("street", company.CEO.Contact.Address.Street),
					slog.String("zipcode", company.CEO.Contact.Address.ZipCode),
				),
			),
		),
		slog.Group("company_address",
			slog.String("province", company.Address.Province),
			slog.String("city", company.Address.City),
			slog.String("street", company.Address.Street),
			slog.String("zipcode", company.Address.ZipCode),
		),
	)

	// 也可以使用 slog.Any 直接打印整个对象
	logger.Info("完整公司对象 slog ", slog.Any("company", company))
	logrus.Info("logrus 完整公司对象:", company)

	// 测试 nil 指针的情况
	emptyCompany := &Company{
		Name:      "空公司",
		CEO:       nil,
		Address:   nil,
		EmployNum: 0,
	}
	logger.Info("空公司信息", slog.Any("empty_company", emptyCompany))
}
