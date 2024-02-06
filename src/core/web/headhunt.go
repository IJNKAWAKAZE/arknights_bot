package web

import (
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

var PoolUP = make(map[int]string)
var Pool = make(map[int]string)

func init() {
	PoolUP[7] = viper.GetString("gacha.pool_up_6_1")
	PoolUP[6] = viper.GetString("gacha.pool_up_6")
	PoolUP[5] = viper.GetString("gacha.pool_up_5")
	Pool[6] = viper.GetString("gacha.pool_6")
	Pool[5] = viper.GetString("gacha.pool_5")
	Pool[4] = viper.GetString("gacha.pool_4")
	Pool[3] = viper.GetString("gacha.pool_3")
}

func Headhunt(r *gin.Engine) {
	r.GET("/headhunt", func(c *gin.Context) {
		r6prob := 2.0
		r5prob := 8.0
		r4prob := 50.0
		r3prob := 40.0
		times := 0
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		key := fmt.Sprintf("headhunt:%d", userId)
		if utils.RedisIsExists(key) {
			times, _ = strconv.Atoi(utils.RedisGet(key))
		}
		var operators []utils.Operator
		for i := 0; i < 10; i++ {
			var operator utils.Operator
			name := ""
			autoProb(&r6prob, &r5prob, &r4prob, &r3prob, &times)
			allPro := r6prob + r5prob + r4prob + r3prob
			rankRand := float64(getRandomInt(1, int64(allPro)))
			if rankRand <= r6prob {
				name = randChar(6)
				reProb(&r6prob, &r5prob, &r4prob, &r3prob, &times)
			} else if rankRand <= r6prob+r5prob {
				name = randChar(5)
			} else if rankRand <= r6prob+r5prob+r4prob {
				name = randChar(4)
			} else if rankRand <= r6prob+r5prob+r4prob+r3prob {
				name = randChar(3)
			}
			char := utils.GetOperatorByName(name)
			operator.Profession = char.Profession
			operator.Rarity = char.Rarity
			operator.ThumbURL = char.ThumbURL
			operator.Name = char.Name
			operators = append(operators, operator)
			times++
		}
		utils.RedisSet(key, strconv.Itoa(times), 0)
		c.HTML(http.StatusOK, "Headhunt.tmpl", operators)
	})
}

// 自动调整6星概率
func autoProb(r6prob, r5prob, r4prob, r3prob *float64, times *int) {
	if *times > 50 {
		probUp := float64((*times - 49) * 2)
		probMultiple := (probUp - *r6prob) / (*r5prob + *r4prob + *r3prob)
		*r6prob = probUp
		*r5prob, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", *r5prob-probMultiple**r5prob), 64)
		*r4prob, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", *r4prob-probMultiple**r4prob), 64)
		*r3prob, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", *r3prob-probMultiple**r3prob), 64)
	}
	if *times >= 100 {
		reProb(r6prob, r5prob, r4prob, r3prob, times)
	}
}

// 随机数
func getRandomInt(min, max int64) int64 {
	r, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return r.Int64() + min
}

// 重置概率
func reProb(r6prob, r5prob, r4prob, r3prob *float64, times *int) {
	*r6prob = 2.0
	*r5prob = 8.0
	*r4prob = 50.0
	*r3prob = 40.0
	*times = 0
}

// 随机干员
func randChar(rank int) string {
	upWeight := 50
	if rank == 6 {
		upWeight = viper.GetInt("gacha.up_weight")
	}
	charaProb := int64(upWeight * 100)
	charRand := getRandomInt(1, 10000)
	getChar := ""
	if charRand <= charaProb && (rank == 5 || rank == 6) {
		charUpName := PoolUP[rank]
		getChar = randomChar(charUpName)
	} else {
		if PoolUP[7] != "" && rank == 6 {
			charRand := getRandomInt(1, 10000-charaProb)
			if charRand <= (10000-charaProb)/6*5 {
				charUpName := PoolUP[7]
				getChar = randomChar(charUpName)
			} else {
				charUpName := Pool[rank]
				getChar = randomChar(charUpName)
			}
		} else {
			charName := Pool[rank]
			getChar = randomChar(charName)
		}
	}
	return getChar
}

func randomChar(charStr string) string {
	chars := strings.Split(charStr, "/")
	r := getRandomInt(0, int64(len(chars)-1))
	return chars[r]
}
