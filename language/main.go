package language

type Locale int64

const (
	English Locale = iota
	Afrikaans
	Albanian
	ArabicAlgeria
	ArabicBahrain
	ArabicEgypt
	ArabicIraq
	ArabicJordan
	ArabicKuwait
	ArabicLebanon
	ArabicLibya
	ArabicMorocco
	ArabicOman
	ArabicQatar
	ArabicSaudiArabia
	ArabicSyria
	ArabicTunisia
	ArabicUAE
	ArabicYemen
	Basque
	Belarusian
	Bulgarian
	Catalan
	ChineseHongKong
	ChinesePRC
	ChineseSingapore
	ChineseTaiwan
	Croatian
	Czech
	Danish
	DutchBelgium
	DutchStandard
	EnglishAustralia
	EnglishBelize
	EnglishCanada
	EnglishIreland
	EnglishJamaica
	EnglishNewZealand
	EnglishSouthAfrica
	EnglishTrinidad
	EnglishUnitedKingdom
	EnglishUnitedStates
	Estonian
	Faeroese
	Farsi
	Finnish
	FrenchBelgium
	FrenchCanada
	FrenchLuxembourg
	FrenchStandard
	FrenchSwitzerland
	GaelicScotland
	GermanAustria
	GermanLiechtenstein
	GermanLuxembourg
	GermanStandard
	GermanSwitzerland
	Greek
	Hebrew
	Hindi
	Hungarian
	Icelandic
	Indonesian
	Irish
	ItalianStandard
	ItalianSwitzerland
	Japanese
	Korean
	Kurdish
	Latvian
	Lithuanian
	MacedonianFYROM
	Malayalam
	Malaysian
	Maltese
	Norwegian
	NorwegianBokmal
	NorwegianNynorsk
	Polish
	PortugueseBrazil
	PortuguesePortugal
	Punjabi
	RhaetoRomanic
	Romanian
	RomanianRepublicofMoldova
	Russian
	RussianRepublicofMoldova
	Serbian
	Slovak
	Slovenian
	Sorbian
	SpanishArgentina
	SpanishBolivia
	SpanishChile
	SpanishColombia
	SpanishCostaRica
	SpanishDominicanRepublic
	SpanishEcuador
	SpanishElSalvador
	SpanishGuatemala
	SpanishHonduras
	SpanishMexico
	SpanishNicaragua
	SpanishPanama
	SpanishParaguay
	SpanishPeru
	SpanishPuertoRico
	SpanishSpain
	SpanishUruguay
	SpanishVenezuela
	Swedish
	SwedishFinland
	Thai
	Tsonga
	Tswana
	Turkish
	Ukrainian
	Urdu
	Venda
	Vietnamese
	Welsh
	Xhosa
	Yiddish
	Zulu
)

