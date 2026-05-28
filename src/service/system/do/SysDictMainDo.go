package service

import (
	"ego/src/boot/global/dy"
	"ego/src/cls"
	model "ego/src/model/basic"
	"ego/src/model/msg"
	"ego/src/model/response"
	"ego/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func RunSysDictList(c *gin.Context) {
	ptype := cls.MapSetIDSysGroup(c)
	c.HTML(200, "sysdict/list.htm", gin.H{
		"ptype": ptype,
	})
}

func RunSysDictSet(c *gin.Context) {
	db := utils.GetDB(c)

	name := strings.TrimSpace(c.Query("name"))
	id := cast.ToInt(c.Query("id"))
	bid := cast.ToInt(c.Query("bid"))

	var par model.SysDict

	// 💡 优雅重构：动态构建单次查询，不走两次 First
	query := db.Model(&model.SysDict{})
	if name != "" && id > 0 {
		query = query.Where("pname = ? OR id = ?", name, id)
	} else if name != "" {
		query = query.Where("pname = ?", name)
	} else if id > 0 {
		query = query.Where("id = ?", id)
	} else {
		response.OnFailure(c, msg.QueryParamErr)
		return
	}

	if err := query.First(&par).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogSqlErr(c, msg.SqlFindErr, err)
		}
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	var parValue model.SysDictValue
	if bid > 0 {
		if err := db.Where("id = ?", bid).First(&parValue).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("bid", bid))
			response.OnFailure(c, msg.SqlFindErr)
			return
		}
	}

	var renderValue any = nil

	if parValue.ID > 0 {
		renderValue = parValue
	} else {
		renderValue = map[string]any{}
	}
	c.HTML(http.StatusOK, "sysdict/set.htm", gin.H{
		"parid": par.ID,
		"info": map[string]any{
			"id": par.ID,
		},
		"parval": renderValue,
	})
}

// 列表接口
func RunSysDictGetdata(c *gin.Context) {
	db := utils.GetDB(c)
	var conds []any
	var exprs []string

	pname := strings.TrimSpace(c.Query("pname"))
	if len(pname) > 0 {
		exprs = append(exprs, "pname LIKE ?")
		conds = append(conds, "%"+pname+"%")
	}

	ptype := cast.ToInt(c.Query("ptype"))
	if ptype > 0 {
		exprs = append(exprs, "ptype = ?")
		conds = append(conds, ptype)
	}

	// 将条件切片打平合并到一个通用的 db 句柄中
	tx := db
	if len(exprs) > 0 {
		tx = tx.Where(strings.Join(exprs, " AND "), conds...)
	}

	// 2. 健壮的分页计算
	page := cast.ToInt(c.Query("page"))
	if page < 1 {
		page = 1
	}
	limit := cast.ToInt(c.Query("limit"))
	if limit < 1 || limit > 200 {
		limit = 10
	}
	start := (page - 1) * limit

	var count int64
	if err := tx.Model(&model.SysDict{}).Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}

	// 4. 获取分页列表数据（同样从 tx 派生，保证 Count 的副作用不会带到这里）
	var results []map[string]any
	if err := tx.Model(&model.SysDict{}).
		Order("id DESC").
		Offset(start).
		Limit(limit).
		Find(&results).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	// 5. 内存映射字典准备，避免循环内重复查配置表
	ptypeArr := cls.MapSetIDSysGroup(c)

	// 6. 遍历清洗并补充关联数据
	for key, value := range results {
		var csn int64

		_ = db.Model(&model.SysDictValue{}).Where("dict_id = ?", value["id"]).Count(&csn)
		results[key]["result"] = csn
		results[key]["ptypename"] = ptypeArr[cast.ToInt(value["ptype"])]

		var prn int64
		results[key]["tabnum"] = prn

		var tmparr = []map[string]any{}

		results[key]["list"] = tmparr
	}

	// 7. 统一规范返回
	response.OkTableList(results, count, c)
}

