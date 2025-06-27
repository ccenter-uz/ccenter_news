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
}

type BannerGetAllRes struct {
	Banners []BannerRes `json:"banners"`
	Count   int         `json:"count"`
}

type DeleteImage struct {
	ImgUrl string `json:"img_url"`
}

type Image struct {
	ImgUrl string `json:"img_url"`
}
type ListImages struct {
	Images []Image `json:"images"`
	Count  int     `json:"count"`
}
