package persistence

import (
	"mayfly-go/internal/tag/domain/entity"
	"mayfly-go/internal/tag/domain/repository"
	"mayfly-go/pkg/base"
	"mayfly-go/pkg/gormx"
)

type tagTreeRepoImpl struct {
	base.RepoImpl[*entity.TagTree]
}

func newTagTreeRepo() repository.TagTree {
	return &tagTreeRepoImpl{base.RepoImpl[*entity.TagTree]{M: new(entity.TagTree)}}
}

func (p *tagTreeRepoImpl) SelectByCondition(condition *entity.TagTreeQuery, toEntity any, orderBy ...string) {
	sql := "SELECT p.id, p.type, p.code, p.code_path, p.name, p.remark, p.create_time, p.creator, p.update_time, p.modifier FROM t_tag_tree p WHERE p.is_deleted = 0 "

	params := make([]any, 0)
	if condition.Name != "" {
		sql = sql + " AND p.name LIKE ?"
		params = append(params, "%"+condition.Name+"%")
	}
	if len(condition.Codes) > 0 {
		sql = sql + " AND p.code IN (?)"
		params = append(params, condition.Codes)
	}
	if len(condition.CodePaths) > 0 {
		sql = sql + " AND p.code_path IN (?)"
		params = append(params, condition.CodePaths)
	}
	if condition.CodePathLike != "" {
		sql = sql + " AND p.code_path LIKE ?"
		params = append(params, condition.CodePathLike+"%")
	}
	if condition.Type != 0 {
		sql = sql + " AND p.type = ?"
		params = append(params, condition.Type)
	}
	if len(condition.CodePathLikes) > 0 {
		sql = sql + " AND ("
		for i, v := range condition.CodePathLikes {
			if i == 0 {
				sql = sql + "p.code_path LIKE ?"
			} else {
				sql = sql + " OR p.code_path LIKE ?"
			}
			params = append(params, v+"%")
		}
		sql = sql + ")"
	}
	sql = sql + " ORDER BY p.type, p.code_path"
	gormx.GetListBySql2Model(sql, toEntity, params...)
}
