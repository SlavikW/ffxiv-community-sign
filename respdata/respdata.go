package respdata

type FFXIVCode struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

type PostsList struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Rows []struct {
			PostsId            string  `json:"posts_id"`
			PartName           string  `json:"part_name"`
			PartParentName     string  `json:"part_parent_name"`
			PartParentId       string  `json:"part_parent_id"`
			Uuid               string  `json:"uuid"`
			Avatar             string  `json:"avatar"`
			TestLimitedBadge   string  `json:"test_limited_badge"`
			Posts2CreatorBadge string  `json:"posts2_creator_badge"`
			AdminTag           string  `json:"admin_tag"`
			CharacterName      string  `json:"character_name"`
			AreaId             int     `json:"area_id"`
			AreaName           string  `json:"area_name"`
			GroupName          string  `json:"group_name"`
			Title              string  `json:"title"`
			CoverPic           string  `json:"cover_pic"`
			Type               int     `json:"type"`
			PartId             int     `json:"part_id"`
			ContentPre         string  `json:"content_pre"`
			CreatedAt          string  `json:"created_at"`
			LastCommentTime    *string `json:"last_comment_time"`
			SortUpdatedTime    string  `json:"sort_updated_time"`
			CommentCount       string  `json:"comment_count"`
			LikeCount          string  `json:"like_count"`
			StarCount          string  `json:"star_count"`
			ReadCount          string  `json:"read_count"`
			IsTop              int     `json:"is_top"`
			IsRefine           int     `json:"is_refine"`
			RelayCount         string  `json:"relay_count"`
		} `json:"rows"`
		PageTime string `json:"pageTime"`
	} `json:"data"`
}
