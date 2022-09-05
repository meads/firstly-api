package api

import (
	"context"

	"github.com/meads/firstly-api/db"
	"github.com/pkg/errors"
)

// db   (Queries)
// ^
// |
// api  (FirstlyAPI)
// ^
// |
// http (FirstlyServer)

type Image struct {
	ID      int64  `json:"id"`
	Data    string `json:"data"`
	Created string `json:"created"`
	Deleted bool   `json:"deleted"`
}

type ImageAPI struct {
	ctx   context.Context
	store db.Querier
}

func NewImageAPI(ctx context.Context, store *db.Queries) *ImageAPI {
	return &ImageAPI{
		ctx:   ctx,
		store: store,
	}
}

func toDTO(img *db.Image) *Image {
	deleted := false
	if img.Deleted.Valid && img.Deleted.Int32 == 1 {
		deleted = true
	}
	return &Image{
		ID:      img.ID,
		Data:    img.Data,
		Created: img.Created,
		Deleted: deleted,
	}
}

func (api *ImageAPI) CreateImage(data string) (*Image, error) {
	img, err := api.store.CreateImage(api.ctx, data)
	if err != nil {
		return nil, errors.Wrap(err, "api error creating image")
	}
	return toDTO(&img), nil
}

func (api *ImageAPI) DeleteImage(id int64) error {
	err := api.store.DeleteImage(api.ctx, id)
	if err != nil {
		return errors.Wrap(err, "api error deleting image")
	}
	return nil
}

func (api *ImageAPI) GetImage(id int64) (*Image, error) {
	img, err := api.store.GetImage(api.ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "api error retrieving image")
	}
	return toDTO(&img), nil
}

func (api *ImageAPI) ListImages() ([]*Image, error) {
	imgs, err := api.store.ListImages(api.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "api error retrieving images")
	}

	imagesToReturn := []*Image{}
	for i, _ := range imgs {
		imagesToReturn = append(imagesToReturn, toDTO(&imgs[i]))
	}
	return imagesToReturn, nil
}

func (api *ImageAPI) SoftDeleteImage(id int64) error {
	err := api.store.SoftDeleteImage(api.ctx, id)
	if err != nil {
		return errors.Wrap(err, "api error deleting image")
	}
	return nil
}
