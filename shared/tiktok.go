package shared

const (
	TikTokScriptID = "__NEXT_DATA__"
)

type TikTokPost struct {
	Props struct {
		PageProps struct {
			ItemInfo struct {
				ItemStruct struct {
					Author struct {
						UniqueID string `json:"uniqueId"`
					} `json:"author"`
					Video struct {
						DownloadAddress string `json:"downloadAddr"`
					} `json:"video"`
				} `json:"itemStruct"`
			} `json:"itemInfo"`
		} `json:"pageProps"`
	} `json:"props"`
}
