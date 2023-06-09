package controller

import (
	"GoRedisLearn/DB"
	"GoRedisLearn/RedisUtil"
	"GoRedisLearn/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var RCTX = context.Background()

// LikeBlog 实现点赞功能 使用sorted set 方便查询
func LikeBlog(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	db := DB.GetDB()
	blogId := c.PostForm("id")

	// 1 获取登录用户 TODO 为了方便这里直接指明用户
	var user model.TbUser
	user.Id = 2
	// 2 判断当前登录用户是否已经点赞
	key := "blog:liked:" + blogId
	sismember := rds.ZScore(RCTX, key, strconv.Itoa(user.Id))
	if sismember.Val() == 0 {
		// 3 未点赞，可以点赞
		// 3.1 数据库点赞加一
		tx := db.Debug().Table("tb_blog").Where("id=?", blogId).UpdateColumn("liked", gorm.Expr("liked+?", 1))
		fmt.Println("Tx:", tx)
		if tx.RowsAffected != 0 {
			// 更新成功
			// 3.2 保存用户到redis的scored set集合
			// zadd key score member
			// score 是时间戳 member是用户id
			rds.ZAdd(RCTX, key, redis.Z{Score: float64(time.Now().Unix()), Member: user.Id})
		}
	} else {
		// 4 已点赞 取消点赞
		// 4.1 数据库点赞数减一
		db.Table("tb_blog").Where("id=?", blogId).UpdateColumn("liked", gorm.Expr("liked-?", 1))
		// 4.2把用户从redis的set集合移除
		rds.ZRem(RCTX, key, user.Id)
	}
}

// LikeBlogTop5 获取前五个点赞的
func LikeBlogTop5(c *gin.Context) {
	rds := RedisUtil.RedisUtil
	key := "blog:liked:4"
	// 获取前五个点赞的
	zRange := rds.ZRange(RCTX, key, 0, 5)
	fmt.Println(zRange)
}

// 保存blog到数据库，以及推送消息到其粉丝
func SaveBolg(c *gin.Context) {
	// 当前用户id
	userId := 12

	// 1 获取登录用户
	blog := model.TbBlog{
		UserId:     userId,
		ShopId:     2,
		Title:      "cao",
		Images:     "jin",
		Content:    "bo",
		Liked:      "0",
		Comments:   104,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	db := DB.GetDB()
	rds := RedisUtil.RedisUtil

	// 2 保存探店笔记
	tx := db.Debug().Save(&blog)

	// 保存成功
	if tx.RowsAffected == 0 {
		fmt.Println("保存失败")
	}

	// 查询笔者所有粉丝
	var follows []model.TbFollow
	tx = db.Debug().Where("follow_user_id=?", userId).Find(&follows)

	for _, follow := range follows {
		// 获取粉丝id
		uid := follow.UserId
		// 推送
		// 将消息推送到粉丝收件箱
		key := "feed:" + strconv.Itoa(uid)
		rds.ZADD(key, redis.Z{Score: float64(time.Now().Unix()), Member: blog.Id})
	}
	//5 返回id
	fmt.Println("id:", blog.Id)
}

func QueryBlogOfFollow(c *gin.Context) {

	// 上次查询的最后一个值，便是此次查询的最大值
	max := c.PostForm("lastId")
	// 查询偏移量
	offset, _ := strconv.Atoi(c.PostForm("offset"))

	if offset == 0 {
		// 第一次来的时候 offset没有值 因此给予默认值
		offset = 0
	}
	rds := RedisUtil.RedisUtil
	db := DB.GetDB()

	// 1 获取当前用户 TODO 在此处直接指明用户id
	userId := 1
	// 2 查询收件箱 ZREVRANGEBYSCORE key Max Min LIMIT offset count
	key := "feed:" + strconv.Itoa(userId)
	scores := rds.ZRevRangeByScoreWithScores(RCTX, key, &redis.ZRangeBy{
		Min: "0", Max: max, Offset: int64(offset), Count: 3,
	})
	// TODO 这里应该要加一个非空判断

	// 3 解析数据 blogId minTime(时间戳）offset
	ids := make([]interface{}, 0, len(scores.Val()))
	fmt.Println("len:", len(scores.Val()))
	var minTime float64 = 0
	//offset
	os := 1
	for _, score := range scores.Val() {
		// 获取id
		id := score.Member
		ids = append(ids, id)
		// 获取分数（时间戳）
		t := score.Score
		if minTime == t {
			os++
		} else {
			minTime = t
			os = 1
		}
	}
	// 4 根据id查询blog
	fmt.Println("ids:", ids)
	var blog []model.TbBlog
	db.Debug().Where("id IN ?", ids).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{ids}, WithoutParentheses: true},
	}).Find(&blog)
	// 5 封装并返回
	for _, tbBlog := range blog {
		fmt.Println("blog:", tbBlog)
	}
}

func Caojinbo(c *gin.Context) {
	rds := RedisUtil.RedisUtil

	rds.SET("caojinbo", "lxh", 0)
}
