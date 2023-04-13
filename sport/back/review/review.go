package review

import (
	"fmt"
	"sport/user"
	"time"

	"github.com/google/uuid"
)

type ReviewContent struct {
	Mark   int
	Review string
}

func (or *ReviewContent) Validate() error {
	if or.Mark < 0 || or.Mark > 6 {
		return fmt.Errorf("validate ReviewContent: mark must be between 0 and 6")
	}
	if (or.Mark == 0) && (or.Review == "") {
		return fmt.Errorf("validate ReviewContent: either mark or review must be provided")
	}
	if len(or.Review) > MaxReviewLen {
		return fmt.Errorf("validate ReviewContent: invalid ReviewLen")
	}
	return nil
}

type UpdateReviewRequest struct {
	ReviewContent
	AccessToken string
}

func (ur *UpdateReviewRequest) Validate() error {
	if ur.AccessToken == "" {
		return fmt.Errorf("Validate UpdateReviewRequest: invalid AccessToken")
	}
	return ur.ReviewContent.Validate()
}

type OriginType string

const (
	ReviewOriginUser OriginType = "user"
	ReviewOriginRsv  OriginType = "rsv"
)

// create request
type ReviewRequest struct {
	TrainingID *uuid.UUID
	RsvID      uuid.UUID
	UserID     *uuid.UUID
	Email      string
	UserInfo   user.PubUserInfo
}

// table
type Review struct {
	ID          uuid.UUID
	AccessToken string
	CreatedOn   time.Time

	ReviewRequest

	ReviewContent
}

// review info available for all clients
type PubReview struct {
	CreatedOn time.Time
	UserInfo  user.PubUserInfo
	ReviewContent
}

type ReviewType string

const (
	TokenReviewType   ReviewType = "token"
	ContentReviewType ReviewType = "content"
	AnyReviewType     ReviewType = "any"
)

// reviews for specified item which are available for anyone to read
type ReviewContentResponse struct {
	ID       uuid.UUID
	UserInfo user.PubUserInfo
	ReviewContent
}

// can be used by user to make actual review
type ReviewTokenResponse struct {
	ExpireOn    time.Time
	AccessToken string
}

// gets returned to user after making GET /review/user request
// One of those fields will be filled depending on Type field value
type ReviewResponse struct {
	Type    ReviewType
	Content *ReviewContentResponse
	Token   *ReviewTokenResponse
}

type DeleteReviewRequest struct {
	ID uuid.UUID
}

func (d *DeleteReviewRequest) Validate() error {
	if d.ID == uuid.Nil {
		return fmt.Errorf("validate DeleteReviewRequest: invalid id")
	}
	return nil
}

func (h *ReviewRequest) NewOpinion() *Review {
	return &Review{
		ReviewRequest: *h,
		ID:            uuid.New(),
		AccessToken:   uuid.New().String(),
		CreatedOn:     time.Now().In(time.UTC),
		ReviewContent: ReviewContent{
			Mark:   0,
			Review: "",
		},
	}
}
