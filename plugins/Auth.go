package plugins

import (
	"github.com/dgrijalva/jwt-go"
	"nji"
	"reflect"
	"time"
)

type Auth struct {
	*Claims
}

func (pl *Auth) Exec(c *nji.Context) (err error) {
	token := c.Request.Header.Get("Authorization")[4:] // `JWT YWxhZGRpbjpvcGVuc2VzYW1l`
	pl.Claims, err = pl.ParseToken(token)
	return
}

func (pl *Auth) Support() nji.Method {
	return nji.MethodAny
}

func (pl Auth) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	return func(base nji.ViewAddr, c *nji.Context) {
		c.Error = (*Auth)(base.Offset(offset)).Exec(c)
	}
}

// JWT
var jwtSecret = ""

// Claim是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 根据用户的用户名和密码产生token
func (g *Auth) GenerateToken(username, password string) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func (g Auth) ParseToken(token string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
