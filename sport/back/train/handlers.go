package train

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"sport/helpers"
	"sport/static"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
GET /api/training

DESC
	get instructor training/s

REQUEST
	- Authorization: Bearer <token> in header
	- query string:
		- optional 'id' param. if provided will filter results searching for specified id
			(empty array will be returned if not found along with 200 OK status code)
		- optional 'q' param: if provided then search will be performed.
		  results will be ordered by relevance desc
		- optional 'group_ids' param: json array of group ids ([]string)

RETURNS
	- status code:
		- 200 on success
		- error status code
			- 400 when validation of query string fails
			- 500 on unexpected error
	- on success application/json body of type []TrainingResponse
*/
func (ctx *Ctx) HandlerGetTrainings() gin.HandlerFunc {
	return func(g *gin.Context) {

		var err error
		status := 500
		var ret []*TrainingWithJoins
		var userIDptr *uuid.UUID
		//userID := ctx.MustGet("UserID").(uuid.UUID)
		trainingIDstr := g.Query("id")
		var trainingIDptr *uuid.UUID
		var instructorIDptr *uuid.UUID
		var dcIDptr *uuid.UUID
		instructorIDstr := g.Query("instructor_id")
		var q string
		groupIdsStr := g.Query("group_ids")
		dcIDstr := g.Query("dc_id")
		var smIDptr *uuid.UUID
		smIDstr := g.Query("sm_id")

		var groupIds []string
		if groupIdsStr != "" {
			if err = json.Unmarshal([]byte(groupIdsStr), &groupIds); err != nil {
				status = 400
				goto end
			}
		}

		if trainingIDstr != "" {
			var trainingID uuid.UUID
			if trainingID, err = uuid.Parse(trainingIDstr); err != nil {
				status = 400
				goto end
			}
			trainingIDptr = &trainingID
		}

		if instructorIDstr != "" {
			var iid uuid.UUID
			if iid, err = uuid.Parse(instructorIDstr); err != nil {
				status = 400
				goto end
			} else {
				instructorIDptr = &iid
			}
		}

		if dcIDstr != "" {
			var id uuid.UUID
			id, err = uuid.Parse(dcIDstr)
			if err != nil {
				status = 400
				goto end
			}
			dcIDptr = &id
		}

		if smIDstr != "" {
			var id uuid.UUID
			id, err = uuid.Parse(smIDstr)
			if err != nil {
				status = 400
				goto end
			}
			smIDptr = &id
		}

		if g.GetHeader("Authorization") != "" {
			var userID uuid.UUID
			if userID, err = ctx.Api.AuthorizeUserFromCtx(g); err != nil {
				status = 401
				goto end
			} else {
				userIDptr = &userID
			}
		}

		q = g.Query("q")

		if userIDptr == nil && trainingIDptr == nil && q == "" && smIDptr == nil {
			err = fmt.Errorf("Anon cant scan trainings")
			status = 403
			goto end
		}

		if ret, err = ctx.DalReadTrainings(
			DalReadTrainingsRequest{
				TrainingID:   trainingIDptr,
				UserID:       userIDptr,
				Query:        q,
				InstructorID: instructorIDptr,
				WithOccs:     true,
				// only owned trainings
				WithGroups:    userIDptr != nil && instructorIDptr == nil,
				WithDcs:       userIDptr != nil && instructorIDptr == nil,
				WithSubModels: true,
				GroupIds:      groupIds,
				DcID:          dcIDptr,
				SmID:          smIDptr,
			}); err != nil {
			goto end
		}

		g.AbortWithStatusJSON(200, ret)
		return

	end:
		g.AbortWithError(status, err)
	}
}