func (l Locale) String() string {
	switch l {
	case Afrikaans:
		return "af"
	case Albanian:
		return "sq"
	case ArabicAlgeria:
		return "ar-dz"
	case ArabicBahrain:
		return "ar-bh"
	case ArabicEgypt:
		return "ar-eg"
	case ArabicIraq:
		return "ar-iq"
	case ArabicJordan:
		return "ar-jo"
	case ArabicKuwait:
		return "ar-kw"
	case ArabicLebanon:
		return "ar-lb"
	case ArabicLibya:
		return "ar-ly"
	case ArabicMorocco:
		return "ar-ma"
	case ArabicOman:
		return "ar-om"
	case ArabicQatar:
		return "ar-qa"
	case ArabicSaudiArabia:
		return "ar-sa"
	case ArabicSyria:
		return "ar-sy"
	case ArabicTunisia:
		return "ar-tn"
	case ArabicUAE:
		return "ar-ae"
	case ArabicYemen:
		return "ar-ye"
	case Basque:
		return "eu"
	case Belarusian:
		return "be"
	case Bulgarian:
		return "bg"
	case Catalan:
		return "ca"
	case ChineseHongKong:
		return "zh-hk"
	case ChinesePRC:
		return "zh-cn"
	case ChineseSingapore:
		return "zh-sg"
	case ChineseTaiwan:
		return "zh-tw"
	case Croatian:
		return "hr"
	case Czech:
		return "cs"
	case Danish:
		return "da"
	case DutchBelgium:
		return "nl-be"
	case DutchStandard:
		return "nl"
	case English:
		return "en"
	case EnglishAustralia:
		return "en-au"
	case EnglishBelize:
		return "en-bz"
	case EnglishCanada:
		return "en-ca"
	case EnglishIreland:
		return "en-ie"
	case EnglishJamaica:
		return "en-jm"
	case EnglishNewZealand:
		return "en-nz"
	case EnglishSouthAfrica:
		return "en-za"
	case EnglishTrinidad:
		return "en-tt"
	case EnglishUnitedKingdom:
		return "en-gb"
	case EnglishUnitedStates:
		return "en-us"
	case Estonian:
		return "et"
	case Faeroese:
		return "fo"
	case Farsi:
		return "fa"
	case Finnish:
		return "fi"
	case FrenchBelgium:
		return "fr-be"
	case FrenchCanada:
		return "fr-ca"
	case FrenchLuxembourg:
		return "fr-lu"
	case FrenchStandard:
		return "fr"
	case FrenchSwitzerland:
		return "fr-ch"
	case GaelicScotland:
		return "gd"
	case GermanAustria:
		return "de-at"
	case GermanLiechtenstein:
		return "de-li"
	case GermanLuxembourg:
		return "de-lu"
	case GermanStandard:
		return "de"
	case GermanSwitzerland:
		return "de-ch"
	case Greek:
		return "el"
	case Hebrew:
		return "he"
	case Hindi:
		return "hi"
	case Hungarian:
		return "hu"
	case Icelandic:
		return "is"
	case Indonesian:
		return "id"
	case Irish:
		return "ga"
	case ItalianStandard:
		return "it"
	case ItalianSwitzerland:
		return "it-ch"
	case Japanese:
		return "ja"
	case Korean:
		return "ko"
	case Kurdish:
		return "ku"
	case Latvian:
		return "lv"
	case Lithuanian:
		return "lt"
	case MacedonianFYROM:
		return "mk"
	case Malayalam:
		return "ml"
	case Malaysian:
		return "ms"
	case Maltese:
		return "mt"
	case Norwegian:
		return "no"
	case NorwegianBokmal:
		return "nb"
	case NorwegianNynorsk:
		return "nn"
	case Polish:
		return "pl"
	case PortugueseBrazil:
		return "pt-br"
	case PortuguesePortugal:
		return "pt"
	case Punjabi:
		return "pa"
	case RhaetoRomanic:
		return "rm"
	case Romanian:
		return "ro"
	case RomanianRepublicofMoldova:
		return "ro-md"
	case Russian:
		return "ru"
	case RussianRepublicofMoldova:
		return "ru-md"
	case Serbian:
		return "sr"
	case Slovak:
		return "sk"
	case Slovenian:
		return "sl"
	case Sorbian:
		return "sb"
	case SpanishArgentina:
		return "es-ar"
	case SpanishBolivia:
		return "es-bo"
	case SpanishChile:
		return "es-cl"
	case SpanishColombia:
		return "es-co"
	case SpanishCostaRica:
		return "es-cr"
	case SpanishDominicanRepublic:
		return "es-do"
	case SpanishEcuador:
		return "es-ec"
	case SpanishElSalvador:
		return "es-sv"
	case SpanishGuatemala:
		return "es-gt"
	case SpanishHonduras:
		return "es-hn"
	case SpanishMexico:
		return "es-mx"
	case SpanishNicaragua:
		return "es-ni"
	case SpanishPanama:
		return "es-pa"
	case SpanishParaguay:
		return "es-py"
	case SpanishPeru:
		return "es-pe"
	case SpanishPuertoRico:
		return "es-pr"
	case SpanishSpain:
		return "es"
	case SpanishUruguay:
		return "es-uy"
	case SpanishVenezuela:
		return "es-ve"
	case Swedish:
		return "sv"
	case SwedishFinland:
		return "sv-fi"
	case Thai:
		return "th"
	case Tsonga:
		return "ts"
	case Tswana:
		return "tn"
	case Turkish:
		return "tr"
	case Ukrainian:
		return "ua"
	case Urdu:
		return "ur"
	case Venda:
		return "ve"
	case Vietnamese:
		return "vi"
	case Welsh:
		return "cy"
	case Xhosa:
		return "xh"
	case Yiddish:
		return "ji"
	case Zulu:
		return "zu"
	}
	return "en"
}

