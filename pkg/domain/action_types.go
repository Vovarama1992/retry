package domain

type CheckedActionTypeName string

const (
	ActionClickCtaTop                    CheckedActionTypeName = "click_cta_top"
	ActionClickCtaBottom                 CheckedActionTypeName = "click_cta_bottom"
	ActionScrollDepth                    CheckedActionTypeName = "scroll_depth"
	ActionScrollSectionView              CheckedActionTypeName = "scroll_section_view"
	ActionExternalLinkGarageRaid         CheckedActionTypeName = "external_link_garage_raid"
	ActionExternalLinkRaids              CheckedActionTypeName = "external_link_raids"
	ActionExternalLinkDetails            CheckedActionTypeName = "external_link_details"
	ActionGalleryScrollRight             CheckedActionTypeName = "gallery_scroll_right"
	ActionGalleryScrollLeft              CheckedActionTypeName = "gallery_scroll_left"
	ActionExternalLinkMentorPage         CheckedActionTypeName = "external_link_mentor_page"
	ActionExternalLinkSocial             CheckedActionTypeName = "external_link_social"
	ActionFaqOpenAnswer                  CheckedActionTypeName = "faq_open_answer"
	ActionClickLinksBuyAccess            CheckedActionTypeName = "click_links_buy_access"
	ActionClickLinksTelegram             CheckedActionTypeName = "click_links_telegram"
	ActionClickLinksYoutubeEntertainment CheckedActionTypeName = "click_links_youtube_entertainment"
	ActionClickLinksYoutubeStreams       CheckedActionTypeName = "click_links_youtube_streams"
	ActionClickLinksInstagram            CheckedActionTypeName = "click_links_instagram"
	ActionClickLinksTiktok               CheckedActionTypeName = "click_links_tiktok"
	ActionClickLinksArtstation           CheckedActionTypeName = "click_links_artstation"
	ActionClickLinks3DGuide              CheckedActionTypeName = "click_links_3d_guide"
)

var ValidActionTypes = map[CheckedActionTypeName]struct{}{
	ActionClickCtaTop:                    {},
	ActionClickCtaBottom:                 {},
	ActionScrollDepth:                    {},
	ActionScrollSectionView:              {},
	ActionExternalLinkGarageRaid:         {},
	ActionExternalLinkRaids:              {},
	ActionExternalLinkDetails:            {},
	ActionGalleryScrollRight:             {},
	ActionGalleryScrollLeft:              {},
	ActionExternalLinkMentorPage:         {},
	ActionExternalLinkSocial:             {},
	ActionFaqOpenAnswer:                  {},
	ActionClickLinksBuyAccess:            {},
	ActionClickLinksTelegram:             {},
	ActionClickLinksYoutubeEntertainment: {},
	ActionClickLinksYoutubeStreams:       {},
	ActionClickLinksInstagram:            {},
	ActionClickLinksTiktok:               {},
	ActionClickLinksArtstation:           {},
	ActionClickLinks3DGuide:              {},
}

func IsValidActionType(t string) bool {
	_, ok := ValidActionTypes[CheckedActionTypeName(t)]
	return ok
}