// 返回content里的数据列表
func RunSysDictGetSetdata(c *gin.Context) {

	id := cast.ToInt(c.Query("id"))
	if id < 1 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}
	db := utils.GetDB(c)
	// 1. 安全地查询字典主列表，并规整为通用的 []map[string]any
	var result []map[string]any
	if err := db.Session(&gorm.Session{}).
		Model(&model.SysDictValue{}).
		Where("dict_id = ?", id).
		Order("keyid ASC").
		Find(&result).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err, zap.Int("dict_id", id))
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	// 提取所有子项的 ID 集合
	var valIDs []int
	for _, val := range result {
		if valID := cast.ToInt(val["id"]); valID > 0 {
			valIDs = append(valIDs, valID)
		}
	}

	// 通过 Group 分组一次性拉回所有的计数映射
	svnCountMap := make(map[int]int64)
	if len(valIDs) > 0 {
		type SvnCountRow struct {
			DictValueId int   `gorm:"column:dict_value_id"`
			Total       int64 `gorm:"column:total"`
		}
		var rows []SvnCountRow

		// 执行一次聚合查询：SELECT dict_value_id, count(*) as total FROM ... WHERE dict_value_id IN (...) GROUP BY ...
		err := db.Session(&gorm.Session{}).
			Model(&model.SysDictValueSvn{}).
			Select("dict_value_id, count(*) as total").
			Where("dict_value_id IN ?", valIDs).
			Group("dict_value_id").
			Find(&rows).Error

		if err == nil {
			for _, row := range rows {
				svnCountMap[row.DictValueId] = row.Total
			}
		} else {
			utils.LogSqlErr(c, msg.SqlFindListErr, err)
		}
	}

	// 2. 内存清洗：映射前端标准字段（纯内存操作，效率极高）
	for i, value := range result {
		// 转换前端所需的 key 和 value 映射
		result[i]["key"] = value["keyid"]
		result[i]["value"] = value["defaultval"]

		// 从内存 map 直接读取统计数，默认找不到就是 0
		vID := cast.ToInt(value["id"])
		result[i]["svn"] = svnCountMap[vID]
	}

	count := int64(len(result))
	response.OkTableList(result, count, c)
}

func RunSysDictAdd(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))
	var info model.SysDict
	if id > 0 {
		if err := db.Where("id = ?", id).First(&info).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			response.OnFailure(c, msg.SqlFindErr)
			return
		}
	}
	ptype := cls.MapSetIDSysGroup(c)
	keyType := utils.GetParameter(c, "KeyType")
	c.HTML(200, "sysdict/add.htm", gin.H{
		"info":    info,
		"ptype":   ptype,
		"keyType": keyType,
	})
}

// RunSysDictDel 删除数据字典
func RunSysDictDel(c *gin.Context) {
	db := utils.GetDB(c)

	id := cast.ToInt(c.Query("id"))
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	var data model.SysDict
	if err := db.Where("id = ?", id).First(&data).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	// 开启事务，确保级联删除的原子性
	err := db.Transaction(func(tx *gorm.DB) error {
		// A. 删除主表（数据字典）
		if err := tx.Delete(&data).Error; err != nil {
			return err
		}

		// B. 级联删除子表（字典值表）
		if err := tx.Where("dict_id = ?", id).Delete(&model.SysDictValue{}).Error; err != nil {
			return err
		}

		// C. 级联删除SVN表（字典SVN）
		if err := tx.Where("dict_id = ?", id).Delete(&model.SysDictValueSvn{}).Error; err != nil {
			return err
		}

		return nil
	})

	// 事务执行结果拦截
	if err != nil {
		utils.LogSqlErr(c, msg.SqlDeleteErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlDeleteErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("删除数据字典 [%s]", data.Pname),
		map[string]any{"id": id},
		data,
		nil,
	)

	response.OnSuccess(c)
}

// 新增数据字典接口
func RunSysDictCreate(c *gin.Context) {
	var req struct {
		Pname   string `json:"pname" binding:"required"` // 字典名称必填
		Ptype   int    `json:"ptype" binding:"required"` // 系统组分类必填
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c) // 获取 WithContext 安全会话隔离的 DB 实例
	pname := strings.TrimSpace(req.Pname)

	// 3. 业务防重校验：检查名称是否已经存在
	var dupDict model.SysDict
	err := db.Where("pname = ?", pname).First(&dupDict).Error
	if err == nil {
		response.OnFailure(c, msg.NameExistsErr)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.String("pname", pname))
		response.OnFailure(c, msg.SqlErr)
		return
	}

	// 5. 实例化模型对象
	item := model.SysDict{
		Pname:   pname,
		Ptype:   req.Ptype,
		Remarks: strings.TrimSpace(req.Remarks),
	}

	// 6. 执行写入数据库
	if err := db.Create(&item).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCreateErr, err)
		response.OnFailure(c, msg.SqlCreateErr)
		return
	}

	// 7. 审计日志留痕（对齐模块规范：记录请求报文与新生成的结果对象）
	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("新增了数据字典 [%s]", item.Pname),
		req,
		nil,
		item,
	)

	// 8. 成功返回
	response.OnSuccess(c)
}

