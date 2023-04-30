package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"code":    0, //0 成功 -1失败
		"message": "查询成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")

	//盐值
	rand.Seed(time.Now().Unix()) //随机数种子 真随机数
	salt := fmt.Sprintf("%06d", rand.Int())

	data := models.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0 成功 -1失败
			"message": "用户名或密码不能为空",
			"data":    user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //0 成功 -1失败
			"message": "用户名已注册",
			"data":    user,
		})
		return
	}

	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1, //0 成功 -1失败
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}

	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0 成功 -1失败
		"message": "新增用户成功",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "用户id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0 成功 -1失败
		"message": "删除用户成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param id formData string false "用户id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //0 成功 -1失败
			"message": "修改参数不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, //0 成功 -1失败
			"message": "修改用户成功",
			"data":    user,
		})
	}
}

// FindUserByNameAndPwd
// @Summary 通过用户名和密码获取用户信息
// @Tags 用户模块
// @param name query string false "name"
// @param password query string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	//name := c.Query("name")
	//password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0 成功 -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	fmt.Println(user)

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //0 成功 -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    0, //成功
		"message": "登录成功",
		"data":    data,
	})
}

// 防止跨域站点的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}

	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息：", msg)
		//time.Now().Format("2022-01-28 17:55:54")
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// 新建群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//// LoadCommunity 加载群列表
//func LoadCommunity(c *gin.Context) {
//	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
//	//	name := c.Request.FormValue("name")
//	data, msg := models.LoadCommunity(uint(ownerId))
//	if len(data) != 0 {
//		utils.RespList(c.Writer, 0, data, msg)
//	} else {
//		utils.RespFail(c.Writer, msg)
//	}
//}
//
//// JoinGroups 加入群 userId uint, comId uint
//func JoinGroups(c *gin.Context) {
//	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
//	comId := c.Request.FormValue("comId")
//
//	//	name := c.Request.FormValue("name")
//	data, msg := models.JoinGroup(uint(userId), comId)
//	if data == 0 {
//		utils.RespOK(c.Writer, data, msg)
//	} else {
//		utils.RespFail(c.Writer, msg)
//	}
//}
//
//func FindByID(c *gin.Context) {
//	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
//
//	//	name := c.Request.FormValue("name")
//	data := models.FindByID(uint(userId))
//	utils.RespOK(c.Writer, data, "ok")
//}
//
