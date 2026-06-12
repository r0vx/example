package main

import (
	"fmt"
	"math/rand"
	"time"

	"example/admin"
	"example/models"
)

// 随机生成中文名字
var firstNames = []string{"张", "王", "李", "赵", "刘", "陈", "杨", "黄", "周", "吴", "徐", "孙", "马", "朱", "胡", "郭", "何", "高", "林", "罗"}
var lastNames = []string{"伟", "芳", "娜", "敏", "静", "丽", "强", "磊", "洋", "艳", "勇", "军", "杰", "涛", "明", "超", "秀英", "华", "慧", "建"}

// 随机生成公司名
var companies = []string{"阿里巴巴", "腾讯", "百度", "字节跳动", "美团", "京东", "网易", "华为", "小米", "滴滴", "拼多多", "蚂蚁集团", "快手", "携程", "新浪", "搜狐", "优酷", "哔哩哔哩", "知乎", "微博"}

func randomName() string {
	return firstNames[rand.Intn(len(firstNames))] + lastNames[rand.Intn(len(lastNames))]
}

func randomCompany() string {
	return companies[rand.Intn(len(companies))]
}

func randomStatus() string {
	if rand.Float32() > 0.2 {
		return models.StatusActive
	}
	return models.StatusInactive
}

func randomDate() time.Time {
	// 过去两年内的随机日期
	days := rand.Intn(730)
	return time.Now().AddDate(0, 0, -days)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db := admin.ConnectDB()

	fmt.Println("开始插入 1000 条用户数据...")

	users := make([]models.User, 1000)
	for i := range 1000 {
		name := randomName()
		email := fmt.Sprintf("user%d@test.com", i+1)
		users[i] = models.User{
			Name:             name,
			Company:          randomCompany(),
			Status:           randomStatus(),
			RegistrationDate: randomDate(),
		}
		users[i].Account = email
	}

	// 批量插入
	result := db.CreateInBatches(users, 100)
	if result.Error != nil {
		fmt.Printf("插入失败: %v\n", result.Error)
		return
	}

	fmt.Printf("成功插入 %d 条用户数据\n", result.RowsAffected)
}
