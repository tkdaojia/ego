package model

import (
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type SqlID struct {
	ID int `gorm:"primaryKey;autoIncrement;comment:主键id" json:"id"`
}

type SqlBase struct {
	ID        int            `gorm:"primaryKey;autoIncrement;comment:主键id" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SysGroup struct {
	SqlBase
	Gname    string `gorm:"comment:组名称" json:"gname"`
	Icon     string `gorm:"comment:图标" json:"icon"`
	Indextpl string `gorm:"comment:首页模板" json:"indextpl"`
	State    int    `gorm:"default:1;comment:状态" json:"state"`
	Ordnum   int    `gorm:"comment:排序号" json:"ordnum"`
}

type SysRole struct {
	SqlBase
	Name       string `gorm:"comment:唯一标识" json:"gname"`
	Rname      string `gorm:"comment:名称" json:"rname"`
	Remarks    string `gorm:"comment:备注" json:"remarks"`
	Menucate   string `gorm:"comment:菜单大类" json:"menucate"`
	Menulist   string `gorm:"comment:菜单列表" json:"menulist"`
	Modulelist string `json:"modulelist"`
	Sysgroup   string `gorm:"comment:系统组" json:"sysgroup"`
	Rolelist   string `json:"rolelist"`
}

type SysModule struct {
	SqlID
	Mname   string `json:"mname"`
	Remarks string `json:"remarks"`
	Pid     int    `json:"pid"`
}

type SysMenu struct {
	SqlID
	Mname  string `json:"mname"`
	Mlink  string `json:"mlink"`
	Icon   string `json:"icon"`
	Pid    int    `json:"pid"`
	Ordnum int    `json:"ordnum"`
	Status int    `json:"status"`
	Typeid int    `json:"typeid" gorm:"default:0"`
}

type SysOperationLog struct {
	SqlID
	Username    string    `gorm:"type:varchar(50);not null;default:''" json:"username"`
	Nickname    string    `gorm:"type:varchar(50);not null;default:''" json:"nickname"`
	Module      string    `gorm:"type:varchar(20);not null;default:''" json:"module"`
	Action      string    `gorm:"type:varchar(20);not null;default:''" json:"action"`
	Do          string    `gorm:"type:varchar(30);not null;default:''" json:"do"`
	Description string    `gorm:"type:varchar(255);not null;default:''" json:"description"`
	Url         string    `gorm:"type:varchar(255);not null;default:''" json:"url"`
	Method      string    `gorm:"type:varchar(10);not null;default:''" json:"method"`
	Ip          string    `gorm:"type:varchar(50);not null;default:''" json:"ip"`
	UserAgent   string    `gorm:"type:varchar(500);not null;default:''" json:"userAgent"`
	Param       string    `gorm:"type:text" json:"param"`
	DataOld     string    `gorm:"type:text" json:"dataOld"`
	DataNew     string    `gorm:"type:text" json:"dataNew"`
	Status      int       `gorm:"default:1" json:"status"`
	ErrorMsg    string    `gorm:"type:text" json:"errorMsg"`
	Latency     int64     `gorm:"type:bigint;not null;default:0" json:"latency"` // 耗时(ms)
	CreatedAt   time.Time `gorm:"index:idx_created_at" json:"createdAt"`
}

type SysDict struct {
	SqlID
	Sys     int    `json:"sys" gorm:"default:0"`
	Ptype   int    `json:"ptype" gorm:"default:null"`
	Pname   string `json:"pname" gorm:"default:null"`
	Remarks string `json:"remarks" gorm:"default:null"`
}

type SysDictValue struct {
	SqlID
	DictId     int    `json:"dict_id" gorm:"default:0"`
	Keyid      int    `json:"keyid" gorm:"default:0"`
	Defaultval string `json:"defaultval" gorm:"default:''"`
	Keystr     string `json:"keystr" gorm:"default:null"`
	State      int    `json:"state" gorm:"default:1"`
	Uptime     int64  `json:"uptime" gorm:"default:0"`
}

type SysDictValueSvn struct {
	SqlID
	DictValueId int    `json:"dict_value_id" gorm:"default:0"`
	Oldstr      string `json:"oldstr" gorm:"default:''"`
	Newstr      string `json:"newstr" gorm:"default:''"`
	Uptime      int64  `json:"uptime" gorm:"default:0"`
	Post_uid    int    `json:"post_uid" gorm:"default:0"`
	Post_user   string `json:"post_user"  gorm:"default:''"`
	DictId      int    `json:"dict_id" gorm:"default:0"`
}

type SysAccount struct {
	SqlID
	Account     string `json:"account"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	Truename    string `json:"truename"`
	Role_id     string `json:"role_id"`
	Sysgroup    string `json:"sysgroup"`
	Login_count int    `json:"login_count"`
	Lastlogin   int64  `json:"lastlogin"`
	Lastip      string `json:"lastip"`
	Status      int    `json:"status"`
}

type Online struct {
	ID     int    `json:"id"`
	Uid    int    `json:"uid"`
	Ip     string `json:"ip"`
	Active int64  `json:"active"`
	Status int    `json:"status"`
}

type File struct {
	ID         int    `json:"id"`
	Valid      int    `json:"valid"`
	Filename   string `json:"filename"`
	Filetype   string `json:"filetype"`
	Fileext    string `json:"fileext"`
	Filesize   int64  `json:"filesize"`
	Filesizecn string `json:"filesizecn"`
	Filepath   string `json:"filepath"`
	Thumbpath  string `json:"thumbpath"`
	Userid     int    `json:"userid"`
	Username   string `json:"username"`
	Addtimes   int64  `json:"addtimes"`
	Ip         string `json:"ip"`
	Web        string `json:"web"`
	Mtype      string `json:"mtype"`
	Tab_id     int    `json:"tab_id" gorm:"default:null"`
	Field_id   int    `json:"field_id" gorm:"default:null"`
	Mid        int    `json:"mid"`
	Downcount  int    `json:"downcount"`
}

type C_RedisCache struct {
	REDIS *redis.Client
	CACHE *freecache.Cache
}
