package db

type Pic struct {
	FriendPic []FriendPic `json:"FriendPic"`
	Tips      string      `json:"Tips"`
}
type FriendPic struct {
	FileMd5  string `json:"FileMd5"`
	FileSize int    `json:"FileSize"`
	Path     string `json:"Path"`
	URL      string `json:"Url"`
}
type Voice struct {
	Tips string `json:"Tips"`
	URL  string `json:"Url"`
}
