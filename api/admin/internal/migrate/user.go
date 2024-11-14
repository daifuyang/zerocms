package migrate

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"zerocms/api/model/user"
	"zerocms/utils"
)

// 创建默认管理员

func createAdmin(userModel user.SysUserModel) (err error) {
	_, err = userModel.FindOne(context.Background(), 1)
	if err == sqlc.ErrNotFound {
		// 生成盐值
		var salt string
		salt, err = utils.GenerateSalt(16) // 16字节盐值
		if err != nil {
			return err
		}

		// 生成带盐的哈希
		hash, err := utils.GenerateHash("123456", salt)
		if err != nil {
			return err
		}

		newUserModel := &user.SysUser{
			LoginName: sql.NullString{
				String: "admin",
				Valid:  true, // 表示这是一个有效的字符串值
			},
			Password: sql.NullString{
				String: hash,
				Valid:  true,
			},
			Salt: sql.NullString{
				String: salt,
				Valid:  true,
			},
		}
		_, err = userModel.Insert(context.Background(), newUserModel)
		return err
	}
	return err
}