// 更新数据字典接口
func RunSysDictUpdate(c *gin.Context) {
	// 1. 定义专属更新的局部匿名结构体，ID 设为必填约束
	var req struct {
		ID      int    `json:"id" binding:"required"`    // 更新必须带上主键 ID
		Pname   string `json:"pname" binding:"required"` // 字典名称必填
		Ptype   int    `json:"ptype" binding:"required"` // 系统组分类必填
		Remarks string `json:"remarks"`
	}

	// 2. 拦截参数绑定错误
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c) // 获取 WithContext 安全会话隔离的 DB 实例
	pname := strings.TrimSpace(req.Pname)

	// 3. 防御性检查：验证原本数据是否存在
	var oldDict model.SysDict
	if err := db.Where("id = ?", req.ID).First(&oldDict).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", req.ID))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	if oldDict.Pname != req.Pname {
		response.OnFailure(c, "请不要修改名称,可另外新建")
		return
	}

	// 4. 业务防重校验：检查名称是否重复，但要排除当前编辑的 ID
	var dupDict model.SysDict
	err := db.Where("pname = ? AND id != ?", pname, req.ID).First(&dupDict).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.String("pname", pname))
			response.OnFailure(c, "系统繁忙，数据校验失败")
			return
		}
	} else {
		response.OnFailure(c, msg.NameExistsErr)
		return
	}

	// 6. 拼装更新数据映射 (使用 map 可以完美安全地向数据库写入零值或 NULL)
	updateData := map[string]any{
		"pname":   pname,
		"ptype":   req.Ptype,
		"remarks": strings.TrimSpace(req.Remarks),
	}

	if err := db.Model(&model.SysDict{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", req.ID))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	// 再次从 db 获取最新数据用于日志完整呈现
	var newDict model.SysDict
	_ = db.Where("id = ?", req.ID).First(&newDict)

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("更新了数据字典 [%s]", oldDict.Pname),
		req,
		oldDict, // 修改前的数据
		newDict, // 修改后的新数据
	)

	response.OnSuccess(c)
}

