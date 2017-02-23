package global

import (
	"crypto/md5"
	b64 "encoding/base64"
)

type UserInfo struct {
	API       string
	Name      string
	Org       string
	OrgGUID   string
	OrgURL    string
	Space     string
	SpaceGUID string
	SpaceURL  string
}

func (info UserInfo) IsValid() bool {
	return info.API != "" && info.Name != "" && info.Org != "" &&
		info.OrgURL != "" && info.Space != "" && info.SpaceURL != ""
}

func Md5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return encodeData(hash[:])
}

func encodeData(data []byte) string {
	return b64.URLEncoding.WithPadding(b64.NoPadding).EncodeToString(data)
}

func (info UserInfo) GetAPIHash() string {
	return Md5Hash(info.API)
}

func (info UserInfo) GetNameHash() string {
	return Md5Hash(info.Name)
}

func (info UserInfo) GetOrgHash() string {
	return Md5Hash(info.Org)
}

func (info UserInfo) GetSpaceHash() string {
	return Md5Hash(info.Space)
}

var CurrentUserInfo *UserInfo
