package lang

type LangMap map[string]string

func NewLangMap(pl string, en string) LangMap {
	ret := make(LangMap)
	ret["pl"] = pl
	ret["en"] = en
	return ret
}

type LocaleStruct struct {
	UnreadMsgsFmt  LangMap
	NewChatroomFmt LangMap
}

var Locale = LocaleStruct{
	UnreadMsgsFmt:  NewLangMap("Nowe wiadomośći w '%s'", "New messages in '%s'"),
	NewChatroomFmt: NewLangMap("Nowy kanał: '%s'", "New channel: '%s'"),
}
