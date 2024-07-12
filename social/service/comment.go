package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/social/entity"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/valyala/fasthttp"
)

var seqname = "seq_comment"

// CommentService 同步表结构，服务继承基本服务的方法
type CommentService struct {
	service.OrmBaseService
}

var commentService = &CommentService{}

func GetCommentService() *CommentService {
	return commentService
}

func (svc *CommentService) GetSeqName() string {
	return seqname
}

func (svc *CommentService) NewEntity(data []byte) (interface{}, error) {
	event := &entity.Comment{}
	if data == nil {
		return event, nil
	}
	err := message.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}

	return event, err
}

func (svc *CommentService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Comment, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func FastGet(urlStr string, requestParam interface{}) ([]byte, error) {
	args := stock.ReflectFields(requestParam)
	uri := urlStr + "?" + args
	status, resp, err := fasthttp.Get(nil, uri)
	if err != nil {
		logger.Sugar.Errorf("fasthttp Get fail:", err.Error())
		status, resp, err = fasthttp.Get(nil, uri)
		if err != nil {
			logger.Sugar.Errorf("second fasthttp Get fail:", err.Error())
			return nil, err
		}
	}
	if status != fasthttp.StatusOK {
		logger.Sugar.Errorf("fasthttp Get status:%s", status)
		return nil, errors.New(fmt.Sprint(status))
	}

	return resp, nil
}

// GetComments 获取某个视频的评论数据
func (svc *CommentService) GetComments(awemeId string, from int, limit int) ([]*entity.Comment, error) {
	//head = {
	//	'accept': 'application/json, text/plain, */*',
	//		'accept-encoding': 'gzip, deflate, br',
	//		'accept-language': 'zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7',
	//		'cookie': '换成自己的cookie值',
	//		'referer': 'https://www.douyin.com/',
	//		'sec-ch-ua': '"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"',
	//		'sec-ch-ua-mobile': '?0',
	//		'sec-ch-ua-platform': '"macOS"',
	//		'sec-fetch-dest': 'empty',
	//		'sec-fetch-mode': 'cors',
	//		'sec-fetch-site': 'same-origin',
	//		'user-agent': ua,
	//}
	//params := {
	//	'device_platform': 'webapp',
	//		'aid': 6383,
	//		'channel': 'channel_pc_web',
	//		'aweme_id': video_id,  // 视频id
	//	'cursor': page * 20,
	//		'count': 20,
	//		'item_type': 0,
	//		'insert_ids': '',
	//		'rcFT': '',
	//		'pc_client_type': 1,
	//		'version_code': '170400',
	//		'version_name': '17.4.0',
	//		'cookie_enabled': 'true',
	//		'screen_width': 1440,
	//		'screen_height': 900,
	//		'browser_language': 'zh-CN',
	//		'browser_platform': 'MacIntel',
	//		'browser_name': 'Chrome',
	//		'browser_version': '109.0.0.0',
	//		'browser_online': 'true',
	//		'engine_name': 'Blink',
	//		'engine_version': '109.0.0.0',
	//		'os_name': 'Mac OS',
	//		'os_version': '10.15.7',
	//		'cpu_core_num': 4,
	//		'device_memory': 8,
	//		'platform': 'PC',
	//		'downlink': 1.5,
	//		'effective_type': '4g',
	//		'round_trip_time': 150,
	//		'webid': 7184233910711879229,
	//		'msToken': 'LZ3nJ12qCwmFPM1NgmgYAz73RHVG_5ytxc_EMHr_3Mnc9CxfayXlm2kbvRaaisoAdLjRVPdLx5UDrc0snb5UDyQVRdGpd3qHgk64gLh6Tb6lR16WG7VHZQ==',
	//}
	resp, err := FastGet("https://www.douyin.com/aweme/v1/web/comment/list/", nil)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	comments := make([]*entity.Comment, 0)
	c := &entity.Comment{}
	err = json.Unmarshal(resp, c)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	comments = append(comments, c)

	return comments, nil
}

func init() {
	err := service.GetSession().Sync(new(entity.Comment))
	if err != nil {
		return
	}
	commentService.OrmBaseService.GetSeqName = commentService.GetSeqName
	commentService.OrmBaseService.FactNewEntity = commentService.NewEntity
	commentService.OrmBaseService.FactNewEntities = commentService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("comment", commentService)
}