/*

PUT /api/training/occ

DESC:
	override trainings occurrences

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type PostInstructorTrainingRequest

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 409 if occ already exist (same title)
			- 404 if user was not found
			- 400 when validation of request fails
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPutOccs() gin.HandlerFunc {
	return func(g *gin.Context) {
		userID := g.MustGet("UserID").(uuid.UUID)
		var req PutTrainingOccsRequest
		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, req.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}
		_, err := ctx.DalReadSingleTraining(DalReadTrainingsRequest{
			UserID:     &userID,
			TrainingID: &req.TrainingID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}
		occs := make([]Occ, len(req.Occurrences))
		soccs := make([]SecondaryOcc, 0, len(req.Occurrences))
		for i := 0; i < len(req.Occurrences); i++ {
			occs[i] = req.Occurrences[i].NewOcc(req.TrainingID)
			for j := 0; j < len(req.Occurrences[i].SecondaryOccs); j++ {
				soccs = append(
					soccs,
					req.Occurrences[i].SecondaryOccs[j].NewSecondaryOcc(occs[i].ID, req.TrainingID))
			}
		}

		tx, err := ctx.Dal.Db.BeginTx(context.Background(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		})
		if err != nil {
			g.AbortWithError(500, err)
			return
		}

		if err = ctx.DalDeleteTrainingOccs(req.TrainingID, tx); err != nil {
			tx.Rollback()
			g.AbortWithError(500, err)
			return
		}

		if err = ctx.DalCreateOccs(occs, tx); err != nil {
			tx.Rollback()
			g.AbortWithError(500, err)
			return
		}

		if err = ctx.DalCreate2Occs(soccs, tx); err != nil {
			tx.Rollback()
			g.AbortWithError(500, err)
			return
		}

		if err = tx.Commit(); err != nil {
			//tx.Rollback()
			g.AbortWithError(500, err)
			return
		}

		g.AbortWithStatus(204)
	}
}

/*
POST /api/training

DESC
	create new instructor training

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type CreateTrainingRequest

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 409 if training already exist (same title)
			- 403 if user is not an instructor
			- 400 when validation of request fails
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPostTraining() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400
		var req CreateTrainingRequest
		var t *Training
		var occs []Occ
		var soccs []SecondaryOcc
		userID := g.MustGet("UserID").(uuid.UUID)
		var tx *sql.Tx

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &req, func() error {
				return req.Validate(ctx)
			},
		); err != nil {
			goto end
		}

		if t, err = req.Training.NewTraining(ctx, userID); err != nil {
			status = 403
			goto end
		}

		occs = make([]Occ, len(req.Occurrences))
		soccs = make([]SecondaryOcc, 0, len(req.Occurrences))
		for i := 0; i < len(req.Occurrences); i++ {
			occs[i] = req.Occurrences[i].NewOcc(t.ID)
			for j := 0; j < len(req.Occurrences[i].SecondaryOccs); j++ {
				soccs = append(
					soccs,
					req.Occurrences[i].SecondaryOccs[j].NewSecondaryOcc(occs[i].ID, t.ID))
			}
		}

		tx, err = ctx.Dal.Db.BeginTx(context.Background(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		})
		if err != nil {
			status = 500
			goto end
		}

		if err = ctx.DalCreateTraining(t, userID, tx); err != nil {
			tx.Rollback()
			if helpers.PgIsUqViolation(err) {
				status = 409
			} else {
				status = 500
			}
			goto end
		}

		if err = ctx.DalCreateOccs(occs, tx); err != nil {
			tx.Rollback()
			status = 500
			goto end
		}

		if err = ctx.DalCreate2Occs(soccs, tx); err != nil {
			tx.Rollback()
			status = 500
			goto end
		}

		if err = tx.Commit(); err != nil {
			//tx.Rollback()
			status = 500
			goto end
		}

		if req.ReturnID {
			fmt.Fprintf(g.Writer, "%s", t.ID)
			g.AbortWithStatus(200)
		} else {
			g.AbortWithStatus(204)
		}
		return

	end:
		g.AbortWithError(status, err)
	}
}

/*
PATCH /api/training

DESC
	modify instructor training info
	<<note that this will perform full merge on instructor table with provided type>>

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type UpdateInstructorTrainingRequest

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 404 if training does not exist
			- 400 when validation of request fails
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerPatchTraining() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400
		var request UpdateTrainingRequest
		userID := g.MustGet("UserID").(uuid.UUID)
		var instructorID uuid.UUID

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &request, func() error {
				return request.Validate(ctx)
			},
		); err != nil {
			goto end
		}

		if instructorID, err = ctx.Instr.DalReadInstructorID(userID); err != nil {
			goto end
		}

		if err = ctx.DalUpdateTraining(&request, instructorID, userID); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		g.AbortWithError(status, err)
	}
}

/*
DELETE /api/training

DESC
	delete single instructor training offer

REQUEST
	- Authorization: Bearer <token> in header
	- application/json body of type DeleteInstructorTrainingRequest

RETURNS
	- status code:
		- 204 on success
		- error status code
			- 404 if training does not exist
			- 400 when validation of request fails
			- 500 on unexpected error
*/
func (ctx *Ctx) HandlerDeleteTraining() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		status := 400
		var request ObjectKey
		userID := g.MustGet("UserID").(uuid.UUID)

		if err = helpers.ReadJsonBodyFromReader(
			g.Request.Body, &request, request.Validate,
		); err != nil {
			goto end
		}

		if err = ctx.DalDeleteTraining(request.ID, userID); err != nil {
			if helpers.IsENF(err) {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		g.AbortWithStatus(204)
		return

	end:
		g.AbortWithError(status, err)
	}
}

//

