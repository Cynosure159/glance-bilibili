# bilibili 请求视频列表接口

### 查询用户投稿视频明细

> https://api.bilibili.com/x/space/wbi/arc/search

> ~~https://api.bilibili.com/x/space/arc/search~~ （已废弃，保留是方便遇到问题的人搜索到此处）

*请求方式：GET*

鉴权方式：[Wbi 签名](../misc/sign/wbi.md)

另见 [根据关键词查找视频](../video/collection.md#根据关键词查找视频), 功能基本相同, 暂未发现风控校验

**url参数：**

| 参数名  | 类型 | 内容         | 必要性 | 备注                                                                          |
| ------- | ---- | ------------ | ------ | ----------------------------------------------------------------------------- |
| mid     | num  | 目标用户mid  | 必要   |                                                                               |
| order   | str  | 排序方式     | 非必要 | 默认为pubdate<br />最新发布：pubdate<br />最多播放：click<br />最多收藏：stow |
| tid     | num  | 筛选目标分区 | 非必要 | 默认为0<br />0：不进行分区筛选<br />分区tid为所筛选的分区                     |
| keyword | str  | 关键词筛选   | 非必要 | 用于使用关键词搜索该UP主视频稿件                                              |
| pn      | num  | 页码         | 非必要 | 默认为 `1`                                                                    |
| ps      | num  | 每页项数     | 非必要 | 默认为 `30`                                                                   |

**json回复：**

根对象：

| 字段    | 类型 | 内容     | 备注                        |
| ------- | ---- | -------- | --------------------------- |
| code    | num  | 返回值   | 0：成功<br />-400：请求错误<br />-412：请求被拦截 |
| message | str  | 错误信息 | 默认为0                     |
| ttl     | num  | 1        |                 |
| data    | obj  | 信息本体 |                             |

`data`对象：

| 字段            | 类型 | 内容           | 备注 |
| --------------- | ---- | -------------- | ---- |
| list            | obj  | 列表信息       |      |
| page            | obj  | 页面信息       |      |
| episodic_button | obj  | “播放全部“按钮 |      |
| is_risk         | bool |                |      |
| gaia_res_type   | num  |                |      |
| gaia_data       | obj  |                |      |

`data`中的`list`对象：

| 字段  | 类型   | 内容             | 备注 |
| ----- | ------ | ---------------- | ---- |
| slist | array  | 空数组           |      |
| tlist | obj    | 投稿视频分区索引 |      |
| vlist | array | 投稿视频列表     |      |

`list`中的`tlist`对象：

| 字段  | 类型 | 内容         | 备注                  |
| ----- | ---- | ------------ | --------------------- |
| {tid} | obj  | 该分区的详情 | 字段名为存在的分区tid |
| ……    | obj  | ……           | 向下扩展              |

`tlist`中的`{tid}`对象：

| 字段  | 类型 | 内容                 | 备注 |
| ----- | ---- | -------------------- | ---- |
| count | num  | 投稿至该分区的视频数 |      |
| name  | str  | 该分区名称           |      |
| tid   | num  | 该分区tid            |      |

`list`中的`vlist`数组：

| 项   | 类型 | 内容            | 备注 |
| ---- | ---- | --------------- | ---- |
| 0    | obj  | 投稿视频1       |      |
| n    | obj  | 投稿视频（n+1） |      |
| ……   | obj  | ……              | ……   |

`list`中的`vlist`数组中的对象：

| 字段               | 类型 | 内容           | 备注                         |
| ------------------ | ---- | -------------- | ---------------------------- |
| aid                | num  | 稿件avid       |                              |
| attribute          | num  |                |                              |
| author             | str  | 视频UP主       | 不一定为目标用户（合作视频） |
| bvid               | str  | 稿件bvid       |                              |
| comment            | num  | 视频评论数     |                              |
| copyright          | str  | 视频版权类型   |                              |
| created            | num  | 投稿时间       | 时间戳                       |
| description        | str  | 视频简介       |                              |
| elec_arc_type      | num  | 充电为1，否则0 | 可能还有其他情况             |
| enable_vt          | num  | 0              | 作用尚不明确                 |
| hide_click         | bool | false          | 作用尚不明确                 |
| is_avoided         | num  | 0              | 作用尚不明确                 |
| is_charging_arc    | bool | 是否为充电视频 |                              |
| is_lesson_video    | num  | 是否为课堂视频 | 0：否<br />1：是             |
| is_lesson_finished | num  | 课堂是否已完结 | 0：否<br />1：是             |
| is_live_playback   | num  | 是否为直播回放 | 0：否<br />1：是             |
| is_pay             | num  | 0              | 作用尚不明确                 |
| is_self_view       | bool | 是否仅自己可见 |                           |
| is_steins_gate     | num  | 是否为互动视频 | 0：否<br />1：是             |
| is_union_video     | num  | 是否为合作视频 | 0：否<br />1：是             |
| jump_url           | str  | 跳转链接       | 跳转到课堂的链接，否则为""   |
| length             | str  | 视频长度       | MM:SS                        |
| mid                | num  | 视频UP主mid    | 不一定为目标用户（合作视频） |
| meta               | obj  | 所属合集或课堂 | 无数据时为 null              |
| pic                | str  | 视频封面       |                              |
| play               | num  | 视频播放次数   |                              |
| playback_position  | num  | 百分比播放进度 | 封面下方显示的粉色条         |
| review             | num  | 0              | 作用尚不明确                 |
| season_id          | num  | 合集或课堂编号 | 都不属于时为0                |
| subtitle           | str  | 空             | 作用尚不明确                 |
| title              | str  | 视频标题       |                              |
| typeid             | num  | 视频分区tid    |                              |
| video_review       | num  | 视频弹幕数     |                              |
| vt                 | num  | 0              | 作用尚不明确                 |
| vt_display         | str  | 空             | 作用尚不明确                 |

`list`中的`vlist`数组中的对象中的`meta`对象：

| 字段       | 类型 | 内容         | 备注             |
| ---------- | ---- | ------------ | ---------------- |
| attribute  | num  | 0            | 作用尚不明确     |
| cover      | str  | 合集封面URL  |                  |
| ep_count   | num  | 合集视频数量 |                  |
| ep_num     | num  | 合集视频数量 |                  |
| first_aid  | num  | 首个视频av号 |                  |
| id         | num  | 合集id       |                  |
| intro      | str  | 合集介绍     |                  |
| mid        | num  | UP主uid      | 若为课堂，则为0  |
| ptime      | num  | unix时间(s)  | 最后更新时间     |
| sign_state | num  | 0            | 作用尚不明确     |
| stat       | obj  | 合集统计数据 |                  |
| title      | str  | 合集名称     |                  |

`list`中的`vlist`数组中的对象中的`meta`对象中的`stat`对象：

| 字段       | 类型 | 内容         | 备注                 |
| ---------- | ---- | ------------ | -------------------- |
| coin       | num  | 合集总投币数 |                      |
| danmaku    | num  | 合集总弹幕数 |                      |
| favorite   | num  | 合集总收藏数 |                      |
| like       | num  | 合集总点赞数 |                      |
| mtime      | num  | unix时间(s)  | 其他统计数据更新时间 |
| reply      | num  | 合集总评论数 |                      |
| season_id  | num  | 合集id       |                      |
| share      | num  | 合集总分享数 |                      |
| view       | num  | 合集总播放量 |                      |
| vt         | num  | 0            | 作用尚不明确         |
| vv         | num  | 0            | 作用尚不明确         |

`data`中的`page`对象：

| 字段  | 类型 | 内容       | 备注 |
| ----- | ---- | ---------- | ---- |
| count | num  | 总计稿件数 |      |
| pn    | num  | 当前页码   |      |
| ps    | num  | 每页项数   |      |

`data`中的`episodic_button`对象：

| 字段 | 类型 | 内容          | 备注 |
| ---- | ---- | ------------- | ---- |
| text | str  | 按钮文字      |      |
| uri  | str  | 全部播放页url |      |

**示例：**

`pn`（页码）和`ps`（每页项数）只改变`vlist`中成员的多少与内容

以每页2项查询用户`mid=53456`的第1页投稿视频明细

```shell
curl -G 'https://api.bilibili.com/x/space/arc/search' \
--data-urlencode 'mid=53456' \
--data-urlencode 'ps=2' \
--data-urlencode 'pn=1'
```

<details>
<summary>查看响应示例：</summary>

```json
{
	"code": 0,
	"message": "0",
	"ttl": 1,
	"data": {
		"list": {
			"slist": [],
			"tlist": {
				"1": {
					"tid": 1,
					"count": 3,
					"name": "动画"
				},
				"129": {
					"tid": 129,
					"count": 1,
					"name": "舞蹈"
				},
				"160": {
					"tid": 160,
					"count": 96,
					"name": "生活"
				},
				"177": {
					"tid": 177,
					"count": 4,
					"name": "纪录片"
				},
				"181": {
					"tid": 181,
					"count": 50,
					"name": "影视"
				},
				"188": {
					"tid": 188,
					"count": 444,
					"name": "科技"
				},
				"196": {
					"tid": 196,
					"count": 2,
					"name": "课堂"
				}
			},
			"vlist": [{
				"comment": 985,
				"typeid": 250,
				"play": 224185,
				"pic": "http://i0.hdslb.com/bfs/archive/5e56c10a9bd67f2fcac46fdd0fc2caa8769700c8.jpg",
				"subtitle": "",
				"description": "这一次，我们的样片日记首次来到了西藏，在桃花季开启了藏东样片之旅！这趟“开荒”之旅我们跋山涉水，一路硬刚，多亏有路虎卫士这样的神队友撑全场！这次的素材我们也上传到了官网（ysjf.com/material），欢迎大家去看看~如果你喜欢这期视频，请多多支持我们，并把视频分享给你的朋友们一起看看！",
				"copyright": "1",
				"title": "和朋友去西藏拍样片日记……",
				"review": 0,
				"author": "影视飓风",
				"mid": 946974,
				"created": 1745290800,
				"length": "22:11",
				"video_review": 2365,
				"aid": 114375683741573,
				"bvid": "BV1ac5yzhE94",
				"hide_click": false,
				"is_pay": 0,
				"is_union_video": 1,
				"is_steins_gate": 0,
				"is_live_playback": 0,
				"is_lesson_video": 0,
				"is_lesson_finished": 0,
				"lesson_update_info": "",
				"jump_url": "",
				"meta": {
					"id": 2046621,
					"title": "样片日记",
					"cover": "https://archive.biliimg.com/bfs/archive/e2ca3e5a6672cf35c9e61ac02e8d739cc0aafa8b.jpg",
					"mid": 946974,
					"intro": "",
					"sign_state": 0,
					"attribute": 140,
					"stat": {
						"season_id": 2046621,
						"view": 31755096,
						"danmaku": 171253,
						"reply": 33685,
						"favorite": 409505,
						"coin": 935105,
						"share": 199467,
						"like": 1791607,
						"mtime": 1745309513,
						"vt": 0,
						"vv": 0
					},
					"ep_count": 13,
					"first_aid": 238588630,
					"ptime": 1745290800,
					"ep_num": 13
				},
				"is_avoided": 0,
				"season_id": 2046621,
				"attribute": 16793984,
				"is_charging_arc": false,
				"elec_arc_type": 0,
				"vt": 0,
				"enable_vt": 0,
				"vt_display": "",
				"playback_position": 0,
				"is_self_view": false
			}, {
				"comment": 0,
				"typeid": 197,
				"play": 8506,
				"pic": "https://archive.biliimg.com/bfs/archive/489f3df26a190a152ad479bfe50a73f1cd4c43c5.jpg",
				"subtitle": "",
				"description": "8节课，Tim和青青带你用iPhone拍出电影感",
				"copyright": "1",
				"title": "【影视飓风】只看8节课，用iPhone拍出电影感",
				"review": 0,
				"author": "影视飓风",
				"mid": 946974,
				"created": 1744865737,
				"length": "00:00",
				"video_review": 9,
				"aid": 114351440726681,
				"bvid": "BV1WB5ezxEnz",
				"hide_click": false,
				"is_pay": 0,
				"is_union_video": 0,
				"is_steins_gate": 0,
				"is_live_playback": 0,
				"is_lesson_video": 1,
				"is_lesson_finished": 1,
				"lesson_update_info": "8",
				"jump_url": "https://www.bilibili.com/cheese/play/ss190402215",
				"meta": {
					"id": 190402215,
					"title": "【影视飓风】只看8节课，用iPhone拍出电影感",
					"cover": "https://archive.biliimg.com/bfs/archive/489f3df26a190a152ad479bfe50a73f1cd4c43c5.jpg",
					"mid": 0,
					"intro": "",
					"sign_state": 0,
					"attribute": 0,
					"stat": {
						"season_id": 190402215,
						"view": 1111222,
						"danmaku": 1853,
						"reply": 0,
						"favorite": 0,
						"coin": 0,
						"share": 0,
						"like": 0,
						"mtime": 0,
						"vt": 0,
						"vv": 0
					},
					"ep_count": 0,
					"ptime": 1744865737,
					"ep_num": 0
				},
				"is_avoided": 0,
				"season_id": 190402215,
				"attribute": 1073758592,
				"is_charging_arc": false,
				"elec_arc_type": 0,
				"vt": 0,
				"enable_vt": 0,
				"vt_display": "",
				"playback_position": 0,
				"is_self_view": false
			}]
		},
		"page": {
			"pn": 1,
			"ps": 42,
			"count": 786
		},
		"episodic_button": {
			"text": "播放全部",
			"uri": "//www.bilibili.com/medialist/play/946974?from=space"
		},
		"is_risk": false,
		"gaia_res_type": 0,
		"gaia_data": null
	}
}
```

</details>
