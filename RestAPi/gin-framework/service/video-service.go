package service

import "github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/gin-framework/entity"

type VideoService interface {
	FindAll() []entity.Video
	Save(entity.Video) entity.Video
}

type videoService struct {
	videos []entity.Video
}

func New() VideoService {
	return &videoService{}
}

func (sv *videoService) FindAll() []entity.Video {
	return sv.videos
}

func (sv *videoService) Save(video entity.Video) entity.Video {
	sv.videos = append(sv.videos, video)
	return video
}