/*
desc to be done <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

POST /api/training/img

*/
func (ctx *Ctx) HandlerPostTrainingImage() gin.HandlerFunc {
	return func(g *gin.Context) {

		var file *multipart.FileHeader
		var err error
		var status = 400
		userID := g.MustGet("UserID").(uuid.UUID)
		fpath := uuid.New().String()
		var ext = ""
		var fullpath = ""
		var trainingID uuid.UUID
		var trainingIDstr string
		var isMain bool
		var tr *TrainingWithJoins

		// if err := ctx.Request.ParseForm(); err != nil {
		// 	goto end
		// }

		file, err = g.FormFile("image")
		if err != nil {
			goto end
		}

		trainingIDstr = g.Request.FormValue("training_id")
		if trainingIDstr == "" {
			status = 404
			err = fmt.Errorf("no training_id was provided")
			goto end
		}

		if trainingID, err = uuid.Parse(trainingIDstr); err != nil {
			goto end
		}

		if tr, err = ctx.DalReadSingleTraining(
			DalReadTrainingsRequest{
				UserID:     &userID,
				TrainingID: &trainingID,
			}); err != nil {
			if err == sql.ErrNoRows {
				status = 404
			} else {
				status = 500
			}
			goto end
		}

		if len(tr.Training.SecondaryImgUrls) >= ctx.Config.MaxSecondaryTrImages {
			err = fmt.Errorf("user exceeded allowed number of training images")
			status = 429
			goto end
		}

		if g.Request.FormValue("main") != "" {
			isMain = true
		}

		if file.Size == 0 {
			err = fmt.Errorf("no image file was provided")
			goto end
		}

		if err = static.ValidateImage(file); err != nil {
			goto end
		}

		ext = path.Ext(file.Filename)

		if err = static.ValidateExt(ext); err != nil {
			goto end
		}

		fpath += ext
		fpath = path.Join(TrainingImgPath, fpath)
		fullpath = path.Join(
			ctx.Static.Config.Basepath, fpath)

		// path is taken from config, filename is {uuid}.{ext}
		if err = g.SaveUploadedFile(file, fullpath); err != nil {
			status = 500
			goto end
		}

		if isMain {

			if err = ctx.DalUpdateTrainingMainImage(trainingID, fpath); err != nil {
				goto end
			}

			// delete previous image
			pp := tr.Training.MainImgID
			if pp != "" {
				pp = path.Join(
					ctx.Static.Config.Basepath, pp)
				if err = os.Remove(pp); err != nil {
					// log error but dont interrupt the flow
					fmt.Println(err)
					err = nil
				}
			}

		} else {

			if err = ctx.DalAddTrainingSecondaryImage(trainingID, fpath); err != nil {
				goto end
			}

		}

		g.AbortWithStatus(204)
		return

	end:
		_ = g.AbortWithError(status, err)
	}
}

type DeleteTrainingImageRequest struct {
	ID         string
	TrainingID uuid.UUID
}

func (d *DeleteTrainingImageRequest) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("Validate DeleteTrainingImageRequest: invalid ID")
	}
	if d.TrainingID == uuid.Nil {
		return fmt.Errorf("Validate DeleteTrainingImageRequest: invalid TrainingID")
	}
	return nil
}

/*
desc to be done <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

DELETE /api/training/img

*/
func (ctx *Ctx) HandlerDeleteTrainingImage() gin.HandlerFunc {
	return func(g *gin.Context) {

		var request DeleteTrainingImageRequest
		userID := g.MustGet("UserID").(uuid.UUID)

		if err := helpers.ReadJsonBodyFromReader(
			g.Request.Body, &request, request.Validate,
		); err != nil {
			g.AbortWithError(400, err)
			return
		}

		t, err := ctx.DalReadSingleTraining(
			DalReadTrainingsRequest{
				UserID:     &userID,
				TrainingID: &request.TrainingID,
			})
		if err != nil {
			if err == sql.ErrNoRows {
				g.AbortWithError(404, err)
			} else {
				g.AbortWithError(500, err)
			}
			return
		}

		// find target image

		path := path.Join(
			ctx.Static.Config.Basepath, request.ID)

		if t.Training.MainImgID != "" && t.Training.MainImgID == request.ID {
			if err := ctx.DalUpdateTrainingMainImage(request.TrainingID, ""); err != nil {
				g.AbortWithError(500, err)
				return
			}
			if err := os.Remove(path); err != nil {
				g.AbortWithError(500, err)
				return
			}
			g.AbortWithStatus(204)
			return
		}

		for i := range t.Training.SecondaryImgIDs {
			if t.Training.SecondaryImgIDs[i] != "" &&
				t.Training.SecondaryImgIDs[i] == request.ID {
				if err := ctx.DalRemoveTrainingSecondaryImage(
					request.TrainingID, request.ID,
				); err != nil {
					g.AbortWithError(500, err)
					return
				}
				if err := os.Remove(path); err != nil {
					g.AbortWithError(500, err)
					return
				}
				g.AbortWithStatus(204)
				return
			}
		}

		g.AbortWithError(404, fmt.Errorf("update target not found"))
		return
	}
}
