package types

const (
	Instagram = "instagram"
	Highlight = "highlight"
	Story     = "story"
	TikTok    = "tiktok"
	VSCO      = "vsco"
)

var MediaTypes = []string{Highlight, Instagram, Story, TikTok, VSCO}

func ValidMediaType(media string) bool {
	return media == Instagram || media == Highlight || media == Story || media == VSCO || media == TikTok
}

func ValidNetworkType(media string) bool {
	return media == Instagram || media == VSCO || media == TikTok
}