func (l Locale) FromString(locale string) Locale {
	switch locale {
	case "af":
		return Afrikaans
	case "sq":
		return Albanian
	case "ar-dz":
		return ArabicAlgeria
	case "ar-bh":
		return ArabicBahrain
	case "ar-eg":
		return ArabicEgypt
	case "ar-iq":
		return ArabicIraq
	case "ar-jo":
		return ArabicJordan
	case "ar-kw":
		return ArabicKuwait
	case "ar-lb":
		return ArabicLebanon
	case "ar-ly":
		return ArabicLibya
	case "ar-ma":
		return ArabicMorocco
	case "ar-om":
		return ArabicOman
	case "ar-qa":
		return ArabicQatar
	case "ar-sa":
		return ArabicSaudiArabia
	case "ar-sy":
		return ArabicSyria
	case "ar-tn":
		return ArabicTunisia
	case "ar-ae":
		return ArabicUAE
	case "ar-ye":
		return ArabicYemen
	case "eu":
		return Basque
	case "be":
		return Belarusian
	case "bg":
		return Bulgarian
	case "ca":
		return Catalan
	case "zh-hk":
		return ChineseHongKong
	case "zh-cn":
		return ChinesePRC
	case "zh-sg":
		return ChineseSingapore
	case "zh-tw":
		return ChineseTaiwan
	case "hr":
		return Croatian
	case "cs":
		return Czech
	case "da":
		return Danish
	case "nl-be":
		return DutchBelgium
	case "nl":
		return DutchStandard
	case "en":
		return English
	case "en-au":
		return EnglishAustralia
	case "en-bz":
		return EnglishBelize
	case "en-ca":
		return EnglishCanada
	case "en-ie":
		return EnglishIreland
	case "en-jm":
		return EnglishJamaica
	case "en-nz":
		return EnglishNewZealand
	case "en-za":
		return EnglishSouthAfrica
	case "en-tt":
		return EnglishTrinidad
	case "en-gb":
		return EnglishUnitedKingdom
	case "en-us":
		return EnglishUnitedStates
	case "et":
		return Estonian
	case "fo":
		return Faeroese
	case "fa":
		return Farsi
	case "fi":
		return Finnish
	case "fr-be":
		return FrenchBelgium
	case "fr-ca":
		return FrenchCanada
	case "fr-lu":
		return FrenchLuxembourg
	case "fr":
		return FrenchStandard
	case "fr-ch":
		return FrenchSwitzerland
	case "gd":
		return GaelicScotland
	case "de-at":
		return GermanAustria
	case "de-li":
		return GermanLiechtenstein
	case "de-lu":
		return GermanLuxembourg
	case "de":
		return GermanStandard
	case "de-ch":
		return GermanSwitzerland
	case "el":
		return Greek
	case "he":
		return Hebrew
	case "hi":
		return Hindi
	case "hu":
		return Hungarian
	case "is":
		return Icelandic
	case "id":
		return Indonesian
	case "ga":
		return Irish
	case "it":
		return ItalianStandard
	case "it-ch":
		return ItalianSwitzerland
	case "ja":
		return Japanese
	case "ko":
		return Korean
	case "ku":
		return Kurdish
	case "lv":
		return Latvian
	case "lt":
		return Lithuanian
	case "mk":
		return MacedonianFYROM
	case "ml":
		return Malayalam
	case "ms":
		return Malaysian
	case "mt":
		return Maltese
	case "no":
		return Norwegian
	case "nb":
		return NorwegianBokmal
	case "nn":
		return NorwegianNynorsk
	case "pl":
		return Polish
	case "pt-br":
		return PortugueseBrazil
	case "pt":
		return PortuguesePortugal
	case "pa":
		return Punjabi
	case "rm":
		return RhaetoRomanic
	case "ro":
		return Romanian
	case "ro-md":
		return RomanianRepublicofMoldova
	case "ru":
		return Russian
	case "ru-md":
		return RussianRepublicofMoldova
	case "sr":
		return Serbian
	case "sk":
		return Slovak
	case "sl":
		return Slovenian
	case "sb":
		return Sorbian
	case "es-ar":
		return SpanishArgentina
	case "es-bo":
		return SpanishBolivia
	case "es-cl":
		return SpanishChile
	case "es-co":
		return SpanishColombia
	case "es-cr":
		return SpanishCostaRica
	case "es-do":
		return SpanishDominicanRepublic
	case "es-ec":
		return SpanishEcuador
	case "es-sv":
		return SpanishElSalvador
	case "es-gt":
		return SpanishGuatemala
	case "es-hn":
		return SpanishHonduras
	case "es-mx":
		return SpanishMexico
	case "es-ni":
		return SpanishNicaragua
	case "es-pa":
		return SpanishPanama
	case "es-py":
		return SpanishParaguay
	case "es-pe":
		return SpanishPeru
	case "es-pr":
		return SpanishPuertoRico
	case "es":
		return SpanishSpain
	case "es-uy":
		return SpanishUruguay
	case "es-ve":
		return SpanishVenezuela
	case "sv":
		return Swedish
	case "sv-fi":
		return SwedishFinland
	case "th":
		return Thai
	case "ts":
		return Tsonga
	case "tn":
		return Tswana
	case "tr":
		return Turkish
	case "ua":
		return Ukrainian
	case "ur":
		return Urdu
	case "ve":
		return Venda
	case "vi":
		return Vietnamese
	case "cy":
		return Welsh
	case "xh":
		return Xhosa
	case "ji":
		return Yiddish
	case "zu":
		return Zulu
	}
	return English
}
