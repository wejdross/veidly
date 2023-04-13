package instr

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sport/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileSection struct {
	Title   string
	Content string
}

type ProfileSectionArr []ProfileSection

// used to convert golang struct into db composite type representation
func (ud *ProfileSectionArr) Value() (driver.Value, error) {
	return json.Marshal(ud)
}

// used to convert DB composite type representation into golang struct
func (ud *ProfileSectionArr) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ud)
}

const MaxExtraProfileImgs = 6

func (ctx *Ctx) DalUpdateProfileImg(tx *sql.Tx, path string, uid uuid.UUID) error {
	const q = `update instructors set bg_img_path = $1 where user_id = $2`
	res, err := tx.Exec(q, path, uid)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalAddExtraProfileImg(tx *sql.Tx, path string, uid uuid.UUID) error {
	const q = `update instructors set 
		extra_img_paths = array_append(extra_img_paths, $1)
	where user_id = $2 and (extra_img_paths is null or cardinality(extra_img_paths) < $3)`
	res, err := tx.Exec(q, path, uid, MaxExtraProfileImgs)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

func (ctx *Ctx) DalRemoveExtraProfileImg(tx *sql.Tx, path string, uid uuid.UUID) error {
	const q = `update instructors set 
		extra_img_paths = array_remove(extra_img_paths, $1)
	where user_id = $2`
	res, err := tx.Exec(q, path, uid)
	if err != nil {
		return err
	}
	return helpers.PgMustBeOneRow(res)
}

/*
	403 if not an instructor
	413 if user hit his limit on images
	400 on invalid request
*/
func (ctx *Ctx) HandlerPostProfileImg() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		primary := g.Query("primary")

		i, err := ctx.DalReadInstructor(userID, UserID)
		if err != nil {
			g.AbortWithError(403, err)
			return
		}

		h, err := g.FormFile("image")
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		fpath, err := ctx.Static.ValidateImgAndGetRelpath(h, InstrProfileImgDir)
		if err != nil {
			g.AbortWithError(400, err)
			return
		}

		tx, err := ctx.Dal.Db.Begin()
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		oldImgPath := ""

		if primary != "" {
			if err := ctx.DalUpdateProfileImg(
				tx, fpath, userID,
			); err != nil {
				g.AbortWithError(500, err)
				tx.Rollback()
				return
			}
			oldImgPath = i.BgImgPath
		} else {
			if err := ctx.DalAddExtraProfileImg(
				tx, fpath, userID,
			); err != nil {
				if err == sql.ErrNoRows {
					g.AbortWithError(413, err)
				} else {
					g.AbortWithError(500, err)
				}
				tx.Rollback()
				return
			}
		}

		if err := ctx.Static.UpsertImage(h, fpath, oldImgPath); err != nil {
			g.AbortWithError(500, err)
			tx.Rollback()
			return
		}

		if err := tx.Commit(); err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatus(204)
	}
}

type DeleteProfileImgRequest struct {
	Path string
}

func (p *DeleteProfileImgRequest) Validate() error {
	if p.Path == "" {
		return fmt.Errorf("validate DeleteProfileImgRequest: invalid path")
	}
	return nil
}

func (ctx *Ctx) HandlerDeleteProfileImg() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		primary := g.Query("primary")
		var req DeleteProfileImgRequest

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		i, err := ctx.DalReadInstructor(userID, UserID)
		if err != nil {
			g.AbortWithError(403, err)
			return
		}

		tx, err := ctx.Dal.Db.Begin()
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		if primary != "" {
			if req.Path != i.BgImgPath {
				g.AbortWithError(403, fmt.Errorf("%s tried to delete invalid img: %s",
					userID.String(), req.Path))
				tx.Rollback()
				return
			}
			if err := ctx.DalUpdateProfileImg(
				tx, "", userID,
			); err != nil {
				g.AbortWithError(500, err)
				tx.Rollback()
				return
			}
		} else {
			found := false
			for _, p := range i.ExtraImgPaths {
				if p == req.Path {
					found = true
					break
				}
			}
			if !found {
				g.AbortWithError(403, fmt.Errorf("%s tried to delete invalid img: %s",
					userID.String(), req.Path))
				tx.Rollback()
				return
			}
			if err := ctx.DalRemoveExtraProfileImg(tx, req.Path, userID); err != nil {
				g.AbortWithError(500, err)
				tx.Rollback()
				return
			}
		}

		if err := ctx.Static.DeleteImg(req.Path); err != nil {
			g.AbortWithError(500, err)
			tx.Rollback()
			return
		}

		if err := tx.Commit(); err != nil {
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatus(204)
	}
}