func RunSysDictSetPost(c *gin.Context) {
	var req struct {
		ID         int    `json:"id"`                            // 对应前端隐藏域 id
		DictID     int    `json:"dict_id"`                       // 对应前端隐藏域 dict_id
		KeyID      int    `json:"keyid" binding:"required"`      // 对应前端 key
		KeyStr     string `json:"keystr"`                        // 对应前端 keystr
		DefaultVal string `json:"defaultval" binding:"required"` // 对应前端 value
		State      int    `json:"state"`                         // 对应前端 state
	}

	// 2. 拦截 JSON 解析错误
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	inputVal := strings.TrimSpace(req.DefaultVal)
	inputKeyStr := strings.TrimSpace(req.KeyStr)

	// 3. 基本业务校验
	if req.KeyID < 1 {
		response.OnFailure(c, "KEY 不能为空或小于 1")
		return
	}
	if len(inputVal) < 1 {
		response.OnFailure(c, "VALUE 不能为空")
		return
	}
	state := 0
	if req.State < 2 {
		state = 1
	} else {
		state = 2
	}

	db := utils.GetDB(c)

	var one model.SysDict
	// 💡 提示：前端传过来的 dict_id 才是主表的 ID，前端的 id 是子表值表的 ID
	if err := db.Where("id = ?", req.DictID).Limit(1).Find(&one).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("dict_id", req.DictID))
		response.OnFailure(c, "数据字典校验异常")
		return
	}
	if one.ID < 1 {
		response.OnFailure(c, "主表数据不存在")
		return
	}

	var parValue model.SysDictValue
	if err := db.Where("dict_id = ? AND keyid = ?", req.DictID, req.KeyID).Limit(1).Find(&parValue).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err)
		response.OnFailure(c, "查询历史配置异常")
		return
	}

	var svn = dy.Map{}
	var newsvn = dy.Map{}
	isUpdate := parValue.ID > 0

	if isUpdate {
		svn["key"] = parValue.Keyid
		svn["value"] = parValue.Defaultval
		svn["str"] = parValue.Keystr

		newsvn["key"] = req.KeyID
		newsvn["value"] = inputVal
		newsvn["str"] = inputKeyStr
	}

	// 7. 写入/更新值表数据
	parValue.DictId = req.DictID
	parValue.Keyid = req.KeyID
	parValue.Defaultval = inputVal
	parValue.State = state
	parValue.Keystr = inputKeyStr

	if err := db.Save(&parValue).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err)
		response.OnFailure(c, "保存配置失败")
		return
	}

	if isUpdate {
		uid := utils.LibGetUid(c)
		var dsvn model.SysDictValueSvn
		dsvn.DictValueId = parValue.ID
		dsvn.Uptime = utils.GetTimestamp()
		dsvn.Post_uid = uid

		dsvn.Post_user = utils.LibGetUname(c)

		// 序列化快照
		oldJson, _ := json.Marshal(svn)
		dsvn.Oldstr = string(oldJson)
		newJson, _ := json.Marshal(newsvn)
		dsvn.Newstr = string(newJson)

		dsvn.DictId = req.DictID

		if err := db.Save(&dsvn).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlCreateErr, err)
		}
	}

	response.OnSuccess(c)
}

func RunSysDictSvnGetdata(c *gin.Context) {
	db := utils.GetDB(c) // 拿到的已经是安全的克隆实例
	parValueID := cast.ToInt(c.Query("parvalid"))

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))
	if limit < 1 || limit > 200 {
		limit = 10
	}
	start := 0
	if page > 1 {
		start = (page - 1) * limit
	}

	baseQuery := db.Model(&model.SysDictValueSvn{}).Where("dict_value_id = ?", parValueID)

	var count int64
	if err := baseQuery.Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err, zap.Int("parvalid", parValueID))
		response.OnFailure(c, msg.SqlCountErr)
		return
	}

	// 熔断保护：如果没有数据直接返回
	if count == 0 {
		response.OkTableList([]map[string]any{}, 0, c)
		return
	}

	var results []map[string]any
	if err := baseQuery.Order("id DESC").Offset(start).Limit(limit).Find(&results).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err, zap.Int("parvalid", parValueID))
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	for key, value := range results {
		if uptime, exists := value["uptime"]; exists {
			results[key]["ptime"] = utils.Pack.DateTime.TimestampToDate(uptime)
		} else {
			results[key]["ptime"] = ""
		}
	}

	response.OkTableList(results, count, c)
}

func RunSysDictGetSetdataTable(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))

	if id < 1 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	var par model.SysDict
	if err := db.Where("id = ?", id).Limit(1).Find(&par).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	var results []map[string]any
	var count int64

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))
	if limit < 1 || limit > 200 {
		limit = 30
	}
	start := 0
	if page > 1 {
		start = (page - 1) * limit
	}

	query := db.Model(&model.SysDictValue{}).Where("dict_id = ? AND state = 1", id)

	if err := query.Count(&count).Error; err != nil {
		utils.LogSqlErr(c, "查询静态字典值Count失败", err, zap.Int("id", id))
		response.OnFailure(c, "数据统计异常")
		return
	}

	if count > 0 {
		// 将静态字典的字段映射为统一的 key 和 value 返回给前端表格
		err := query.Select("keyid AS `key`, defaultval AS `value`, keystr").
			Offset(start).
			Limit(limit).
			Order("keyid ASC").
			Find(&results).Error

		if err != nil {
			utils.LogSqlErr(c, "查询静态字典值Find失败", err, zap.Int("id", id))
			response.OnFailure(c, "数据加载失败")
			return
		}
	}

	response.OkTableList(results, count, c)
}
