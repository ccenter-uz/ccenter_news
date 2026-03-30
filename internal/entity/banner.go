package entity

type BannerCreate struct {
	Text     MultilingualField `json:"text"`
	Title    MultilingualField `json:"title"`
	Label    MultilingualField `json:"label"`
	Date     string            `json:"date"`
	ImgUrl   string            `json:"img_url"`
	FileLink string            `json:"file_link"`
	HrefName string            `json:"href_name"`
	Type     string            `json:"type"`
	Order    int               `json:"order"`
	Markdown MultilingualField `json:"markdown"`
}

type BannerUpdate struct {
	Id       string            `json:"id"`
	Text     MultilingualField `json:"text"`
	Title    MultilingualField `json:"title"`
	Label    MultilingualField `json:"label"`
	Date     string            `json:"date"`
	ImgUrl   string            `json:"img_url"`
	FileLink string            `json:"file_link"`
	HrefName string            `json:"href_name"`
	Type     string            `json:"type"`
	Order    int               `json:"order"`
	Markdown MultilingualField `json:"markdown"`
}

type BannerRes struct {
	Id        string            `json:"id"`
	Text      MultilingualField `json:"text"`
	Title     MultilingualField `json:"title"`
	Label     MultilingualField `json:"label"`
	Date      string            `json:"date"`
	ImgUrl    string            `json:"img_url"`
	FileLink  *string           `json:"file_link"`
	HrefName  *string           `json:"href_name"`
	Type      string            `json:"type"`
	Order     int               `json:"order"`
	CreatedAt string            `json:"created_at"`
	Markdown  MultilingualField `json:"markdown"`
}

type BannerGetAllRes struct {
	Banners []BannerRes `json:"banners"`
	Count   int         `json:"count"`
}

type DeleteFile struct {
	Url string `json:"url"`
}

type Url struct {
	Url string `json:"url"`
}
